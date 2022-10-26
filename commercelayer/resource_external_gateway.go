package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceExternalGateway() *schema.Resource {
	return &schema.Resource{
		Description: `Price lists are collections of SKU prices, 
		defined by currency and market. When a list of SKUs is fetched, 
		only SKUs with a price defined in the market's price list and at least 
		a stock item in one of the market stock locations will be returned. 
		A user can create price lists to manage international business or B2B/B2C models.`,
		ReadContext:   resourceExternalGatewayReadFunc,
		CreateContext: resourceExternalGatewayCreateFunc,
		UpdateContext: resourceExternalGatewayUpdateFunc,
		DeleteContext: resourceExternalGatewayDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The external gateway unique identifier",
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
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The payment gateway's internal name.",
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
						"authorize_url": {
							Description: "The endpoint used by the external gateway to authorize payments.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"capture_url": {
							Description: "The endpoint used by the external gateway to capture payments.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"void_url": {
							Description: "The endpoint used by the external gateway to void payments.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"refund_url": {
							Description: "The endpoint used by the external gateway to refund payments.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"token_url": {
							Description: "The endpoint used by the external gateway to create a customer payment token.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceExternalGatewayReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.ExternalGatewaysApi.GETExternalGatewaysExternalGatewayId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	externalGateway, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(externalGateway.GetId())

	return nil
}

func resourceExternalGatewayCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	externalGatewayCreate := commercelayer.ExternalGatewayCreate{
		Data: commercelayer.ExternalGatewayCreateData{
			Type: externalGatewayType,
			Attributes: commercelayer.POSTExternalGateways201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
				AuthorizeUrl:    stringRef(attributes["authorize_url"]),
				CaptureUrl:      stringRef(attributes["capture_url"]),
				VoidUrl:         stringRef(attributes["void_url"]),
				TokenUrl:        stringRef(attributes["token_url"]),
				RefundUrl:       stringRef(attributes["refund_url"]),
			},
		},
	}

	err := d.Set("type", externalGatewayType)
	if err != nil {
		return diagErr(err)
	}

	externalGateway, _, err := c.ExternalGatewaysApi.POSTExternalGateways(ctx).ExternalGatewayCreate(externalGatewayCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*externalGateway.Data.Id)

	return nil
}

func resourceExternalGatewayDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.ExternalGatewaysApi.DELETEExternalGatewaysExternalGatewayId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceExternalGatewayUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var externalGatewayUpdate = commercelayer.ExternalGatewayUpdate{
		Data: commercelayer.ExternalGatewayUpdateData{
			Type: externalGatewayType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHExternalGatewaysExternalGatewayId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"].(string)),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
				AuthorizeUrl:    stringRef(attributes["authorize_url"]),
				CaptureUrl:      stringRef(attributes["capture_url"]),
				VoidUrl:         stringRef(attributes["void_url"]),
				TokenUrl:        stringRef(attributes["token_url"]),
				RefundUrl:       stringRef(attributes["refund_url"]),
			},
		},
	}

	_, _, err := c.ExternalGatewaysApi.PATCHExternalGatewaysExternalGatewayId(ctx, d.Id()).ExternalGatewayUpdate(externalGatewayUpdate).Execute()

	return diag.FromErr(err)
}
