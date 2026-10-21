package setting

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/suzuki-shunsuke/flute/v2/flute"

	"github.com/zahiar/terraform-provider-graylog/graylog/testutil"
)

func TestAccLDAPSetting(t *testing.T) {
	if err := testutil.SetEnv(); err != nil {
		t.Fatal(err)
	}

	var backendBody string

	resourceName := "graylog_ldap_setting.system"

	roleBody := `{
  "total": 1,
  "page": 1,
  "per_page": 100,
  "count": 1,
  "roles": [
    {
      "id": "6ad77e0659a8b019391944c9",
      "name": "Reader",
      "description": "Grants basic permissions for every user (built-in)",
      "permissions": ["messages:read"],
      "read_only": true
    }
  ],
  "query": "",
  "grand_total": 1,
  "context": {}
}`

	getRoute := flute.Route{
		Name: "get backends",
		Matcher: flute.Matcher{
			Method: "GET",
		},
		Tester: flute.Tester{
			Path:         "/api/system/authentication/services/backends",
			PartOfHeader: testutil.Header(),
		},
		Response: flute.Response{
			Response: func(req *http.Request) (*http.Response, error) {
				if backendBody == "" {
					backendBody = `[
  {
    "id": "test-backend-id",
    "title": "LDAP",
    "type": "ldap",
    "enabled": false,
    "default_roles": ["6ad77e0659a8b019391944c9"],
    "config": {
      "servers": [{"host": "localhost", "port": 389}],
      "transport_security": "none",
      "type": "ldap",
      "user_search_base": "",
      "user_search_pattern": "",
      "user_name_attribute": ""
    },
    "authentication_dn": ""
  }
]`
				}
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(backendBody)),
				}, nil
			},
		},
	}

	rolesRoute := flute.Route{
		Name: "get roles",
		Matcher: flute.Matcher{
			Method: "GET",
		},
		Tester: flute.Tester{
			Path:         "/api/authz/roles",
			PartOfHeader: testutil.Header(),
		},
		Response: flute.Response{
			Response: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       ioutil.NopCloser(strings.NewReader(roleBody)),
				}, nil
			},
		},
	}

	postRoute := flute.Route{
		Name: "create a backend",
		Matcher: flute.Matcher{
			Method: "POST",
		},
		Tester: flute.Tester{
			Path:         "/api/system/authentication/services/backends",
			PartOfHeader: testutil.Header(),
		},
		Response: flute.Response{
			Base: http.Response{
				StatusCode: 201,
			},
		},
	}

	deleteRoute := flute.Route{
		Name: "delete a backend",
		Matcher: flute.Matcher{
			Method: "DELETE",
		},
		Tester: flute.Tester{
			Path:         "/api/system/authentication/services/backends/test-backend-id",
			PartOfHeader: testutil.Header(),
		},
		Response: flute.Response{
			Base: http.Response{
				StatusCode: 204,
			},
		},
	}

	createStep := resource.TestStep{
		ResourceName: resourceName,
		PreConfig: func() {
			testutil.SetHTTPClient(t, getRoute, rolesRoute, postRoute, deleteRoute)
		},
		Config: `
resource "graylog_ldap_setting" "system" {
  title                     = "LDAP"
  server_url                = "ldap://localhost:389"
  authentication_username   = ""
  authentication_password   = ""
  user_search_base          = ""
  user_search_pattern       = ""
  user_name_attribute       = ""
  default_role              = "Reader"
  enabled                   = false
}
`,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, "title", "LDAP"),
			resource.TestCheckResourceAttr(resourceName, "server_url", "ldap://localhost:389"),
			resource.TestCheckResourceAttr(resourceName, "default_role", "Reader"),
			resource.TestCheckResourceAttr(resourceName, "enabled", "false"),
		),
	}

	updateRoute := flute.Route{
		Name: "update a backend",
		Matcher: flute.Matcher{
			Method: "PUT",
		},
		Tester: flute.Tester{
			Path:         "/api/system/authentication/services/backends/test-backend-id",
			PartOfHeader: testutil.Header(),
		},
		Response: flute.Response{
			Base: http.Response{
				StatusCode: 200,
			},
		},
	}

	updateStep := resource.TestStep{
		ResourceName: resourceName,
		PreConfig: func() {
			testutil.SetHTTPClient(t, getRoute, rolesRoute, updateRoute, deleteRoute)
		},
		Config: `
resource "graylog_ldap_setting" "system" {
  title                     = "LDAP"
  server_url                = "ldaps://localhost:636"
  authentication_username   = ""
  authentication_password   = ""
  user_search_base          = ""
  user_search_pattern       = ""
  user_name_attribute       = ""
  default_role              = "Reader"
  enabled                   = false
}
`,
		Check: resource.ComposeTestCheckFunc(
			resource.TestCheckResourceAttr(resourceName, "title", "LDAP"),
			resource.TestCheckResourceAttr(resourceName, "server_url", "ldaps://localhost:636"),
			resource.TestCheckResourceAttr(resourceName, "default_role", "Reader"),
		),
	}

	resource.Test(t, resource.TestCase{
		Providers: testutil.SingleResourceProviders("graylog_ldap_setting", Resource()),
		Steps: []resource.TestStep{
			createStep,
			updateStep,
		},
	})
}
