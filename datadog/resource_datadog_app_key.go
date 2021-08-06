package datadog

import (
	"context"

	"github.com/terraform-providers/terraform-provider-datadog/datadog/internal/utils"

	datadogV2 "github.com/DataDog/datadog-api-client-go/api/v2/datadog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDatadogAppKey() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides a Datadog APP Key resource. This can be used to create and manage Datadog APP Keys.",
		CreateContext: resourceDatadogAppKeyCreate,
		ReadContext:   resourceDatadogAppKeyRead,
		UpdateContext: resourceDatadogAppKeyUpdate,
		DeleteContext: resourceDatadogAppKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "Name for APP Key.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"key": {
				Description: "The value of the APP Key.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func buildDatadogAppKeyCreateV2Struct(d *schema.ResourceData) *datadogV2.ApplicationKeyCreateRequest {
	appKeyAttributes := datadogV2.NewApplicationKeyCreateAttributes(d.Get("name").(string))
	appKeyData := datadogV2.NewApplicationKeyCreateData(*appKeyAttributes, datadogV2.APPLICATIONKEYSTYPE_APPLICATION_KEYS)
	appKeyRequest := datadogV2.NewApplicationKeyCreateRequest(*appKeyData)

	return appKeyRequest
}

func buildDatadogAppKeyUpdateV2Struct(d *schema.ResourceData) *datadogV2.ApplicationKeyUpdateRequest {
	appKeyAttributes := datadogV2.NewApplicationKeyUpdateAttributes(d.Get("name").(string))
	appKeyData := datadogV2.NewApplicationKeyUpdateData(*appKeyAttributes, d.Id(), datadogV2.APPLICATIONKEYSTYPE_APPLICATION_KEYS)
	appKeyRequest := datadogV2.NewApplicationKeyUpdateRequest(*appKeyData)

	return appKeyRequest
}

func updateAppKeyState(d *schema.ResourceData, appKeyData *datadogV2.FullApplicationKey) diag.Diagnostics {
	appKeyAttributes := appKeyData.GetAttributes()

	if err := d.Set("name", appKeyAttributes.GetName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("key", appKeyAttributes.GetKey()); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceDatadogAppKeyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConf := meta.(*ProviderConfiguration)
	datadogClientV2 := providerConf.DatadogClientV2
	authV2 := providerConf.AuthV2

	resp, httpResponse, err := datadogClientV2.KeyManagementApi.CreateCurrentUserApplicationKey(authV2, *buildDatadogAppKeyCreateV2Struct(d))
	if err != nil {
		return utils.TranslateClientErrorDiag(err, httpResponse, "error creating app key")
	}

	appKeyData := resp.GetData()
	d.SetId(appKeyData.GetId())

	return updateAppKeyState(d, &appKeyData)
}

func resourceDatadogAppKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConf := meta.(*ProviderConfiguration)
	datadogClientV2 := providerConf.DatadogClientV2
	authV2 := providerConf.AuthV2

	resp, httpResponse, err := datadogClientV2.KeyManagementApi.GetCurrentUserApplicationKey(authV2, d.Id())
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return utils.TranslateClientErrorDiag(err, httpResponse, "error getting app key")
	}
	appKeyData := resp.GetData()
	return updateAppKeyState(d, &appKeyData)
}

func resourceDatadogAppKeyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConf := meta.(*ProviderConfiguration)
	datadogClientV2 := providerConf.DatadogClientV2
	authV2 := providerConf.AuthV2

	resp, httpResponse, err := datadogClientV2.KeyManagementApi.UpdateCurrentUserApplicationKey(authV2, d.Id(), *buildDatadogAppKeyUpdateV2Struct(d))
	if err != nil {
		return utils.TranslateClientErrorDiag(err, httpResponse, "error updating app key")
	}
	appKeyData := resp.GetData()
	return updateAppKeyState(d, &appKeyData)
}

func resourceDatadogAppKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerConf := meta.(*ProviderConfiguration)
	datadogClientV2 := providerConf.DatadogClientV2
	authV2 := providerConf.AuthV2

	if httpResponse, err := datadogClientV2.KeyManagementApi.DeleteCurrentUserApplicationKey(authV2, d.Id()); err != nil {
		return utils.TranslateClientErrorDiag(err, httpResponse, "error deleting app key")
	}

	return nil
}
