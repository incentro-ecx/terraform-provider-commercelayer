package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourcePaypalGateway() *schema.Resource {
	return &schema.Resource{
		Description: "Configuring a PayPal payment gateway for a market lets you safely process payments " +
			"through PayPal.To create a PayPal gateway choose a meaningful name that helps you identify it within " +
			"your organization and connect your PayPal account by adding your client ID and secret " +
			"(contact PayPal's support if you are not sure about the requested data).",
		ReadContext:   resourcePaypalGatewayReadFunc,
		CreateContext: resourcePaypalGatewayCreateFunc,
		UpdateContext: resourcePaypalGatewayUpdateFunc,
		DeleteContext: resourcePaypalGatewayDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The paypal payment unique identifier",
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
						"client_id": {
							Description: "The gateway client ID.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"client_secret": {
							Description: "The gateway client secret.",
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

func resourcePaypalGatewayReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.PaypalGatewaysApi.GETPaypalGatewaysPaypalGatewayId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	paypalGateway, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(paypalGateway.GetId())

	return nil
}

func resourcePaypalGatewayCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	paypalGatewayCreate := commercelayer.PaypalGatewayCreate{
		Data: commercelayer.PaypalGatewayCreateData{
			Type: paypalGatewaysType,
			Attributes: commercelayer.POSTPaypalGateways201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				ClientId:        attributes["client_id"].(string),
				ClientSecret:    attributes["client_secret"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	err := d.Set("type", paypalGatewaysType)
	if err != nil {
		return diagErr(err)
	}

	paypalGateway, _, err := c.PaypalGatewaysApi.POSTPaypalGateways(ctx).PaypalGatewayCreate(paypalGatewayCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*paypalGateway.Data.Id)

	return nil
}

func resourcePaypalGatewayDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.PaypalGatewaysApi.DELETEPaypalGatewaysPaypalGatewayId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourcePaypalGatewayUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var paypalGatewayUpdate = commercelayer.PaypalGatewayUpdate{
		Data: commercelayer.PaypalGatewayUpdateData{
			Type: paypalGatewaysType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHPaypalGatewaysPaypalGatewayId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"]),
				ClientId:        stringRef(attributes["client_id"]),
				ClientSecret:    stringRef(attributes["client_secret"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.PaypalGatewaysApi.PATCHPaypalGatewaysPaypalGatewayId(ctx, d.Id()).
		PaypalGatewayUpdate(paypalGatewayUpdate).Execute()

	return diag.FromErr(err)
}
