package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceInventoryReturnLocation() *schema.Resource {
	return &schema.Resource{
		Description: "Inventory return locations build a hierarchy of stock locations within an inventory " +
			"model, determining the available options for the returns. In the case a SKU is available in more stock " +
			"locations, it gets returned to those with the highest priority.",
		ReadContext:   resourceInventoryReturnLocationReadFunc,
		CreateContext: resourceInventoryReturnLocationCreateFunc,
		UpdateContext: resourceInventoryReturnLocationUpdateFunc,
		DeleteContext: resourceInventoryReturnLocationDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The inventory return location unique identifier",
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
						"priority": {
							Description: "The inventory model's internal name.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"reference": {
							Description: "A string that you can use to add any external identifier to the resource. This " +
								"can be useful for integrating the resource to an external system, like an ERP, a " +
								"InventoryReturnLocationing tool, a CRM, or whatever.",
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
			"relationships": {
				Description: "Resource relationships",
				Type:        schema.TypeList,
				MaxItems:    1,
				MinItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"stock_location_id": {
							Description: "The associated stock location id.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"inventory_model_id": {
							Description: "The associated inventory model id.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceInventoryReturnLocationReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.InventoryReturnLocationsApi.GETInventoryReturnLocationsInventoryReturnLocationId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	inventoryModel, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(inventoryModel.GetId().(string))

	return nil
}

func resourceInventoryReturnLocationCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	inventoryModelCreate := commercelayer.InventoryReturnLocationCreate{
		Data: commercelayer.InventoryReturnLocationCreateData{
			Type: inventoryReturnLocationsType,
			Attributes: commercelayer.POSTInventoryReturnLocations201ResponseDataAttributes{
				Priority:        int32(attributes["priority"].(int)),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.InventoryReturnLocationCreateDataRelationships{
				StockLocation: commercelayer.DeliveryLeadTimeCreateDataRelationshipsStockLocation{
					Data: commercelayer.DeliveryLeadTimeDataRelationshipsStockLocationData{
						Type: stringRef(stockLocationType),
						Id:   stringRef(relationships["stock_location_id"]),
					},
				},
				InventoryModel: commercelayer.InventoryReturnLocationCreateDataRelationshipsInventoryModel{
					Data: commercelayer.InventoryReturnLocationDataRelationshipsInventoryModelData{
						Type: stringRef(inventoryModelType),
						Id:   stringRef(relationships["inventory_model_id"]),
					},
				},
			},
		},
	}

	err := d.Set("type", inventoryReturnLocationsType)
	if err != nil {
		return diagErr(err)
	}

	inventoryModel, _, err := c.InventoryReturnLocationsApi.POSTInventoryReturnLocations(ctx).InventoryReturnLocationCreate(inventoryModelCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(inventoryModel.Data.GetId().(string))

	return nil
}

func resourceInventoryReturnLocationDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.InventoryReturnLocationsApi.DELETEInventoryReturnLocationsInventoryReturnLocationId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceInventoryReturnLocationUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	var inventoryModelUpdate = commercelayer.InventoryReturnLocationUpdate{
		Data: commercelayer.InventoryReturnLocationUpdateData{
			Type: inventoryReturnLocationsType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHInventoryReturnLocationsInventoryReturnLocationId200ResponseDataAttributes{
				Priority:        intToInt32Ref(attributes["priority"]),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.InventoryReturnLocationUpdateDataRelationships{
				StockLocation: &commercelayer.DeliveryLeadTimeCreateDataRelationshipsStockLocation{
					Data: commercelayer.DeliveryLeadTimeDataRelationshipsStockLocationData{
						Type: stringRef(stockLocationType),
						Id:   stringRef(relationships["stock_location_id"]),
					},
				},
				InventoryModel: &commercelayer.InventoryReturnLocationCreateDataRelationshipsInventoryModel{
					Data: commercelayer.InventoryReturnLocationDataRelationshipsInventoryModelData{
						Type: stringRef(inventoryModelType),
						Id:   stringRef(relationships["inventory_model_id"]),
					},
				},
			},
		},
	}

	_, _, err := c.InventoryReturnLocationsApi.PATCHInventoryReturnLocationsInventoryReturnLocationId(ctx, d.Id()).
		InventoryReturnLocationUpdate(inventoryModelUpdate).Execute()

	return diag.FromErr(err)
}
