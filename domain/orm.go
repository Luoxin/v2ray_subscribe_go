package domain

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/mcuadros/go-defaults"

	"github.com/pkg/errors"
)

func Scan(src interface{}, dst interface{}) error {
	x := func(buf []byte) error {
		bufLen := len(buf)
		if bufLen >= 2 && buf[0] == '{' && buf[bufLen-1] == '}' {
			return json.Unmarshal(buf, dst)
		} else if bufLen > 0 {
			return json.Unmarshal(buf, dst)
		} else {
			defaults.SetDefaults(dst)
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
}

func Value(m interface{}) (driver.Value, error) {
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
