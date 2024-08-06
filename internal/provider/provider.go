package provider

import (
	"context"
	"os"

	"github.com/zeiss/terraform-provider-openfga/sdk"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure OpenFGA satisfies various provider interfaces.
var (
	_ provider.Provider              = &openfgaProvider{}
	_ provider.ProviderWithFunctions = &openfgaProvider{}
)

// openfgaProvider defines the provider implementation.
type openfgaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// OpenFGAProvider defines the provider implementation.
type OpenFGAProvider struct {
	ApiURL types.String `tfsdk:"api_url"`
}

func (p *openfgaProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "openfga"
	resp.Version = p.version
}

func (p *openfgaProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				MarkdownDescription: "OpenFGA API URL",
				Optional:            true,
			},
		},
	}
}

func (p *openfgaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data OpenFGAProvider

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.ApiURL.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Unknown OpenFGA API URL",
			"The provider cannot create the OpenFGA API client as there is an unknown configuration value for the OpenFGA API URL. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the OPENFGA_API_URL environment variable.",
		)
	}

	api_url := os.Getenv("OPENFGA_API_URL")

	if !data.ApiURL.IsNull() {
		api_url = data.ApiURL.ValueString()
	}

	client, err := sdk.NewClient(api_url)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create OpenFGA client", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *openfgaProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewStoreResource,
		NewModelResource,
		NewTupleResource,
	}
}

func (p *openfgaProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *openfgaProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &openfgaProvider{
			version: version,
		}
	}
}
