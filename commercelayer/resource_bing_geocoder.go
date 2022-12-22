package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceBingGeocoders() *schema.Resource {
	return &schema.Resource{
		Description: "Geocoders lets you convert an address in text form into geographic coordinates " +
			"(like latitude and longitude). A geocoder can be optionally associated with a market. " +
			"By connecting a geocoder to a market, all the shipping and billing addresses belonging " +
			"to that market orders will be geocoded",
		ReadContext:   resourceBingGeocodersReadFunc,
		CreateContext: resourceBingGeocodersCreateFunc,
		UpdateContext: resourceBingGeocodersUpdateFunc,
		DeleteContext: resourceBingGeocodersDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The bing geocoder unique identifier",
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
							Description: "The geocoder's internal name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"reference": {
							Description: "A string that you can use to add any external identifier to the resource. " +
								"This can be useful for integrating the resource to an external system, like an ERP, a marketing tool, a CRM, or whatever.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"reference_origin": {
							Description: "Any identifier of the third party system that defines the reference code",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"key": {
							Description: "The Bing Virtualearth key.",
							Type:        schema.TypeString,
							Required:    true,
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

func resourceBingGeocodersReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.BingGeocodersApi.GETBingGeocodersBingGeocoderId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	bingGeocoder, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(bingGeocoder.GetId())

	return nil
}

func resourceBingGeocodersCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	bingGeocoderCreate := commercelayer.BingGeocoderCreate{
		Data: commercelayer.BingGeocoderCreateData{
			Type: bingGeocodersType,
			Attributes: commercelayer.POSTBingGeocoders201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
				Key:             attributes["key"].(string),
			},
		},
	}

	err := d.Set("type", bingGeocodersType)
	if err != nil {
		return diagErr(err)
	}

	bingGeocoders, _, err := c.BingGeocodersApi.POSTBingGeocoders(ctx).BingGeocoderCreate(bingGeocoderCreate).Execute()

	if err != nil {
		return diagErr(err)
	}

	d.SetId(*bingGeocoders.Data.Id)

	return nil
}

func resourceBingGeocodersDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.BingGeocodersApi.DELETEBingGeocodersBingGeocoderId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceBingGeocodersUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var bingGeocodersUpdate = commercelayer.BingGeocoderUpdate{
		Data: commercelayer.BingGeocoderUpdateData{
			Type: bingGeocodersType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHBingGeocodersBingGeocoderId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
				Key:             stringRef(attributes["key"]),
			},
		},
	}

	_, _, err := c.BingGeocodersApi.PATCHBingGeocodersBingGeocoderId(ctx, d.Id()).BingGeocoderUpdate(bingGeocodersUpdate).Execute()

	return diag.FromErr(err)
}
