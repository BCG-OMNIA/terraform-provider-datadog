package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/terraform-providers/terraform-provider-datadog/datadog"

	datadogV2 "github.com/DataDog/datadog-api-client-go/api/v2/datadog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDatadogAppKey_Update(t *testing.T) {
	t.Parallel()
	ctx, accProviders := testAccProviders(context.Background(), t)
	accProvider := testAccProvider(t, accProviders)
	appKeyName := uniqueEntityName(ctx, t)
	appKeyNameUpdate := appKeyName + "-2"
	resourceName := "datadog_app_key.foo"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: accProviders,
		CheckDestroy:      testAccCheckDatadogAppKeyDestroy(accProvider),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDatadogAppKeyConfigRequired(appKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatadogAppKeyExists(accProvider, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", appKeyName),
					testAccCheckDatadogAppKeyValueMatches(accProvider, resourceName),
				),
			},
			{
				Config: testAccCheckDatadogAppKeyConfigRequired(appKeyNameUpdate),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatadogAppKeyExists(accProvider, resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", appKeyNameUpdate),
					testAccCheckDatadogAppKeyValueMatches(accProvider, resourceName),
				),
			},
		},
	})
}

func testAccCheckDatadogAppKeyConfigRequired(uniq string) string {
	return fmt.Sprintf(`
resource "datadog_app_key" "foo" {
  name = "%s"
}`, uniq)
}

func testAccCheckDatadogAppKeyExists(accProvider func() (*schema.Provider, error), n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		provider, _ := accProvider()
		providerConf := provider.Meta().(*datadog.ProviderConfiguration)
		datadogClientV2 := providerConf.DatadogClientV2
		authV2 := providerConf.AuthV2

		if err := datadogAppKeyExistsHelper(authV2, s, datadogClientV2, n); err != nil {
			return err
		}
		return nil
	}
}

func datadogAppKeyExistsHelper(ctx context.Context, s *terraform.State, client *datadogV2.APIClient, name string) error {
	id := s.RootModule().Resources[name].Primary.ID
	if _, _, err := client.KeyManagementApi.GetCurrentUserApplicationKey(ctx, id); err != nil {
		return fmt.Errorf("received an error retrieving app key %s", err)
	}
	return nil
}

func testAccCheckDatadogAppKeyValueMatches(accProvider func() (*schema.Provider, error), n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		provider, _ := accProvider()
		providerConf := provider.Meta().(*datadog.ProviderConfiguration)
		datadogClientV2 := providerConf.DatadogClientV2
		authV2 := providerConf.AuthV2

		if err := datadogAppKeyValueMatches(authV2, s, datadogClientV2, n); err != nil {
			return err
		}
		return nil
	}
}

func datadogAppKeyValueMatches(ctx context.Context, s *terraform.State, client *datadogV2.APIClient, name string) error {
	primaryResource := s.RootModule().Resources[name].Primary
	id := primaryResource.ID
	expectedKey := primaryResource.Attributes["key"]
	resp, _, err := client.KeyManagementApi.GetCurrentUserApplicationKey(ctx, id)
	if err != nil {
		return fmt.Errorf("received an error retrieving app key %s", err)
	}
	actualKey := resp.Data.Attributes.GetKey()
	if expectedKey != actualKey {
		return fmt.Errorf("app key value does not match")
	}
	return nil
}

func testAccCheckDatadogAppKeyDestroy(accProvider func() (*schema.Provider, error)) func(*terraform.State) error {
	return func(s *terraform.State) error {
		provider, _ := accProvider()
		providerConf := provider.Meta().(*datadog.ProviderConfiguration)
		datadogClientV2 := providerConf.DatadogClientV2
		authV2 := providerConf.AuthV2

		if err := datadogAppKeyDestroyHelper(authV2, s, datadogClientV2); err != nil {
			return err
		}
		return nil
	}
}

func datadogAppKeyDestroyHelper(ctx context.Context, s *terraform.State, client *datadogV2.APIClient) error {
	for _, r := range s.RootModule().Resources {
		if r.Type != "datadog_app_key" {
			continue
		}

		id := r.Primary.ID
		_, httpResponse, err := client.KeyManagementApi.GetCurrentUserApplicationKey(ctx, id)

		if err != nil {
			if httpResponse.StatusCode == 404 {
				continue
			}
			return fmt.Errorf("received an error retrieving app key %s", err)
		}

		return fmt.Errorf("app key still exists")
	}

	return nil
}
