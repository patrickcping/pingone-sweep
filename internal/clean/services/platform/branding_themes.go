package platform

import (
	"context"
	"net/http"

	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/logger"
)

var (
	BootstrapBrandingThemeNames = []string{
		"Ping Default",
	}
)

type CleanEnvironmentPlatformBrandingThemesConfig struct {
	Environment                 clean.CleanEnvironmentConfig
	BootstrapBrandingThemeNames []string
}

func (c *CleanEnvironmentPlatformBrandingThemesConfig) Clean(ctx context.Context) error {
	l := logger.Get()

	configKey := "Branding Themes"

	l.Debug().Msgf(`[%s] Cleaning bootstrap config for environment ID "%s"..`, configKey, c.Environment.EnvironmentID)

	if len(c.BootstrapBrandingThemeNames) == 0 {
		l.Info().Msgf("[%s] No bootstrap names configured - skipping", configKey)
		return nil
	}

	var response *management.EntityArray
	err := clean.ReadAllConfig(
		ctx,
		configKey,
		c.Environment,
		func() (any, *http.Response, error) {
			return c.Environment.Client.ManagementAPIClient.BrandingThemesApi.ReadBrandingThemes(ctx, c.Environment.EnvironmentID).Execute()
		},
		&response,
	)
	if err != nil {
		return err
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasThemes() {

		l.Debug().Msgf("[%s] Configuration items found, looping..", configKey)
		for _, theme := range embedded.GetThemes() {

			err := clean.TryCleanConfig(
				ctx,
				configKey,
				c.Environment,
				clean.ConfigItem{
					IdentifierToEvaluate: *theme.Configuration.Name,
					Id:                   *theme.Id,
					Default:              &theme.Default,
				},
				clean.ConfigItemEval{
					IdentifierListToSearch: c.BootstrapBrandingThemeNames,
					StartsWithStringMatch:  false,
				},
				func() (any, *http.Response, error) {
					fR, fErr := c.Environment.Client.ManagementAPIClient.BrandingThemesApi.DeleteBrandingTheme(ctx, c.Environment.EnvironmentID, theme.GetId()).Execute()
					return nil, fR, fErr
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
