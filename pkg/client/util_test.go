package client

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testTransport struct {
	roundtripfunc func(req *http.Request) (*http.Response, error)
}

func (transport testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return transport.roundtripfunc(req)
}

func TestGetMetricDataSucessful(t *testing.T) {
	transport := testTransport{
		roundtripfunc: func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, "/api/v1/storage", r.URL.Path)
			assert.Equal(t, "GET", r.Method)
			return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(`{"storage_layer":{"chunks": {"total_chunks": 2}}}`))}, nil
		},
	}
	FBClient := FluentBitClient{FBHost: "fake", FBPort: 1234, HTTPClient: http.Client{Transport: transport}}
	expected := Response{StorageLayer: StorageLayer{Chunks: ChunksStorage{TotalChunks: 2}}}
	resp, err := FBClient.GetMetricData()
	assert.NoError(t, err)
	assert.Equal(t, expected, *resp)
}

func TestGetMetricDataWithNot200Response(t *testing.T) {
	transport := testTransport{
		roundtripfunc: func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, "/api/v1/storage", r.URL.Path)
			assert.Equal(t, "GET", r.Method)
			return &http.Response{StatusCode: http.StatusInternalServerError, Body: ioutil.NopCloser(strings.NewReader(`INTERNAL SERVER ERRROR`))}, nil
		},
	}
	FBClient := FluentBitClient{FBHost: "fake", FBPort: 1234, HTTPClient: http.Client{Transport: transport}}
	_, err := FBClient.GetMetricData()

	assert.Error(t, err)
}

func TestGetMetricWithNonExpectedBody(t *testing.T) {
	transport := testTransport{
		roundtripfunc: func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, "/api/v1/storage", r.URL.Path)
			assert.Equal(t, "GET", r.Method)
			return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(strings.NewReader(`{"unknown_field":{"unknown_field": {"unknown_field": 2}}}}`))}, nil
		},
	}
	FBClient := FluentBitClient{FBHost: "fake", FBPort: 1234, HTTPClient: http.Client{Transport: transport}}
	_, err := FBClient.GetMetricData()

	assert.Error(t, err)
}
