package proxy

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
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

// func ParseHttpLink(link string) (*Http, error) {
// 	u, err := url.Parse(link)
// 	if err != nil {
// 	    return nil, err
// 	}
//
// 	h := Http{
// 		Base:     Base{
// 			Name:    "",
// 			Server:  u.Host,
// 			Path:    u.Path,
// 			Type:    "http",
// 			UDP:     false,
// 			Useable: false,
// 		},
// 		Username: "",
// 		Password: "",
// 		Tls:      false,
// 	}
// }
