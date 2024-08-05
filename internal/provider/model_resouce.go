package provider

import (
	"context"
	"fmt"

	"github.com/zeiss/terraform-provider-openfga/sdk"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &modelResource{}
	_ resource.ResourceWithImportState = &modelResource{}
)

func NewModelResource() resource.Resource {
	return &modelResource{}
}

type modelResource struct {
	client *sdk.Client
}

// ModelResouce describes the resource data model.
type ModelResouce struct {
	ID    types.String        `tfsdk:"id"`
	Spec  types.String        `tfsdk:"spec"`
	Store *StoreResourceEmbed `tfsdk:"store"`
}

// StoreResourceEmbed describes the nested store resource data model.
type StoreResourceEmbed struct {
	ID types.String `tfsdk:"id"`
}

func (r *modelResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_model"
}

func (r *modelResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Model resource",
		Attributes: map[string]schema.Attribute{
			"store": schema.SingleNestedAttribute{
				MarkdownDescription: "Model's store",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Store identifier",
						Required:            true,
					},
				},
			},
			"spec": schema.StringAttribute{
				MarkdownDescription: "Model specification",
				Required:            true,
			},
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Model identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *modelResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*sdk.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *modelResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ModelResouce

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createdModel, err := r.client.CreateModel(ctx, data.Store.ID.ValueString(), data.Spec.ValueString())

	tflog.Info(ctx, fmt.Sprintf("Read a model with id %s", createdModel.ID))

	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create model, got error: %s", err.Error()),
		)
		return
	}

	// Set the created models's ID in the Terraform state
	data.ID = types.StringValue(createdModel.ID)

	tflog.Info(ctx, fmt.Sprintf("Model with id %s got created", createdModel.ID))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *modelResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ModelResouce

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read a model with id %s", data.ID.ValueString()))

	// Retrieve the model using the GetModel method
	model, err := r.client.GetAuthorizationModel(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read model, got error: %s", err.Error()),
		)
		return
	}

	// Update the data model with the retrieved model information
	data.ID = types.StringValue(model.ID)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *modelResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ModelResouce

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Update a model with id %s", data.ID))
}

func (r *modelResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ModelResouce

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
}

func (r *modelResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
