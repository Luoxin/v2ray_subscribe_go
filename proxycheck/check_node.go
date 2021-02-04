package proxycheck

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Dreamacro/clash/adapters/outbound"
	"github.com/Dreamacro/clash/constant"
	C "github.com/Dreamacro/clash/constant"
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
	"github.com/thedevsaddam/retry"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"subsrcibe/utils"
	"time"
)

type ProxyCheck struct {
	pool *ants.Pool
}

func NewProxyCheck() *ProxyCheck {
	return &ProxyCheck{}
}

func (p *ProxyCheck) Start() error {
	pool, err := ants.NewPool(10)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	p.pool = pool
	return nil
}

func (p *ProxyCheck) Add(nodeUrl string, logic func(err error, delay, speed float64) error) error {
	err := p.pool.Submit(func() {
		delay, speed, err := p.Check(nodeUrl)
		err = retry.DoFunc(5, 500*time.Millisecond, func() error {
			err = logic(err, delay, speed)
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

func (p *ProxyCheck) Check(nodeUrl string) (float64, float64, error) {
	proxyConfig := utils.ParseProxy(nodeUrl)

	j, err := json.Marshal(proxyConfig)
	if err != nil {
		log.Errorf("err:%v", err)
		return 0, 0, err
	}

	var proxyItem map[string]interface{}
	err = json.Unmarshal(j, &proxyItem)
	if err != nil {
		log.Errorf("err:%v", err)
		return 0, 0, err
	}

	proxyItem["name"] = "test proxy"

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

	req, err := http.NewRequest(http.MethodHead, url, nil)
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
		log.Errorf("err:%v", err)
		return
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	delay = time.Now().Sub(start)
	speed = float64(resp.ContentLength) / float64(delay.Milliseconds())

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
