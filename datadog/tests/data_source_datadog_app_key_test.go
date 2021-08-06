package test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatadogAppKeyDatasource_matchId(t *testing.T) {
	_, accProviders := testAccProviders(context.Background(), t)
	ctx, accProviders := testAccProviders(context.Background(), t)
	appKeyName := uniqueEntityName(ctx, t)
	accProvider := testAccProvider(t, accProviders)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: accProviders,
		CheckDestroy:      testAccCheckDatadogAppKeyDestroy(accProvider),
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAppKeyIdConfig(appKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatadogAppKeyExists(accProvider, "datadog_app_key.app_key_1"),
					resource.TestCheckResourceAttr("datadog_app_key.app_key_1", "name", fmt.Sprintf("%s 1", appKeyName)),
					resource.TestCheckResourceAttr("datadog_app_key.app_key_2", "name", fmt.Sprintf("%s 2", appKeyName)),
					resource.TestCheckResourceAttr("data.datadog_app_key.app_key", "name", fmt.Sprintf("%s 1", appKeyName)),
				),
			},
		},
	})
}

func TestAccDatadogAppKeyDatasource_matchName(t *testing.T) {
	_, accProviders := testAccProviders(context.Background(), t)
	ctx, accProviders := testAccProviders(context.Background(), t)
	appKeyName := uniqueEntityName(ctx, t)
	accProvider := testAccProvider(t, accProviders)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: accProviders,
		CheckDestroy:      testAccCheckDatadogAppKeyDestroy(accProvider),
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceAppKeyNameConfig(appKeyName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatadogAppKeyExists(accProvider, "datadog_app_key.app_key_1"),
					resource.TestCheckResourceAttr("datadog_app_key.app_key_1", "name", fmt.Sprintf("%s 1", appKeyName)),
					resource.TestCheckResourceAttr("datadog_app_key.app_key_2", "name", fmt.Sprintf("%s 2", appKeyName)),
					resource.TestCheckResourceAttr("data.datadog_app_key.app_key", "name", fmt.Sprintf("%s 1", appKeyName)),
				),
			},
		},
	})
}

func TestAccDatadogAppKeyDatasource_matchIdError(t *testing.T) {
	_, accProviders := testAccProviders(context.Background(), t)
	accProvider := testAccProvider(t, accProviders)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: accProviders,
		CheckDestroy:      testAccCheckDatadogAppKeyDestroy(accProvider),
		Steps: []resource.TestStep{
			{
				Config:      testAccDatasourceAppKeyIdOnlyConfig("11111111-2222-3333-4444-555555555555"),
				ExpectError: regexp.MustCompile("error getting app key"),
			},
		},
	})
}

func TestAccDatadogAppKeyDatasource_matchNameError(t *testing.T) {
	_, accProviders := testAccProviders(context.Background(), t)
	ctx, accProviders := testAccProviders(context.Background(), t)
	appKeyName := uniqueEntityName(ctx, t)
	accProvider := testAccProvider(t, accProviders)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: accProviders,
		CheckDestroy:      testAccCheckDatadogAppKeyDestroy(accProvider),
		Steps: []resource.TestStep{
			{
				Config:      testAccDatasourceAppKeyNameOnlyConfig(appKeyName),
				ExpectError: regexp.MustCompile("your query returned no result, please try a less specific search criteria"),
			},
		},
	})
}

func TestAccDatadogAppKeyDatasource_missingParametersError(t *testing.T) {
	_, accProviders := testAccProviders(context.Background(), t)
	accProvider := testAccProvider(t, accProviders)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: accProviders,
		CheckDestroy:      testAccCheckDatadogAppKeyDestroy(accProvider),
		Steps: []resource.TestStep{
			{
				Config:      testAccDatasourceAppKeyMissingParametersConfig(),
				ExpectError: regexp.MustCompile("missing id or name parameter"),
			},
		},
	})
}

func testAccAppKeyConfig(uniq string) string {
	return fmt.Sprintf(`
resource "datadog_app_key" "app_key_1" {
  name = "%s 1"
}

resource "datadog_app_key" "app_key_2" {
  name = "%s 2"
}`, uniq, uniq)
}

func testAccDatasourceAppKeyIdConfig(uniq string) string {
	return fmt.Sprintf(`
%s
data "datadog_app_key" "app_key" {
  depends_on = [
    datadog_app_key.app_key_1,
    datadog_app_key.app_key_2,
  ]
  id = datadog_app_key.app_key_1.id
}`, testAccAppKeyConfig(uniq))
}

func testAccDatasourceAppKeyIdOnlyConfig(uniq string) string {
	return fmt.Sprintf(`
data "datadog_app_key" "app_key" {
  id = "%s"
}`, uniq)
}

func testAccDatasourceAppKeyNameConfig(uniq string) string {
	return fmt.Sprintf(`
%s
data "datadog_app_key" "app_key" {
  depends_on = [
    datadog_app_key.app_key_1,
    datadog_app_key.app_key_2,
  ]
  name = datadog_app_key.app_key_1.name
}`, testAccAppKeyConfig(uniq))
}

func testAccDatasourceAppKeyNameOnlyConfig(uniq string) string {
	return fmt.Sprintf(`
data "datadog_app_key" "app_key" {
  name = "%s"
}`, uniq)
}

func testAccDatasourceAppKeyMissingParametersConfig() string {
	return fmt.Sprintf(`
data "datadog_app_key" "app_key" {
}`)
}
