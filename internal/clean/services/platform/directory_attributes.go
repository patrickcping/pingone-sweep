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

	l.Debug().Msgf(`Cleaning bootstrap directory attributes for environment ID "%s"..`, c.Environment.EnvironmentID)

	if len(c.BootstrapAttributeNames) == 0 {
		l.Warn().Msgf("No bootstrap directory attribute names configured - skipping")
		return nil
	}

	var schemaName string
	if c.SchemaName != nil {
		schemaName = *c.SchemaName
	} else {
		schemaName = "User"
	}

	l.Debug().Msgf(`Fetching ID for schema "%s"..`, schemaName)
	schema, err := fetchDefaultSchema(ctx, c.Environment.Client, c.Environment.EnvironmentID, schemaName)
	if err != nil {
		return err
	}

	if schema == nil {
		return fmt.Errorf("No schema found - the API responded with no data")
	}

	l.Debug().Msgf(`Schema ID found as "%s"`, schema.GetId())

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
		return fmt.Errorf("No directory attributes found - the API responded with no data")
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasAttributes() {

		l.Debug().Msg("Directory attributes found, looping..")
		for _, attributeInstance := range embedded.GetAttributes() {

			attribute := attributeInstance.SchemaAttribute

			l.Debug().Msgf(`Looping bootstrapped directory attribute names for "%s"..`, attribute.GetName())
			for _, bootstrapAttributeName := range c.BootstrapAttributeNames {

				if attribute.GetName() == bootstrapAttributeName {
					l.Debug().Msgf(`Found "%s" directory attribute`, bootstrapAttributeName)

					if attribute.GetEnabled() {
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
							l.Info().Msgf(`Directory attribute "%s" disabled`, bootstrapAttributeName)
						} else {
							l.Warn().Msgf(`Dry run: directory attribute "%s" with ID "%s" would be disabled`, bootstrapAttributeName, attribute.GetId())
						}
					} else {
						l.Info().Msgf(`Directory attribute "%s" already disabled - no action taken`, bootstrapAttributeName)
					}

					break
				}
			}
		}

	} else {
		l.Debug().Msg("No directory attributes found in the target environment")
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
