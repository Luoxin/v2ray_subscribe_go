package proxy

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Socket struct {
	Base
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s Socket) ToSurge() string {
	panic("implement me")
}

func (s Socket) String() string {
	data, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(data)
}

func (s Socket) ToClash() string {
	data, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(data)
}

func (s Socket) Link() string {
	u := url.URL{
		Scheme:     "socket",
		User:       url.UserPassword(s.Username, s.Password),
		Path:       s.Path,
		ForceQuery: false,
	}

	switch s.Type {
	case "socket4":
		u.Scheme = "socket4"
	case "socket5":
		u.Scheme = "socket5"
	}

	host := s.Server
	if s.Port != 0 {
		host = fmt.Sprintf("%v:%v", host, s.Port)
	}
	u.Host = host

	return u.String()
}

func (s Socket) Identifier() string {
	return net.JoinHostPort(s.Server, strconv.Itoa(s.Port))
}

func (s Socket) Clone() Proxy {
	return &s
}

func ParseHttpLink(link string) (*Http, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	h := Http{
		Base: Base{
			Path:    u.Path,
			Type:    "http",
			UDP:     false,
			Useable: true,
		},
		Tls: false,
	}

	h.Port, err = strconv.Atoi(u.Port())
	if err != nil {
		log.Errorf("err:%v", err)
		h.Port = 80
	}
	h.Server = strings.TrimSuffix(u.Host, fmt.Sprintf(":%v", h.Port))

	h.Username = u.User.Username()
	h.Password, _ = u.User.Password()

	if u.Scheme == "https" {
		h.Tls = true
	}

	return &h, nil
}
