package cmd

import (
	"context"

	"github.com/zeiss/terraform-provider-openfga/internal/cfg"
	"github.com/zeiss/terraform-provider-openfga/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/spf13/cobra"
	"github.com/zeiss/pkg/server"
)

var (
	config  *cfg.Config
	version string = "dev"
)

func init() {
	config = cfg.New()

	Root.PersistentFlags().BoolVar(&config.Flags.Debug, "debug", config.Flags.Debug, "debug")

	Root.SilenceUsage = true
	Root.SilenceErrors = true
}

var Root = &cobra.Command{
	Version: version,
	RunE: func(cmd *cobra.Command, args []string) error {
		p := NewProvider(config)

		s, _ := server.WithContext(cmd.Context())
		s.Listen(p, false)

		return s.Wait()
	},
}

var _ server.Listener = (*Provider)(nil)

// Provider is the server that implements the Noop interface.
type Provider struct {
	cfg *cfg.Config
}

// NewProvider returns a new instance of Provider.
func NewProvider(cfg *cfg.Config) *Provider {
	return &Provider{cfg}
}

// Start starts the server.
func (s *Provider) Start(ctx context.Context, ready server.ReadyFunc, run server.RunFunc) func() error {
	return func() error {
		opts := providerserver.ServeOpts{
			Address: "registry.terraform.io/zeiss/openfga",
			Debug:   s.cfg.Flags.Debug,
		}

		err := providerserver.Serve(ctx, provider.New(version), opts)
		if err != nil {
			return err
		}

		return nil
	}
}
