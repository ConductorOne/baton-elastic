package elastic

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
