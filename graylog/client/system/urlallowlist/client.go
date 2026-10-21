package urlallowlist

import (
	"context"
	"errors"
	"net/http"

	"github.com/zahiar/terraform-provider-graylog/graylog/client/httpclient"
)

type Client struct {
	Client httpclient.Client
}

func (cl Client) Get(ctx context.Context) (map[string]interface{}, *http.Response, error) {
	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/system/urlallowlist",
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Upsert(
	ctx context.Context, data map[string]interface{},
) (map[string]interface{}, *http.Response, error) {
	if data == nil {
		return nil, nil, errors.New("request body is nil")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "PUT",
		Path:         "/system/urlallowlist",
		RequestBody:  data,
		ResponseBody: &body,
	})
	return body, resp, err
}

func (cl Client) Delete(ctx context.Context) (*http.Response, error) {
	body := map[string]interface{}{
		"entries":  []interface{}{},
		"disabled": false,
	}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "PUT",
		Path:         "/system/urlallowlist",
		RequestBody:  body,
		ResponseBody: &body,
	})
	return resp, err
}
