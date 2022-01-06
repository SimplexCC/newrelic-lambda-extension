package checks

import (
	"bytes"
	"errors"
	"github.com/newrelic/newrelic-lambda-extension/config"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockClientError struct{}

func (c *mockClientError) Get(string) (*http.Response, error) {
	return nil, errors.New("Something went wrong")
}

type mockClientRedirect struct{}

func (c *mockClientRedirect) Get(string) (*http.Response, error) {
	body := ioutil.NopCloser(bytes.NewBufferString("Hello World"))
	return &http.Response{Body: body, StatusCode: 301}, nil
}

type mockClientSuccess struct{}

func (c *mockClientSuccess) Get(string) (*http.Response, error) {
	body := ioutil.NopCloser(bytes.NewBufferString("<html><body>You are being <a href=\"https://github.com/newrelic/node-newrelic/releases/tag/v8.5.0\">redirected</a>.</body></html>"))
	return &http.Response{Body: body, StatusCode: 302, Header: map[string][]string{"Location": {"https://github.com/newrelic/node-newrelic/releases/tag/v8.5.0"}}}, nil
}

func TestRuntimeCheck(t *testing.T) {
	dirname, err := os.MkdirTemp("", "")
	assert.Nil(t, err)
	defer os.RemoveAll(dirname)

	oldPath := runtimeLookupPath
	defer func() {
		runtimeLookupPath = oldPath
	}()
	runtimeLookupPath = filepath.Join(dirname, runtimeLookupPath)

	os.MkdirAll(filepath.Join(runtimeLookupPath, "node"), os.ModePerm)
	conf := config.Configuration{}
	client = &mockClientSuccess{}
	r, err := checkAndReturnRuntime(&conf)
	assert.Equal(t, runtimeConfigs[Node].language, r.language)
	assert.Equal(t, "v8.5.0", r.AgentVersion)
	assert.Nil(t, err)
}

func TestRuntimeCheckNil(t *testing.T) {
	conf := config.Configuration{}
	r, err := checkAndReturnRuntime(&conf)
	assert.Equal(t, runtimeConfig{}, r)
	assert.Nil(t, err)
}

func TestLatestAgentTag(t *testing.T) {
	client = &mockClientError{}
	assert.Nil(t, latestAgentTag(&runtimeConfig{}))

	client = &mockClientRedirect{}
	assert.Nil(t, latestAgentTag(&runtimeConfig{}))

	client = &mockClientSuccess{}
	r := runtimeConfig{}
	assert.Nil(t, latestAgentTag(&r))
	assert.Equal(t, "v8.5.0", r.AgentVersion)
}
