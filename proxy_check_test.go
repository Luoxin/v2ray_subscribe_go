package main

import (
	"subsrcibe/subscription"
	"testing"
)

func Test_checkNode(t *testing.T) {
	s = &State{
		Config: &Config{
			ProxyCheckUrl: "https://www.google.com/generate_204",
		},
	}

	checkNode(&subscription.ProxyNode{
		ProxyNodeType: 1,
		NodeDetail: &subscription.ProxyNode_NodeDetail{
			Buf: "vmess://ew0KICAidiI6ICIyIiwNCiAgInBzIjogIiIsDQogICJhZGQiOiAiNDcuNzUuNDkuMyIsDQogICJwb3J0IjogIjM2NjQ0IiwNCiAgImlkIjogIjdmMTg5YWQ2LTE2MGYtNGExZi1hODgwLWVhN2Y4NzZhZWFmYiIsDQogICJhaWQiOiAiMjMzIiwNCiAgIm5ldCI6ICJ3cyIsDQogICJ0eXBlIjogIm5vbmUiLA0KICAiaG9zdCI6ICIiLA0KICAicGF0aCI6ICIiLA0KICAidGxzIjogIiINCn0=",
		},
	})
}
