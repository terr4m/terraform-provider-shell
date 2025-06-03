package provider

import (
	"context"
	"fmt"
	"runtime"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/terr4m/terraform-provider-shell/internal/tfdynamic"
)

var (
	_ resource.Resource              = &ScriptResource{}
	_ resource.ResourceWithConfigure = &ScriptResource{}
)

// NewScriptResource creates a new resource resource.
func NewScriptResource() resource.Resource {
	return &ScriptResource{}
}

// ScriptResource defines the resource implementation.
type ScriptResource struct {
	providerData *ShellProviderData
}

// ScriptResourceModel describes the resource data model.
type ScriptResourceModel struct {
	Environment      types.Map      `tfsdk:"environment"`
	WorkingDirectory types.String   `tfsdk:"working_directory"`
	OSCommands       types.Map      `tfsdk:"os_commands"`
	Output           types.Dynamic  `tfsdk:"output"`
	Timeouts         timeouts.Value `tfsdk:"timeouts"`
}

// CRUDCommandsModel describes a set of CRUD commands.
type CRUDCommandsModel struct {
	Create CommandModel `tfsdk:"create"`
	Read   CommandModel `tfsdk:"read"`
	Update CommandModel `tfsdk:"update"`
	Delete CommandModel `tfsdk:"delete"`
}

// CommandModel describes an interpreter and a command string.
type CommandModel struct {
	Interpreter types.List   `tfsdk:"interpreter"`
	Command     types.String `tfsdk:"command"`
}

// Metadata returns the resource metadata.
func (d *ScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_script", req.ProviderTypeName)
}

// Schema returns the resource schema.
func (r *ScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "The Shell script resource allows you to execute arbitrary commands as part of a Terraform lifecycle.",
		MarkdownDescription: "The _Shell_ script resource (`shell_script`) allows you to execute arbitrary commands as part of a _Terraform_ lifecycle. All commands must output a JSON string to the file defined by the `TF_SCRIPT_OUTPUT` environment variable and the file must be consistent on re-reading. You can access the output value in state in the read, update and delete commands via the `TF_SCRIPT_STATE_OUTPUT` environment variable.",
		Attributes: map[string]schema.Attribute{
			"environment": schema.MapAttribute{
				Description:         "The environment variables to set when executing commands; to be combined with the OS environment and the provider environment.",
				MarkdownDescription: "The environment variables to set when executing commands; to be combined with the OS environment and the provider environment.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"working_directory": schema.StringAttribute{
				Description:         "The working directory to use when executing the commands; this will default to the Terraform working directory.",
				MarkdownDescription: "The working directory to use when executing the commands; this will default to the _Terraform_ working directory.",
				Optional:            true,
			},
			"os_commands": schema.MapNestedAttribute{
				Description:         "A map of commands to run as part of the Terraform lifecycle where the map key is the GOOS value or default; default must be provided.",
				MarkdownDescription: "A map of commands to run as part of the Terraform lifecycle where the map key is the `GOOS` value or `default`; `default` must be provided.",
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"create": schema.SingleNestedAttribute{
							MarkdownDescription: "The create command configuration.",
							Required:            true,
							Attributes: map[string]schema.Attribute{
								"interpreter": schema.ListAttribute{
									MarkdownDescription: "The interpreter to use for executing the create command; if not set the platform default interpreter will be used.",
									ElementType:         types.StringType,
									Optional:            true,
									Validators: []validator.List{
										listvalidator.SizeAtLeast(1),
									},
								},
								"command": schema.StringAttribute{
									MarkdownDescription: "The create command to execute.",
									Required:            true,
								},
							},
						},
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
						"update": schema.SingleNestedAttribute{
							MarkdownDescription: "The update command configuration.",
							Required:            true,
							Attributes: map[string]schema.Attribute{
								"interpreter": schema.ListAttribute{
									MarkdownDescription: "The interpreter to use for executing the update command; if not set the platform default interpreter will be used.",
									ElementType:         types.StringType,
									Optional:            true,
									Validators: []validator.List{
										listvalidator.SizeAtLeast(1),
									},
								},
								"command": schema.StringAttribute{
									MarkdownDescription: "The update command to execute.",
									Required:            true,
								},
							},
						},
						"delete": schema.SingleNestedAttribute{
							MarkdownDescription: "The delete command configuration.",
							Required:            true,
							Attributes: map[string]schema.Attribute{
								"interpreter": schema.ListAttribute{
									MarkdownDescription: "The interpreter to use for executing the delete command; if not set the platform default interpreter will be used.",
									ElementType:         types.StringType,
									Optional:            true,
									Validators: []validator.List{
										listvalidator.SizeAtLeast(1),
									},
								},
								"command": schema.StringAttribute{
									MarkdownDescription: "The delete command to execute.",
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
				Create:            true,
				CreateDescription: "Timeout for creating the resource; this defaults to the provider value if not set. This should be a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).",
				Read:              true,
				ReadDescription:   "Timeout for reading the resource; this defaults to the provider value if not set. This should be a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).",
				Update:            true,
				UpdateDescription: "Timeout for updating the resource; this defaults to the provider value if not set. This should be a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).",
				Delete:            true,
				DeleteDescription: "Timeout for deleting the resource; this defaults to the provider value if not set. This should be a string that can be [parsed as a duration](https://pkg.go.dev/time#ParseDuration) consisting of numbers and unit suffixes, such as `30s` or `2h45m`. Valid time units are `s` (seconds), `m` (minutes), `h` (hours).",
			}),
		},
	}
}

// Configure configures the resource.
func (r *ScriptResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*ShellProviderData)
	if !ok {
		resp.Diagnostics.AddError("Unexpected resource provider data.", fmt.Sprintf("expected *ShellProviderData, got: %T", req.ProviderData))
		return
	}

	r.providerData = providerData
}

// Create creates the resource.
func (r *ScriptResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ScriptResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	var commands map[string]CRUDCommandsModel
	if resp.Diagnostics.Append(data.OSCommands.ElementsAs(ctx, &commands, false)...); resp.Diagnostics.HasError() {
		return
	}

	var command CRUDCommandsModel
	command, ok := commands[runtime.GOOS]
	if !ok {
		command = commands["default"]
	}

	timeout, diags := data.Timeouts.Create(ctx, r.providerData.DefaultTimeouts.Create)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	raw, diags := runCommand(ctx, r.providerData, command.Create.Interpreter, data.Environment, data.WorkingDirectory, command.Create.Command, TFLifecycleCreate, nil, true)
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

// Read reads the resource state.
func (r *ScriptResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ScriptResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	var commands map[string]CRUDCommandsModel
	if resp.Diagnostics.Append(data.OSCommands.ElementsAs(ctx, &commands, false)...); resp.Diagnostics.HasError() {
		return
	}

	var command CRUDCommandsModel
	command, ok := commands[runtime.GOOS]
	if !ok {
		command = commands["default"]
	}

	timeout, diags := data.Timeouts.Read(ctx, r.providerData.DefaultTimeouts.Read)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var stateData ScriptResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...); resp.Diagnostics.HasError() {
		return
	}

	stateOutput, err := tfdynamic.EncodeDynamic(ctx, stateData.Output)
	if err != nil {
		resp.Diagnostics.AddError("Failed to encode the state output.", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	raw, diags := runCommand(ctx, r.providerData, command.Read.Interpreter, data.Environment, data.WorkingDirectory, command.Read.Command, TFLifecycleRead, stateOutput, true)
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

// Update updates the resource.
func (r *ScriptResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ScriptResourceModel
	if resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	var commands map[string]CRUDCommandsModel
	if resp.Diagnostics.Append(data.OSCommands.ElementsAs(ctx, &commands, false)...); resp.Diagnostics.HasError() {
		return
	}

	var command CRUDCommandsModel
	command, ok := commands[runtime.GOOS]
	if !ok {
		command = commands["default"]
	}

	timeout, diags := data.Timeouts.Update(ctx, r.providerData.DefaultTimeouts.Update)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var stateData ScriptResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...); resp.Diagnostics.HasError() {
		return
	}

	stateOutput, err := tfdynamic.EncodeDynamic(ctx, stateData.Output)
	if err != nil {
		resp.Diagnostics.AddError("Failed to encode the state output.", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	raw, diags := runCommand(ctx, r.providerData, command.Update.Interpreter, data.Environment, data.WorkingDirectory, command.Update.Command, TFLifecycleUpdate, stateOutput, true)
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

// Delete deletes the resource.
func (r *ScriptResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ScriptResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &data)...); resp.Diagnostics.HasError() {
		return
	}

	var commands map[string]CRUDCommandsModel
	if resp.Diagnostics.Append(data.OSCommands.ElementsAs(ctx, &commands, false)...); resp.Diagnostics.HasError() {
		return
	}

	var command CRUDCommandsModel
	command, ok := commands[runtime.GOOS]
	if !ok {
		command = commands["default"]
	}

	timeout, diags := data.Timeouts.Delete(ctx, r.providerData.DefaultTimeouts.Delete)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	var stateData ScriptResourceModel
	if resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...); resp.Diagnostics.HasError() {
		return
	}

	stateOutput, err := tfdynamic.EncodeDynamic(ctx, stateData.Output)
	if err != nil {
		resp.Diagnostics.AddError("Failed to encode the state output.", err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	_, diags = runCommand(ctx, r.providerData, command.Delete.Interpreter, data.Environment, data.WorkingDirectory, command.Delete.Command, TFLifecycleDelete, stateOutput, false)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
}
