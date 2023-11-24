package davinci

import (
	"context"
	"net/http"

	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/logger"
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

type CleanEnvironmentDaVinciFormsConfig struct {
	Environment               clean.CleanEnvironmentConfig
	BootstrapDaVinciFormNames []string
}

func (c *CleanEnvironmentDaVinciFormsConfig) Clean(ctx context.Context) error {
	l := logger.Get()

	configKey := "DaVinci Forms"

	l.Debug().Msgf(`[%s] Cleaning bootstrap config for environment ID "%s"..`, configKey, c.Environment.EnvironmentID)

	if len(c.BootstrapDaVinciFormNames) == 0 {
		l.Info().Msgf("[%s] No bootstrap names configured - skipping", configKey)
		return nil
	}

	ok, err := clean.BillOfMaterialsHasService(ctx, configKey, c.Environment, management.ENUMPRODUCTTYPE_ONE_DAVINCI)
	if err != nil {
		return err
	}

	if !ok {
		l.Info().Msgf("[%s] Bill of materials does not contain applicable service %s - skipping", configKey, management.ENUMPRODUCTTYPE_ONE_DAVINCI)
		return nil
	}

	var response *management.EntityArray
	err = clean.ReadAllConfig(
		ctx,
		configKey,
		c.Environment,
		func() (any, *http.Response, error) {
			return c.Environment.Client.ManagementAPIClient.FormManagementApi.ReadAllForms(ctx, c.Environment.EnvironmentID).Execute()
		},
		&response,
	)
	if err != nil {
		return err
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasForms() {

		l.Debug().Msgf("[%s] Configuration items found, looping..", configKey)
		for _, form := range embedded.GetForms() {

			err := clean.TryCleanConfig(
				ctx,
				configKey,
				c.Environment,
				clean.ConfigItem{
					IdentifierToEvaluate: form.Name,
					Id:                   *form.Id,
				},
				clean.ConfigItemEval{
					IdentifierListToSearch: c.BootstrapDaVinciFormNames,
					StartsWithStringMatch:  false,
				},
				func() (any, *http.Response, error) {
					r, err := c.Environment.Client.ManagementAPIClient.FormManagementApi.DeleteForm(ctx, c.Environment.EnvironmentID, form.GetId()).Execute()
					return nil, r, err
				},
				nil,
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
