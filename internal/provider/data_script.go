package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terr4m/terraform-provider-shell/internal/tfdynamic"
)

var _ datasource.DataSource = &ScriptDataSource{}
var _ datasource.DataSourceWithConfigure = &ScriptDataSource{}

// NewScriptDataSource creates a new consistent hash data source.
func NewScriptDataSource() datasource.DataSource {
	return &ScriptDataSource{}
}

// ScriptDataSource defines the data source implementation.
type ScriptDataSource struct {
	providerData *ShellProviderData
}

// ScriptDataSourceModel describes the data source data model.
type ScriptDataSourceModel struct {
	Interpreter      types.List    `tfsdk:"interpreter"`
	Environment      types.Map     `tfsdk:"environment"`
	WorkingDirectory types.String  `tfsdk:"working_directory"`
	Command          types.String  `tfsdk:"command"`
	Output           types.Dynamic `tfsdk:"output"`
}

// Metadata returns the data source metadata.
func (d *ScriptDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_script", req.ProviderTypeName)
}

// Schema returns the data source schema.
func (d *ScriptDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "The Shell script data source allows you to execute an arbitrary command as the read part of a Terraform lifecycle.",
		MarkdownDescription: "The _Shell_ script data source allows you to execute an arbitrary command as the read part of a _Terraform_ lifecycle. The script must output a JSON string to the file defined by the `TF_SCRIPT_OUTPUT` environment variable and the file must be consistent on re-reading.",
		Attributes: map[string]schema.Attribute{
			"interpreter": schema.ListAttribute{
				Description:         "The interpreter to use for executing the command; if not set the provider interpreter will be used.",
				MarkdownDescription: "The interpreter to use for executing the command; if not set the provider interpreter will be used.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"environment": schema.MapAttribute{
				Description:         "The environment variables to set when executing command; to be combined with the OS environment and the provider environment.",
				MarkdownDescription: "The environment variables to set when executing command; to be combined with the OS environment and the provider environment.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"working_directory": schema.StringAttribute{
				Description:         "The working directory to use when executing the command; this will default to the Terraform working directory.",
				MarkdownDescription: "The working directory to use when executing the command; this will default to the _Terraform_ working directory..",
				Optional:            true,
			},
			"command": schema.StringAttribute{
				Description:         "The command to run when reading the data source.",
				MarkdownDescription: "The command to run when reading the data source; this must write a JSON string to the file defined by the `TF_SCRIPT_OUTPUT` environment variable.",
				Required:            true,
			},
			"output": schema.DynamicAttribute{
				Description:         "The output of the script as a structured type.",
				MarkdownDescription: "The output of the script as a structured type.",
				Computed:            true,
			},
		},
	}
}

// Configure configures the data source.
func (d *ScriptDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*ShellProviderData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected data source provider data.", fmt.Sprintf("expected *ShellProviderData, got: %T", req.ProviderData))
		return
	}

	d.providerData = providerData
}

// Read reads the data source.
func (d *ScriptDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ScriptDataSourceModel
	if resp.Diagnostics.Append(req.Config.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	raw, diags := runCommand(ctx, d.providerData, data.Interpreter, data.Environment, data.WorkingDirectory, data.Command, true)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	out, diags := tfdynamic.Decode(ctx, raw)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Output = out

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
