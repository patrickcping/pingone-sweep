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

	l.Debug().Msgf(`Cleaning bootstrap DaVinci forms for environment ID "%s"..`, c.Environment.EnvironmentID)

	if len(c.BootstrapDaVinciFormNames) == 0 {
		l.Warn().Msgf("No bootstrap DaVinci form names configured - skipping")
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
		return fmt.Errorf("No DaVinci forms found - the API responded with no data")
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasForms() {

		l.Debug().Msg("DaVinci forms found, looping..")
		for _, form := range embedded.GetForms() {

			l.Debug().Msgf(`Looping bootstrapped form names for "%s"..`, form.GetName())
			for _, defaultFormName := range c.BootstrapDaVinciFormNames {

				if form.GetName() == defaultFormName {
					l.Debug().Msgf(`Found "%s" DaVinci form`, defaultFormName)

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
						l.Info().Msgf(`DaVinci form "%s" deleted`, defaultFormName)
					} else {
						l.Warn().Msgf(`Dry run: DaVinci form "%s" with ID "%s" would be deleted`, defaultFormName, form.GetId())
					}

					break
				}
			}
		}
		l.Debug().Msg("DaVinci forms done")

	} else {
		l.Debug().Msg("No DaVinci forms found in the target environment")
	}

	return nil
}
