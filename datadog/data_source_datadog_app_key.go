package datadog

import (
	"context"

	"github.com/terraform-providers/terraform-provider-datadog/datadog/internal/utils"

	datadogV2 "github.com/DataDog/datadog-api-client-go/api/v2/datadog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatadogAppKey() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to retrieve information about an existing app key.",
		ReadContext: dataSourceDatadogAppKeyRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Id for APP Key.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Description: "Name for APP Key.",
				Type:        schema.TypeString,
				Optional:    true,
			},

			// Computed values
			"key": {
				Description: "The value of the APP Key.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func dataSourceDatadogAppKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConf := meta.(*ProviderConfiguration)
	datadogClientV2 := providerConf.DatadogClientV2
	authV2 := providerConf.AuthV2

	if id := d.Get("id").(string); id != "" {
		resp, httpResponse, err := datadogClientV2.KeyManagementApi.GetCurrentUserApplicationKey(authV2, id)
		if err != nil {
			return utils.TranslateClientErrorDiag(err, httpResponse, "error getting app key")
		}
		appKeyData := resp.GetData()
		d.SetId(appKeyData.GetId())
		return updateAppKeyState(d, &appKeyData)
	} else if name := d.Get("name").(string); name != "" {
		optionalParams := datadogV2.NewListCurrentUserApplicationKeysOptionalParameters()
		optionalParams.WithFilter(name)

		appKeysResponse, httpResponse, err := datadogClientV2.KeyManagementApi.ListCurrentUserApplicationKeys(authV2, *optionalParams)
		if err != nil {
			return utils.TranslateClientErrorDiag(err, httpResponse, "error getting app keys")
		}

		appKeysData := appKeysResponse.GetData()

		if len(appKeysData) > 1 {
			return diag.Errorf("your query returned more than one result, please try a more specific search criteria")
		}
		if len(appKeysData) == 0 {
			return diag.Errorf("your query returned no result, please try a less specific search criteria")
		}

		appKeyPartialData := appKeysData[0]

		id := appKeyPartialData.GetId()
		appKeyResponse, httpResponse, err := datadogClientV2.KeyManagementApi.GetCurrentUserApplicationKey(authV2, id)
		if err != nil {
			return utils.TranslateClientErrorDiag(err, httpResponse, "error getting app key")
		}
		appKeyFullData := appKeyResponse.GetData()
		d.SetId(appKeyFullData.GetId())
		return updateAppKeyState(d, &appKeyFullData)
	}

	return diag.Errorf("missing id or name parameter")
}
