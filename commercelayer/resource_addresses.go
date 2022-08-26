package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceAddress() *schema.Resource {
	return &schema.Resource{
		Description: "Get notified when specific events occur on a Commercelayer store. For more information, see " +
			"Addresss Overview.",
		ReadContext:   resourceAddressReadFunc,
		CreateContext: resourceAddressCreateFunc,
		UpdateContext: resourceAddressUpdateFunc,
		DeleteContext: resourceAddressDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The address unique identifier",
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
						"business": {
							Description: "Indicates if it's a business or a personal address",
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
						"first_name": {
							Description: "Address first name",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"last_name": {
							Description: "Address last name",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"company": {
							Description: "Address company name",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"line_1": {
							Description: "Address line 1, i.e. Street address, PO Box",
							Type:        schema.TypeString,
							Required:    true,
						},
						"line_2": {
							Description: "Address line 2, i.e. Apartment, Suite, Building",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"city": {
							Description: "Address city",
							Type:        schema.TypeString,
							Required:    true,
						},
						"zip_code": {
							Description: "ZIP or postal code",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"state_code": {
							Description: "State, province or region code",
							Type:        schema.TypeString,
							Required:    true,
						},
						"country_code": {
							Description: "The international 2-letter country code as defined by the ISO 3166-1 standard",
							Type:        schema.TypeString,
							Required:    true,
						},
						"phone": {
							Description: "Phone number (including extension).",
							Type:        schema.TypeString,
							Required:    true,
						},
						"email": {
							Description: "Email address",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"notes": {
							Description: "A free notes attached to the address. When used as a shipping address, this " +
								"can be useful to let the customers add specific delivery instructions.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"lat": {
							Description: "The address geocoded latitude. This is automatically generated when " +
								"creating a shipping/billing address for an order and a valid geocoder is attached to " +
								"the order's market.",
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"lng": {
							Description: "The address geocoded longitude. This is automatically generated when " +
								"creating a shipping/billing address for an order and a valid geocoder is attached " +
								"to the order's market.",
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"billing_info": {
							Description: "Customer's billing information (i.e. VAT number, codice fiscale)",
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
						//	TODO: implement geocoder relation
					},
				},
			},
		},
	}
}

func resourceAddressReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.AddressesApi.GETAddressesAddressId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	address, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(address.GetId())

	return nil
}

func resourceAddressCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := d.Get("attributes").([]interface{})[0].(map[string]interface{})

	addressCreate := commercelayer.AddressCreate{
		Data: commercelayer.AddressCreateData{
			Type: addressType,
			Attributes: commercelayer.POSTAddresses201ResponseDataAttributes{
				Business:        boolRef(attributes["business"]),
				FirstName:       stringRef(attributes["first_name"]),
				LastName:        stringRef(attributes["last_name"]),
				Company:         stringRef(attributes["company"]),
				Line1:           attributes["line_1"].(string),
				Line2:           stringRef(attributes["line_2"]),
				City:            attributes["city"].(string),
				ZipCode:         stringRef(attributes["zip_code"]),
				StateCode:       attributes["state_code"].(string),
				CountryCode:     attributes["country_code"].(string),
				Phone:           attributes["phone"].(string),
				Email:           stringRef(attributes["email"]),
				Notes:           stringRef(attributes["notes"]),
				Lat:             float64ToFloat32Ref(attributes["lat"]),
				Lng:             float64ToFloat32Ref(attributes["lng"]),
				BillingInfo:     stringRef(attributes["billing_info"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			Relationships: nil,
		},
	}

	address, _, err := c.AddressesApi.POSTAddresses(ctx).AddressCreate(addressCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*address.Data.Id)

	return nil
}

func resourceAddressDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.AddressesApi.DELETEAddressesAddressId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceAddressUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := d.Get("attributes").([]interface{})[0].(map[string]interface{})

	var addressUpdate = commercelayer.AddressUpdate{
		Data: commercelayer.AddressUpdateData{
			Type: addressType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHAddressesAddressId200ResponseDataAttributes{
				Business:        boolRef(attributes["business"]),
				FirstName:       stringRef(attributes["first_name"]),
				LastName:        stringRef(attributes["last_name"]),
				Company:         stringRef(attributes["company"]),
				Line1:           stringRef(attributes["line_1"]),
				Line2:           stringRef(attributes["line_2"]),
				City:            stringRef(attributes["city"]),
				ZipCode:         stringRef(attributes["zip_code"]),
				StateCode:       stringRef(attributes["state_code"]),
				CountryCode:     stringRef(attributes["country_code"]),
				Phone:           stringRef(attributes["phone"]),
				Email:           stringRef(attributes["email"]),
				Notes:           stringRef(attributes["notes"]),
				Lat:             float64ToFloat32Ref(attributes["lat"]),
				Lng:             float64ToFloat32Ref(attributes["lng"]),
				BillingInfo:     stringRef(attributes["billing_info"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			Relationships: nil,
		},
	}

	_, _, err := c.AddressesApi.PATCHAddressesAddressId(ctx, d.Id()).AddressUpdate(addressUpdate).Execute()

	return diag.FromErr(err)

}
