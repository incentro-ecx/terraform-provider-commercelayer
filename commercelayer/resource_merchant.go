package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceMerchant() *schema.Resource {
	return &schema.Resource{
		Description: "A merchant is the fiscal representative that is selling in a specific market. Tax calculators " +
			"use the merchant's address (and the shipping address) to determine the tax rate for an order.",
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
			"type": {
				Description: "The resource type",
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
							Description: "The merchant's internal name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"reference": {
							Description: "A string that you can use to add any external identifier to the resource. This " +
								"can be useful for integrating the resource to an external system, like an ERP, a " +
								"marketing tool, a CRM, or whatever.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"reference_origin": {
							Description: "Any identifier of the third party system that defines the reference code",
							Type:        schema.TypeString,
							Optional:    true,
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
						"address_id": {
							Description: "The associated address id.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceMerchantReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.MerchantsApi.GETMerchantsMerchantId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	merchant, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(merchant.GetId())

	return nil
}

func resourceMerchantCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	merchantCreate := commercelayer.MerchantCreate{
		Data: commercelayer.MerchantCreateData{
			Type: merchantType,
			Attributes: commercelayer.POSTMerchants201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.MerchantCreateDataRelationships{
				Address: commercelayer.CustomerAddressCreateDataRelationshipsAddress{
					Data: commercelayer.BingGeocoderDataRelationshipsAddressesData{
						Type: stringRef(addressType),
						Id:   stringRef(relationships["address_id"]),
					},
				},
			},
		},
	}

	err := d.Set("type", merchantType)
	if err != nil {
		return diagErr(err)
	}

	merchant, _, err := c.MerchantsApi.POSTMerchants(ctx).MerchantCreate(merchantCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*merchant.Data.Id)

	return nil
}

func resourceMerchantDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.MerchantsApi.DELETEMerchantsMerchantId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceMerchantUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	var merchantUpdate = commercelayer.MerchantUpdate{
		Data: commercelayer.MerchantUpdateData{
			Type: merchantType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHMerchantsMerchantId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"].(string)),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.MerchantUpdateDataRelationships{
				Address: &commercelayer.CustomerAddressCreateDataRelationshipsAddress{
					Data: commercelayer.BingGeocoderDataRelationshipsAddressesData{
						Type: stringRef(addressType),
						Id:   stringRef(relationships["address_id"]),
					},
				},
			},
		},
	}

	_, _, err := c.MerchantsApi.PATCHMerchantsMerchantId(ctx, d.Id()).MerchantUpdate(merchantUpdate).Execute()

	return diag.FromErr(err)
}
