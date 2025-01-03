// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	hreq "github.com/imroc/req/v3"
)

// Ensure LeeIOProvider satisfies various provider interfaces.
var _ provider.Provider = &LeeIOProvider{}
var _ provider.ProviderWithFunctions = &LeeIOProvider{}

// LeeIOProvider defines the provider implementation.
type LeeIOProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// LeeIOProviderModel describes the provider data model.
type LeeIOProviderModel struct{}

func (p *LeeIOProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "leeio"
	resp.Version = p.version
}

func (p *LeeIOProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{},
	}
}

func (p *LeeIOProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data LeeIOProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	client := hreq.C()
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *LeeIOProvider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

func (p *LeeIOProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewIPDataSource,
	}
}

func (p *LeeIOProvider) Functions(ctx context.Context) []func() function.Function {
	return nil
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &LeeIOProvider{
			version: version,
		}
	}
}
