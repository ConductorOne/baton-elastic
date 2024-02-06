package connector

import (
	"context"
	"fmt"

	"github.com/conductorone/baton-elastic/pkg/elastic"
	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/annotations"
	"github.com/conductorone/baton-sdk/pkg/pagination"
	ent "github.com/conductorone/baton-sdk/pkg/types/entitlement"
	grant "github.com/conductorone/baton-sdk/pkg/types/grant"
	rs "github.com/conductorone/baton-sdk/pkg/types/resource"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const roleMembership = "member"

type roleBuilder struct {
	resourceType         *v2.ResourceType
	client               *elastic.Client
	shouldSyncDeployment bool
}

func (r *roleBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return r.resourceType
}

// Create a new connector resource for Elastic deployment role.
func deploymentRoleResource(role string) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"role_name": role,
		// no ID in api response
		"role_id": role,
	}

	roleOptions := []rs.RoleTraitOption{
		rs.WithRoleProfile(profile),
	}

	ret, err := rs.NewRoleResource(
		role,
		deploymentRoleResourceType,
		role,
		roleOptions,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// List returns all the roles from the database as resource objects.
// Roles include a RoleTrait because they are the 'shape' of a standard role.
func (r *roleBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if !r.shouldSyncDeployment {
		return nil, "", nil, nil
	}

	roles, err := r.client.ListDeploymentRoles(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error listing roles: %w", err)
	}

	var rv []*v2.Resource
	for key := range roles {
		ur, err := deploymentRoleResource(key)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error creating role resource for role %s: %w", key, err)
		}
		rv = append(rv, ur)
	}

	return rv, "", nil, nil
}

func (r *roleBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement
	assignmentOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(deploymentUserResourceType),
		ent.WithDisplayName(fmt.Sprintf("%s Role %s", resource.DisplayName, roleMembership)),
		ent.WithDescription(fmt.Sprintf("Member of %s elasticsearch role", resource.DisplayName)),
	}

	rv = append(rv, ent.NewAssignmentEntitlement(
		resource,
		roleMembership,
		assignmentOptions...,
	))

	return rv, "", nil, nil
}

func (r *roleBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	users, err := r.client.ListDeploymentUsers(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	var rv []*v2.Grant
	for _, user := range users {
		userCopy := user
		if hasRole(resource.Id.Resource, user.Roles) {
			ur, err := deploymentUserResource(&userCopy)
			if err != nil {
				return nil, "", nil, fmt.Errorf("error creating user resource for role %s: %w", resource.Id.Resource, err)
			}
			gr := grant.NewGrant(resource, roleMembership, ur.Id)
			rv = append(rv, gr)
		}
	}

	return rv, "", nil, nil
}

func (r *roleBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)

	if principal.Id.ResourceType != deploymentUserResourceType.Id {
		l.Warn(
			"baton-elastic: only users can be granted role membership",
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("principal_id", principal.Id.Resource),
		)
		return nil, fmt.Errorf("baton-elastic: only users can be granted role membership")
	}

	user, err := r.client.GetDeploymentUser(ctx, principal.Id.Resource)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	newUser := user[principal.Id.Resource]
	newUser.Roles = append(newUser.Roles, entitlement.Resource.Id.Resource)

	err = r.client.UpdateUser(ctx, principal.Id.Resource, newUser)
	if err != nil {
		return nil, fmt.Errorf("baton-elastic: failed to grant role to user: %w", err)
	}

	return nil, nil
}

func (r *roleBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	principal := grant.Principal
	entitlement := grant.Entitlement

	if principal.Id.ResourceType != deploymentUserResourceType.Id {
		l.Warn(
			"baton-elastic: only users can have role membership revoked",
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("principal_id", principal.Id.Resource),
		)
		return nil, fmt.Errorf("baton-elastic: only users can have role membership revoked")
	}

	user, err := r.client.GetDeploymentUser(ctx, principal.Id.Resource)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	var roles []string
	for _, role := range user[principal.Id.Resource].Roles {
		if role != entitlement.Resource.Id.Resource {
			roles = append(roles, role)
		}
	}

	newUser := user[principal.Id.Resource]
	newUser.Roles = roles
	err = r.client.UpdateUser(ctx, principal.Id.Resource, newUser)
	if err != nil {
		return nil, fmt.Errorf("baton-elastic: failed to revoke user role: %w", err)
	}

	return nil, nil
}

func newDeploymentRoleBuilder(client *elastic.Client, shouldSyncDeployment bool) *roleBuilder {
	return &roleBuilder{
		resourceType:         deploymentRoleResourceType,
		client:               client,
		shouldSyncDeployment: shouldSyncDeployment,
	}
}

func hasRole(target string, roles []string) bool {
	for _, role := range roles {
		if target == role {
			return true
		}
	}
	return false
}
