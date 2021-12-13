package httpcli

import (
	"fmt"
	"testing"
)

func TestIsHttps(t *testing.T) {
	us := []string{
		"Https://1232.22.22.22/xx/xx/xx",
		"Http://1232.22.22.22/xx/xx/xx",
		"https://1232.22.22.22/xx/xx/xx",
		"http://1232.22.22.22/xx/xx/xx",
		"HTTPS://1232.22.22.22/xx/xx/xx",
	}

	for _, u := range us {
		fmt.Println(u, IsHttps(u))
	}
}
