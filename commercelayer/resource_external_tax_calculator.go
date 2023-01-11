package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceExternalTaxCalculator() *schema.Resource {
	return &schema.Resource{
		Description: "Create an external tax calculator to delegate tax calculation" +
			" logic to the specified external service. Use the order payload to compute " +
			"your own logic and return the tax rate to be applied to the order.",
		ReadContext:   resourceExternalTaxCalculatorReadFunc,
		CreateContext: resourceExternalTaxCalculatorCreateFunc,
		UpdateContext: resourceExternalTaxCalculatorUpdateFunc,
		DeleteContext: resourceExternalTaxCalculatorDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The external tax calculator unique identifier",
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
							Description: "The external tax calculator's internal name",
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
						"tax_calculator_url": {
							Description: "The URL to the service that will compute the taxes.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceExternalTaxCalculatorReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.ExternalTaxCalculatorsApi.GETExternalTaxCalculatorsExternalTaxCalculatorId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	externalTaxCalculator, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(externalTaxCalculator.GetId())

	return nil
}

func resourceExternalTaxCalculatorCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	externalTaxCalculatorCreate := commercelayer.ExternalTaxCalculatorCreate{
		Data: commercelayer.ExternalTaxCalculatorCreateData{
			Type: externalTaxCalculatorType,
			Attributes: commercelayer.POSTExternalTaxCalculators201ResponseDataAttributes{
				Name:             attributes["name"].(string),
				Reference:        stringRef(attributes["reference"]),
				ReferenceOrigin:  stringRef(attributes["reference_origin"]),
				Metadata:         keyValueRef(attributes["metadata"]),
				TaxCalculatorUrl: attributes["tax_calculator_url"].(string),
			},
		},
	}

	err := d.Set("type", externalTaxCalculatorType)
	if err != nil {
		return diagErr(err)
	}

	externalTaxCalculator, _, err := c.ExternalTaxCalculatorsApi.POSTExternalTaxCalculators(ctx).ExternalTaxCalculatorCreate(externalTaxCalculatorCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*externalTaxCalculator.Data.Id)

	return nil
}

func resourceExternalTaxCalculatorDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.ExternalTaxCalculatorsApi.DELETEExternalTaxCalculatorsExternalTaxCalculatorId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceExternalTaxCalculatorUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var ExternalTaxCalculatorUpdate = commercelayer.ExternalTaxCalculatorUpdate{
		Data: commercelayer.ExternalTaxCalculatorUpdateData{
			Type: externalTaxCalculatorType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHExternalTaxCalculatorsExternalTaxCalculatorId200ResponseDataAttributes{
				Name:             stringRef(attributes["name"].(string)),
				Reference:        stringRef(attributes["reference"]),
				ReferenceOrigin:  stringRef(attributes["reference_origin"]),
				Metadata:         keyValueRef(attributes["metadata"]),
				TaxCalculatorUrl: stringRef(attributes["tax_calculator_url"].(string)),
			},
		},
	}

	_, _, err := c.ExternalTaxCalculatorsApi.PATCHExternalTaxCalculatorsExternalTaxCalculatorId(ctx, d.Id()).ExternalTaxCalculatorUpdate(ExternalTaxCalculatorUpdate).Execute()

	return diag.FromErr(err)
}
