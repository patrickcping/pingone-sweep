package cmd

import (
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

var cleanBrandingThemesCmd = &cobra.Command{
	Use:   "branding-themes",
	Short: "Clean unwanted demo branding theme configuration",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep branding-themes --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-sweep branding-themes --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --theme-name "Ping Default" --theme-name "Ping Default 2" --dry-run
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool("dry-run")
		themeNames := viper.GetStringSlice("pingone.services.platform.branding-themes.names")

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
				EnvironmentID: viper.GetString("pingone.target-environment-id"),
				DryRun:        viper.GetBool("dry-run"),
				Client:        apiClient.API,
			},
			BootstrapBrandingThemeNames: themeNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanBrandingThemesCmd.PersistentFlags().StringArrayVar(&themeNames, "theme-name", platform.BootstrapBrandingThemeNames, "The list of theme names to search for to delete.")
	viper.BindPFlag("pingone.services.platform.branding-themes.names", cleanBrandingThemesCmd.PersistentFlags().Lookup("theme-name"))
}
