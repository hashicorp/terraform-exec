package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func main() {
	providerserver.Serve(context.Background(), New, providerserver.ServeOpts{
		Address: "terraform-exec/registry/test",
	})
}

type testProvider struct {
}

func New() provider.Provider {
	return &testProvider{}
}

func (p *testProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "test"
	resp.Version = "latest"
}

// Schema defines the provider-level schema for configuration data.
func (p *testProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Provider for testing",
	}
}

func (p *testProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, "Configuring test provider")
}

func (p *testProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *testProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLogResource,
	}
}

func NewLogResource() resource.Resource {
	return &createIndefinitelyResource{}
}

type createIndefinitelyResource struct {
}

func (p *createIndefinitelyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_create_indefinitely"
}

func (p *createIndefinitelyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resourceschema.Schema{
		MarkdownDescription: "a resource that cannot be created as it hangs indefinitely on creation",
	}
}

func (r *createIndefinitelyResource) Configure(ctx context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	// do nothing
}

func (p *createIndefinitelyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// do nothing
}

func (p *createIndefinitelyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	for {
		select {
		case <-ctx.Done():
			tflog.Info(ctx, "context is done, stopping creation")
			return
		default:
			tflog.Info(ctx, "creating log resource (will create until context is canceled)")
		}
	}
}
func (p *createIndefinitelyResource) Update(ctx context.Context, request resource.UpdateRequest, response *resource.UpdateResponse) {
	// do nothing
}

func (p *createIndefinitelyResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	// do nothing
}
