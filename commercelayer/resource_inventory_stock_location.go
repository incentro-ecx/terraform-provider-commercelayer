package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceInventoryStockLocation() *schema.Resource {
	return &schema.Resource{
		Description: "Inventory stock locations build a hierarchy of stock locations within an inventory " +
			"model, determining the availability of SKU's that are being purchased. In the case a SKU is available " +
			"in more stock locations, it gets shipped from those with the highest priority.",
		ReadContext:   resourceInventoryStockLocationReadFunc,
		CreateContext: resourceInventoryStockLocationCreateFunc,
		UpdateContext: resourceInventoryStockLocationUpdateFunc,
		DeleteContext: resourceInventoryStockLocationDeleteFunc,
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
							Description: "The stock location priority within the associated inventory model.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"on_hold": {
							Description: "Indicates if the shipment should be put on hold if fulfilled from the " +
								"associated stock location. This is useful to manage use cases like back-orders, " +
								"pre-orders or personalized orders that need to be customized before being fulfilled.",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"reference": {
							Description: "A string that you can use to add any external identifier to the resource. This " +
								"can be useful for integrating the resource to an external system, like an ERP, a " +
								"InventoryStockLocationing tool, a CRM, or whatever.",
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

func resourceInventoryStockLocationReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.InventoryStockLocationsApi.GETInventoryStockLocationsInventoryStockLocationId(ctx, d.Id()).Execute()
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

func resourceInventoryStockLocationCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	inventoryModelCreate := commercelayer.InventoryStockLocationCreate{
		Data: commercelayer.InventoryStockLocationCreateData{
			Type: inventoryStockLocationsType,
			Attributes: commercelayer.POSTInventoryStockLocations201ResponseDataAttributes{
				Priority:        int32(attributes["priority"].(int)),
				OnHold:          boolRef(attributes["on_hold"]),
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

	err := d.Set("type", inventoryStockLocationsType)
	if err != nil {
		return diagErr(err)
	}

	inventoryModel, _, err := c.InventoryStockLocationsApi.POSTInventoryStockLocations(ctx).InventoryStockLocationCreate(inventoryModelCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(inventoryModel.Data.GetId().(string))

	return nil
}

func resourceInventoryStockLocationDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.InventoryStockLocationsApi.DELETEInventoryStockLocationsInventoryStockLocationId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceInventoryStockLocationUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	var inventoryModelUpdate = commercelayer.InventoryStockLocationUpdate{
		Data: commercelayer.InventoryStockLocationUpdateData{
			Type: inventoryStockLocationsType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHInventoryStockLocationsInventoryStockLocationId200ResponseDataAttributes{
				Priority:        intToInt32Ref(attributes["priority"]),
				OnHold:          boolRef(attributes["on_hold"]),
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

	_, _, err := c.InventoryStockLocationsApi.PATCHInventoryStockLocationsInventoryStockLocationId(ctx, d.Id()).
		InventoryStockLocationUpdate(inventoryModelUpdate).Execute()

	return diag.FromErr(err)
}
