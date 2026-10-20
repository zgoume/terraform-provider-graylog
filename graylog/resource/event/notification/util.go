package notification

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/zahiar/terraform-provider-graylog/graylog/convert"
)

const (
	keyID        = "id"
	keyConfig    = "config"
	keyAPISecret = "api_secret"
	keyBasicAuth = "basic_auth"
)

func getDataFromResourceData(d *schema.ResourceData) (map[string]interface{}, error) {
	data, err := convert.GetFromResourceData(d, Resource())
	if err != nil {
		return nil, err
	}

	if err := convert.JSONToData(data, keyConfig); err != nil {
		return nil, err
	}

	return data, nil
}

func setDataToResourceData(d *schema.ResourceData, data map[string]interface{}) error {
	// The API may return api_secret/basic_auth (empty) and *_is_set booleans.
	// Delete these so SetResourceData doesn't overwrite our state values.
	delete(data, keyAPISecret)
	delete(data, keyBasicAuth)
	delete(data, keyAPISecret+"_is_set")
	delete(data, keyBasicAuth+"_is_set")

	if err := convert.DataToJSON(data, keyConfig); err != nil {
		return err
	}

	if err := convert.SetResourceData(d, Resource(), data); err != nil {
		return err
	}

	d.SetId(data[keyID].(string))
	return nil
}

func getSecretPayloadFromResourceData(d *schema.ResourceData) map[string]interface{} {
	payload := make(map[string]interface{})
	if v, ok := d.GetOk(keyAPISecret); ok && v.(string) != "" {
		payload[keyAPISecret] = v.(string)
	}
	if v, ok := d.GetOk(keyBasicAuth); ok && v.(string) != "" {
		payload[keyBasicAuth] = v.(string)
	}
	return payload
}

func getSecretSetPayload(d *schema.ResourceData) map[string]interface{} {
	payload := make(map[string]interface{})
	for _, key := range []string{keyAPISecret, keyBasicAuth} {
		newVal, oldVal := d.GetChange(key)
		newStr := newVal.(string)
		oldStr := oldVal.(string)
		if newStr != "" && newStr != oldStr {
			payload[key] = newStr
		} else if newStr == "" && oldStr != "" {
			payload[key] = ""
		}
	}
	return payload
}

func mergeSecretsIntoConfig(config map[string]interface{}, secrets map[string]interface{}) {
	for k, v := range secrets {
		config[k] = v
	}
}

func removeSecretKeys(data map[string]interface{}) {
	delete(data, keyAPISecret)
	delete(data, keyBasicAuth)
}

func handleSecretRead(d *schema.ResourceData, apiData map[string]interface{}) {
	config, ok := apiData[keyConfig].(map[string]interface{})
	if !ok {
		return
	}
	delete(config, keyAPISecret)
	delete(config, keyBasicAuth)
}
