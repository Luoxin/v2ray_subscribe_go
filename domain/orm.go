package domain

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

func Scan(src interface{}, dst proto.Message) error {
	x := func(buf []byte) error {
		bufLen := len(buf)

		if bufLen >= 2 && buf[0] == '{' && buf[bufLen-1] == '}' {
			u := &jsonpb.Unmarshaler{AllowUnknownFields: true}
			return u.Unmarshal(bytes.NewBuffer(buf), dst)
		} else if bufLen > 0 {
			return proto.Unmarshal(buf, dst)
		} else {
			dst.Reset()
		}
		return nil
	}

	switch r := src.(type) {
	case []byte:
		buf := src.([]byte)
		return x(buf)
	case string:
		buf := []byte(src.(string))
		return x(buf)
	default:
		return errors.New(
			fmt.Sprintf("unknown type %v %s to scan", r, reflect.ValueOf(src).String()))
	}
	return nil
}

func Value(m proto.Message) (driver.Value, error) {
	if m == nil {
		return nil, nil
	}

	return json.Marshal(m)
}

func (m *CrawlerConf_Rule) Scan(value interface{}) error {
	return Scan(value, m)
}

func (m *CrawlerConf_Rule) Value() (driver.Value, error) {
	return Value(m)
}
