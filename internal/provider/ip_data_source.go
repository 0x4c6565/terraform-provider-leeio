// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	hreq "github.com/imroc/req/v3"
)

var _ datasource.DataSource = &IPDataSource{}

func NewIPDataSource() datasource.DataSource {
	return &IPDataSource{}
}

// IPDataSource defines the data source implementation.
type IPDataSource struct {
	client *hreq.Client
}

// IPDataSourceModel describes the data source data model.
type IPDataSourceModel struct {
	Address     types.String `tfsdk:"address"`
	AddressIPv4 types.String `tfsdk:"address_ipv4"`
	AddressIPv6 types.String `tfsdk:"address_ipv6"`
}

func (d *IPDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip"
}

func (d *IPDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This data source returns the source IP address of the execution environment from where Terraform is running",

		Attributes: map[string]schema.Attribute{
			"address": schema.StringAttribute{
				MarkdownDescription: "IP address",
				Computed:            true,
			},
			"address_ipv4": schema.StringAttribute{
				MarkdownDescription: "IPv4 address. Null if not routable",
				Computed:            true,
			},
			"address_ipv6": schema.StringAttribute{
				MarkdownDescription: "IPv6 address. Null if not routable",
				Computed:            true,
			},
		},
	}
}

func (d *IPDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client, _ = req.ProviderData.(*hreq.Client)
}

func (d *IPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data IPDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var address string

	hresp, err := d.client.R().
		SetSuccessResult(&address).
		SetHeader("Accept", "application/json").
		EnableDump().
		Get("https://ip.lee.io")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read IP address, got error: %s", err))
		tflog.Debug(ctx, hresp.Dump())
		return
	}

	data.Address = types.StringValue(address)

	var ipv4Address string
	hresp, err = d.client.R().
		SetSuccessResult(&ipv4Address).
		SetHeader("Accept", "application/json").
		EnableDump().
		Get("https://ipv4.lee.io")
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to retrieve IPv4 address, got error: %s", err))
		tflog.Debug(ctx, hresp.Dump())
	} else {
		data.AddressIPv4 = types.StringValue(ipv4Address)
	}

	var ipv6Address string
	hresp, err = d.client.R().
		SetSuccessResult(&ipv6Address).
		SetHeader("Accept", "application/json").
		EnableDump().
		Get("https://ipv6.lee.io")
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to retrieve IPv6 address, got error: %s", err))
		tflog.Debug(ctx, hresp.Dump())
	} else {
		data.AddressIPv6 = types.StringValue(ipv6Address)
	}

	tflog.Trace(ctx, "read a data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
