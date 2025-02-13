package provider

import (
	"context"
	"fmt"

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
	Interpreter      types.List    `tfsdk:"interpreter"`
	Environment      types.Map     `tfsdk:"environment"`
	WorkingDirectory types.String  `tfsdk:"working_directory"`
	Commands         *ScriptsModel `tfsdk:"commands"`
	Output           types.Dynamic `tfsdk:"output"`
}

// ScriptsModel describes the scripts data model.
type ScriptsModel struct {
	Create types.String `tfsdk:"create"`
	Read   types.String `tfsdk:"read"`
	Update types.String `tfsdk:"update"`
	Delete types.String `tfsdk:"delete"`
}

// Metadata returns the resource metadata.
func (d *ScriptResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_script", req.ProviderTypeName)
}

// Schema returns the resource schema.
func (r *ScriptResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "The Shell script resource allows you to execute arbitrary commands as part of a Terraform lifecycle.",
		MarkdownDescription: "The _Shell_ script resource allows you to execute arbitrary commands as part of a _Terraform_ lifecycle. All commands must output a JSON string to the file defined by the `TF_SCRIPT_OUTPUT` environment variable and the file must be consistent on re-reading.",
		Attributes: map[string]schema.Attribute{
			"interpreter": schema.ListAttribute{
				Description:         "The interpreter to use for executing the commands; if not set the provider interpreter will be used.",
				MarkdownDescription: "The interpreter to use for executing the commands; if not set the provider interpreter will be used.",
				ElementType:         types.StringType,
				Optional:            true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"environment": schema.MapAttribute{
				Description:         "The environment variables to set when executing commands; to be combined with the OS environment and the provider environment.",
				MarkdownDescription: "The environment variables to set when executing commands; to be combined with the OS environment and the provider environment.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"working_directory": schema.StringAttribute{
				Description:         "The working directory to use when executing the commands; this will default to the Terraform working directory.",
				MarkdownDescription: "The working directory to use when executing the commands; this will default to the _Terraform_ working directory..",
				Optional:            true,
			},
			"commands": schema.SingleNestedAttribute{
				Description:         "The commands to run as part of the Terraform lifecycle.",
				MarkdownDescription: "The commands to run as part of the _Terraform_ lifecycle. All commands must write a JSON string to the file defined by the `TF_SCRIPT_OUTPUT` environment variable.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"create": schema.StringAttribute{
						MarkdownDescription: "The command to execute when creating the resource.",
						Required:            true,
					},
					"read": schema.StringAttribute{
						MarkdownDescription: "The command to execute when reading the resource.",
						Required:            true,
					},
					"update": schema.StringAttribute{
						MarkdownDescription: "The command to execute when updating the resource.",
						Required:            true,
					},
					"delete": schema.StringAttribute{
						MarkdownDescription: "The command to execute when deleting the resource.",
						Required:            true,
					},
				},
			},
			"output": schema.DynamicAttribute{
				Description:         "The output of the script as a structured type.",
				MarkdownDescription: "The output of the script as a structured type.",
				Computed:            true,
			},
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

	raw, diags := runCommand(ctx, r.providerData, data.Interpreter, data.Environment, data.WorkingDirectory, data.Commands.Create, true)
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

	raw, diags := runCommand(ctx, r.providerData, data.Interpreter, data.Environment, data.WorkingDirectory, data.Commands.Read, true)
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

	raw, diags := runCommand(ctx, r.providerData, data.Interpreter, data.Environment, data.WorkingDirectory, data.Commands.Update, true)
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

	_, diags := runCommand(ctx, r.providerData, data.Interpreter, data.Environment, data.WorkingDirectory, data.Commands.Delete, false)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}
}
