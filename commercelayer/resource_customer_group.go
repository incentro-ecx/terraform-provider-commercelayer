package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceCustomerGroup() *schema.Resource {
	return &schema.Resource{
		Description: `A customer group is a resource that can be used to organize customers into groups. 
		When you associate a customer group to a market, that market becomes private and can be accessed
		 only by the customers belonging to the group. You can use customer groups to manage B2B customers, 
		 B2C loyalty programs, private sales, and more.`,
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
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The customer group's internal name",
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
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.CustomerGroupsApi.GETCustomerGroupsCustomerGroupId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	customer_group, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(customer_group.GetId())

	return nil
}

func resourceCustomerGroupCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := d.Get("attributes").([]interface{})[0].(map[string]interface{})

	customerGroupCreate := commercelayer.CustomerGroupCreate{
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

	customer_group, _, err := c.CustomerGroupsApi.POSTCustomerGroups(ctx).CustomerGroupCreate(customerGroupCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*customer_group.Data.Id)

	return nil
}

func resourceCustomerGroupDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.CustomerGroupsApi.DELETECustomerGroupsCustomerGroupId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceCustomerGroupUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := d.Get("attributes").([]interface{})[0].(map[string]interface{})

	var customerGroupUpdate = commercelayer.CustomerGroupUpdate{
		Data: commercelayer.CustomerGroupUpdateData{
			Type: customerGroupType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHCustomerGroupsCustomerGroupId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"].(string)),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.CustomerGroupsApi.PATCHCustomerGroupsCustomerGroupId(ctx, d.Id()).CustomerGroupUpdate(customerGroupUpdate).Execute()

	return diag.FromErr(err)
}
