package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceShippingZone() *schema.Resource {
	return &schema.Resource{
		Description: "Shipping zones determine the available shipping methods for a given shipping address. The " +
			"match is evaluated against a set of regular expressions on the address country, state or zip code.",
		ReadContext:   resourceShippingZoneReadFunc,
		CreateContext: resourceShippingZoneCreateFunc,
		UpdateContext: resourceShippingZoneUpdateFunc,
		DeleteContext: resourceShippingZoneDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The shipping zone unique identifier",
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
							Description: "The shipping zone's internal name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"country_code_regex": {
							Description: "The regex that will be evaluated to match the shipping address country code.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"not_country_code_regex": {
							Description: "The regex that will be evaluated as negative match for the shipping " +
								"address country code.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"state_code_regex": {
							Description: "The regex that will be evaluated to match the shipping address state code.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"not_state_code_regex": {
							Description: "The regex that will be evaluated as negative match for the shipping " +
								"address state code.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"zip_code_regex": {
							Description: "The regex that will be evaluated to match the shipping address zip code.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"not_zip_code_regex": {
							Description: "The regex that will be evaluated as negative match for the shipping zip " +
								"country code.",
							Type:     schema.TypeString,
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
		},
	}
}

func resourceShippingZoneReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.ShippingZonesApi.GETShippingZonesShippingZoneId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	shippingZone, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(shippingZone.GetId().(string))

	return nil
}

func resourceShippingZoneCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	shippingZoneCreate := commercelayer.ShippingZoneCreate{
		Data: commercelayer.ShippingZoneCreateData{
			Type: shippingZoneType,
			Attributes: commercelayer.POSTShippingZones201ResponseDataAttributes{
				Name:                attributes["name"].(string),
				CountryCodeRegex:    stringRef(attributes["country_code_regex"]),
				NotCountryCodeRegex: stringRef(attributes["not_country_code_regex"]),
				StateCodeRegex:      stringRef(attributes["state_code_regex"]),
				NotStateCodeRegex:   stringRef(attributes["not_state_code_regex"]),
				ZipCodeRegex:        stringRef(attributes["zip_code_regex"]),
				NotZipCodeRegex:     stringRef(attributes["not_zip_code_regex"]),
				Reference:           stringRef(attributes["reference"]),
				ReferenceOrigin:     stringRef(attributes["reference_origin"]),
				Metadata:            keyValueRef(attributes["metadata"]),
			},
		},
	}

	err := d.Set("type", shippingZoneType)
	if err != nil {
		return diagErr(err)
	}

	shippingZone, _, err := c.ShippingZonesApi.POSTShippingZones(ctx).ShippingZoneCreate(shippingZoneCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(shippingZone.Data.GetId().(string))

	return nil
}

func resourceShippingZoneDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.ShippingZonesApi.DELETEShippingZonesShippingZoneId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceShippingZoneUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var shippingZoneUpdate = commercelayer.ShippingZoneUpdate{
		Data: commercelayer.ShippingZoneUpdateData{
			Type: shippingZoneType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHShippingZonesShippingZoneId200ResponseDataAttributes{
				Name:                stringRef(attributes["name"]),
				CountryCodeRegex:    stringRef(attributes["country_code_regex"]),
				NotCountryCodeRegex: stringRef(attributes["not_country_code_regex"]),
				StateCodeRegex:      stringRef(attributes["state_code_regex"]),
				NotStateCodeRegex:   stringRef(attributes["not_state_code_regex"]),
				ZipCodeRegex:        stringRef(attributes["zip_code_regex"]),
				NotZipCodeRegex:     stringRef(attributes["not_zip_code_regex"]),
				Reference:           stringRef(attributes["reference"]),
				ReferenceOrigin:     stringRef(attributes["reference_origin"]),
				Metadata:            keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.ShippingZonesApi.PATCHShippingZonesShippingZoneId(ctx, d.Id()).ShippingZoneUpdate(shippingZoneUpdate).Execute()

	return diag.FromErr(err)
}
