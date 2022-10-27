package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceInventoryModel() *schema.Resource {
	return &schema.Resource{
		Description: "An inventory model defines a list of stock locations ordered by priority. The priority and " +
			"cutoff determine how the availability of SKU's gets calculated within a market.",
		ReadContext:   resourceInventoryModelReadFunc,
		CreateContext: resourceInventoryModelCreateFunc,
		UpdateContext: resourceInventoryModelUpdateFunc,
		DeleteContext: resourceInventoryModelDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The inventory model unique identifier",
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
							Description: "The inventory model's internal name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"strategy": {
							Description: "The inventory model's shipping strategy: one between 'no_split' (default), " +
								"'split_shipments', 'ship_from_primary' and 'ship_from_first_available_or_primary'.",
							Type:             schema.TypeString,
							Default:          "no_split",
							Optional:         true,
							ValidateDiagFunc: inventoryModelStrategyValidation,
						},
						"stock_locations_cutoff": {
							Description: "The maximum number of stock locations used for inventory computation",
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     2,
						},
						"reference": {
							Description: "A string that you can use to add any external identifier to the resource. This " +
								"can be useful for integrating the resource to an external system, like an ERP, a " +
								"InventoryModeling tool, a CRM, or whatever.",
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

func resourceInventoryModelReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.InventoryModelsApi.GETInventoryModelsInventoryModelId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	inventoryModel, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(inventoryModel.GetId())

	return nil
}

func resourceInventoryModelCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	inventoryModelCreate := commercelayer.InventoryModelCreate{
		Data: commercelayer.InventoryModelCreateData{
			Type: inventoryModelType,
			Attributes: commercelayer.POSTInventoryModels201ResponseDataAttributes{
				Name:                 attributes["name"].(string),
				Strategy:             stringRef(attributes["strategy"]),
				StockLocationsCutoff: intToInt32Ref(attributes["stock_locations_cutoff"]),
				Reference:            stringRef(attributes["reference"]),
				ReferenceOrigin:      stringRef(attributes["reference_origin"]),
				Metadata:             keyValueRef(attributes["metadata"]),
			},
		},
	}

	err := d.Set("type", inventoryModelType)
	if err != nil {
		return diagErr(err)
	}

	inventoryModel, _, err := c.InventoryModelsApi.POSTInventoryModels(ctx).InventoryModelCreate(inventoryModelCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(*inventoryModel.Data.Id)

	return nil
}

func resourceInventoryModelDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.InventoryModelsApi.DELETEInventoryModelsInventoryModelId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceInventoryModelUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var inventoryModelUpdate = commercelayer.InventoryModelUpdate{
		Data: commercelayer.InventoryModelUpdateData{
			Type: inventoryModelType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHInventoryModelsInventoryModelId200ResponseDataAttributes{
				Name:                 stringRef(attributes["name"]),
				Strategy:             stringRef(attributes["strategy"]),
				StockLocationsCutoff: intToInt32Ref(attributes["stock_locations_cutoff"]),
				Reference:            stringRef(attributes["reference"]),
				ReferenceOrigin:      stringRef(attributes["reference_origin"]),
				Metadata:             keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.InventoryModelsApi.PATCHInventoryModelsInventoryModelId(ctx, d.Id()).
		InventoryModelUpdate(inventoryModelUpdate).Execute()

	return diag.FromErr(err)
}
