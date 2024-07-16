package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-elastic/pkg/elastic"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type roleMappingBuilder struct {
	resourceType         *v2.ResourceType
	client               *elastic.Client
	shouldSyncDeployment bool
}

func (d *roleMappingBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return d.resourceType
}

// Create a new connector resource for Elastic deployment user.
func roleMappingResource(role string) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"role_mapping_id":   role,
		"role_mapping_name": role,
	}

	status := v2.UserTrait_Status_STATUS_UNSPECIFIED
	userTraitOptions := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithStatus(status),
	}

	ret, err := rs.NewUserResource(
		role,
		roleMappingResourceType,
		role,
		userTraitOptions,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// List returns all the role mappings.
func (d *roleMappingBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if !d.shouldSyncDeployment {
		return nil, "", nil, nil
	}

	roles, err := d.client.ListDeploymentRoleMapping(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error listing deployment users: %w", err)
	}

	var rv []*v2.Resource
	for role := range roles {
		ur, err := roleMappingResource(role)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error creating role mapping resource %s: %w", role, err)
		}
		rv = append(rv, ur)
	}

	return rv, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (d *roleMappingBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (d *roleMappingBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newRoleMappingBuilder(client *elastic.Client, shouldSyncDeployment bool) *roleMappingBuilder {
	return &roleMappingBuilder{
		resourceType:         roleMappingResourceType,
		client:               client,
		shouldSyncDeployment: shouldSyncDeployment,
	}
}
