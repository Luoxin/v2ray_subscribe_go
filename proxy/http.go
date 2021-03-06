package proxy

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/Luoxin/Eutamias/utils/json"

	log "github.com/sirupsen/logrus"
)

type Http struct {
	Base
	Username string `json:"username"`
	Password string `json:"password"`
	Tls      bool   `json:"tls"`
}

func (h Http) ToSurge() string {
	panic("implement me")
}

func (h Http) String() string {
	data, err := json.Marshal(h)
	if err != nil {
		return ""
	}
	return string(data)
}

func (h Http) ToClash() string {
	data, err := json.Marshal(h)
	if err != nil {
		return ""
	}
	return string(data)
}

func (h Http) Link() string {
	u := url.URL{
		User:       url.UserPassword(h.Username, h.Password),
		Path:       h.Path,
		ForceQuery: false,
	}
	if h.Tls {
		u.Scheme = "https"
	} else {
		u.Scheme = "http"
	}

	host := h.Server
	if h.Port != 0 {
		host = fmt.Sprintf("%v:%v", host, h.Port)
	}
	u.Host = host

	return u.String()
}

func (h Http) Identifier() string {
	return net.JoinHostPort(h.Server, strconv.Itoa(h.Port))
}

func (h Http) Clone() Proxy {
	return &h
}

func ParseSocketLink(link string) (*Http, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	h := Http{
		Base: Base{
			Path:    u.Path,
			Type:    u.Scheme,
			Useable: true,
		},
		Tls: false,
	}

	h.Port, err = strconv.Atoi(u.Port())
	if err != nil {
		log.Errorf("err:%v", err)
	}
	h.Server = strings.TrimSuffix(u.Host, fmt.Sprintf(":%v", h.Port))

	h.Username = u.User.Username()
	h.Password, _ = u.User.Password()

	return &h, nil
}
