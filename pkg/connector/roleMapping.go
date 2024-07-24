package connector

import (
	"context"
	"fmt"
	"slices"

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

type roleMappingBuilder struct {
	resourceType         *v2.ResourceType
	client               *elastic.Client
	shouldSyncDeployment bool
}

const NF = -1

func (r *roleMappingBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return r.resourceType
}

// Create a new connector resource for role mapping.
func roleMappingResource(role string) (*v2.Resource, error) {
	profile := map[string]interface{}{
		"role_mapping_id":   role,
		"role_mapping_name": role,
	}

	roleOptions := []rs.RoleTraitOption{
		rs.WithRoleProfile(profile),
	}

	ret, err := rs.NewRoleResource(
		role,
		roleMappingResourceType,
		role,
		roleOptions,
	)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// List returns all the role mappings.
func (r *roleMappingBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	if !r.shouldSyncDeployment {
		return nil, "", nil, nil
	}

	roles, err := r.client.ListDeploymentRoleMapping(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error listing role mappings: %w", err)
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
func (d *roleMappingBuilder) Entitlements(ctx context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement
	roles, err := d.client.ListDeploymentRoleMapping(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for roleMappingName, role := range roles {
		for _, mappingRole := range role.Roles {
			assignmentOptions := []ent.EntitlementOption{
				ent.WithGrantableTo(roleMappingResourceType),
				ent.WithDisplayName(fmt.Sprintf("%s Role %s", resource.DisplayName, roleMappingName)),
				ent.WithDescription(fmt.Sprintf("Member of %s elasticsearch role", resource.DisplayName)),
			}

			rv = append(rv, ent.NewAssignmentEntitlement(
				resource,
				mappingRole,
				assignmentOptions...,
			))
		}
	}

	return rv, "", nil, nil
}

// GetRoleMappingUsers returns role mapping users.
func (r *roleMappingBuilder) GetRoleMappingUsers(ctx context.Context, name string) ([]string, error) {
	var users []string
	roles, err := r.client.GetDeploymentRoleMapping(ctx, name)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if field := role.Rules.(map[string]any)["field"]; field != nil {
			userData := Utility{
				Data: fmt.Sprintf("%s", field.(map[string]any)["username"]),
			}
			users = userData.TrimPrefix("[").TrimSuffix("]").Split(" ")
		}
	}

	return users, nil
}

// Grants always returns an empty slice for users since they don't have any entitlements.
func (r *roleMappingBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	var rv []*v2.Grant
	roles, err := r.client.ListDeploymentRoleMapping(ctx)
	if err != nil {
		return nil, "", nil, err
	}

	for roleMappingName, role := range roles {
		if roleMappingName != resource.Id.Resource {
			continue
		}

		if field := role.Rules.(map[string]any)["field"]; field != nil {
			userData := Utility{
				Data: fmt.Sprintf("%s", field.(map[string]any)["username"]),
			}
			users := userData.TrimPrefix("[").TrimSuffix("]").Split(" ")
			for _, userName := range users {
				ur, err := deploymentUserResource(&elastic.DeploymentUser{
					Username: userName,
				})
				if err != nil {
					return nil, "", nil, fmt.Errorf("error creating role mapping resource for user %s: %w", resource.Id.Resource, err)
				}

				for _, roleName := range role.Roles {
					gr := grant.NewGrant(resource, roleName, ur.Id)
					rv = append(rv, gr)
				}
			}
		}
	}

	return rv, "", nil, nil
}

func (r *roleMappingBuilder) Grant(ctx context.Context, principal *v2.Resource, entitlement *v2.Entitlement) (annotations.Annotations, error) {
	l := ctxzap.Extract(ctx)
	if principal.Id.ResourceType != deploymentUserResourceType.Id {
		l.Warn(
			"baton-elastic: only users can be granted role mapping membership",
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("principal_id", principal.Id.Resource),
		)
		return nil, fmt.Errorf("baton-elastic: only users can be granted role mapping membership")
	}

	user, err := r.client.GetDeploymentUser(ctx, principal.Id.Resource)
	if err != nil {
		return nil, fmt.Errorf("error fetching user: %w", err)
	}

	newUser := user[principal.Id.Resource]
	roleMappingName := entitlement.Resource.Id.Resource
	users, err := r.GetRoleMappingUsers(ctx, roleMappingName)
	if err != nil {
		return nil, err
	}

	userPos := slices.IndexFunc(users, func(c string) bool {
		return c == newUser.Username
	})
	if userPos != NF {
		l.Warn(
			"baton-elastic: user already has this role mapping",
			zap.String("principal_id", principal.Id.String()),
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("roleMappingName", roleMappingName),
		)
		return nil, fmt.Errorf("baton-elastic: user %s already has this role mapping", principal.DisplayName)
	}

	users = append(users, newUser.Username)
	data := elastic.MappingRolesBody{
		Roles:   newUser.Roles,
		Enabled: true,
		Rules: elastic.Rule{
			Field: elastic.Field{
				Username: users,
			},
		},
	}
	err = r.client.UpdateUserMappingRole(ctx, data, roleMappingName)
	if err != nil {
		return nil, fmt.Errorf("baton-elastic: failed to grant role mapping to user: %w", err)
	}

	l.Warn("Role Mapping Membership has been created.",
		zap.String("roleMappingName", roleMappingName),
		zap.String("User", newUser.Username),
	)

	return nil, nil
}

func (r *roleMappingBuilder) Revoke(ctx context.Context, grant *v2.Grant) (annotations.Annotations, error) {
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

	newUser := user[principal.Id.Resource]
	roleMappingName := entitlement.Resource.Id.Resource
	users, err := r.GetRoleMappingUsers(ctx, roleMappingName)
	if err != nil {
		return nil, err
	}

	userPos := slices.IndexFunc(users, func(c string) bool {
		return c == newUser.Username
	})
	if userPos == NF {
		l.Warn(
			"baton-elastic: user does not have this role mapping",
			zap.String("principal_id", principal.Id.String()),
			zap.String("principal_type", principal.Id.ResourceType),
			zap.String("roleMappingName", roleMappingName),
		)
		return nil, fmt.Errorf("baton-elastic: user %s does not have this role mapping", principal.DisplayName)
	}

	users = append(users[:userPos], users[userPos+1:]...)
	data := elastic.MappingRolesBody{
		Roles:   newUser.Roles,
		Enabled: true,
		Rules: elastic.Rule{
			Field: elastic.Field{
				Username: users,
			},
		},
	}
	err = r.client.UpdateUserMappingRole(ctx, data, roleMappingName)
	if err != nil {
		return nil, fmt.Errorf("baton-elastic: failed to revoke role mapping to user: %w", err)
	}

	l.Warn("Role Membership has been revoked.",
		zap.String("role Mapping", roleMappingName),
		zap.String("User", newUser.Username),
	)

	return nil, nil
}

func newRoleMappingBuilder(client *elastic.Client, shouldSyncDeployment bool) *roleMappingBuilder {
	return &roleMappingBuilder{
		resourceType:         roleMappingResourceType,
		client:               client,
		shouldSyncDeployment: shouldSyncDeployment,
	}
}
