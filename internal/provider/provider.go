package provider

import (
	"context"
	"os"
	"terraform-provider-violet/internal/violet"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &violetProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &violetProvider{
			version: version,
		}
	}
}

// violetProvider is the provider implementation.
type violetProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *violetProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "violet"
	resp.Version = p.version
}

type violetProviderModel struct {
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
	AppId     types.String `tfsdk:"app_id"`
	AppSecret types.String `tfsdk:"app_secret"`
}

// Schema defines the provider-level schema for configuration data.
func (p *violetProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"app_id": schema.StringAttribute{
				Optional: true,
			},
			"app_secret": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a violet API client for data sources and resources.
func (p *violetProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Violet client")
	var config violetProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Violet username",
			"The provider cannot create the Violet API client as there is an unknown configuration value for the Violet username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VIOLET_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Violet password",
			"The provider cannot create the Violet API client as there is an unknown configuration value for the Violet password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VIOLET_PASSWORD environment variable.",
		)
	}

	if config.AppId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("app_id"),
			"Unknown Violet app_id",
			"The provider cannot create the Violet API client as there is an unknown configuration value for the Violet app_id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VIOLET_APP_ID environment variable.",
		)
	}

	if config.AppSecret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("app_secret"),
			"Unknown Violet app_secret",
			"The provider cannot create the Violet API client as there is an unknown configuration value for the Violet app_secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VIOLET_APP_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	username := os.Getenv("VIOLET_USERNAME")
	password := os.Getenv("VIOLET_PASSWORD")
	appId := os.Getenv("VIOLET_APP_ID")
	appSecret := os.Getenv("VIOLET_APP_SECRET")

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.AppId.IsNull() {
		appId = config.AppId.ValueString()
	}

	if !config.AppSecret.IsNull() {
		appSecret = config.AppSecret.ValueString()
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Violet username",
			"The provider cannot create the Violet API client as there is an unknown configuration value for the Violet username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VIOLET_USERNAME environment variable.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Violet password",
			"The provider cannot create the Violet API client as there is an unknown configuration value for the Violet password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VIOLET_PASSWORD environment variable.",
		)
	}

	if appId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("app_id"),
			"Unknown Violet app_id",
			"The provider cannot create the Violet API client as there is an unknown configuration value for the Violet app_id. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VIOLET_APP_ID environment variable.",
		)
	}

	if appSecret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("app_secret"),
			"Unknown Violet app_secret",
			"The provider cannot create the Violet API client as there is an unknown configuration value for the Violet app_secret. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the VIOLET_APP_SECRET environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := violet.VioletClient{
		Username:  username,
		Password:  password,
		AppId:     appId,
		AppSecret: appSecret,
		BaseUrl:   "https://sandbox-api.violet.io/v1/",
	}
	client.Login(ctx)

	resp.DataSourceData = &client
	resp.ResourceData = &client
}

// DataSources defines the data sources implemented in the provider.
func (p *violetProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		WebhookDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *violetProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWebhookResource,
	}
}
