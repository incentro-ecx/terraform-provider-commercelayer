package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceAdyenGateway() *schema.Resource {
	return &schema.Resource{
		Description: "Configuring a Adyen payment gateway for a market lets you safely process payments through Adyen. " +
			"The Adyen gateway is compliant with the PSD2 European regulation so that you can implement a payment flow " +
			"that supports SCA and 3DS2 by using the Adyen's official JS SDK and libraries." +
			"To create a Adyen gateway choose a meaningful name that helps you identify it within your organization and gather all the credentials requested " +
			"(like secret and publishable keys, etc. — contact Adyen's support if you are not sure about the requested data).",
		ReadContext:   resourceAdyenGatewayReadFunc,
		CreateContext: resourceAdyenGatewayCreateFunc,
		UpdateContext: resourceAdyenGatewayUpdateFunc,
		DeleteContext: resourceAdyenGatewayDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The adyen payment unique identifier",
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
						"merchant_account": {
							Description: "The gateway merchant account.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"api_key": {
							Description: "The gateway API key.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"public_key": {
							Description: "The public key linked to your API credential.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"live_url_prefix": {
							Description: "The prefix of the endpoint used for live transactions.",
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

func resourceAdyenGatewayReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.AdyenGatewaysApi.GETAdyenGatewaysAdyenGatewayId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	adyenGateway, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(adyenGateway.GetId())

	return nil
}

func resourceAdyenGatewayCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	adyenGatewayCreate := commercelayer.AdyenGatewayCreate{
		Data: commercelayer.AdyenGatewayCreateData{
			Type: adyenGatewaysType,
			Attributes: commercelayer.POSTAdyenGateways201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				MerchantAccount: attributes["merchant_account"].(string),
				ApiKey:          attributes["api_key"].(string),
				LiveUrlPrefix:   attributes["live_url_prefix"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	err := d.Set("type", adyenGatewaysType)
	if err != nil {
		return diagErr(err)
	}

	adyenGateway, _, err := c.AdyenGatewaysApi.POSTAdyenGateways(ctx).AdyenGatewayCreate(adyenGatewayCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*adyenGateway.Data.Id)

	return nil
}

func resourceAdyenGatewayDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.AdyenGatewaysApi.DELETEAdyenGatewaysAdyenGatewayId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceAdyenGatewayUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var adyenGatewayUpdate = commercelayer.AdyenGatewayUpdate{
		Data: commercelayer.AdyenGatewayUpdateData{
			Type: adyenGatewaysType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHAdyenGatewaysAdyenGatewayId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"]),
				MerchantAccount: stringRef(attributes["merchant_account"]),
				ApiKey:          stringRef(attributes["api_key"]),
				LiveUrlPrefix:   stringRef(attributes["live_url_prefix"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.AdyenGatewaysApi.PATCHAdyenGatewaysAdyenGatewayId(ctx, d.Id()).
		AdyenGatewayUpdate(adyenGatewayUpdate).Execute()

	return diag.FromErr(err)
}
