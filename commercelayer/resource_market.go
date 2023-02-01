package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceMarket() *schema.Resource {
	return &schema.Resource{
		Description: "A market is made of a merchant, an inventory model, and a price list (plus an optional " +
			"customer group, geocoder, and tax calculator)",
		ReadContext:   resourceMarketReadFunc,
		CreateContext: resourceMarketCreateFunc,
		UpdateContext: resourceMarketUpdateFunc,
		DeleteContext: resourceMarketDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The market unique identifier",
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
							Description: "The Market's internal name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"facebook_pixel_id": {
							Description: "The Facebook Pixed ID",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"checkout_url": {
							Description: "The checkout URL for this market",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"external_prices_url": {
							Description: "The URL used to fetch prices from an external source",
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
						"merchant_id": {
							Description: "The associated merchant id.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"price_list_id": {
							Description: "The associated price list id.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"inventory_model_id": {
							Description: "The associated inventory model id.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"customer_group_id": {
							Description: "The associated customer group id.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"tax_calculator_id": {
							Description: "The associated tax calculator id.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceMarketReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.MarketsApi.GETMarketsMarketId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	Market, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(Market.GetId())

	return nil
}

func resourceMarketCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	marketCreate := commercelayer.MarketCreate{
		Data: commercelayer.MarketCreateData{
			Type: marketType,
			Attributes: commercelayer.POSTMarkets201ResponseDataAttributes{
				Name:              attributes["name"].(string),
				FacebookPixelId:   stringRef(attributes["facebook_pixel_id"]),
				CheckoutUrl:       stringRef(attributes["checkout_url"]),
				ExternalPricesUrl: stringRef(attributes["external_prices_url"]),
				Reference:         stringRef(attributes["reference"]),
				ReferenceOrigin:   stringRef(attributes["reference_origin"]),
				Metadata:          keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.MarketCreateDataRelationships{
				Merchant: commercelayer.MarketCreateDataRelationshipsMerchant{
					Data: commercelayer.MarketDataRelationshipsMerchantData{
						Type: stringRef(merchantType),
						Id:   stringRef(relationships["merchant_id"]),
					},
				},
				PriceList: commercelayer.MarketCreateDataRelationshipsPriceList{
					Data: commercelayer.MarketDataRelationshipsPriceListData{
						Type: stringRef(priceListType),
						Id:   stringRef(relationships["price_list_id"]),
					},
				},
				InventoryModel: commercelayer.InventoryReturnLocationCreateDataRelationshipsInventoryModel{
					Data: commercelayer.InventoryReturnLocationDataRelationshipsInventoryModelData{
						Type: stringRef(inventoryModelType),
						Id:   stringRef(relationships["inventory_model_id"]),
					},
				},
			},
		},
	}

	taxCalculatorId := stringRef(relationships["tax_calculator_id"])
	if taxCalculatorId != nil {
		marketCreate.Data.Relationships.TaxCalculator = &commercelayer.MarketCreateDataRelationshipsTaxCalculator{
			Data: commercelayer.MarketDataRelationshipsTaxCalculatorData{
				Type: stringRef(taxCalculatorType),
				Id:   taxCalculatorId,
			}}
	}

	customerGroupId := stringRef(relationships["customer_group_id"])
	if customerGroupId != nil {
		marketCreate.Data.Relationships.CustomerGroup = &commercelayer.CustomerCreateDataRelationshipsCustomerGroup{
			Data: commercelayer.CustomerDataRelationshipsCustomerGroupData{
				Type: stringRef(customerGroupType),
				Id:   customerGroupId,
			}}
	}

	err := d.Set("type", marketType)
	if err != nil {
		return diagErr(err)
	}

	market, _, err := c.MarketsApi.POSTMarkets(ctx).MarketCreate(marketCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*market.Data.Id)

	return nil
}

func resourceMarketDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.MarketsApi.DELETEMarketsMarketId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceMarketUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	var marketUpdate = commercelayer.MarketUpdate{
		Data: commercelayer.MarketUpdateData{
			Type: marketType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHMarketsMarketId200ResponseDataAttributes{
				Name:              stringRef(attributes["name"]),
				FacebookPixelId:   stringRef(attributes["facebook_pixel_id"]),
				CheckoutUrl:       stringRef(attributes["checkout_url"]),
				ExternalPricesUrl: stringRef(attributes["external_prices_url"]),
				Reference:         stringRef(attributes["reference"]),
				ReferenceOrigin:   stringRef(attributes["reference_origin"]),
				Metadata:          keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.MarketUpdateDataRelationships{
				Merchant: &commercelayer.MarketCreateDataRelationshipsMerchant{
					Data: commercelayer.MarketDataRelationshipsMerchantData{
						Type: stringRef(merchantType),
						Id:   stringRef(relationships["merchant_id"]),
					},
				},
				PriceList: &commercelayer.MarketCreateDataRelationshipsPriceList{
					Data: commercelayer.MarketDataRelationshipsPriceListData{
						Type: stringRef(priceListType),
						Id:   stringRef(relationships["price_list_id"]),
					},
				},
				InventoryModel: &commercelayer.InventoryReturnLocationCreateDataRelationshipsInventoryModel{
					Data: commercelayer.InventoryReturnLocationDataRelationshipsInventoryModelData{
						Type: stringRef(inventoryModelType),
						Id:   stringRef(relationships["inventory_model_id"]),
					},
				},
			},
		},
	}

	taxCalculatorId := stringRef(relationships["tax_calculator_id"])
	if taxCalculatorId != nil {
		marketUpdate.Data.Relationships.TaxCalculator = &commercelayer.MarketCreateDataRelationshipsTaxCalculator{
			Data: commercelayer.MarketDataRelationshipsTaxCalculatorData{
				Type: stringRef(taxCalculatorType),
				Id:   stringRef(taxCalculatorId),
			}}
	}

	customerGroupId := stringRef(relationships["customer_group_id"])
	if customerGroupId != nil {
		marketUpdate.Data.Relationships.CustomerGroup = &commercelayer.CustomerCreateDataRelationshipsCustomerGroup{
			Data: commercelayer.CustomerDataRelationshipsCustomerGroupData{
				Type: stringRef(customerGroupType),
				Id:   stringRef(customerGroupId),
			}}
	}

	_, _, err := c.MarketsApi.PATCHMarketsMarketId(ctx, d.Id()).MarketUpdate(marketUpdate).Execute()

	return diag.FromErr(err)
}
