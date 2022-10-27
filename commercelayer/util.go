package commercelayer

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func diagErr(err error) diag.Diagnostics {
	apiErr, ok := err.(*commercelayer.GenericOpenAPIError)
	if ok {
		return diag.Errorf("%s: %s", apiErr.Error(), string(apiErr.Body()))
	}
	return diag.FromErr(err)
}

func stringRef(val interface{}) *string {
	if val == nil {
		return nil
	}
	ref := val.(string)
	if ref == "" {
		return nil
	}
	return &ref
}

func intToInt32Ref(val interface{}) *int32 {
	if val == nil {
		return nil
	}
	ref := int32(val.(int))
	if ref == 0 {
		return nil
	}
	return &ref
}

func keyValueRef(val interface{}) map[string]interface{} {
	if val == nil {
		return map[string]interface{}{}
	}
	return val.(map[string]interface{})
}

func boolRef(val interface{}) *bool {
	if val == nil {
		return nil
	}

	ref := val.(bool)

	return &ref
}

func float64ToFloat32Ref(val interface{}) *float32 {
	if val == nil {
		return nil
	}

	ref := float32(val.(float64))

	if ref == 0 {
		return nil
	}

	return &ref
}

func stringSliceValueRef(val interface{}) []string {
	if val == nil {
		return []string{}
	}

	var s []string

	for _, v := range val.([]interface{}) {
		s = append(s, v.(string))
	}

	return s
}

func nestedMap(val interface{}) map[string]any {
	if val == nil {
		return map[string]any{}
	}
	valMap := val.([]any)
	if len(valMap) == 0 {
		return map[string]any{}
	}

	return valMap[0].(map[string]any)
}
