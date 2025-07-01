package elastic

import (
	"encoding/json"
	"fmt"
)

type DeploymentUser struct {
	Username string      `json:"username"`
	Roles    []string    `json:"roles"`
	FullName string      `json:"full_name"`
	Email    string      `json:"email"`
	Enabled  bool        `json:"enabled"`
	Metadata interface{} `json:"metadata"`
}

type DeploymentRole struct {
	Cluster      []string      `json:"cluster"`
	Applications []interface{} `json:"applications"`
	RunAs        []string      `json:"run_as"`
}

type User struct {
	Email          string `json:"email"`
	MemberSince    string `json:"member_since"`
	Name           string `json:"name"`
	OrganizationID string `json:"organization_id"`
	UserID         string `json:"user_id"`
}

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type MappingRolesResponse struct {
	Roles        []string `json:"roles,omitempty"`
	Enabled      bool     `json:"enabled,omitempty"`
	Rules        Rule     `json:"rules,omitempty"`
	RuleTemplate any      `json:"role_templates,omitempty"`
	Metadata     any      `json:"metadata,omitempty"`
}

type Rules struct {
	Field map[string][]any `json:"field,omitempty"`
	All   any              `json:"all,omitempty"`
}

type Rule struct {
	Field Field `json:"field,omitempty"`
}

type StringOrSlice []string

func (s *StringOrSlice) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		*s = []string{single}
		return nil
	}

	var slice []string
	if err := json.Unmarshal(data, &slice); err == nil {
		*s = slice
		return nil
	}

	return fmt.Errorf("StringOrSlice: unable to unmarshal %s", string(data))
}

type MappingRolesBody struct {
	Roles   []string `json:"roles,omitempty"`
	Enabled bool     `json:"enabled,omitempty"`
	Rules   Rule     `json:"rules,omitempty"`
}

type Field struct {
	Username StringOrSlice `json:"username,omitempty"`
}

type RequestRoleBody struct {
	Cluster      []string       `json:"cluster,omitempty"`
	Indices      []Indices      `json:"indices,omitempty"`
	Applications []Applications `json:"applications,omitempty"`
	RunAs        []string       `json:"run_as,omitempty"`
	Metadata     Metadata       `json:"metadata,omitempty"`
}

type Metadata struct {
	Version int `json:"version,omitempty"`
}

type Applications struct {
	Application string   `json:"application,omitempty"`
	Privileges  []string `json:"privileges,omitempty"`
	Resources   []string `json:"resources,omitempty"`
}

type FieldSecurity struct {
	Grant []string `json:"grant,omitempty"`
}

type Indices struct {
	Names         []string      `json:"names,omitempty"`
	Privileges    []string      `json:"privileges,omitempty"`
	FieldSecurity FieldSecurity `json:"field_security,omitempty"`
	Query         string        `json:"query,omitempty"`
}

type UserBody struct {
	Password string       `json:"password,omitempty"`
	Roles    []string     `json:"roles,omitempty"`
	FullName string       `json:"full_name,omitempty"`
	Email    string       `json:"email,omitempty"`
	Metadata UserMetadata `json:"metadata,omitempty"`
}

type UserMetadata struct {
	Intelligence int `json:"intelligence,omitempty"`
}
