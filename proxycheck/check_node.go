package proxycheck

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Dreamacro/clash/adapters/outbound"
	"github.com/Dreamacro/clash/constant"
	C "github.com/Dreamacro/clash/constant"
	"github.com/luoxin/v2ray_subscribe_go/faker"
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
	"github.com/thedevsaddam/retry"

	"github.com/luoxin/v2ray_subscribe_go/subscribe/proxy"
)

type ProxyCheck struct {
	pool *ants.Pool
	w    sync.WaitGroup

	maxSize int
	faker   *faker.Faker
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
		maxSize: 10,
		faker:   faker.New(),
	}
}

func (p *ProxyCheck) Init() error {
	if p.pool == nil {
		pool, err := ants.NewPool(10)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		p.pool = pool
	} else {
		p.pool.Tune(p.maxSize)
	}

	return nil
}

func (p *ProxyCheck) SetMaxSize(size int) error {
	p.maxSize = size

	if p.pool != nil {
		p.pool.Tune(size)
	}

	return nil
}

func (p *ProxyCheck) AddWithClash(nodeUrl string, logic func(result Result) error) error {
	p.w.Add(1)
	err := p.pool.Submit(func() {
		defer p.w.Done()
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
			p.w.Done()
			log.Errorf("err:%v", err)
			return
		}
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func (p *ProxyCheck) AddWithLink(nodeUrl string, logic func(result Result) error) error {
	p.w.Add(1)
	err := p.pool.Submit(func() {
		defer p.w.Done()

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
	if err != nil {
		p.w.Done()
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func (p *ProxyCheck) CheckWithClash(clashConfig string) (float64, float64, error) {
	var proxyItem map[string]interface{}
	err := json.Unmarshal([]byte(clashConfig), &proxyItem)
	if err != nil {
		log.Errorf("err:%v", err)
		return 0, 0, err
	}

	proxyItem["name"] = "test proxy"

	coverFloat2int := func(field string) {
		if x, ok := proxyItem[field].(float64); ok {
			proxyItem[field] = int(x)
		}
	}

	coverFloat2int("port")
	coverFloat2int("alterId")

	proxy, err := outbound.ParseProxy(proxyItem)
	if err != nil {
		return 0, 0, err
	}

	time.Sleep(time.Second * 1)

	delay, speed, err := p.URLTest(proxy, "https://www.google.com")
	//delay, err := URLTest(proxy, "http://www.gstatic.com/generate_204")
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
	p.w.Wait()
}

func (p *ProxyCheck) URLTest(proxy constant.Proxy, url string) (delay time.Duration, speed float64, err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)

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
	req = req.WithContext(ctx)
	req.Header.Set("User-Agent", p.faker.UserAgent())

	transport := &http.Transport{
		Dial: func(string, string) (net.Conn, error) {
			return instance, nil
		},
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false, // 不要忽略ssl校验
		},
	}

	client := http.Client{
		Transport: transport,
		//CheckRedirect: func(req *http.Request, via []*http.Request) error {
		//	return http.ErrUseLastResponse
		//},
		Timeout: time.Minute,
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

	delay = time.Now().Sub(start)

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
