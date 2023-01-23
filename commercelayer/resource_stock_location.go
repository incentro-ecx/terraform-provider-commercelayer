package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceStockLocation() *schema.Resource {
	return &schema.Resource{
		Description: "Shipping zones determine the available shipping methods for a given shipping address. The " +
			"match is evaluated against a set of regular expressions on the address country, state or zip code.",
		ReadContext:   resourceStockLocationReadFunc,
		CreateContext: resourceStockLocationCreateFunc,
		UpdateContext: resourceStockLocationUpdateFunc,
		DeleteContext: resourceStockLocationDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The stock location unique identifier",
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
							Description: "The stock location's internal name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"label_format": {
							Description: "The shipping label format for this stock location. Can be one of 'PDF'" +
								", 'ZPL', 'EPL2', or 'PNG'",
							Type:     schema.TypeString,
							Optional: true,
						},
						"suppress_etd": {
							Description: "Flag it if you want to skip the electronic invoice creation when " +
								"generating the customs info for this stock location shipments.",
							Type:     schema.TypeBool,
							Optional: true,
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

func resourceStockLocationReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.StockLocationsApi.GETStockLocationsStockLocationId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	stockLocation, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(stockLocation.GetId())

	return nil
}

func resourceStockLocationCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	stockLocationCreate := commercelayer.StockLocationCreate{
		Data: commercelayer.StockLocationCreateData{
			Type: stockLocationType,
			Attributes: commercelayer.POSTStockLocations201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				LabelFormat:     stringRef(attributes["label_format"]),
				SuppressEtd:     boolRef(attributes["suppress_etd"]),
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

	err := d.Set("type", stockLocationType)
	if err != nil {
		return diagErr(err)
	}

	stockLocation, _, err := c.StockLocationsApi.POSTStockLocations(ctx).StockLocationCreate(stockLocationCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*stockLocation.Data.Id)

	return nil
}

func resourceStockLocationDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.StockLocationsApi.DELETEStockLocationsStockLocationId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceStockLocationUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	var stockLocationUpdate = commercelayer.StockLocationUpdate{
		Data: commercelayer.StockLocationUpdateData{
			Type: stockLocationType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHStockLocationsStockLocationId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"]),
				LabelFormat:     stringRef(attributes["label_format"]),
				SuppressEtd:     boolRef(attributes["suppress_etd"]),
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

	_, _, err := c.StockLocationsApi.PATCHStockLocationsStockLocationId(ctx, d.Id()).StockLocationUpdate(stockLocationUpdate).Execute()

	return diag.FromErr(err)
}
