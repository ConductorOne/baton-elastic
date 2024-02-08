package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const baseUrl = "https://api.elastic-cloud.com/"

type Client struct {
	httpClient         *http.Client
	apiKey             string
	organizationID     string
	deploymentApiKey   string
	deploymentEndpoint string
}

func NewClient(httpClient *http.Client, deploymentApiKey, deploymentEndpoint, apiKey, organizationID string) *Client {
	return &Client{
		httpClient:         httpClient,
		apiKey:             apiKey,
		organizationID:     organizationID,
		deploymentApiKey:   deploymentApiKey,
		deploymentEndpoint: deploymentEndpoint,
	}
}

// ListOrganizations returns a list of all Elastic organizations.
func (c *Client) ListOrganizations(ctx context.Context) ([]Organization, error) {
	var res struct {
		Organizations []Organization `json:"organizations"`
	}

	orgUrl, _ := url.JoinPath(baseUrl, "api/v1/organizations")
	if err := c.doRequest(ctx, orgUrl, &res, http.MethodGet, nil); err != nil {
		return nil, err
	}

	return res.Organizations, nil
}

// ListOrgMembers returns a list of all Elastic organization members.
func (c *Client) ListOrgMembers(ctx context.Context, orgId string) ([]User, error) {
	var res struct {
		Members []User `json:"members"`
	}

	if c.organizationID != "" {
		orgId = c.organizationID
	}

	orgUrl, _ := url.JoinPath(baseUrl, "api/v1/organizations", orgId, "members")
	if err := c.doRequest(ctx, orgUrl, &res, http.MethodGet, nil); err != nil {
		return nil, err
	}

	return res.Members, nil
}

// ListDeploymentUsers returns a list of all Elastic deployment users.
func (c *Client) ListDeploymentUsers(ctx context.Context) (map[string]DeploymentUser, error) {
	res := make(map[string]DeploymentUser)

	usersUrl, _ := url.JoinPath(c.deploymentEndpoint, "_security/user")
	if err := c.doRequest(ctx, usersUrl, &res, http.MethodGet, nil); err != nil {
		return nil, err
	}

	return res, nil
}

// GetDeploymentUser returns a single user from Elastic deployment.
func (c *Client) GetDeploymentUser(ctx context.Context, username string) (map[string]DeploymentUser, error) {
	res := make(map[string]DeploymentUser)
	usersUrl, _ := url.JoinPath(c.deploymentEndpoint, "_security/user", username)
	if err := c.doRequest(ctx, usersUrl, &res, http.MethodGet, nil); err != nil {
		return nil, err
	}

	return res, nil
}

// ListDeploymentRoles returns a list of all Elastic roles on deployment.
func (c *Client) ListDeploymentRoles(ctx context.Context) (map[string]DeploymentRole, error) {
	res := make(map[string]DeploymentRole)
	usersUrl, _ := url.JoinPath(c.deploymentEndpoint, "_security/role")
	if err := c.doRequest(ctx, usersUrl, &res, http.MethodGet, nil); err != nil {
		return nil, err
	}

	return res, nil
}

// DeploymentAuth returns info about user that is authenticated with the deployment api key.
func (c *Client) DeploymentAuth(ctx context.Context) error {
	var res struct {
		Username string `json:"username"`
		Enabled  bool   `json:"enabled"`
	}

	authUrl, _ := url.JoinPath(c.deploymentEndpoint, "_security/_authenticate")
	if err := c.doRequest(ctx, authUrl, &res, http.MethodGet, nil); err != nil {
		return err
	}

	if res.Username == "" && !res.Enabled {
		return fmt.Errorf("invalid deployment api key")
	}

	return nil
}

// UpdateUser updates user. Used to grant or revoke user roles.
func (c *Client) UpdateUser(ctx context.Context, username string, user DeploymentUser) error {
	url, _ := url.JoinPath(c.deploymentEndpoint, "_security/user/", username)
	requestBody, err := json.Marshal(user)
	if err != nil {
		return err
	}

	var res struct {
		Created bool `json:"created"`
	}

	e := c.doRequest(ctx, url, &res, http.MethodPost, requestBody)
	if e != nil {
		return fmt.Errorf("error updating user: %w", e)
	}

	return nil
}

func (c *Client) doRequest(ctx context.Context, url string, res interface{}, method string, payload []byte) error {
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if c.deploymentEndpoint != "" && strings.Contains(url, c.deploymentEndpoint) {
		req.Header.Add("Authorization", fmt.Sprintf("ApiKey %s", c.deploymentApiKey))
	} else {
		req.Header.Add("Authorization", fmt.Sprintf("ApiKey %s", c.apiKey))
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}

	return nil
}
