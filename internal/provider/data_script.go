package provider

import (
	"context"
	"fmt"
	"runtime"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
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
	Environment      types.Map      `tfsdk:"environment"`
	WorkingDirectory types.String   `tfsdk:"working_directory"`
	Inputs           types.Dynamic  `tfsdk:"inputs"`
	OSCommands       types.Map      `tfsdk:"os_commands"`
	Output           types.Dynamic  `tfsdk:"output"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

// ReadCommandModel describes a set of CRUD commands.
type ReadCommandModel struct {
	Read CommandModel `tfsdk:"read"`
}

// Metadata returns the data source metadata.
func (d *ScriptDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_script", req.ProviderTypeName)
}

// Schema returns the data source schema.
func (d *ScriptDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "The Shell script data source allows you to execute an arbitrary command as the read part of a Terraform lifecycle.",
		MarkdownDescription: "The _Shell_ script data source (`shell_script`) allows you to execute an arbitrary command as the read part of a _Terraform_ lifecycle. The script must output a JSON string to the file defined by the `TF_SCRIPT_OUTPUT` environment variable and the file must be consistent on re-reading. If the script exits with a non-zero code the provider will ready any text from the file defined by the `TF_SCRIPT_ERROR` environment variable and return it as part of the error diagnostics.",
		Attributes: map[string]schema.Attribute{
			"environment": schema.MapAttribute{
				Description:         "The environment variables to set when executing command; to be combined with the OS environment and the provider environment.",
				MarkdownDescription: "The environment variables to set when executing command; to be combined with the OS environment and the provider environment.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"working_directory": schema.StringAttribute{
				Description:         "The working directory to use when executing the command; this will default to the Terraform working directory.",
				MarkdownDescription: "The working directory to use when executing the command; this will default to the _Terraform_ working directory.",
				Optional:            true,
			},
			"inputs": schema.DynamicAttribute{
				Description:         "Inputs to be made available to the script.",
				MarkdownDescription: "Inputs to be made available to the script; these can be accessed as JSON via the `TF_SCRIPT_INPUTS` environment variable.",
				Optional:            true,
			},
			"os_commands": schema.MapNestedAttribute{
				Description:         "A map of commands to run as part of the Terraform lifecycle where the map key is the GOOS value or default; default must be provided.",
				MarkdownDescription: "A map of commands to run as part of the Terraform lifecycle where the map key is the `GOOS` value or `default`; `default` must be provided.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"read": schema.SingleNestedAttribute{
							MarkdownDescription: "The read command configuration.",
							Required:            true,
							Attributes: map[string]schema.Attribute{
								"interpreter": schema.ListAttribute{
									MarkdownDescription: "The interpreter to use for executing the read command; if not set the platform default interpreter will be used.",
									ElementType:         types.StringType,
									Optional:            true,
									Validators: []validator.List{
										listvalidator.SizeAtLeast(1),
									},
								},
								"command": schema.StringAttribute{
									MarkdownDescription: "The read command to execute.",
									Required:            true,
								},
							},
						},
					},
				},
			},
			"output": schema.DynamicAttribute{
				Description:         "The output of the script as a structured type.",
				MarkdownDescription: "The output of the script as a structured type.",
				Computed:            true,
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read:            true,
				ReadDescription: "Timeout for reading the data source; this defaults to the provider value if not set. This should be a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).",
			}),
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

	var commands map[string]ReadCommandModel
	if resp.Diagnostics.Append(data.OSCommands.ElementsAs(ctx, &commands, false)...); resp.Diagnostics.HasError() {
		return
	}

	var command ReadCommandModel
	command, ok := commands[runtime.GOOS]
	if !ok {
		command = commands["default"]
	}

	timeout, diags := data.Timeouts.Read(ctx, d.providerData.DefaultTimeouts.Read)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	inputs, err := tfdynamic.EncodeDynamic(ctx, data.Inputs)
	if err != nil {
		resp.Diagnostics.AddError("Failed to encode the inputs.", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	raw, diags := runCommand(ctx, d.providerData, command.Read.Interpreter, data.Environment, data.WorkingDirectory, command.Read.Command, TFLifecycleRead, inputs, nil, true)
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
