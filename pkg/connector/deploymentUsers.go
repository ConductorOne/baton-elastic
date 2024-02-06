package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-elastic/pkg/elastic"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/helpers"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
)

type deploymentUserBuilder struct {
	resourceType         *v2.ResourceType
	client               *elastic.Client
	shouldSyncDeployment bool
}

func (d *deploymentUserBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return d.resourceType
}

// Create a new connector resource for Elastic deployment user.
func deploymentUserResource(user *elastic.DeploymentUser) (*v2.Resource, error) {
	firstname, lastname := helpers.SplitFullName(user.FullName)
	profile := map[string]interface{}{
		"first_name": firstname,
		"last_name":  lastname,
		"login":      user.Email,
		// username is an ID
		"user_id": user.Username,
	}

	var status v2.UserTrait_Status_Status
	status = v2.UserTrait_Status_STATUS_UNSPECIFIED

	if user.Enabled {
		status = v2.UserTrait_Status_STATUS_ENABLED
	}

	userTraitOptions := []rs.UserTraitOption{
		rs.WithUserProfile(profile),
		rs.WithEmail(user.Email, true),
		rs.WithStatus(status),
	}

	ret, err := rs.NewUserResource(
		user.FullName,
		deploymentUserResourceType,
		user.Username,
		userTraitOptions,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// List returns all the users from the database as resource objects.
// Users include a UserTrait because they are the 'shape' of a standard user.
func (d *deploymentUserBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if !d.shouldSyncDeployment {
		return nil, "", nil, nil
	}

	users, err := d.client.ListDeploymentUsers(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error listing deployment users: %w", err)
	}

	var rv []*v2.Resource
	for key := range users {
		userCopy := users[key]
		ur, err := deploymentUserResource(&userCopy)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error creating user resource for deployment user %s: %w", key, err)
		}
		rv = append(rv, ur)
	}

	return rv, "", nil, nil
}

// Entitlements always returns an empty slice for users.
func (d *deploymentUserBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (d *deploymentUserBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	return nil, "", nil, nil
}

func newDeploymentUserBuilder(client *elastic.Client, shouldSyncDeployment bool) *deploymentUserBuilder {
	return &deploymentUserBuilder{
		resourceType:         deploymentUserResourceType,
		client:               client,
		shouldSyncDeployment: shouldSyncDeployment,
	}
}
