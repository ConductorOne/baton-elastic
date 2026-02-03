package connector

import (
	"testing"

	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/stretchr/testify/assert"
)

func TestRoleMappingBuilderGrants(t *testing.T) {
	if apiKey == "" && organizationID == "" && deploymentApiKey == "" && deploymentEndpoint == "" {
		t.Skip()
	}

	opts := rs.SyncOpAttrs{}
	cli := getClientForTesting(ctx)
	assert.Nil(t, cli)

	p := &roleMappingBuilder{
		resourceType:         roleMappingResourceType,
		client:               cli,
		shouldSyncDeployment: true,
	}

	roleMapping := "mapping7"
	resource, err := roleMappingResource(roleMapping)
	assert.Nil(t, err)

	_, _, err1 := p.Grants(ctx, resource, opts)
	assert.Nil(t, err1)
}

func TestGetUsers(t *testing.T) {
	if apiKey == "" && organizationID == "" && deploymentApiKey == "" && deploymentEndpoint == "" {
		t.Skip()
	}

	cli := getClientForTesting(ctx)
	assert.Nil(t, cli)

	p := &roleMappingBuilder{
		resourceType:         roleMappingResourceType,
		client:               cli,
		shouldSyncDeployment: true,
	}

	roleMapping := "mapping7"
	_, err := p.GetRoleMappingUsers(ctx, roleMapping)
	assert.Nil(t, err)
}
