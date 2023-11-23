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

	l.Debug().Msgf(`Cleaning bootstrap branding themes for environment ID "%s"..`, c.Environment.EnvironmentID)

	if len(c.BootstrapBrandingThemeNames) == 0 {
		l.Warn().Msgf("No bootstrap branding theme names configured - skipping")
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
		return fmt.Errorf("No branding themes found - the API responded with no data")
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasThemes() {

		l.Debug().Msg("Branding themes found, looping..")
		for _, theme := range embedded.GetThemes() {

			l.Debug().Msgf(`Looping bootstrapped theme names for "%s"..`, theme.Configuration.GetName())
			for _, defaultThemeName := range c.BootstrapBrandingThemeNames {

				if theme.Configuration.GetName() == defaultThemeName {
					l.Debug().Msgf(`Found "%s" branding theme`, defaultThemeName)

					if theme.GetDefault() {
						l.Warn().Msgf(`The "%s" branding theme is set as the default theme - this will not be deleted`, defaultThemeName)

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
						l.Info().Msgf(`Branding theme "%s" deleted`, defaultThemeName)
					} else {
						l.Warn().Msgf(`Dry run: branding theme "%s" with ID "%s" would be deleted`, defaultThemeName, theme.GetId())
					}

					break
				}
			}
		}
		l.Debug().Msg("Branding themes done")

	} else {
		l.Debug().Msg("No Branding themes found in the target environment")
	}

	return nil
}
