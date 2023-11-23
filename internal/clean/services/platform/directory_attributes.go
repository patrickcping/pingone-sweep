package platform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/patrickcping/pingone-clean-config/internal/clean"
	"github.com/patrickcping/pingone-clean-config/internal/logger"
	"github.com/patrickcping/pingone-clean-config/internal/sdk"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
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

	debugModule := "Directory Attributes"

	l.Debug().Msgf(`[%s] Cleaning bootstrap config for environment ID "%s"..`, debugModule, c.Environment.EnvironmentID)

	if len(c.BootstrapAttributeNames) == 0 {
		l.Warn().Msgf("[%s] No bootstrap names configured - skipping", debugModule)
		return nil
	}

	var schemaName string
	if c.SchemaName != nil {
		schemaName = *c.SchemaName
	} else {
		schemaName = "User"
	}

	l.Debug().Msgf(`[%s] Fetching ID for schema "%s"..`, debugModule, schemaName)
	schema, err := fetchDefaultSchema(ctx, c.Environment.Client, c.Environment.EnvironmentID, schemaName)
	if err != nil {
		return err
	}

	if schema == nil {
		return fmt.Errorf("[%s] No schema found - the API responded with no data", debugModule)
	}

	l.Debug().Msgf(`[%s] Schema ID found as "%s"`, debugModule, schema.GetId())

	var response *management.EntityArray
	err = sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return c.Environment.Client.SchemasApi.ReadAllSchemaAttributes(ctx, c.Environment.EnvironmentID, schema.GetId()).Execute()
		},
		"ReadAllSchemaAttributes",
		sdk.DefaultCreateReadRetryable,
		&response,
	)

	if err != nil {
		return err
	}

	if response == nil {
		return fmt.Errorf("[%s] No configuration items found - the API responded with no data", debugModule)
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasAttributes() {

		l.Debug().Msgf("[%s] Configuration items found, looping..", debugModule)
		for _, attributeInstance := range embedded.GetAttributes() {

			attribute := attributeInstance.SchemaAttribute

			l.Debug().Msgf(`[%s] Looping names for "%s"..`, debugModule, attribute.GetName())
			for _, bootstrapAttributeName := range c.BootstrapAttributeNames {

				if attribute.GetName() == bootstrapAttributeName {
					l.Debug().Msgf(`[%s] Found "%s"`, debugModule, bootstrapAttributeName)

					if !attribute.GetEnabled() {
						l.Info().Msgf(`[%s] "%s" is already disabled - no action taken`, debugModule, bootstrapAttributeName)

						break
					}

					attributeUpdate := management.NewSchemaAttributePatch()
					attributeUpdate.SetEnabled(false)

					if !c.Environment.DryRun {
						err := sdk.ParseResponse(
							ctx,

							func() (any, *http.Response, error) {
								return c.Environment.Client.SchemasApi.UpdateAttributePatch(ctx, c.Environment.EnvironmentID, schema.GetId(), attribute.GetId()).SchemaAttributePatch(*attributeUpdate).Execute()
							},
							"UpdateAttributePatch",
							sdk.DefaultCreateReadRetryable,
							nil,
						)

						if err != nil {
							return err
						}
						l.Info().Msgf(`[%s] "%s" disabled`, debugModule, attribute.GetName())
					} else {
						l.Warn().Msgf(`[%s] Dry run: "%s" with ID "%s" would be disabled`, debugModule, attribute.GetName(), attribute.GetId())
					}

					break
				}
			}
		}
		l.Debug().Msgf("[%s] Done", debugModule)

	} else {
		l.Debug().Msgf("[%s] No configuration items found in the target environment", debugModule)
	}

	return nil
}

func fetchDefaultSchema(ctx context.Context, apiClient *management.APIClient, environmentId string, schemaName string) (*management.Schema, error) {
	var schema management.Schema

	// Run the API call
	var entityArray *management.EntityArray
	err := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return apiClient.SchemasApi.ReadAllSchemas(ctx, environmentId).Execute()
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
			return nil, fmt.Errorf("Cannot find schema from name - The schema %s for environment %s cannot be found", schemaName, environmentId)
		}

	}

	return &schema, nil
}
