package geolite

import (
	"fmt"
	"testing"
)

func TestGetCountry(t *testing.T) {
	x, _ := GetCountry("141.101.115.17")
	fmt.Printf("%+v", x)
}
