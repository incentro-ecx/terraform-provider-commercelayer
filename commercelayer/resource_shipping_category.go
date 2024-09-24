package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceShippingCategory() *schema.Resource {
	return &schema.Resource{
		Description: "Shipping categories determine which shipping methods are available for the associated " +
			"SKU's. Unless the selected inventory model strategy is no_split, if an order contains line items " +
			"belonging to more than one shipping category it is split into more shipments.",
		ReadContext:   resourceShippingCategoryReadFunc,
		CreateContext: resourceShippingCategoryCreateFunc,
		UpdateContext: resourceShippingCategoryUpdateFunc,
		DeleteContext: resourceShippingCategoryDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The shipping category unique identifier",
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
							Description: "The shipping category's internal name.",
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

func resourceShippingCategoryReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.ShippingCategoriesApi.GETShippingCategoriesShippingCategoryId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	shippingCategory, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(shippingCategory.GetId().(string))

	return nil
}

func resourceShippingCategoryCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	shippingCategoryCreate := commercelayer.ShippingCategoryCreate{
		Data: commercelayer.ShippingCategoryCreateData{
			Type: shippingCategoryType,
			Attributes: commercelayer.POSTShippingCategories201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	err := d.Set("type", shippingCategoryType)
	if err != nil {
		return diagErr(err)
	}

	shippingCategory, _, err := c.ShippingCategoriesApi.POSTShippingCategories(ctx).ShippingCategoryCreate(shippingCategoryCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(shippingCategory.Data.GetId().(string))

	return nil
}

func resourceShippingCategoryDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.ShippingCategoriesApi.DELETEShippingCategoriesShippingCategoryId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceShippingCategoryUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var shippingCategoryUpdate = commercelayer.ShippingCategoryUpdate{
		Data: commercelayer.ShippingCategoryUpdateData{
			Type: shippingCategoryType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHShippingCategoriesShippingCategoryId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.ShippingCategoriesApi.PATCHShippingCategoriesShippingCategoryId(ctx, d.Id()).ShippingCategoryUpdate(shippingCategoryUpdate).Execute()

	return diag.FromErr(err)
}
