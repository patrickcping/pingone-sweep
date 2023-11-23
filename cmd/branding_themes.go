package cmd

import (
	"strings"

	"github.com/patrickcping/pingone-clean-config/internal/clean"
	"github.com/patrickcping/pingone-clean-config/internal/clean/services/platform"
	"github.com/patrickcping/pingone-clean-config/internal/logger"
	"github.com/spf13/cobra"
)

var (
	themeNames []string
)

var cleanBrandingThemesCmd = &cobra.Command{
	Use:   "branding-themes",
	Short: "Clean unwanted demo branding theme configuration",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-cleanconfig branding-themes --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-cleanconfig branding-themes --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --theme-name "Ping Default" --theme-name "Ping Default 2" --dry-run
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()
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
				Client:        apiClient.API.ManagementAPIClient,
				EnvironmentID: environmentID,
				DryRun:        dryRun,
			},
			BootstrapBrandingThemeNames: themeNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanBrandingThemesCmd.PersistentFlags().StringArrayVar(&themeNames, "theme-name", platform.BootstrapBrandingThemeNames, "The list of theme names to search for to delete.")
}
