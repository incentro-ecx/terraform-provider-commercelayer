package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commercelayer "github.com/incentro-dc/go-commercelayer-sdk/api"
)

func resourceWebhook() *schema.Resource {
	return &schema.Resource{
		Description: "A webhook object is returned as part of the response body of each successful list, retrieve, " +
			"create or update API call to the /api/webhooks endpoint.",
		ReadContext:   resourceWebhookReadFunc,
		CreateContext: resourceWebhookCreateFunc,
		UpdateContext: resourceWebhookUpdateFunc,
		DeleteContext: resourceWebhookDeleteFunc,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Description: "The webhook unique identifier",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: "The resource type",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"shared_secret": {
				Description: "The shared secret used to sign the external request payload.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
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
							Description: "Unique name for the webhook.",
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "webhook",
						},
						"topic": {
							Description: "The identifier of the resource/event that will trigger the webhook.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"callback_url": {
							Description: "URI where the webhook subscription should send the POST request when the " +
								"event occurs.",
							Type:     schema.TypeString,
							Required: true,
						},
						"include_resources": {
							Description: "List of related commercelayer_inventory_stock_location that should be included in the webhook body.",
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional: true,
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

func resourceWebhookReadFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	resp, _, err := c.WebhooksApi.GETWebhooksWebhookId(ctx, d.Id()).Execute()
	if err != nil {
		return diagErr(err)
	}

	webhook, ok := resp.GetDataOk()
	if !ok {
		d.SetId("")
		return nil
	}

	d.SetId(webhook.GetId().(string))

	return nil
}

func resourceWebhookCreateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	webhookCreate := commercelayer.WebhookCreate{
		Data: commercelayer.WebhookCreateData{
			Type: webhookType,
			Attributes: commercelayer.POSTWebhooks201ResponseDataAttributes{
				Name:             stringRef(attributes["name"]),
				Topic:            attributes["topic"].(string),
				CallbackUrl:      attributes["callback_url"].(string),
				IncludeResources: stringSliceValueRef(attributes["include_resources"]),
				Reference:        stringRef(attributes["reference"]),
				ReferenceOrigin:  stringRef(attributes["reference_origin"]),
				Metadata:         keyValueRef(attributes["metadata"]),
			},
		},
	}

	err := d.Set("type", webhookType)
	if err != nil {
		return diagErr(err)
	}

	webhook, _, err := c.WebhooksApi.POSTWebhooks(ctx).WebhookCreate(webhookCreate).Execute()
	if err != nil {
		return diagErr(err)
	}

	d.SetId(webhook.Data.GetId().(string))

	//Fetch the shared secret (this is a work-around because the create does not return it)
	resp, _, err := c.WebhooksApi.GETWebhooksWebhookId(ctx, webhook.Data.GetId().(string)).Execute()
	if err != nil {
		return diagErr(err)
	}

	getWebhook := resp.GetData()

	err = d.Set("shared_secret", &getWebhook.Attributes.SharedSecret)
	if err != nil {
		return diagErr(err)
	}

	return nil
}

func resourceWebhookDeleteFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)
	_, err := c.WebhooksApi.DELETEWebhooksWebhookId(ctx, d.Id()).Execute()
	return diag.FromErr(err)
}

func resourceWebhookUpdateFunc(ctx context.Context, d *schema.ResourceData, i interface{}) diag.Diagnostics {
	c := i.(*commercelayer.APIClient)

	attributes := nestedMap(d.Get("attributes"))

	var webhookUpdate = commercelayer.WebhookUpdate{
		Data: commercelayer.WebhookUpdateData{
			Type: webhookType,
			Id:   d.Id(),
			Attributes: commercelayer.PATCHWebhooksWebhookId200ResponseDataAttributes{
				Name:             stringRef(attributes["name"]),
				Topic:            stringRef(attributes["topic"]),
				CallbackUrl:      stringRef(attributes["callback_url"]),
				IncludeResources: stringSliceValueRef(attributes["include_resources"]),
				Reference:        stringRef(attributes["reference"]),
				ReferenceOrigin:  stringRef(attributes["reference_origin"]),
				Metadata:         keyValueRef(attributes["metadata"]),
			},
		},
	}

	_, _, err := c.WebhooksApi.PATCHWebhooksWebhookId(ctx, d.Id()).WebhookUpdate(webhookUpdate).Execute()

	return diag.FromErr(err)

}
