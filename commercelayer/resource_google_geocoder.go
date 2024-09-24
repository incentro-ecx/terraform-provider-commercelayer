package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceGoogleGeocoders() *schema.Resource {
	return &schema.Resource{
		Description: "Geocoders lets you convert an address in text form into geographic coordinates " +
			"(like latitude and longitude). A geocoder can be optionally associated with a market. " +
			"By connecting a geocoder to a market, all the shipping and billing addresses belonging " +
			"to that market orders will be geocoded",
		ReadContext:   resourceGoogleGeocodersReadFunc,
		CreateContext: resourceGoogleGeocodersCreateFunc,
		UpdateContext: resourceGoogleGeocodersUpdateFunc,
		DeleteContext: resourceGoogleGeocodersDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The google geocoder unique identifier",
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
						"api_key": {
							Description: "The Google Map API key",
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

func resourceGoogleGeocodersReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.GoogleGeocodersApi.GETGoogleGeocodersGoogleGeocoderId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	googleGeocoder, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(googleGeocoder.GetId().(string))

	return nil
}

func resourceGoogleGeocodersCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	googleGeocoderCreate := commercelayer.GoogleGeocoderCreate{
		Data: commercelayer.GoogleGeocoderCreateData{
			Type: googleGeocodersType,
			Attributes: commercelayer.POSTGoogleGeocoders201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
				ApiKey:          attributes["api_key"].(string),
			},
		},
	}

	err := d.Set("type", googleGeocodersType)
	if err != nil {
		return diagErr(err)
	}

	googleGeocoders, _, err := c.GoogleGeocodersApi.POSTGoogleGeocoders(ctx).GoogleGeocoderCreate(googleGeocoderCreate).Execute()

	if err != nil {
		return diagErr(err)
	}

	d.SetId(googleGeocoders.Data.GetId().(string))

	return nil
}

func resourceGoogleGeocodersDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.GoogleGeocodersApi.DELETEGoogleGeocodersGoogleGeocoderId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceGoogleGeocodersUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var googleGeocodersUpdate = commercelayer.GoogleGeocoderUpdate{
		Data: commercelayer.GoogleGeocoderUpdateData{
			Type: googleGeocodersType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHGoogleGeocodersGoogleGeocoderId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
				ApiKey:          stringRef(attributes["api_key"]),
			},
		},
	}

	_, _, err := c.GoogleGeocodersApi.PATCHGoogleGeocodersGoogleGeocoderId(ctx, d.Id()).GoogleGeocoderUpdate(googleGeocodersUpdate).Execute()

	return diag.FromErr(err)
}
