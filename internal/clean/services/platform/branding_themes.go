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

	debugModule := "Branding Themes"

	l.Debug().Msgf(`[%s] Cleaning bootstrap config for environment ID "%s"..`, debugModule, c.Environment.EnvironmentID)

	if len(c.BootstrapBrandingThemeNames) == 0 {
		l.Warn().Msgf("[%s] No bootstrap names configured - skipping", debugModule)
		return nil
	}

	var response *management.EntityArray
	err := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return c.Environment.Client.BrandingThemesApi.ReadBrandingThemes(ctx, c.Environment.EnvironmentID).Execute()
		},
		"ReadBrandingThemes",
		sdk.DefaultCreateReadRetryable,
		&response,
	)

	if err != nil {
		return err
	}

	if response == nil {
		return fmt.Errorf("[%s] No configuration items found - the API responded with no data", debugModule)
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasThemes() {

		l.Debug().Msgf("[%s] Configuration items found, looping..", debugModule)
		for _, theme := range embedded.GetThemes() {

			l.Debug().Msgf(`[%s] Looping names for "%s"..`, debugModule, theme.Configuration.GetName())
			for _, defaultThemeName := range c.BootstrapBrandingThemeNames {

				if theme.Configuration.GetName() == defaultThemeName {
					l.Debug().Msgf(`[%s] Found "%s"`, debugModule, defaultThemeName)

					if theme.GetDefault() {
						l.Warn().Msgf(`[%s] "%s" is set as the environment default - this configuration will not be deleted`, debugModule, theme.Configuration.GetName())

						break
					}

					if !c.Environment.DryRun {
						err := sdk.ParseResponse(
							ctx,

							func() (any, *http.Response, error) {
								r, err := c.Environment.Client.BrandingThemesApi.DeleteBrandingTheme(ctx, c.Environment.EnvironmentID, theme.GetId()).Execute()
								return nil, r, err
							},
							"DeleteBrandingTheme",
							sdk.DefaultCreateReadRetryable,
							nil,
						)

						if err != nil {
							return err
						}
						l.Info().Msgf(`[%s] "%s" deleted`, debugModule, theme.Configuration.GetName())
					} else {
						l.Warn().Msgf(`[%s] Dry run: "%s" with ID "%s" would be deleted`, debugModule, theme.Configuration.GetName(), theme.GetId())
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
