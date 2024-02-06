package main

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-sdk/pkg/cli"
	"github.com/spf13/cobra"
)

// config defines the external configuration required for the connector to run.
type config struct {
	cli.BaseConfig `mapstructure:",squash"` // Puts the base config options in the same place as the connector options

	ApiKey             string `mapstructure:"api-key"`
	OrganizationID     string `mapstructure:"organization-id,omitempty"`
	DeploymentApiKey   string `mapstructure:"deployment-api-key"`
	DeploymentEndpoint string `mapstructure:"deployment-endpoint"`
}

// validateConfig is run after the configuration is loaded, and should return an error if it isn't valid.
func validateConfig(ctx context.Context, cfg *config) error {
	if cfg.ApiKey == "" {
		return fmt.Errorf("api key is missing, please provide it via --api-key flag or $BATON_API_KEY environment variable")
	}

	return nil
}

// cmdFlags sets the cmdFlags required for the connector.
func cmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().String("api-key", "", "Elastic API key used to communicate with Elastic cloud API. ($BATON_API_KEY)")
	cmd.PersistentFlags().String("organization-id", "", "Optional. Provide your Elastic organization ID if you want to sync members of a single organization. ($BATON_ORGANIZATION_ID)")
	cmd.PersistentFlags().String("deployment-api-key", "", "API key of your elasticsearch deployment. ($BATON_DEPLOYMENT_API_KEY)")
	cmd.PersistentFlags().String("deployment-endpoint", "", "Elasticsearch endpoint used to sync deployment resources. ($BATON_DEPLOYMENT_ENDPOINT)")

	cmd.MarkFlagsRequiredTogether("deployment-api-key", "deployment-endpoint")
}
