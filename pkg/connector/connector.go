package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-elastic/pkg/elastic"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/connectorbuilder"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
)

type Connector struct {
	client               *elastic.Client
	shouldSyncDeployment bool
}

// ResourceSyncers returns a ResourceSyncer for each resource type that should be synced from the upstream service.
func (d *Connector) ResourceSyncers(ctx context.Context) []connectorbuilder.ResourceSyncer {
	return []connectorbuilder.ResourceSyncer{
		newOrganizationBuilder(d.client),
		newUserBuilder(d.client),
		newDeploymentRoleBuilder(d.client, d.shouldSyncDeployment),
		newDeploymentUserBuilder(d.client, d.shouldSyncDeployment),
		newRoleMappingBuilder(d.client, d.shouldSyncDeployment),
	}
}

// Metadata returns metadata about the connector.
func (d *Connector) Metadata(ctx context.Context) (*v2.ConnectorMetadata, error) {
	return &v2.ConnectorMetadata{
		DisplayName: "Elastic connector",
		Description: "Connector syncing users and roles from Elastic cloud and optionally from elasticsearch deployment.",
	}, nil
}

// Validate is called to ensure that the connector is properly configured. It should exercise any API credentials
// to be sure that they are valid.
func (d *Connector) Validate(ctx context.Context) (annotations.Annotations, error) {
	_, err := d.client.ListOrganizations(ctx)
	if err != nil {
		return nil, fmt.Errorf("error validating elastic cloud credentials: %w", err)
	}

	if d.shouldSyncDeployment {
		err := d.client.DeploymentAuth(ctx)
		if err != nil {
			return nil, fmt.Errorf("error validating elasticsearch deployment credentials: %w", err)
		}
	}

	return nil, nil
}

// New returns a new instance of the connector.
func New(ctx context.Context, deploymentApiKey, deploymentEndpoint, apiKey, organizationID string) (*Connector, error) {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil, err
	}

	shouldSyncDeployment := false
	if deploymentEndpoint != "" {
		shouldSyncDeployment = true
	}

	return &Connector{
		client:               elastic.NewClient(httpClient, deploymentApiKey, deploymentEndpoint, apiKey, organizationID),
		shouldSyncDeployment: shouldSyncDeployment,
	}, nil
}
