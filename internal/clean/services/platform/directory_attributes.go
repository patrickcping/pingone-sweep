package platform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/patrickcping/pingone-sweep/internal/sdk"
)

var (
	BootstrapDirectoryAttributeNames = []string{
		"accountId",
		"address",
		"email",
		"externalId",
		"locale",
		"mobilePhone",
		"name",
		"nickname",
		"photo",
		"preferredLanguage",
		"primaryPhone",
		"timezone",
		"title",
		"type",
	}
)

type CleanEnvironmentPlatformDirectoryAttributeConfig struct {
	Environment             clean.CleanEnvironmentConfig
	BootstrapAttributeNames []string
	SchemaName              *string
}

func (c *CleanEnvironmentPlatformDirectoryAttributeConfig) Clean(ctx context.Context) error {
	l := logger.Get()

	configKey := "Directory Attributes"

	l.Debug().Msgf(`[%s] Cleaning bootstrap config for environment ID "%s"..`, configKey, c.Environment.EnvironmentID)

	if len(c.BootstrapAttributeNames) == 0 {
		l.Warn().Msgf("[%s] No bootstrap names configured - skipping", configKey)
		return nil
	}

	var schemaName string
	if c.SchemaName != nil {
		schemaName = *c.SchemaName
	} else {
		schemaName = "User"
	}

	l.Debug().Msgf(`[%s] Fetching ID for schema "%s"..`, configKey, schemaName)
	schema, err := fetchDefaultSchema(ctx, c.Environment, schemaName)
	if err != nil {
		return err
	}

	if schema == nil {
		return fmt.Errorf("[%s] No schema found - the API responded with no data", configKey)
	}

	l.Debug().Msgf(`[%s] Schema ID found as "%s"`, configKey, schema.GetId())

	var response *management.EntityArray
	err = clean.ReadAllConfig(
		ctx,
		configKey,
		c.Environment,
		func() (any, *http.Response, error) {
			return c.Environment.Client.ManagementAPIClient.SchemasApi.ReadAllSchemaAttributes(ctx, c.Environment.EnvironmentID, schema.GetId()).Execute()
		},
		&response,
	)
	if err != nil {
		return err
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasAttributes() {

		l.Debug().Msgf("[%s] Configuration items found, looping..", configKey)
		for _, attributeInstance := range embedded.GetAttributes() {

			attribute := attributeInstance.SchemaAttribute

			err := clean.TryCleanConfig(
				ctx,
				configKey,
				c.Environment,
				clean.ConfigItem{
					IdentifierToEvaluate: attribute.Name,
					Id:                   *attribute.Id,
					Enabled:              &attribute.Enabled,
				},
				clean.ConfigItemEval{
					IdentifierListToSearch: c.BootstrapAttributeNames,
					StartsWithStringMatch:  false,
				},
				nil,
				func() (any, *http.Response, error) {
					attributeUpdate := management.NewSchemaAttributePatch()
					attributeUpdate.SetEnabled(false)
					return c.Environment.Client.ManagementAPIClient.SchemasApi.UpdateAttributePatch(ctx, c.Environment.EnvironmentID, schema.GetId(), attribute.GetId()).SchemaAttributePatch(*attributeUpdate).Execute()
				},
			)

			if err != nil {
				return err
			}

		}
		l.Debug().Msgf("[%s] Done", configKey)

	} else {
		l.Debug().Msgf("[%s] No configuration items found in the target environment", configKey)
	}

	return nil
}

func fetchDefaultSchema(ctx context.Context, env clean.CleanEnvironmentConfig, schemaName string) (*management.Schema, error) {
	var schema management.Schema

	// Run the API call
	var entityArray *management.EntityArray
	err := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return env.Client.ManagementAPIClient.SchemasApi.ReadAllSchemas(ctx, env.EnvironmentID).Execute()
		},
		"ReadAllSchemas",
		sdk.DefaultCreateReadRetryable,
		&entityArray,
	)
	if err != nil {
		return nil, err
	}

	if schemas, ok := entityArray.Embedded.GetSchemasOk(); ok {

		found := false
		for _, schemaItem := range schemas {

			if schemaItem.GetName() == schemaName {
				schema = schemaItem
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("Cannot find schema from name - The schema %s for environment %s cannot be found", schemaName, env.EnvironmentID)
		}

	}

	return &schema, nil
}
