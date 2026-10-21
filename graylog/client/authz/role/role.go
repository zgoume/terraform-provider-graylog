package role

import (
	"context"
	"net/http"

	"github.com/zahiar/terraform-provider-graylog/graylog/client/httpclient"
)

type Role struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
	ReadOnly    bool     `json:"read_only"`
}

type RolesResponse struct {
	Total      int                    `json:"total"`
	Page       int                    `json:"page"`
	PerPage    int                    `json:"per_page"`
	Count      int                    `json:"count"`
	Roles      []Role                 `json:"roles"`
	Query      string                 `json:"query"`
	GrandTotal int                    `json:"grand_total"`
	Context    map[string]interface{} `json:"context"`
}

type Client struct {
	Client httpclient.Client
}

func (cl Client) GetAll(ctx context.Context) (*RolesResponse, *http.Response, error) {
	body := &RolesResponse{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/api/authz/roles?page=1&per_page=100",
		ResponseBody: &body,
	})
	return body, resp, err
}
