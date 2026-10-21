package setting

import (
	"context"
	"fmt"
	"net/http"

	"github.com/zahiar/terraform-provider-graylog/graylog/client/httpclient"
)

type Client struct {
	Client httpclient.Client
}

func (cl Client) ListBackends(ctx context.Context) ([]map[string]interface{}, *http.Response, error) {
	body := []map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/api/system/authentication/services/backends",
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) GetBackend(ctx context.Context, backendID string) (map[string]interface{}, *http.Response, error) {
	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/api/system/authentication/services/backends/" + backendID,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) CreateBackend(ctx context.Context, data map[string]interface{}) (*http.Response, error) {
	if data == nil {
		return nil, fmt.Errorf("request body is nil")
	}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:      "POST",
		Path:        "/api/system/authentication/services/backends",
		RequestBody: data,
	})
	return resp, err
}

func (cl Client) UpdateBackend(ctx context.Context, backendID string, data map[string]interface{}) (*http.Response, error) {
	if data == nil {
		return nil, fmt.Errorf("request body is nil")
	}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:      "PUT",
		Path:        "/api/system/authentication/services/backends/" + backendID,
		RequestBody: data,
	})
	return resp, err
}

func (cl Client) DeleteBackend(ctx context.Context, backendID string) (*http.Response, error) {
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method: "DELETE",
		Path:   "/api/system/authentication/services/backends/" + backendID,
	})
	return resp, err
}
