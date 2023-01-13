package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceTaxjarAccount() *schema.Resource {
	return &schema.Resource{
		Description: "Configure your TaxJar account to automatically compute tax calculations " +
			"for the orders of the associated market.",
		ReadContext:   resourceTaxjarAccountReadFunc,
		CreateContext: resourceTaxjarAccountCreateFunc,
		UpdateContext: resourceTaxjarAccountUpdateFunc,
		DeleteContext: resourceTaxjarAccountDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The taxjar account unique identifier",
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
							Description: "The tax calculator's internal name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"api_key": {
							Description: "The TaxJar account API key.",
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
		},
	}
}

func resourceTaxjarAccountReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.TaxjarAccountsApi.GETTaxjarAccountsTaxjarAccountId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	taxjarAccount, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(taxjarAccount.GetId())

	return nil
}

func resourceTaxjarAccountCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	taxjarAccountCreate := commercelayer.TaxjarAccountCreate{
		Data: commercelayer.TaxjarAccountCreateData{
			Type: taxjarAccountsType,
			Attributes: commercelayer.POSTTaxjarAccounts201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				ApiKey:          attributes["api_key"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	err := d.Set("type", taxjarAccountsType)
	if err != nil {
		return diagErr(err)
	}

	taxjarAccount, _, err := c.TaxjarAccountsApi.POSTTaxjarAccounts(ctx).TaxjarAccountCreate(taxjarAccountCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*taxjarAccount.Data.Id)

	return nil
}

func resourceTaxjarAccountDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.TaxjarAccountsApi.DELETETaxjarAccountsTaxjarAccountId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceTaxjarAccountUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var taxjarAccountUpdate = commercelayer.TaxjarAccountUpdate{
		Data: commercelayer.TaxjarAccountUpdateData{
			Type: taxjarAccountsType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHTaxjarAccountsTaxjarAccountId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.TaxjarAccountsApi.PATCHTaxjarAccountsTaxjarAccountId(ctx, d.Id()).
		TaxjarAccountUpdate(taxjarAccountUpdate).Execute()

	return diag.FromErr(err)
}
