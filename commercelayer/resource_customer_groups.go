package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceCustomerGroup() *schema.Resource {
	return &schema.Resource{
		Description: "A CustomerGroup is the fiscal representative that is selling in a specific market. Tax calculators " +
			"use the CustomerGroup's address (and the shipping address) to determine the tax rate for an order.",
		ReadContext:   resourceCustomerGroupReadFunc,
		CreateContext: resourceCustomerGroupCreateFunc,
		UpdateContext: resourceCustomerGroupUpdateFunc,
		DeleteContext: resourceCustomerGroupDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The CustomerGroup unique identifier",
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
							Description: "The CustomerGroup's internal name",
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

func resourceCustomerGroupReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	return diag.Errorf("Not implemented")
}

func resourceCustomerGroupCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := d.Get("attributes").([]interface{})[0].(map[string]interface{})

	CustomerGroupCreate := commercelayer.CustomerGroupCreate{
		Data: commercelayer.CustomerGroupCreateData{
			Type: customerGroupType,
			Attributes: commercelayer.POSTCustomerGroups201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
				},
			},
		}

	CustomerGroup, _, err := c.CustomerGroupsApi.POSTCustomerGroups(ctx).CustomerGroupCreate(CustomerGroupCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*CustomerGroup.Data.Id)

	return nil
}

func resourceCustomerGroupDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	return diag.Errorf("Not implemented")
}

func resourceCustomerGroupUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	return diag.Errorf("Not implemented")
}
