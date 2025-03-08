package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

var _ datasource.DataSource = &CodeDataSource{}

func NewCodeDataSource() datasource.DataSource {
	return &CodeDataSource{}
}

type CodeDataSource struct {
}

type CodeDataSourceModel struct {
	Code    types.String     `tfsdk:"code"`
	Options *EvaluateOptions `tfsdk:"options"`
	Output  types.String     `tfsdk:"output"`
}

func (d *CodeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_code"
}

func (d *CodeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Evaluates the provided code string as Jsonnet",

		Attributes: map[string]schema.Attribute{
			"code": schema.StringAttribute{
				MarkdownDescription: "The Jsonnet code to be evaluated",
				Required:            true,
			},
			"options": schema.SingleNestedAttribute{
				MarkdownDescription: "Additional options to be passed to the Jsonnet command",
				Attributes: map[string]schema.Attribute{
					"jpaths": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
					},
				},
				Optional: true,
			},
			"output": schema.StringAttribute{
				MarkdownDescription: "The evaluated Jsonnet, as a string",
				Computed:            true,
			},
		},
	}
}

func (d *CodeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
}

func (d *CodeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CodeDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var options []EvaluateOptions
	if data.Options != nil {
		options = append(options, *data.Options)
	}
	output, err := evaluate(data.Code.ValueString(), options)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to evaluate Jsonnet, got error: %s", err))
		return
	}

	data.Output = types.StringValue(strings.TrimSpace(output))

	tflog.Trace(ctx, "read a data source")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
