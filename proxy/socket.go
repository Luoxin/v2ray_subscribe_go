package proxy

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
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
