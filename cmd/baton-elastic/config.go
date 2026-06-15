package main

import (
	"github.com/conductorone/baton-sdk/pkg/field"
)

var (
	APIKeyField = field.StringField(
		"api-key",
		field.WithDescription("Elastic Cloud API key used to communicate with Elastic Cloud API. Can be set via $BATON_API_KEY."),
		field.WithIsSecret(true),
	)

	OrganizationIDField = field.StringField(
		"organization-id",
		field.WithDescription("Optional Elastic Cloud organization ID to sync members of a single org."),
	)

	DeploymentAPIKeyField = field.StringField(
		"deployment-api-key",
		field.WithDescription("API key of your Elasticsearch deployment."),
		field.WithIsSecret(true),
	)

	DeploymentEndpointField = field.StringField(
		"deployment-endpoint",
		field.WithDescription("Elasticsearch deployment endpoint URL, e.g. http://localhost:9200."),
	)

	fieldRelationships = []field.SchemaFieldRelationship{
		field.FieldsAtLeastOneUsed(
			APIKeyField,
			DeploymentAPIKeyField,
		),
		field.FieldsRequiredTogether(
			DeploymentAPIKeyField,
			DeploymentEndpointField,
		),
	}

	ConfigurationSchema = field.NewConfiguration(
		[]field.SchemaField{
			APIKeyField,
			OrganizationIDField,
			DeploymentAPIKeyField,
			DeploymentEndpointField,
		},
		field.WithConstraints(fieldRelationships...),
		field.WithConnectorDisplayName("Elastic"),
		field.WithHelpUrl("/docs/baton/elastic"),
		field.WithIconUrl("/static/app-icons/elastic.svg"),
	)
)
