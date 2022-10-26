package commercelayer

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/incentro-dc/go-commercelayer-sdk/api"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

var baseProvider = &schema.Provider{
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
		"commercelayer_address":                 resourceAddress(),
		"commercelayer_merchant":                resourceMerchant(),
		"commercelayer_price_list":              resourcePriceList(),
		"commercelayer_customer_group":          resourceCustomerGroup(),
		"commercelayer_webhook":                 resourceWebhook(),
		"commercelayer_external_gateway":        resourceExternalGateway(),
		"commercelayer_external_tax_calculator": resourceExternalTaxCalculator(),
	},
}

type Configuration struct {
	cached        bool
	cacheLocation string
}

type ProviderOption func(configuration *Configuration)

func WithCachedToken(cacheLocation string) ProviderOption {
	return func(c *Configuration) {
		c.cached = true
		c.cacheLocation = cacheLocation
	}
}

func Provider(opts ...ProviderOption) plugin.ProviderFunc {
	c := Configuration{
		cached: false,
	}

	for _, opt := range opts {
		opt(&c)
	}

	return func() *schema.Provider {
		baseProvider.ConfigureContextFunc = c.configureFunc
		return baseProvider
	}
}

func (c *Configuration) configureFunc(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
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

	newCtx := context.Background()

	var tokenSource = credentials.TokenSource(newCtx)
	if c.cached {
		tokenSource = newCachedTokenSource(tokenSource, c.cacheLocation)
		tokenSource.Token()
	}

	httpClient := oauth2.NewClient(newCtx, tokenSource)

	commercelayerClient := api.NewAPIClient(&api.Configuration{
		HTTPClient: httpClient,
		Debug:      true,
		Servers: []api.ServerConfiguration{
			{URL: apiEndpoint},
		},
	})

	return commercelayerClient, nil
}
