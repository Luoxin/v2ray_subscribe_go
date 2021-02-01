package proxycheck

import (
	"context"
	"encoding/json"
	"github.com/Dreamacro/clash/adapters/outbound"
	"github.com/panjf2000/ants/v2"
	log "github.com/sirupsen/logrus"
	"github.com/thedevsaddam/retry"
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

	ctx, _ := context.WithTimeout(context.Background(), time.Minute)

	// x, err := proxy.URLTest(ctx, "https://www.google.com/generate_204")
	delay, err := proxy.URLTest(ctx, "http://www.gstatic.com/generate_204")
	if err != nil {
		log.Errorf("err:%v", err)
		return 0, 0, err
	}

	return float64(delay), 0, nil
}
