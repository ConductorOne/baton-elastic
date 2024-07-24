package connector

import (
	"context"
	"os"
	"testing"

	"github.com/conductorone/baton-elastic/pkg/elastic"
	"github.com/conductorone/baton-sdk/pkg/uhttp"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/stretchr/testify/assert"
)

var (
	apiKey             = os.Getenv("BATON_API_KEY=")
	organizationID     = os.Getenv("BATON_ORGANIZATION_ID")
	deploymentApiKey   = os.Getenv("BATON_DEPLOYMENT_API_KEY==")
	deploymentEndpoint = os.Getenv("BATON_DEPLOYMENT_ENDPOINT")
	ctx                = context.Background()
)

func TestClientListDeploymentRoleMapping(t *testing.T) {
	if apiKey == "" && organizationID == "" && deploymentApiKey == "" && deploymentEndpoint == "" {
		t.Skip()
	}

	cli := getClientForTesting(ctx)
	assert.Nil(t, cli)
	res, err := cli.ListDeploymentRoleMapping(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestClientGetDeploymentRoleMapping(t *testing.T) {
	if apiKey == "" && organizationID == "" && deploymentApiKey == "" && deploymentEndpoint == "" {
		t.Skip()
	}

	cli := getClientForTesting(ctx)
	assert.Nil(t, cli)
	res, err := cli.GetDeploymentRoleMapping(ctx, "mapping7")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func TestClientUpdateUsersWithinRoleMapping(t *testing.T) {
	if apiKey == "" && organizationID == "" && deploymentApiKey == "" && deploymentEndpoint == "" {
		t.Skip()
	}

	cli := getClientForTesting(ctx)
	assert.Nil(t, cli)
	body := elastic.MappingRolesBody{
		Roles: []string{"superuser",
			"viewer",
			"kibana_admin",
			"transform_admin",
			"inference_admin",
			"user",
			"admin",
		},
		Enabled: true,
		Rules: elastic.Rule{
			Field: elastic.Field{
				Username: []string{"jacknich", "bryancooper"},
			},
		},
	}
	err := cli.UpdateUserMappingRole(ctx, body, "mapping7")
	assert.Nil(t, err)
}

func TestClientAddDeploymentRole(t *testing.T) {
	if apiKey == "" && organizationID == "" && deploymentApiKey == "" && deploymentEndpoint == "" {
		t.Skip()
	}

	cli := getClientForTesting(ctx)
	assert.Nil(t, cli)
	body := elastic.RequestRoleBody{
		Cluster: []string{"all"},
		Indices: []elastic.Indices{
			{
				Names:      []string{"index1", "index2"},
				Privileges: []string{"all"},
				FieldSecurity: elastic.FieldSecurity{
					Grant: []string{"title", "body"},
				},
				Query: "{\"match\": {\"title\": \"foo\"}}",
			},
		},
		Applications: []elastic.Applications{
			{
				Application: "myapp",
				Privileges:  []string{"admin", "read"},
				Resources:   []string{"*"},
			},
		},
		RunAs: []string{"other_user"},
		Metadata: elastic.Metadata{
			Version: 1,
		},
	}
	err := cli.AddDeploymentRole(ctx, body, "my_admin_role")
	assert.Nil(t, err)
}

func TestClientAddDeploymentRoleV2(t *testing.T) {
	if apiKey == "" && organizationID == "" && deploymentApiKey == "" && deploymentEndpoint == "" {
		t.Skip()
	}

	cli := getClientForTesting(ctx)
	assert.Nil(t, cli)
	body := elastic.RequestRoleBody{
		RunAs:   []string{"clicks_watcher_1"},
		Cluster: []string{"monitor"},
		Indices: []elastic.Indices{
			{
				Names:      []string{"events-*"},
				Privileges: []string{"read"},
				FieldSecurity: elastic.FieldSecurity{
					Grant: []string{"category", "@timestamp", "message"},
				},
				Query: "{\"match\": {\"category\": \"click\"}}",
			},
		},
	}
	err := cli.AddDeploymentRole(ctx, body, "clicks_admin")
	assert.Nil(t, err)
}

func getClientForTesting(ctx context.Context) *elastic.Client {
	httpClient, err := uhttp.NewClient(ctx, uhttp.WithLogger(true, ctxzap.Extract(ctx)))
	if err != nil {
		return nil
	}

	return elastic.NewClient(
		httpClient,
		deploymentApiKey,
		deploymentEndpoint,
		apiKey,
		organizationID,
	)
}

func TestAddUsers(t *testing.T) {
	if apiKey == "" && organizationID == "" && deploymentApiKey == "" && deploymentEndpoint == "" {
		t.Skip()
	}

	cli := getClientForTesting(ctx)
	assert.Nil(t, cli)
	// It adds deployment roles when adding users
	body := elastic.UserBody{
		Password: "secretpwd",
		Roles: []string{"superuser",
			"viewer",
			"kibana_admin",
			"transform_admin",
			"inference_admin",
		},
		FullName: "Jack Nicholson",
		Email:    "jacknich@example.com",
		Metadata: elastic.UserMetadata{
			Intelligence: 7,
		},
	}
	err := cli.AddUsersWithRoles(ctx, body, "bryancooper")
	assert.Nil(t, err)
}

func TestListDeploymentUsers(t *testing.T) {
	if apiKey == "" && organizationID == "" && deploymentApiKey == "" && deploymentEndpoint == "" {
		t.Skip()
	}

	cli := getClientForTesting(ctx)
	assert.Nil(t, cli)
	res, _ := cli.ListDeploymentUsers(ctx)
	assert.Nil(t, res)
}

func TestDeleteDeploymentRoleMapping(t *testing.T) {
	if apiKey == "" && organizationID == "" && deploymentApiKey == "" && deploymentEndpoint == "" {
		t.Skip()
	}

	cli := getClientForTesting(ctx)
	assert.Nil(t, cli)

	err := cli.DeleteDeploymentRoleMapping(ctx, "mapping2")
	assert.Nil(t, err)
}
