package commercelayer

import (
	"bytes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"os"
	"testing"
	"text/template"
)

func init() {
	testAccProviderCommercelayer = Provider()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"commercelayer": func() (*schema.Provider, error) {
			return testAccProviderCommercelayer, nil
		},
	}
}

var testAccProviderCommercelayer *schema.Provider
var testAccProviderFactories = map[string]func() (*schema.Provider, error){}

func TestProvider(t *testing.T) {
	provider := Provider()
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	requiredEnvs := []string{
		"COMMERCELAYER_CLIENT_ID",
		"COMMERCELAYER_CLIENT_SECRET",
		"COMMERCELAYER_API_ENDPOINT",
		"COMMERCELAYER_AUTH_ENDPOINT",
	}
	for _, val := range requiredEnvs {
		if os.Getenv(val) == "" {
			t.Fatalf("%v must be set for acceptance tests", val)
		}
	}
}

func hclTemplate(data string, params map[string]any) string {
	var out bytes.Buffer
	tmpl := template.Must(template.New("hcl").Parse(data))
	err := tmpl.Execute(&out, params)
	if err != nil {
		log.Fatal(err)
	}
	return out.String()
}
