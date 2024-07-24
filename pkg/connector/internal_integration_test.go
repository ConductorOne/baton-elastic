package connector

import (
	"testing"

	"github.com/conductorone/baton-sdk/pkg/pagination"
	"github.com/stretchr/testify/assert"
)

func TestRoleMappingBuilderGrants(t *testing.T) {
	if apiKey == "" && organizationID == "" && deploymentApiKey == "" && deploymentEndpoint == "" {
		t.Skip()
	}

	pToken := &pagination.Token{}
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

	_, _, _, err1 := p.Grants(ctx, resource, pToken)
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
