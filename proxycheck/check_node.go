package proxycheck

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/Dreamacro/clash/adapters/outbound"
	"github.com/Dreamacro/clash/constant"
	C "github.com/Dreamacro/clash/constant"
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
	"github.com/thedevsaddam/retry"

	"subsrcibe/proxy"
)

type ProxyCheck struct {
	pool *ants.Pool
	w    sync.WaitGroup

	maxSize int
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
	proxyConfig, err := proxy.ParseProxyToClash(nodeUrl)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = p.AddWithClash(proxyConfig, logic)
	if err != nil {
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
		log.Errorf("err:%v", err)
		return 0, 0, err
	}

	delay, speed, err := URLTest(proxy, "https://www.google.com")
	//delay, err := URLTest(proxy, "http://www.gstatic.com/generate_204")
	if err != nil {
		log.Errorf("err:%v", err)
		return 0, 0, err
	}

	return float64(delay.Milliseconds()), speed, nil
}

func (p *ProxyCheck) CheckWithLink(nodeUrl string) (float64, float64, error) {
	proxyConfig, err := proxy.ParseProxyToClash(nodeUrl)
	if err != nil {
		log.Errorf("err:%v", err)
		return 0, 0, err
	}

	return p.CheckWithClash(proxyConfig)
}

func (p *ProxyCheck) WaitFinish() {
	p.w.Done()
}

func URLTest(p constant.Proxy, url string) (delay time.Duration, speed float64, err error) {
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)

	addr, err := urlToMetadata(url)
	if err != nil {
		return
	}

	instance, err := p.DialContext(ctx, &addr)
	if err != nil {
		return
	}
	defer instance.Close()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req = req.WithContext(ctx)

	transport := &http.Transport{
		Dial: func(string, string) (net.Conn, error) {
			return instance, nil
		},
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		//TLSHandshakeTimeout: time.Minute,
		//DisableKeepAlives: true,
		//DisableCompression: true,
		//MaxIdleConns:          10,
		//IdleConnTimeout:       time.Minute,
		//ResponseHeaderTimeout: time.Minute,
		//ExpectContinueTimeout: time.Minute,
		//TLSNextProto:          nil,
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

	if len(body) > 0 {
		speed = float64(len(body)) / float64(delay.Milliseconds())
	} else if resp.ContentLength > 0 {
		speed = float64(resp.ContentLength) / float64(delay.Milliseconds())
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
