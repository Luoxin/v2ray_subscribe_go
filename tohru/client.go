package tohru

import (
	"errors"
	"time"

	"github.com/Luoxin/Eutamias/conf"
	"github.com/Luoxin/Eutamias/proxynode"
	"github.com/Luoxin/Eutamias/utils"
	"github.com/Luoxin/Eutamias/utils/json"
	"github.com/Luoxin/Eutamias/version"
	"github.com/elliotchance/pie/pie"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"go.uber.org/atomic"
)

const Hello = "B3vUNO|I,|\"FAco9b<fIPj0K:r,Zsj\"?KFOA}.z1N&LZOP1GYq"

type tohru struct {
	client *resty.Client

	tplToken atomic.String
}

func newTohru() *tohru {
	return &tohru{}
}

var Tohru = newTohru()

func (p *tohru) Init() error {
	p.client = resty.New().
		SetTimeout(time.Second * 5).
		SetRetryMaxWaitTime(time.Second * 5).
		SetRetryWaitTime(time.Second).
		SetLogger(log.New())

	if conf.Config.Debug {
		p.client.SetDebug(true).EnableTrace()
	}

	return nil
}

func (p *tohru) Start() error {
	sync := func() error {
		return p.SyncNode()
	}

	if conf.Config.IsTohru() {
		err := p.CheckUsable()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = sync()
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		go func() {
			for {
				err = sync()
				if err != nil {
					log.Errorf("err:%v", err)
				}
				time.Sleep(time.Minute * 30)
			}

		}()
	}

	return nil
}

func (p *tohru) GenEncryptionPassword(password string) string {
	return utils.Sha384(password)
}

func (p *tohru) CheckUsable() error {
	b, err := json.Marshal(UserInfo{
		Hello:         Hello,
		TohruKey:      conf.Config.Base.TohruKey,
		TohruPassword: p.GenEncryptionPassword(conf.Config.Base.TohruPassword),
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	hello, err := conf.Ecc.ECCEncrypt4Str(string(b))
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var rsp CheckUsableRsp
	resp, err := p.client.R().
		SetBody(CheckUsableReq{
			Version: version.Version,
			Hello:   hello,
		}).
		SetResult(&rsp).
		Post(conf.Config.Base.KobayashiSanAddr + "/tohru/CheckUsable")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New(resp.String())
	}

	p.client.SetHeader(TokenKey, rsp.Token)
	p.tplToken.Store(rsp.Token)

	for _, cookie := range resp.Cookies() {
		p.client.SetCookie(cookie)
	}

	return nil
}

func (p *tohru) SyncNode() error {
	nodeList, err := proxynode.GetNode4Tohru(100)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var rsp SyncNodeRsp
	err = p.DoRequest("/tohru/SyncNode", &SyncNodeReq{
		NodeList: nodeList,
	}, &rsp)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	b, err := conf.Ecc.ECCDecrypt(rsp.NodeList)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var nodeUrlList pie.Strings
	err = json.Unmarshal([]byte(b), &nodeUrlList)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	nodeUrlList.Each(func(s string) {
		_, err = proxynode.AddNodeWithUrl(s)
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	})

	return nil
}

func (p *tohru) Registration(key, pwd string) error {
	var rsp RegistrationRsp
	resp, err := p.client.R().
		SetBody(&RegistrationReq{
			TohruKey:      key,
			TohruPassword: p.GenEncryptionPassword(pwd),
		}).
		SetResult(&rsp).
		SetBasicAuth(conf.Config.Base.KobayashiSanUserName, conf.Config.Base.KobayashiSanPassword).
		Post(conf.Config.Base.KobayashiSanAddr + "/tohru/Registration")
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New(resp.String())
	}

	return nil
}

func (p *tohru) DoRequest(path string, req, rsp interface{}) error {
	var lastErr error

	err := validate.Struct(req)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	for i := 0; i < 3; i++ {
		resp, err := p.client.R().
			SetHeader(TokenKey, p.tplToken.Load()).
			SetBody(req).
			SetResult(rsp).
			Post(conf.Config.Base.KobayashiSanAddr + path)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		switch resp.StatusCode() {
		case 200:
			return nil
		case 403:
			lastErr = errors.New(resp.String())

			err = p.CheckUsable()
			if err != nil {
				log.Errorf("err:%v", err)
				return err
			}
		default:
			return errors.New(resp.String())
		}
	}

	return lastErr
}

func (p *tohru) ChangedPassword(username string, oldPwd, newPwd string) error {
	var rsp ChangePasswordRsp
	err := p.DoRequest("/tohru/ChangePassword", &ChangePasswordReq{
		TohruKey:         username,
		OldTohruPassword: p.GenEncryptionPassword(oldPwd),
		NewTohruPassword: p.GenEncryptionPassword(newPwd),
	}, &rsp)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
