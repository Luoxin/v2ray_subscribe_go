package geolite

import (
	"fmt"
	"testing"
)

func TestGetCountry(t *testing.T) {
	x, _ := GetCountry("google.com")
	fmt.Printf("%+v", x)
}
