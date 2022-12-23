package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceManualGateway() *schema.Resource {
	return &schema.Resource{
		Description: "An manual payment defines a list of stock locations ordered by priority. The priority and " +
			"cutoff determine how the availability of SKU's gets calculated within a market.",
		ReadContext:   resourceManualGatewayReadFunc,
		CreateContext: resourceManualGatewayCreateFunc,
		UpdateContext: resourceManualGatewayUpdateFunc,
		DeleteContext: resourceManualGatewayDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The manual payment unique identifier",
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
						"reference": {
							Description: "A string that you can use to add any external identifier to the resource. This " +
								"can be useful for integrating the resource to an external system, like an ERP, a " +
								"ManualPaymentGatewaying tool, a CRM, or whatever.",
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

func resourceManualGatewayReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.ManualGatewaysApi.GETManualGatewaysManualGatewayId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	manualGateway, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(manualGateway.GetId())

	return nil
}

func resourceManualGatewayCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	manualGatewayCreate := commercelayer.ManualGatewayCreate{
		Data: commercelayer.ManualGatewayCreateData{
			Type: manualGatewayType,
			Attributes: commercelayer.POSTManualGateways201ResponseDataAttributes{
				Name:            attributes["name"].(string),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	err := d.Set("type", manualGatewayType)
	if err != nil {
		return diagErr(err)
	}

	manualGateway, _, err := c.ManualGatewaysApi.POSTManualGateways(ctx).ManualGatewayCreate(manualGatewayCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*manualGateway.Data.Id)

	return nil
}

func resourceManualGatewayDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.ManualGatewaysApi.DELETEManualGatewaysManualGatewayId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceManualGatewayUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var manualGatewayUpdate = commercelayer.ManualGatewayUpdate{
		Data: commercelayer.ManualGatewayUpdateData{
			Type: manualGatewayType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHManualGatewaysManualGatewayId200ResponseDataAttributes{
				Name:            stringRef(attributes["name"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.ManualGatewaysApi.PATCHManualGatewaysManualGatewayId(ctx, d.Id()).
		ManualGatewayUpdate(manualGatewayUpdate).Execute()

	return diag.FromErr(err)
}
