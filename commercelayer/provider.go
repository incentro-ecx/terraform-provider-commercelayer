package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/incentro-dc/go-commercelayer-sdk/api"
	"golang.org/x/oauth2/clientcredentials"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COMMERCELAYER_CLIENT_ID", nil),
				Description: "The client id of a Commercelayer store",
				Sensitive:   true,
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COMMERCELAYER_CLIENT_SECRET", nil),
				Description: "The client secret of a Commercelayer store",
				Sensitive:   true,
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COMMERCELAYER_API_ENDPOINT", nil),
				Description: "The Commercelayer api endpoint",
			},
			"auth_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("COMMERCELAYER_AUTH_ENDPOINT", nil),
				Description: "The Commercelayer auth endpoint",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"commercelayer_address": resourceAddress(),
		},
		ConfigureContextFunc: providerConfigureFunc,
	}
}

func providerConfigureFunc(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	apiEndpoint := d.Get("api_endpoint").(string)
	authEndpoint := d.Get("auth_endpoint").(string)

	credentials := clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     authEndpoint,
		Scopes:       []string{},
	}

	httpClient := credentials.Client(context.Background())

	commercelayerClient := api.NewAPIClient(&api.Configuration{
		HTTPClient: httpClient,
		Debug:      true,
		Servers: []api.ServerConfiguration{
			{URL: apiEndpoint},
		},
	})

	return commercelayerClient, nil
}
