package cmd

import (
	"fmt"
	"strings"

	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/clean/services/platform"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	themeNames []string
)

const (
	brandingThemesCmdName = "branding-themes"

	brandingThemeNamesParamName      = "theme-name"
	brandingThemeNamesParamConfigKey = "pingone.services.platform.branding-themes.names"
)

var (
	brandingThemesConfigurationParamMapping = map[string]string{
		brandingThemeNamesParamName: brandingThemeNamesParamConfigKey,
	}
)

var cleanBrandingThemesCmd = &cobra.Command{
	Use:   brandingThemesCmdName,
	Short: "Clean unwanted demo branding theme configuration",
	Long: fmt.Sprintf(`Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s "Ping Default" --%s "Ping Default 2" --%s
	
	`, brandingThemesCmdName, environmentIDParamName, dryRunParamName, brandingThemesCmdName, environmentIDParamName, brandingThemeNamesParamName, brandingThemeNamesParamName, dryRunParamName),
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool(dryRunParamConfigKey)
		themeNames := viper.GetStringSlice(brandingThemeNamesParamConfigKey)

		l.Debug().Msgf("Clean Command called for branding themes.")
		l.Debug().Msgf("Dry run setting: %t", dryRun)
		l.Debug().Msgf(`Theme names: "%s"`, strings.Join(themeNames, `", "`))

		var err error
		apiClient, err = initApiClient(cmd.Context(), cmd.Version)
		if err != nil {
			return err
		}

		cleanConfig := platform.CleanEnvironmentPlatformBrandingThemesConfig{
			Environment: clean.CleanEnvironmentConfig{
				EnvironmentID: viper.GetString(environmentIDParamConfigKey),
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapBrandingThemeNames: themeNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanBrandingThemesCmd.PersistentFlags().StringArrayVar(&themeNames, brandingThemeNamesParamName, platform.BootstrapBrandingThemeNames, "The list of theme names to search for to delete.")

	// Do the binds
	for k, v := range brandingThemesConfigurationParamMapping {
		viper.BindPFlag(v, cleanBrandingThemesCmd.PersistentFlags().Lookup(k))
	}
}
