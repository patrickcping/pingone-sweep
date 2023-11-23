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
	BootstrapDaVinciFormNames = []string{
		"Example - Password Recovery",
		"Example - Password Recovery User Lookup",
		"Example - Password Reset",
		"Example - Registration",
		"Example - Sign On",
	}
)

type CleanEnvironmentPlatformDaVinciFormsConfig struct {
	Environment               clean.CleanEnvironmentConfig
	BootstrapDaVinciFormNames []string
}

func (c *CleanEnvironmentPlatformDaVinciFormsConfig) Clean(ctx context.Context) error {
	l := logger.Get()

	debugModule := "DaVinci Forms"

	l.Debug().Msgf(`[%s] Cleaning bootstrap config for environment ID "%s"..`, debugModule, c.Environment.EnvironmentID)

	if len(c.BootstrapDaVinciFormNames) == 0 {
		l.Warn().Msgf("[%s] No bootstrap names configured - skipping", debugModule)
		return nil
	}

	var response *management.EntityArray
	err := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return c.Environment.Client.FormManagementApi.ReadAllForms(ctx, c.Environment.EnvironmentID).Execute()
		},
		"ReadAllForms",
		sdk.DefaultCreateReadRetryable,
		&response,
	)

	if err != nil {
		return err
	}

	if response == nil {
		return fmt.Errorf("[%s] No configuration items found - the API responded with no data", debugModule)
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasForms() {

		l.Debug().Msgf("[%s] Configuration items found, looping..", debugModule)
		for _, form := range embedded.GetForms() {

			l.Debug().Msgf(`[%s] Looping names for "%s"..`, debugModule, form.GetName())
			for _, defaultFormName := range c.BootstrapDaVinciFormNames {

				if form.GetName() == defaultFormName {
					l.Debug().Msgf(`[%s] Found "%s"`, debugModule, defaultFormName)

					if !c.Environment.DryRun {
						err := sdk.ParseResponse(
							ctx,

							func() (any, *http.Response, error) {
								r, err := c.Environment.Client.FormManagementApi.DeleteForm(ctx, c.Environment.EnvironmentID, form.GetId()).Execute()
								return nil, r, err
							},
							"DeleteForm",
							sdk.DefaultCreateReadRetryable,
							nil,
						)

						if err != nil {
							return err
						}
						l.Info().Msgf(`[%s] "%s" deleted`, debugModule, form.GetName())
					} else {
						l.Warn().Msgf(`[%s] Dry run: "%s" with ID "%s" would be deleted`, debugModule, form.GetName(), form.GetId())
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
