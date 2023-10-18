package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/rutkowskib/terraform-provider-violet/internal/violet"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &WebhookResource{}
	_ resource.ResourceWithConfigure   = &WebhookResource{}
	_ resource.ResourceWithImportState = &WebhookResource{}
)

// NewWebhookResource is a helper function to simplify the provider implementation.
func NewWebhookResource() resource.Resource {
	return &WebhookResource{}
}

type WebhookResource struct {
	client *violet.VioletClient
}

// Metadata returns the resource type name.
func (r *WebhookResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webhook"
}

// Configure adds the provider configured client to the resource.
func (r *WebhookResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*violet.VioletClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

type WebhookResourceModel struct {
	Id               types.Int64  `tfsdk:"id"`
	AppId            types.Int64  `tfsdk:"app_id"`
	Event            types.String `tfsdk:"event"`
	RemoteEndpoint   types.String `tfsdk:"remote_endpoint"`
	Status           types.String `tfsdk:"status"`
	DateCreated      types.String `tfsdk:"date_created"`
	DateLastModified types.String `tfsdk:"date_last_modified"`
}

// Schema defines the schema for the resource.
func (r *WebhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Resource to manage Violet webhook",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:    true,
				Description: "Webhook id",
			},
			"app_id": schema.Int64Attribute{
				Computed:    true,
				Description: "App Id of application this webhook belongs to",
			},
			"event": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Event webhook will be subscribed to",
			},
			"remote_endpoint": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Endpoint that webhook will be publishing to",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Status of webhook",
			},
			"date_created": schema.StringAttribute{
				Computed:    true,
				Description: "Creation date of webhook",
			},
			"date_last_modified": schema.StringAttribute{
				Computed:    true,
				Description: "Date of last modification of the webhook",
			},
		},
	}
}

func (r *WebhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	id, _ := strconv.Atoi(req.ID)
	state := WebhookResourceModel{
		Id: types.Int64Value(int64(id)),
	}

	resp.State.Set(ctx, state)
}

// Create creates the resource and sets the initial Terraform state.
func (r *WebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan WebhookResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "create", map[string]interface{}{
		"event":           plan.Event.ValueString(),
		"remote_endpoint": plan.RemoteEndpoint.ValueString(),
	})

	input := violet.CreateWebhookInput{
		Event:          plan.Event.ValueString(),
		RemoteEndpoint: plan.RemoteEndpoint.ValueString(),
	}
	err, webhook := r.client.CreateWebhook(ctx, input)

	if err != nil {
		tflog.Error(ctx, "Error creating webhook data", map[string]interface{}{
			"err": err.Error(),
		})
		resp.Diagnostics.AddError(
			"Error creating Violet webhook",
			"Creating webhook failed: "+err.Error(),
		)
		return
	}

	state := WebhookResourceModel{
		Id:               types.Int64Value(webhook.Id),
		AppId:            types.Int64Value(webhook.AppId),
		Event:            types.StringValue(webhook.Event),
		RemoteEndpoint:   types.StringValue(webhook.RemoteEndpoint),
		Status:           types.StringValue(webhook.Status),
		DateCreated:      types.StringValue(webhook.DateCreated),
		DateLastModified: types.StringValue(webhook.DateLastModified),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *WebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var oldState WebhookResourceModel
	diags := req.State.Get(ctx, &oldState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := oldState.Id.ValueInt64()

	tflog.Info(ctx, "Read webhook resource", map[string]interface{}{
		"id": id,
	})

	err, webhook := r.client.GetWebhook(ctx, id)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error reading Violet webhook id: %d", id),
			"Get webhook failed: "+err.Error(),
		)
		return
	}

	state := WebhookResourceModel{
		Id:               types.Int64Value(webhook.Id),
		AppId:            types.Int64Value(webhook.AppId),
		Event:            types.StringValue(webhook.Event),
		RemoteEndpoint:   types.StringValue(webhook.RemoteEndpoint),
		Status:           types.StringValue(webhook.Status),
		DateCreated:      types.StringValue(webhook.DateCreated),
		DateLastModified: types.StringValue(webhook.DateLastModified),
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *WebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *WebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WebhookResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueInt64()

	tflog.Info(ctx, "Delete webhook resource", map[string]interface{}{
		"id": id,
	})

	err := r.client.DeleteWebhook(ctx, id)

	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error deleting Violet webhook id: %d", id),
			"Delete webhook failed: "+err.Error(),
		)
	}
}
