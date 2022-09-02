package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceMerchant() *schema.Resource {
	return &schema.Resource{
		Description:   "",
		ReadContext:   resourceMerchantReadFunc,
		CreateContext: resourceMerchantCreateFunc,
		UpdateContext: resourceMerchantUpdateFunc,
		DeleteContext: resourceMerchantDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The merchant unique identifier",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"attributes": {
				Description: "Resource attributes",
				Type:        schema.TypeList,
				MaxItems:    1,
				MinItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "",
							Type:        schema.TypeString,
							Required:    true,
						},
						"metadata": {
							Description: "Set of key-value pairs that you can attach to the resource. This can be useful " +
								"for storing additional information about the resource in a structured format",
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
						},
					},
				},
			},
			"relationships": {
				Description: "Resource relationships",
				Type:        schema.TypeList,
				MaxItems:    1,
				MinItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						//	TODO: implement geocoder relation
					},
				},
			},
		},
	}
}

func resourceMerchantReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	return diag.Errorf("Not implemented")
}

func resourceMerchantCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := d.Get("attributes").([]interface{})[0].(map[string]interface{})

	merchantCreate := commercelayer.MerchantCreate{
		Data: commercelayer.MerchantCreateData{
			Type: merchantType,
			Attributes: commercelayer.POSTMerchants201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			Relationships: nil,
		},
	}

	merchant, _, err := c.MerchantsApi.POSTMerchants(ctx).MerchantCreate(merchantCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*merchant.Data.Id)

	return nil
}

func resourceMerchantDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	return diag.Errorf("Not implemented")
}

func resourceMerchantUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	return diag.Errorf("Not implemented")
}
