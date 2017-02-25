package integration

import (
	"os"
	"testing"
)

type configStruct struct {
	AccessToken string
}

func config(t *testing.T) configStruct {
	c := configStruct{AccessToken: os.Getenv("ACCESS_TOKEN")}
	if c.AccessToken == "" {
		t.Error("You must set ACCESS_TOKEN environment variable for " +
			"integration testing. " +
			"https://tech.yandex.com/oauth/doc/dg/tasks/get-oauth-token-docpage/")
	}
	return c
}
