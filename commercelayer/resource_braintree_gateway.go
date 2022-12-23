package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceBraintreeGateway() *schema.Resource {
	return &schema.Resource{
		Description: "Configuring a Braintree payment gateway for a market lets you safely process payments " +
			"through Braintree. The Braintree gateway is compliant with the PSD2 European regulation " +
			"so that you can implement a payment flow that supports SCA and 3DS2 by using the " +
			"Braintree official JS SDK and libraries.",
		ReadContext:   resourceBraintreeGatewayReadFunc,
		CreateContext: resourceBraintreeGatewayCreateFunc,
		UpdateContext: resourceBraintreeGatewayUpdateFunc,
		DeleteContext: resourceBraintreeGatewayDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The braintree payment unique identifier",
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
						"merchant_account_id": {
							Description: "The gateway merchant account ID.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"merchant_id": {
							Description: "The gateway merchant ID.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"public_key": {
							Description: "The gateway API public key.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"private_key": {
							Description: "The gateway API private key.",
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

func resourceBraintreeGatewayReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.BraintreeGatewaysApi.GETBraintreeGatewaysBraintreeGatewayId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	braintreeGateway, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(braintreeGateway.GetId())

	return nil
}

func resourceBraintreeGatewayCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	braintreeGatewayCreate := commercelayer.BraintreeGatewayCreate{
		Data: commercelayer.BraintreeGatewayCreateData{
			Type: braintreeGatewaysType,
			Attributes: commercelayer.POSTBraintreeGateways201ResponseDataAttributes{
				Name:              attributes["name"].(string),
				MerchantAccountId: attributes["merchant_account_id"].(string),
				MerchantId:        attributes["merchant_id"].(string),
				PublicKey:         attributes["public_key"].(string),
				PrivateKey:        attributes["private_key"].(string),
				Reference:         stringRef(attributes["reference"]),
				ReferenceOrigin:   stringRef(attributes["reference_origin"]),
				Metadata:          keyValueRef(attributes["metadata"]),
			},
		},
	}

	err := d.Set("type", braintreeGatewaysType)
	if err != nil {
		return diagErr(err)
	}

	braintreeGateway, _, err := c.BraintreeGatewaysApi.POSTBraintreeGateways(ctx).BraintreeGatewayCreate(braintreeGatewayCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*braintreeGateway.Data.Id)

	return nil
}

func resourceBraintreeGatewayDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.BraintreeGatewaysApi.DELETEBraintreeGatewaysBraintreeGatewayId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceBraintreeGatewayUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var braintreeGatewayUpdate = commercelayer.BraintreeGatewayUpdate{
		Data: commercelayer.BraintreeGatewayUpdateData{
			Type: braintreeGatewaysType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHBraintreeGatewaysBraintreeGatewayId200ResponseDataAttributes{
				Name:              stringRef(attributes["name"]),
				MerchantAccountId: stringRef(attributes["merchant_account_id"]),
				MerchantId:        stringRef(attributes["merchant_id"]),
				PublicKey:         stringRef(attributes["public_key"]),
				PrivateKey:        stringRef(attributes["private_key"]),
				Reference:         stringRef(attributes["reference"]),
				ReferenceOrigin:   stringRef(attributes["reference_origin"]),
				Metadata:          keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.BraintreeGatewaysApi.PATCHBraintreeGatewaysBraintreeGatewayId(ctx, d.Id()).
		BraintreeGatewayUpdate(braintreeGatewayUpdate).Execute()

	return diag.FromErr(err)
}
