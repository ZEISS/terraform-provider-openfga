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
	_ resource.Resource                = &tupleresource{}
	_ resource.ResourceWithImportState = &tupleresource{}
)

func NewTupleResource() resource.Resource {
	return &tupleresource{}
}

type tupleresource struct {
	client *sdk.Client
}

// TupleResource describes the resource data model.
type TupleResource struct {
	ID       types.String        `tfsdk:"id"`
	User     types.String        `tfsdk:"user"`
	Relation types.String        `tfsdk:"relation"`
	Document types.String        `tfsdk:"document"`
	Store    *TupleResourceEmbed `tfsdk:"store"`
}

// TupleResourceEmbed describes the nested store resource data model.
type TupleResourceEmbed struct {
	ID    types.String `tfsdk:"id"`
	Model types.String `tfsdk:"model"`
}

func (r *tupleresource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tuple"
}

func (r *tupleresource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Tuple resource",
		Attributes: map[string]schema.Attribute{
			"store": schema.SingleNestedAttribute{
				MarkdownDescription: "Tuple's model",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "Store identifier",
						Required:            true,
					},
					"model": schema.StringAttribute{
						MarkdownDescription: "Model identifier",
						Required:            true,
					},
				},
			},
			"user": schema.StringAttribute{
				MarkdownDescription: "User identifier",
				Required:            true,
			},
			"relation": schema.StringAttribute{
				MarkdownDescription: "Relation identifier",
				Required:            true,
			},
			"document": schema.StringAttribute{
				MarkdownDescription: "Document identifier",
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

func (r *tupleresource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *tupleresource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TupleResource

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createdTuple, err := r.client.CreateTuple(ctx, data.Store.ID.ValueString(), data.Store.Model.ValueString(), data.User.ValueString(), data.Relation.ValueString(), data.Document.ValueString())

	tflog.Info(ctx, fmt.Sprintf("Read a tuple with id %s", createdTuple.ID))

	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to create tuple, got error: %s", err.Error()),
		)
		return
	}

	// Set the created models's ID in the Terraform state
	data.ID = types.StringValue(createdTuple.ID)

	tflog.Info(ctx, fmt.Sprintf("Model with id %s got created", createdTuple.ID))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *tupleresource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TupleResource

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read a tuple with id %s", data.ID.ValueString()))

	// Retrieve the model using the GetTuple method
	tuple, err := r.client.GetTuple(ctx, data.Store.ID.ValueString(), data.Store.Model.ValueString(), data.User.ValueString(), data.Relation.ValueString(), data.Document.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read tuple, got error: %s", err.Error()),
		)
		return
	}

	// Update the data model with the retrieved model information
	data.ID = types.StringValue(tuple.ID)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *tupleresource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ModelResouce

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Update a model with id %s", data.ID))
}

func (r *tupleresource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TupleResource

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Delete a tuple with id %s", data.ID.ValueString()))

	err := r.client.DeleteTuple(ctx, data.Store.ID.ValueString(), data.Store.Model.ValueString(), data.User.ValueString(), data.Relation.ValueString(), data.Document.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to delete tuple, got error: %s", err.Error()),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Tuple with id %s got deleted", data.ID.ValueString()))
}

func (r *tupleresource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
