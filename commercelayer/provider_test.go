package commercelayer

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"net/http"
	"os"
	"testing"
	"text/template"
	"time"
)

var testAccProviderCommercelayer *schema.Provider
var testAccProviderFactories = map[string]func() (*schema.Provider, error){}

type AcceptanceSuite struct {
	suite.Suite
}

func (s *AcceptanceSuite) SetupSuite() {
	credentials := clientcredentials.Config{
		ClientID:     os.Getenv("COMMERCELAYER_CLIENT_ID"),
		ClientSecret: os.Getenv("COMMERCELAYER_CLIENT_SECRET"),
		TokenURL:     os.Getenv("COMMERCELAYER_AUTH_ENDPOINT"),
		Scopes:       []string{},
	}

	token, err := credentials.Token(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	tokenSource := oauth2.StaticTokenSource(token)

	testAccProviderCommercelayer = Provider(WithTokenSource(tokenSource))()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"commercelayer": func() (*schema.Provider, error) {
			return testAccProviderCommercelayer, nil
		},
	}
}

func TestAcceptanceSuite(t *testing.T) {
	if os.Getenv("TF_ACC") == "1" {
		suite.Run(t, new(AcceptanceSuite))
	}
}

func TestProvider(t *testing.T) {
	provider := Provider()()
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(s *AcceptanceSuite) {
	requiredEnvs := []string{
		"COMMERCELAYER_CLIENT_ID",
		"COMMERCELAYER_CLIENT_SECRET",
		"COMMERCELAYER_API_ENDPOINT",
		"COMMERCELAYER_AUTH_ENDPOINT",
	}
	for _, val := range requiredEnvs {
		if os.Getenv(val) == "" {
			s.Failf("%v must be set for acceptance tests", val)
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

func retryRemoval(times int, callable func() (*http.Response, error)) error {
	for retries := 1; retries < times; retries++ {
		resp, err := callable()
		if resp.StatusCode == 404 {
			return nil
		}
		if err != nil {
			return err
		}

		if resp.StatusCode == 200 {
			log.Println("retrying removal")
			time.Sleep(time.Second)
			continue
		}

		return fmt.Errorf("received response code with status %d", resp.StatusCode)
	}

	return nil
}
