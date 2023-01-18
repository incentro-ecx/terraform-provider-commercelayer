package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceSku() *schema.Resource {
	return &schema.Resource{
		Description:   "",
		ReadContext:   resourceSkuReadFunc,
		CreateContext: resourceSkuCreateFunc,
		UpdateContext: resourceSkuUpdateFunc,
		DeleteContext: resourceSkuDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The SKU unique identifier",
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
							Description: "The internal name of the SKU.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"code": {
							Description: "The SKU code, that uniquely identifies the SKU within the organization.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"description": {
							Description: "An internal description of the SKU.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"image_url": {
							Description: "The URL of an image that represents the SKU.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"pieces_per_pack": {
							Description: "The number of pieces that compose the SKU. This is useful to describe sets and bundles.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"weight": {
							Description: "The weight of the SKU. If present, it will be used to calculate the shipping rates.",
							Type:        schema.TypeFloat,
							Optional:    true,
						},
						"unit_of_weight": {
							Description: "Can be one of 'gr', 'lb', or 'oz'",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"hs_tariff_number": {
							Description: "The Harmonized System Code used by customs to identify the products shipped across international borders.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"do_not_ship": {
							Description: "Indicates if the SKU doesn't generate shipments.",
							Type:        schema.TypeBool,
							Optional:    true,
						},
						"do_not_track": {
							Description: "Indicates if the SKU doesn't track the stock inventory.",
							Type:        schema.TypeBool,
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
						"shipping_category_id": {
							Description: "The shipping category id.",
							Type:        schema.TypeString,
							Required:    true,
						}}}},
		},
	}
}

func resourceSkuReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.SkusApi.GETSkusSkuId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	sku, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(sku.GetId())

	return nil
}

func resourceSkuCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	skuCreate := commercelayer.SkuCreate{
		Data: commercelayer.SkuCreateData{
			Type: skusType,
			Attributes: commercelayer.POSTSkus201ResponseDataAttributes{
				Description:     stringRef(attributes["description"].(string)),
				ImageUrl:        stringRef(attributes["image_url"].(string)),
				PiecesPerPack:   intToInt32Ref(attributes["pieces_per_pack"]),
				Weight:          float64ToFloat32Ref(attributes["weight"]),
				UnitOfWeight:    stringRef(attributes["unit_of_weight"].(string)),
				HsTariffNumber:  stringRef(attributes["hs_tariff_number"].(string)),
				DoNotShip:       boolRef(attributes["do_not_ship"]),
				DoNotTrack:      boolRef(attributes["do_not_track"]),
				Name:            attributes["name"].(string),
				Code:            attributes["name"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.SkuCreateDataRelationships{
				ShippingCategory: commercelayer.ShipmentDataRelationshipsShippingCategory{
					Data: commercelayer.ShipmentDataRelationshipsShippingCategoryData{
						Type: stringRef(shippingCategoryType),
						Id:   stringRef(relationships["shipping_category_id"]),
					},
				},
			},
		},
	}

	err := d.Set("type", skusType)
	if err != nil {
		return diagErr(err)
	}

	sku, _, err := c.SkusApi.POSTSkus(ctx).SkuCreate(skuCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*sku.Data.Id)

	return nil
}

func resourceSkuDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.SkusApi.DELETESkusSkuId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceSkuUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	var skuUpdate = commercelayer.SkuUpdate{
		Data: commercelayer.SkuUpdateData{
			Type: skusType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHSkusSkuId200ResponseDataAttributes{
				Description:     stringRef(attributes["description"].(string)),
				ImageUrl:        stringRef(attributes["image_url"].(string)),
				PiecesPerPack:   intToInt32Ref(attributes["pieces_per_pack"]),
				Weight:          float64ToFloat32Ref(attributes["weight"]),
				UnitOfWeight:    stringRef(attributes["unit_of_weight"].(string)),
				HsTariffNumber:  stringRef(attributes["hs_tariff_number"].(string)),
				DoNotShip:       boolRef(attributes["do_not_ship"]),
				DoNotTrack:      boolRef(attributes["do_not_track"]),
				Name:            stringRef(attributes["name"].(string)),
				Code:            stringRef(attributes["name"].(string)),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.SkuUpdateDataRelationships{
				ShippingCategory: &commercelayer.ShipmentDataRelationshipsShippingCategory{
					Data: commercelayer.ShipmentDataRelationshipsShippingCategoryData{
						Type: stringRef(shippingCategoryType),
						Id:   stringRef(relationships["shipping_category_id"]),
					},
				},
			},
		},
	}

	_, _, err := c.SkusApi.PATCHSkusSkuId(ctx, d.Id()).SkuUpdate(skuUpdate).Execute()

	return diag.FromErr(err)
}
