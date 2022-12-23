package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceStripeGateway() *schema.Resource {
	return &schema.Resource{
		Description: "Configuring a Stripe payment gateway for a market lets you safely process payments through Stripe. " +
			"The Stripe gateway is compliant with the PSD2 European regulation so that you can implement a payment flow " +
			"that supports SCA and 3DS2 by using the Stripe's official JS SDK and libraries." +
			"To create a Stripe gateway choose a meaningful name that helps you identify it within your organization and gather all the credentials requested " +
			"(like secret and publishable keys, etc. — contact Stripe's support if you are not sure about the requested data).",
		ReadContext:   resourceStripeGatewayReadFunc,
		CreateContext: resourceStripeGatewayCreateFunc,
		UpdateContext: resourceStripeGatewayUpdateFunc,
		DeleteContext: resourceStripeGatewayDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The stripe payment unique identifier",
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
						"login": {
							Description: "The gateway login.",
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

func resourceStripeGatewayReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.StripeGatewaysApi.GETStripeGatewaysStripeGatewayId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	stripeGateway, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(stripeGateway.GetId())

	return nil
}

func resourceStripeGatewayCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	stripeGatewayCreate := commercelayer.StripeGatewayCreate{
		Data: commercelayer.StripeGatewayCreateData{
			Type: stripeGatewaysType,
			Attributes: commercelayer.POSTStripeGateways201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				Login:           attributes["login"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	err := d.Set("type", stripeGatewaysType)
	if err != nil {
		return diagErr(err)
	}

	stripeGateway, _, err := c.StripeGatewaysApi.POSTStripeGateways(ctx).StripeGatewayCreate(stripeGatewayCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*stripeGateway.Data.Id)

	return nil
}

func resourceStripeGatewayDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.StripeGatewaysApi.DELETEStripeGatewaysStripeGatewayId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceStripeGatewayUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var stripeGatewayUpdate = commercelayer.StripeGatewayUpdate{
		Data: commercelayer.StripeGatewayUpdateData{
			Type: stripeGatewaysType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHStripeGatewaysStripeGatewayId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.StripeGatewaysApi.PATCHStripeGatewaysStripeGatewayId(ctx, d.Id()).
		StripeGatewayUpdate(stripeGatewayUpdate).Execute()

	return diag.FromErr(err)
}
