package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceShippingMethod() *schema.Resource {
	return &schema.Resource{
		Description:   `Shipping methods are used to provide customers with different delivery options.`,
		ReadContext:   resourceShippingMethodReadFunc,
		CreateContext: resourceShippingMethodCreateFunc,
		UpdateContext: resourceShippingMethodUpdateFunc,
		DeleteContext: resourceShippingMethodDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The shipping method unique identifier",
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
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The shipping method's name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"scheme": {
							Description: "The shipping method's scheme, one of 'flat' or 'weight_tiered'.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"currency_code": {
							Description: "The international 3-letter currency code as defined by the ISO " +
								"4217 standard.",
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: currencyCodeValidation,
						},
						"price_amount_cents": {
							Description: "The price of this shipping method, in cents.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"free_over_amount_cents": {
							Description: "Apply free shipping if the order amount is over this value, in cents.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"min_weight": {
							Description: "The minimum weight for which this shipping method is available.",
							Type:        schema.TypeFloat,
							Optional:    true,
						},
						"max_weight": {
							Description: "The maximum weight for which this shipping method is available.",
							Type:        schema.TypeFloat,
							Optional:    true,
						},
						"unit_of_weight": {
							Description: "Can be one of 'gr', 'lb', or 'oz'",
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
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"market_id": {
							Description: "The associated market id.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"shipping_zone_id": {
							Description: "The shipping zone that is used to match the order shipping address.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"shipping_category_id": {
							Description: "The shipping category for which this shipping method is available.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"stock_location_id": {
							Description: "The stock location for which this shipping method is available.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"shipping_method_tier_ids": {
							Description: "The associated shipping method tiers (meaningful when " +
								"billing_scheme != 'flat').",
							Type: schema.TypeList,
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

func resourceShippingMethodReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.ShippingMethodsApi.GETShippingMethodsShippingMethodId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	shippingMethod, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(shippingMethod.GetId())

	return nil
}

func resourceShippingMethodCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	shippingMethodCreate := commercelayer.ShippingMethodCreate{
		Data: commercelayer.ShippingMethodCreateData{
			Type: shippingMethodType,
			Attributes: commercelayer.POSTShippingMethods201ResponseDataAttributes{
				Name:                attributes["name"].(string),
				Scheme:              stringRef(attributes["scheme"]),
				CurrencyCode:        stringRef(attributes["currency_code"]),
				PriceAmountCents:    int32(attributes["price_amount_cents"].(int)),
				FreeOverAmountCents: intToInt32Ref(attributes["free_over_amount_cents"]),
				MinWeight:           float64ToFloat32Ref(attributes["min_weight"]),
				MaxWeight:           float64ToFloat32Ref(attributes["max_weight"]),
				UnitOfWeight:        stringRef(attributes["unit_of_weight"]),
				Reference:           stringRef(attributes["reference"]),
				ReferenceOrigin:     stringRef(attributes["reference_origin"]),
				Metadata:            keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.ShippingMethodCreateDataRelationships{},
		},
	}

	marketId := stringRef(relationships["market_id"])
	if marketId != nil {
		shippingMethodCreate.Data.Relationships.Market = &commercelayer.BillingInfoValidationRuleCreateDataRelationshipsMarket{
			Data: commercelayer.AvalaraAccountDataRelationshipsMarketsData{
				Type: stringRef(marketType),
				Id:   marketId,
			}}
	}

	shippingZoneId := stringRef(relationships["shipping_zone_id"])
	if shippingZoneId != nil {
		shippingMethodCreate.Data.Relationships.ShippingZone = &commercelayer.ShippingMethodCreateDataRelationshipsShippingZone{
			Data: commercelayer.ShippingMethodDataRelationshipsShippingZoneData{
				Type: stringRef(shippingZoneType),
				Id:   shippingZoneId,
			}}
	}

	shippingCategoryId := stringRef(relationships["shipping_category_id"])
	if shippingCategoryId != nil {
		shippingMethodCreate.Data.Relationships.ShippingCategory = &commercelayer.ShippingMethodCreateDataRelationshipsShippingCategory{
			Data: commercelayer.ShipmentDataRelationshipsShippingCategoryData{
				Type: stringRef(shippingCategoryType),
				Id:   shippingCategoryId,
			}}
	}

	stockLocationId := stringRef(relationships["stock_location_id"])
	if stockLocationId != nil {
		shippingMethodCreate.Data.Relationships.StockLocation = &commercelayer.DeliveryLeadTimeCreateDataRelationshipsStockLocation{
			Data: commercelayer.DeliveryLeadTimeDataRelationshipsStockLocationData{
				Type: stringRef(stockLocationType),
				Id:   stockLocationId,
			}}
	}

	// TODO: fix shipping method tiers
	//shippingMethodTierIds := stringRef(relationships["shipping_method_tier_ids"])
	//if shippingMethodTierIds != nil {
	//	shippingMethodCreate.Data.Relationships.ShippingMethodTiers = &commercelayer.ShippingMethodDataRelationshipsShippingMethodTiers{
	//		Data: commercelayer.ShippingMethodDataRelationshipsShippingMethodTiersData{
	//			Type: stringRef(shippingMethodTierType),
	//			Id:   shippingMethodTierIds,
	//		}}
	//}

	err := d.Set("type", shippingMethodType)
	if err != nil {
		return diagErr(err)
	}

	shippingMethod, _, err := c.ShippingMethodsApi.POSTShippingMethods(ctx).ShippingMethodCreate(shippingMethodCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*shippingMethod.Data.Id)

	return nil
}

func resourceShippingMethodDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.ShippingMethodsApi.DELETEShippingMethodsShippingMethodId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceShippingMethodUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	var shippingMethodUpdate = commercelayer.ShippingMethodUpdate{
		Data: commercelayer.ShippingMethodUpdateData{
			Type: shippingMethodType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHShippingMethodsShippingMethodId200ResponseDataAttributes{
				Name:                stringRef(attributes["name"]),
				Scheme:              stringRef(attributes["scheme"]),
				CurrencyCode:        stringRef(attributes["currency_code"]),
				PriceAmountCents:    intToInt32Ref(attributes["price_amount_cents"]),
				FreeOverAmountCents: intToInt32Ref(attributes["free_over_amount_cents"]),
				MinWeight:           float64ToFloat32Ref(attributes["min_weight"]),
				MaxWeight:           float64ToFloat32Ref(attributes["max_weight"]),
				UnitOfWeight:        stringRef(attributes["unit_of_weight"]),
				Reference:           stringRef(attributes["reference"]),
				ReferenceOrigin:     stringRef(attributes["reference_origin"]),
				Metadata:            keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.ShippingMethodCreateDataRelationships{},
		},
	}

	marketId := stringRef(relationships["market_id"])
	if marketId != nil {
		shippingMethodUpdate.Data.Relationships.Market = &commercelayer.BillingInfoValidationRuleCreateDataRelationshipsMarket{
			Data: commercelayer.AvalaraAccountDataRelationshipsMarketsData{
				Type: stringRef(marketType),
				Id:   marketId,
			}}
	}

	shippingZoneId := stringRef(relationships["shipping_zone_id"])
	if shippingZoneId != nil {
		shippingMethodUpdate.Data.Relationships.ShippingZone = &commercelayer.ShippingMethodCreateDataRelationshipsShippingZone{
			Data: commercelayer.ShippingMethodDataRelationshipsShippingZoneData{
				Type: stringRef(shippingZoneType),
				Id:   shippingZoneId,
			}}
	}

	shippingCategoryId := stringRef(relationships["shipping_category_id"])
	if shippingCategoryId != nil {
		shippingMethodUpdate.Data.Relationships.ShippingCategory = &commercelayer.ShippingMethodCreateDataRelationshipsShippingCategory{
			Data: commercelayer.ShipmentDataRelationshipsShippingCategoryData{
				Type: stringRef(shippingCategoryType),
				Id:   shippingCategoryId,
			}}
	}

	stockLocationId := stringRef(relationships["stock_location_id"])
	if stockLocationId != nil {
		shippingMethodUpdate.Data.Relationships.StockLocation = &commercelayer.DeliveryLeadTimeCreateDataRelationshipsStockLocation{
			Data: commercelayer.DeliveryLeadTimeDataRelationshipsStockLocationData{
				Type: stringRef(stockLocationType),
				Id:   stockLocationId,
			}}
	}

	// TODO: fix shipping method tiers
	//shippingMethodTierIds := stringSliceValueRef(relationships["shipping_method_tier_ids"])
	//if shippingMethodTierIds != nil {
	//	shippingMethodUpdate.Data.Relationships.ShippingMethodTiers = &commercelayer.ShippingMethodDataRelationshipsShippingMethodTiers{
	//		Data: commercelayer.ShippingMethodDataRelationshipsShippingMethodTiersData{
	//			Type: stringRef(shippingMethodTierType),
	//			Id:   shippingMethodTierIds,
	//		}}
	//}

	_, _, err := c.ShippingMethodsApi.PATCHShippingMethodsShippingMethodId(ctx, d.Id()).ShippingMethodUpdate(shippingMethodUpdate).Execute()

	return diag.FromErr(err)
}
