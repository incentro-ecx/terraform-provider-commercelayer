package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceTaxCategory() *schema.Resource {
	return &schema.Resource{
		Description: "Create a tax category for an SKU that has special taxation. " +
			"Specify a valid tax code for the associated tax calculator.",
		ReadContext:   resourceTaxCategoryReadFunc,
		CreateContext: resourceTaxCategoryCreateFunc,
		UpdateContext: resourceTaxCategoryUpdateFunc,
		DeleteContext: resourceTaxCategoryDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The tax category unique identifier",
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
							Description: "The external tax calculator's internal name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"code": {
							Description: "The external tax calculator's internal name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"sku_code": {
							Description: "The external tax calculator's internal name",
							Type:        schema.TypeString,
							Optional:    true,
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
						"sku_id": {
							Description: "The associated address id.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"tax_calculator_id": {
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

func resourceTaxCategoryReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.TaxCategoriesApi.GETTaxCategoriesTaxCategoryId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	externalTaxCalculator, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(externalTaxCalculator.GetId())

	return nil
}

func resourceTaxCategoryCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	externalTaxCalculatorCreate := commercelayer.TaxCategoryCreate{
		Data: commercelayer.TaxCategoryCreateData{
			Type: taxCategoriesType,
			Attributes: commercelayer.POSTTaxCategories201ResponseDataAttributes{
				Code:            attributes["code"].(string),
				SkuCode:         stringRef(attributes["sku_code"].(string)),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.TaxCategoryCreateDataRelationships{
				Sku: commercelayer.BundleDataRelationshipsSkus{
					Data: commercelayer.BundleDataRelationshipsSkusData{
						Type: stringRef(skusType),
						Id:   stringRef(relationships["sku_id"]),
					},
				},
				//TODO: fill in default tax calculator
				TaxCalculator: commercelayer.TaxCategoryDataRelationshipsTaxCalculator{
					AvalaraAccount:        nil,
					ExternalTaxCalculator: nil,
					ManualTaxCalculator:   nil,
					TaxjarAccount:         nil,
				},
			},
		},
	}

	err := d.Set("type", externalTaxCalculatorType)
	if err != nil {
		return diagErr(err)
	}

	externalTaxCalculator, _, err := c.TaxCategoriesApi.POSTTaxCategories(ctx).TaxCategoryCreate(externalTaxCalculatorCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*externalTaxCalculator.Data.Id)

	return nil
}

func resourceTaxCategoryDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.TaxCategoriesApi.DELETETaxCategoriesTaxCategoryId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceTaxCategoryUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	var TaxCategoryUpdate = commercelayer.TaxCategoryUpdate{
		Data: commercelayer.TaxCategoryUpdateData{
			Type: externalTaxCalculatorType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHTaxCategoriesTaxCategoryId200ResponseDataAttributes{
				Code:            stringRef(attributes["code"].(string)),
				SkuCode:         stringRef(attributes["code"].(string)),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.TaxCategoryUpdateDataRelationships{
				Sku: &commercelayer.BundleDataRelationshipsSkus{
					Data: commercelayer.BundleDataRelationshipsSkusData{
						Type: stringRef(skusType),
						Id:   stringRef(relationships["sku_id"]),
					},
				},
			},
		},
	}

	_, _, err := c.TaxCategoriesApi.PATCHTaxCategoriesTaxCategoryId(ctx, d.Id()).TaxCategoryUpdate(TaxCategoryUpdate).Execute()

	return diag.FromErr(err)
}
