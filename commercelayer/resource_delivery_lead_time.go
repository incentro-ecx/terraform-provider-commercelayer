package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceDeliveryLeadTime() *schema.Resource {
	return &schema.Resource{
		Description: "Delivery lead times provide customers with detailed information about their shipments. " +
			"This is useful if you ship from many stock locations or offer more shipping method options within a market.",
		ReadContext:   resourceDeliveryLeadTimesReadFunc,
		CreateContext: resourceDeliveryLeadTimesCreateFunc,
		UpdateContext: resourceDeliveryLeadTimesUpdateFunc,
		DeleteContext: resourceDeliveryLeadTimesDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The delivery lead time unique identifier",
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
						"min_hours": {
							Description: "The delivery lead minimum time (in hours) when shipping from the associated stock location with the associated shipping method.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"max_hours": {
							Description: "The delivery lead maximum time (in hours) when shipping from the associated stock location with the associated shipping method.",
							Type:        schema.TypeInt,
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
						"shipping_method_id": {
							Description: "The associated shipping method id.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
		},
	}
}

func resourceDeliveryLeadTimesReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.DeliveryLeadTimesApi.GETDeliveryLeadTimesDeliveryLeadTimeId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	deliveryLeadTime, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(deliveryLeadTime.GetId())

	return nil
}

func resourceDeliveryLeadTimesCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	deliveryLeadTimeCreate := commercelayer.DeliveryLeadTimeCreate{
		Data: commercelayer.DeliveryLeadTimeCreateData{
			Type: deliveryLeadTimesType,
			Attributes: commercelayer.POSTDeliveryLeadTimes201ResponseDataAttributes{
				MinHours:        attributes["min_hours"].(int32),
				MaxHours:        attributes["max_hours"].(int32),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.DeliveryLeadTimeCreateDataRelationships{
				StockLocation: commercelayer.DeliveryLeadTimeDataRelationshipsStockLocation{
					Data: commercelayer.DeliveryLeadTimeDataRelationshipsStockLocationData{
						Type: stringRef(stockLocationType),
						Id:   stringRef(relationships["stock_location_id"]),
					}},
				ShippingMethod: commercelayer.DeliveryLeadTimeDataRelationshipsShippingMethod{Data: commercelayer.DeliveryLeadTimeDataRelationshipsShippingMethodData{
					Type: stringRef(shippingMethodType),
					Id:   stringRef(relationships["shipping_method_id"]),
				}},
			},
		},
	}

	err := d.Set("type", deliveryLeadTimesType)
	if err != nil {
		return diagErr(err)
	}

	deliveryLeadTimes, _, err := c.DeliveryLeadTimesApi.POSTDeliveryLeadTimes(ctx).DeliveryLeadTimeCreate(deliveryLeadTimeCreate).Execute()

	if err != nil {
		return diagErr(err)
	}

	d.SetId(*deliveryLeadTimes.Data.Id)

	return nil
}

func resourceDeliveryLeadTimesDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.DeliveryLeadTimesApi.DELETEDeliveryLeadTimesDeliveryLeadTimeId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceDeliveryLeadTimesUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))
	relationships := nestedMap(d.Get("relationships"))

	var deliveryLeadTimeUpdate = commercelayer.DeliveryLeadTimeUpdate{
		Data: commercelayer.DeliveryLeadTimeUpdateData{
			Type: deliveryLeadTimesType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHDeliveryLeadTimesDeliveryLeadTimeId200ResponseDataAttributes{
				MinHours:        intToInt32Ref(attributes["min_hours"].(int32)),
				MaxHours:        intToInt32Ref(attributes["max_hours"].(int32)),
				Reference:       stringRef(attributes["reference"]),
				ReferenceOrigin: stringRef(attributes["reference_origin"]),
				Metadata:        keyValueRef(attributes["metadata"]),
			},
			Relationships: &commercelayer.DeliveryLeadTimeUpdateDataRelationships{
				StockLocation: &commercelayer.DeliveryLeadTimeDataRelationshipsStockLocation{
					Data: commercelayer.DeliveryLeadTimeDataRelationshipsStockLocationData{
						Type: stringRef(stockLocationType),
						Id:   stringRef(relationships["stock_location_id"]),
					}},
				ShippingMethod: &commercelayer.DeliveryLeadTimeDataRelationshipsShippingMethod{Data: commercelayer.DeliveryLeadTimeDataRelationshipsShippingMethodData{
					Type: stringRef(shippingMethodType),
					Id:   stringRef(relationships["shipping_method_id"]),
				}},
			},
		},
	}

	_, _, err := c.DeliveryLeadTimesApi.PATCHDeliveryLeadTimesDeliveryLeadTimeId(ctx, d.Id()).DeliveryLeadTimeUpdate(deliveryLeadTimeUpdate).Execute()

	return diag.FromErr(err)
}
