package proxycheck

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Dreamacro/clash/adapters/outbound"
	"github.com/Dreamacro/clash/constant"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Luoxin/Eutamias/utils/json"
	"github.com/Luoxin/faker"
	"github.com/alitto/pond"
	log "github.com/sirupsen/logrus"
	"github.com/thedevsaddam/retry"

	"github.com/Luoxin/Eutamias/proxy"
)

// https://github.com/letsfire/factory
// https://github.com/alitto/pond

type ProxyCheck struct {
	maxSize int
	faker   *faker.Faker

	checkUrl string
	pool     *pond.WorkerPool
	timeout  time.Duration
}

//go:generate pie ResultList.*
type ResultList []*Result

type Result struct {
	ProxyUrl     string
	Delay, Speed float64
	Err          error
}

func NewProxyCheck() *ProxyCheck {
	return &ProxyCheck{
		maxSize:  10,
		faker:    faker.New(),
		checkUrl: "https://www.google.com",
		timeout:  time.Minute,
	}
}

func (p *ProxyCheck) Init(maxWorker int) error {
	if maxWorker == 0 {
		maxWorker = 10
	}

	p.pool = pond.New(maxWorker, maxWorker)
	return nil
}

func (p *ProxyCheck) SetCheckUrl(checkUrl string) {
	p.checkUrl = checkUrl
}

func (p *ProxyCheck) SetTimeout(timeout time.Duration) {
	p.timeout = timeout
}

func (p *ProxyCheck) AddWithClash(nodeUrl string, logic func(result Result) error) error {
	p.pool.Submit(func() {
		delay, speed, err := p.CheckWithClash(nodeUrl)
		err = retry.DoFunc(5, 500*time.Millisecond, func() error {
			err = logic(Result{
				ProxyUrl: nodeUrl,
				Delay:    delay,
				Speed:    speed,
				Err:      err,
			})
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			return nil
		})
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	})

	return nil
}

func (p *ProxyCheck) AddWithLink(nodeUrl string, logic func(result Result) error) error {
	p.pool.Submit(func() {
		delay, speed, err := p.CheckWithLink(nodeUrl)
		err = retry.DoFunc(5, 500*time.Millisecond, func() error {
			err = logic(Result{
				ProxyUrl: nodeUrl,
				Delay:    delay,
				Speed:    speed,
				Err:      err,
			})
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
			return nil
		})
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	})

	return nil
}

func (p *ProxyCheck) CheckWithClash(clashConfig string) (float64, float64, error) {
	var proxyItem map[string]interface{}
	err := json.Unmarshal([]byte(clashConfig), &proxyItem)
	if err != nil {
		log.Errorf("err:%v", err)
		return 0, 0, err
	}

	proxyItem["name"] = "eutamias"

	// coverFloat2int := func(field string) {
	// 	if x, ok := proxyItem[field].(float64); ok {
	// 		proxyItem[field] = int(x)
	// 	}
	// }
	//
	// coverFloat2int("port")
	// coverFloat2int("alterId")

	proxyNode, err := outbound.ParseProxy(proxyItem)
	if err != nil {
		return 0, 0, err
	}

	delay, speed, err := p.URLTest(proxyNode, p.checkUrl)
	// delay, err := URLTest(proxy, "http://www.gstatic.com/generate_204")
	if err != nil {
		return 0, 0, err
	}

	return float64(delay.Milliseconds()), speed, nil
}

func (p *ProxyCheck) CheckWithLink(nodeUrl string) (float64, float64, error) {
	proxyConfig, err := proxy.ParseProxyToClash(nodeUrl)
	if err != nil {
		return 0, 0, err
	}

	return p.CheckWithClash(proxyConfig)
}

func (p *ProxyCheck) WaitFinish() {
	p.pool.WaitingTasks()
}

func (p *ProxyCheck) URLTest(proxy constant.Proxy, url string) (delay time.Duration, speed float64, err error) {
	ctx := context.Background()
	ctx, _ = context.WithTimeout(ctx, p.timeout)

	addr, err := urlToMetadata(url)
	if err != nil {
		return
	}

	instance, err := proxy.DialContext(ctx, &addr)
	if err != nil {
		return
	}
	defer instance.Close()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.WithContext(ctx)
	req.Header.Set("User-Agent", p.faker.UserAgent())

	transport := &http.Transport{
		Dial: func(string, string) (net.Conn, error) {
			return instance, nil
		},
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false, // 不要忽略ssl校验
		},
		TLSHandshakeTimeout:   p.timeout,
		DisableKeepAlives:     true,
		DisableCompression:    true,
		IdleConnTimeout:       p.timeout,
		ResponseHeaderTimeout: p.timeout,
		ExpectContinueTimeout: p.timeout,
	}

	client := http.Client{
		Transport: transport,
		Timeout:   p.timeout,
	}

	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		//log.Errorf("err:%v", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		//log.Errorf("err:%v", err)
		return
	}

	delay = time.Since(start)

	if strings.Contains(string(body), "GLaDOS 停止工作") {
		err = errors.New("unusable")
		return
	}

	if resp.ContentLength > 0 {
		speed = float64(resp.ContentLength) / float64(delay.Milliseconds())
	} else {
		speed = float64(len(body)) / float64(delay.Milliseconds())
	}

	return
}

func urlToMetadata(rawURL string) (addr C.Metadata, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return
	}

	port := u.Port()
	if port == "" {
		switch u.Scheme {
		case "https":
			port = "443"
		case "http":
			port = "80"
		default:
			err = fmt.Errorf("%s scheme not Support", rawURL)
			return
		}
	}

	addr = C.Metadata{
		AddrType: C.AtypDomainName,
		Host:     u.Hostname(),
		DstIP:    nil,
		DstPort:  port,
	}
	return
}
