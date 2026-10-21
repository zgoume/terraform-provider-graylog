package urlallowlist

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/suzuki-shunsuke/flute/v2/flute"

	"github.com/zahiar/terraform-provider-graylog/graylog/testutil"
)

func TestAccURLAllowList(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	urlAllowListBody := ""

	resourceURLPath := "/api/system/urlallowlist"
	resourceName := "graylog_url_allow_list.test"

	getRoute := flute.Route{
		Name: "get url allow list",
		Matcher: flute.Matcher{
			Method: "GET",
		},
		Tester: flute.Tester{
			Path:         resourceURLPath,
			PartOfHeader: testutil.Header(),
		},
		Response: flute.Response{
			Response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(urlAllowListBody)),
				}, nil
			},
		},
	}

	postRoute := flute.Route{
		Name: "create url allow list",
		Matcher: flute.Matcher{
			Method: "PUT",
		},
		Tester: flute.Tester{
			Path:         resourceURLPath,
			PartOfHeader: testutil.Header(),
			Test: func(t *testing.T, req *http.Request, svc flute.Service, route flute.Route) {
				urlAllowListBody = `{
  "entries": [
    {
      "id": "076419b6-8b36-4be6-b774-b3d8035640f1",
      "title": "alerta",
      "value": "https://alerta.cloud-temple.sec/api/webhooks/graylog",
      "type": "literal"
    }
  ],
  "disabled": false
}`
			},
		},
		Response: flute.Response{
			Base: http.Response{
				StatusCode: 200,
			},
			BodyString: `{
  "entries": [
    {
      "id": "076419b6-8b36-4be6-b774-b3d8035640f1",
      "title": "alerta",
      "value": "https://alerta.cloud-temple.sec/api/webhooks/graylog",
      "type": "literal"
    }
  ],
  "disabled": false
}`,
		},
	}

	deleteRoute := flute.Route{
		Name: "delete url allow list",
		Matcher: flute.Matcher{
			Method: "PUT",
		},
		Tester: flute.Tester{
			Path:         resourceURLPath,
			PartOfHeader: testutil.Header(),
		},
		Response: flute.Response{
			Base: http.Response{
				StatusCode: 200,
			},
		},
	}

	createStep := resource.TestStep{
		ResourceName: resourceName,
		PreConfig: func() {
			testutil.SetHTTPClient(t, getRoute, postRoute, deleteRoute)
		},
		Config: `
resource "graylog_url_allow_list" "test" {
  entries {
    id = "076419b6-8b36-4be6-b774-b3d8035640f1"
    title = "alerta"
    value = "https://alerta.cloud-temple.sec/api/webhooks/graylog"
    type = "literal"
  }
}
`,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, "entries.0.title", "alerta"),
		),
	}

	updateRoute := flute.Route{
		Name: "update url allow list",
		Matcher: flute.Matcher{
			Method: "PUT",
		},
		Tester: flute.Tester{
			Path:         resourceURLPath,
			PartOfHeader: testutil.Header(),
			Test: func(t *testing.T, req *http.Request, svc flute.Service, route flute.Route) {
				urlAllowListBody = `{
  "entries": [
    {
      "id": "076419b6-8b36-4be6-b774-b3d8035640f1",
      "title": "alerta-updated",
      "value": "https://alerta.cloud-temple.sec/api/webhooks/graylog",
      "type": "literal"
    }
  ],
  "disabled": true
}`
			},
		},
		Response: flute.Response{
			Base: http.Response{
				StatusCode: 200,
			},
			BodyString: `{
  "entries": [
    {
      "id": "076419b6-8b36-4be6-b774-b3d8035640f1",
      "title": "alerta-updated",
      "value": "https://alerta.cloud-temple.sec/api/webhooks/graylog",
      "type": "literal"
    }
  ],
  "disabled": true
}`,
		},
	}

	updateStep := resource.TestStep{
		ResourceName: resourceName,
		PreConfig: func() {
			testutil.SetHTTPClient(t, getRoute, updateRoute, deleteRoute)
		},
		Config: `
resource "graylog_url_allow_list" "test" {
  entries {
    id = "076419b6-8b36-4be6-b774-b3d8035640f1"
    title = "alerta-updated"
    value = "https://alerta.cloud-temple.sec/api/webhooks/graylog"
    type = "literal"
  }
  disabled = true
}
`,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, "entries.0.title", "alerta-updated"),
			resource.TestCheckResourceAttr(resourceName, "disabled", "true"),
		),
	}

	resource.Test(t, resource.TestCase{
		Providers: testutil.SingleResourceProviders("graylog_url_allow_list", Resource()),
		Steps: []resource.TestStep{
			createStep,
			updateStep,
		},
	})
}
