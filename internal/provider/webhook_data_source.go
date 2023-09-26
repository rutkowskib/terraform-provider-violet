package provider

import (
	"context"
	"fmt"
	"github.com/rutkowskib/terraform-provider-violet/internal/violet"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &webhookDataSource{}
	_ datasource.DataSourceWithConfigure = &webhookDataSource{}
)

func WebhookDataSource() datasource.DataSource {
	return &webhookDataSource{}
}

type webhookDataSource struct {
	client *violet.VioletClient
}

func (d *webhookDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*violet.VioletClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected violet.VioletClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *webhookDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "violet_webhook"
}

type webhookModel struct {
	Id               types.Int64  `tfsdk:"id"`
	AppId            types.Int64  `tfsdk:"app_id"`
	Event            types.String `tfsdk:"event"`
	RemoteEndpoint   types.String `tfsdk:"remote_endpoint"`
	Status           types.String `tfsdk:"status"`
	DateCreated      types.String `tfsdk:"date_created"`
	DateLastModified types.String `tfsdk:"date_last_modified"`
}

func (d *webhookDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source to get data of existing Violet webhook",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Required: true,
				Description: "Webhook id",
			},
			"app_id": schema.Int64Attribute{
				Computed: true,
				Description: "App Id of application this webhook belongs to",
			},
			"event": schema.StringAttribute{
				Computed: true,
				Description: "Event webhook will be subscribed to",
			},
			"remote_endpoint": schema.StringAttribute{
				Computed: true,
				Description: "Endpoint that webhook will be publishing to",
			},
			"status": schema.StringAttribute{
				Computed: true,
				Description: "Status of webhook",
			},
			"date_created": schema.StringAttribute{
				Computed: true,
				Description: "Creation date of webhook",
			},
			"date_last_modified": schema.StringAttribute{
				Computed: true,
				Description: "Date of last modification of the webhook",
			},
		},
	}
}

func (d *webhookDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data webhookModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueInt64()

	tflog.Info(ctx, "Read webhookDataSource", map[string]interface{}{
		"id": id,
	})

	webhook := d.client.GetWebhook(ctx, id)

	state := webhookModel{
		Id:               types.Int64Value(webhook.Id),
		AppId:            types.Int64Value(webhook.AppId),
		Event:            types.StringValue(webhook.Event),
		RemoteEndpoint:   types.StringValue(webhook.RemoteEndpoint),
		Status:           types.StringValue(webhook.Status),
		DateCreated:      types.StringValue(webhook.DateCreated),
		DateLastModified: types.StringValue(webhook.DateLastModified),
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
