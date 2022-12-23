package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceCheckoutComGateway() *schema.Resource {
	return &schema.Resource{
		Description: "Configuring a CheckoutCom payment gateway for a market lets you safely process payments through CheckoutCom. " +
			"The CheckoutCom gateway is compliant with the PSD2 European regulation so that you can" +
			"implement a payment flow that supports SCA and 3DS2 by using the CheckoutCom's official " +
			"JS SDK and libraries.",
		ReadContext:   resourceCheckoutComGatewayReadFunc,
		CreateContext: resourceCheckoutComGatewayCreateFunc,
		UpdateContext: resourceCheckoutComGatewayUpdateFunc,
		DeleteContext: resourceCheckoutComGatewayDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The checkoutCom payment unique identifier",
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
						"secret_key": {
							Description: "The gateway secret key.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"public_key": {
							Description: "The gateway public key.",
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

func resourceCheckoutComGatewayReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.CheckoutComGatewaysApi.GETCheckoutComGatewaysCheckoutComGatewayId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	checkoutComGateway, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(checkoutComGateway.GetId())

	return nil
}

func resourceCheckoutComGatewayCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	checkoutComGatewayCreate := commercelayer.CheckoutComGatewayCreate{
		Data: commercelayer.CheckoutComGatewayCreateData{
			Type: checkoutComGatewaysType,
			Attributes: commercelayer.POSTCheckoutComGateways201ResponseDataAttributes{
				Name: attributes["name"].(string),

				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	err := d.Set("type", checkoutComGatewaysType)
	if err != nil {
		return diagErr(err)
	}

	checkoutComGateway, _, err := c.CheckoutComGatewaysApi.POSTCheckoutComGateways(ctx).CheckoutComGatewayCreate(checkoutComGatewayCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*checkoutComGateway.Data.Id)

	return nil
}

func resourceCheckoutComGatewayDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.CheckoutComGatewaysApi.DELETECheckoutComGatewaysCheckoutComGatewayId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceCheckoutComGatewayUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var checkoutComGatewayUpdate = commercelayer.CheckoutComGatewayUpdate{
		Data: commercelayer.CheckoutComGatewayUpdateData{
			Type: checkoutComGatewaysType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHCheckoutComGatewaysCheckoutComGatewayId200ResponseDataAttributes{
				Name: stringRef(attributes["name"]),

				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.CheckoutComGatewaysApi.PATCHCheckoutComGatewaysCheckoutComGatewayId(ctx, d.Id()).
		CheckoutComGatewayUpdate(checkoutComGatewayUpdate).Execute()

	return diag.FromErr(err)
}
