package provider

import (
	"context"
	"runtime"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ShellProvider satisfies various provider interfaces.
var _ provider.Provider = &ShellProvider{}
var _ provider.ProviderWithFunctions = &ShellProvider{}

// New returns a new provider implementation.
func New(version, commit string) func() provider.Provider {
	return func() provider.Provider {
		return &ShellProvider{
			version: version,
			commit:  commit,
		}
	}
}

// ShellProviderData is the data available to the resource and data sources.
type ShellProviderData struct {
	provider    *ShellProvider
	Model       *ShellProviderModel
	Interpreter []string
	Environment map[string]string
	LogOutput   bool
}

// ShellProviderModel describes the provider data model.
type ShellProviderModel struct {
	Interpreter types.List `tfsdk:"interpreter"`
	Environment types.Map  `tfsdk:"environment"`
	LogOutput   types.Bool `tfsdk:"log_output"`
}

// ShellProvider defines the provider implementation.
type ShellProvider struct {
	version string
	commit  string
}

func (p *ShellProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "shell"
	resp.Version = p.version
}

func (p *ShellProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "The Shell provider allows you to execute arbitrary shell scripts and parse their JSON output for use in your Terraform configurations.",
		MarkdownDescription: "The _Shell_ provider allows you to execute arbitrary shell scripts and parse their JSON output for use in your _Terraform_ configurations. This is particularly useful for running scripts that interact with external APIs, or other systems that don't have a native _Terraform_ provider, or for performing complex data transformations.",
		Attributes: map[string]schema.Attribute{
			"interpreter": schema.ListAttribute{
				Description:         "The interpreter to use for executing scripts if not provided by the resource or data source.",
				MarkdownDescription: "The interpreter to use for executing scripts if not provided by the resource or data source. This defaults to `[\"/bin/bash\", \"-c\"]` or `[\"pwsh\", \"-c\"]` on Windows.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"environment": schema.MapAttribute{
				Description:         "The environment variables to set when executing scripts.",
				MarkdownDescription: "The environment variables to set when executing scripts.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"log_output": schema.BoolAttribute{
				Description:         "If true, lines output by the script will be logged at the appropriate level if they have a specific prefix.",
				MarkdownDescription: "If `true`, lines output by the script will be logged at the appropriate level if they start with the `[<LEVEL>]` pattern where `<LEVEL>` can be one of `ERROR`, `WARN`, `INFO`, `DEBUG` & `TRACE`.",
				Optional:            true,
			},
		},
	}
}

// Configure configures the provider.
func (p *ShellProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	if req.ClientCapabilities.DeferralAllowed && !req.Config.Raw.IsFullyKnown() {
		resp.Deferred = &provider.Deferred{
			Reason: provider.DeferredReasonProviderConfigUnknown,
		}
	}

	// Load the provider config
	model := &ShellProviderModel{}
	if resp.Diagnostics.Append(req.Config.Get(ctx, model)...); resp.Diagnostics.HasError() {
		return
	}

	// Set the interpreter
	var interpreter []string
	if !model.Interpreter.IsNull() {
		if resp.Diagnostics.Append(model.Interpreter.ElementsAs(ctx, &interpreter, false)...); resp.Diagnostics.HasError() {
			return
		}
	} else {
		if runtime.GOOS == "windows" {
			interpreter = []string{"pwsh", "-c"}
		} else {
			interpreter = []string{"/bin/bash", "-c"}
		}
	}

	// Set the environment
	environment := map[string]string{}
	if !model.Environment.IsNull() {
		if resp.Diagnostics.Append(model.Environment.ElementsAs(ctx, &environment, false)...); resp.Diagnostics.HasError() {
			return
		}
	}

	// Configure provider data
	providerData := &ShellProviderData{
		provider:    p,
		Model:       model,
		Interpreter: interpreter,
		Environment: environment,
		LogOutput:   model.LogOutput.ValueBool(),
	}

	resp.DataSourceData = providerData
	resp.ResourceData = providerData
}

func (p *ShellProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewScriptDataSource,
	}
}

func (p *ShellProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func (p *ShellProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewScriptResource,
	}
}
