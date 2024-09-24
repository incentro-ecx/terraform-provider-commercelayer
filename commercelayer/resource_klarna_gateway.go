package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceKlarnaGateway() *schema.Resource {
	return &schema.Resource{
		Description: "Configuring a Klarna payment gateway for a market lets you safely process payments through Klarna. " +
			"The Klarna gateway is compliant with the PSD2 European regulation so that you can" +
			"implement a payment flow that supports SCA and 3DS2 by using the Klarna's official " +
			"JS SDK and libraries.",
		ReadContext:   resourceKlarnaGatewayReadFunc,
		CreateContext: resourceKlarnaGatewayCreateFunc,
		UpdateContext: resourceKlarnaGatewayUpdateFunc,
		DeleteContext: resourceKlarnaGatewayDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The klarna payment unique identifier",
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
							Description: "The payment gateway's internal name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"country_code": {
							Description: "The gateway country code one of EU, US, or OC.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"api_key": {
							Description: "The public key linked to your API credential.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"api_secret": {
							Description: "The gateway API key.",
							Type:        schema.TypeString,
							Required:    true,
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

func resourceKlarnaGatewayReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.KlarnaGatewaysApi.GETKlarnaGatewaysKlarnaGatewayId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	klarnaGateway, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(klarnaGateway.GetId().(string))

	return nil
}

func resourceKlarnaGatewayCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	klarnaGatewayCreate := commercelayer.KlarnaGatewayCreate{
		Data: commercelayer.KlarnaGatewayCreateData{
			Type: klarnaGatewaysType,
			Attributes: commercelayer.POSTKlarnaGateways201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				CountryCode:     attributes["country_code"].(string),
				ApiKey:          attributes["api_key"].(string),
				ApiSecret:       attributes["api_secret"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	err := d.Set("type", klarnaGatewaysType)
	if err != nil {
		return diagErr(err)
	}

	klarnaGateway, _, err := c.KlarnaGatewaysApi.POSTKlarnaGateways(ctx).KlarnaGatewayCreate(klarnaGatewayCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(klarnaGateway.Data.GetId().(string))

	return nil
}

func resourceKlarnaGatewayDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.KlarnaGatewaysApi.DELETEKlarnaGatewaysKlarnaGatewayId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceKlarnaGatewayUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var klarnaGatewayUpdate = commercelayer.KlarnaGatewayUpdate{
		Data: commercelayer.KlarnaGatewayUpdateData{
			Type: klarnaGatewaysType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHKlarnaGatewaysKlarnaGatewayId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"]),
				CountryCode:     stringRef(attributes["country_code"]),
				ApiKey:          stringRef(attributes["api_key"]),
				ApiSecret:       stringRef(attributes["api_secret"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.KlarnaGatewaysApi.PATCHKlarnaGatewaysKlarnaGatewayId(ctx, d.Id()).
		KlarnaGatewayUpdate(klarnaGatewayUpdate).Execute()

	return diag.FromErr(err)
}
