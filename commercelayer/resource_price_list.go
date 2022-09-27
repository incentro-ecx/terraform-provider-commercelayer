package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourcePriceList() *schema.Resource {
	return &schema.Resource{
		Description: `Price lists are collections of SKU prices, 
		defined by currency and market. When a list of SKUs is fetched, 
		only SKUs with a price defined in the market's price list and at least 
		a stock item in one of the market stock locations will be returned. 
		A user can create price lists to manage international business or B2B/B2C models.`,
		ReadContext:   resourcePriceListReadFunc,
		CreateContext: resourcePriceListCreateFunc,
		UpdateContext: resourcePriceListUpdateFunc,
		DeleteContext: resourcePriceListDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The PriceList unique identifier",
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
							Description: "The price list's internal name",
							Type:        schema.TypeString,
							Required:    true,
						},
						"currency_code": {
							Description:      "The international 3-letter currency code as defined by the ISO 4217 standard.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: currencyCodeValidation,
						},
						"tax_included": {
							Description: "Indicates if the associated prices include taxes.",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
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

func resourcePriceListReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.PriceListsApi.GETPriceListsPriceListId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	priceList, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(priceList.GetId())

	return nil
}

func resourcePriceListCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := d.Get("attributes").([]interface{})[0].(map[string]interface{})

	priceListCreate := commercelayer.PriceListCreate{
		Data: commercelayer.PriceListCreateData{
			Type: priceListType,
			Attributes: commercelayer.POSTPriceLists201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				CurrencyCode:    attributes["currency_code"].(string),
				TaxIncluded:     boolRef(attributes["tax_included"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	priceList, _, err := c.PriceListsApi.POSTPriceLists(ctx).PriceListCreate(priceListCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*priceList.Data.Id)

	return nil
}

func resourcePriceListDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.PriceListsApi.DELETEPriceListsPriceListId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourcePriceListUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := d.Get("attributes").([]interface{})[0].(map[string]interface{})

	var PriceListUpdate = commercelayer.PriceListUpdate{
		Data: commercelayer.PriceListUpdateData{
			Type: priceListType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHPriceListsPriceListId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"].(string)),
				CurrencyCode:    stringRef(attributes["currency_code"].(string)),
				TaxIncluded:     boolRef(attributes["tax_included"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.PriceListsApi.PATCHPriceListsPriceListId(ctx, d.Id()).PriceListUpdate(PriceListUpdate).Execute()

	return diag.FromErr(err)
}
