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
)

const orgMembership = "member"

type organizationBuilder struct {
	resourceType *v2.ResourceType
	client       *elastic.Client
}

func (r *organizationBuilder) ResourceType(ctx context.Context) *v2.ResourceType {
	return r.resourceType
}

// Create a new connector resource for Elastic organization.
func organizationResource(organization elastic.Organization) (*v2.Resource, error) {
	ret, err := rs.NewResource(
		organization.Name,
		organizationResourceType,
		organization.ID,
		rs.WithAnnotation(
			&v2.ChildResourceType{ResourceTypeId: userResourceType.Id},
		))

	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *organizationBuilder) List(ctx context.Context, parentResourceID *v2.ResourceId, pToken *pagination.Token) ([]*v2.Resource, string, annotations.Annotations, error) {
	orgs, err := r.client.ListOrganizations(ctx)
	if err != nil {
		return nil, "", nil, fmt.Errorf("error listing organizations: %w", err)
	}

	var rv []*v2.Resource
	for _, org := range orgs {
		or, err := organizationResource(org)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error creating organization resource: %w", err)
		}
		rv = append(rv, or)
	}

	return rv, "", nil, nil
}

func (r *organizationBuilder) Entitlements(_ context.Context, resource *v2.Resource, _ *pagination.Token) ([]*v2.Entitlement, string, annotations.Annotations, error) {
	var rv []*v2.Entitlement
	assignmentOptions := []ent.EntitlementOption{
		ent.WithGrantableTo(userResourceType),
		ent.WithDisplayName(fmt.Sprintf("%s Organization %s", resource.DisplayName, orgMembership)),
		ent.WithDescription(fmt.Sprintf("Member of %s Elastic organization", resource.DisplayName)),
	}

	rv = append(rv, ent.NewAssignmentEntitlement(
		resource,
		orgMembership,
		assignmentOptions...,
	))

	return rv, "", nil, nil
}

func (r *organizationBuilder) Grants(ctx context.Context, resource *v2.Resource, pToken *pagination.Token) ([]*v2.Grant, string, annotations.Annotations, error) {
	members, err := r.client.ListOrgMembers(ctx, resource.Id.Resource)
	if err != nil {
		return nil, "", nil, err
	}

	var rv []*v2.Grant
	for _, member := range members {
		memberCopy := member
		ur, err := userResource(&memberCopy, resource.Id)
		if err != nil {
			return nil, "", nil, fmt.Errorf("error creating user resource for organization %s: %w", resource.Id.Resource, err)
		}
		gr := grant.NewGrant(resource, orgMembership, ur.Id)
		rv = append(rv, gr)
	}

	return rv, "", nil, nil
}

func newOrganizationBuilder(client *elastic.Client) *organizationBuilder {
	return &organizationBuilder{
		resourceType: organizationResourceType,
		client:       client,
	}
}
