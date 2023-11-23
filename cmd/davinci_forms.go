package cmd

import (
	"strings"

	"github.com/patrickcping/pingone-clean-config/internal/clean"
	"github.com/patrickcping/pingone-clean-config/internal/clean/services/platform"
	"github.com/patrickcping/pingone-clean-config/internal/logger"
	"github.com/spf13/cobra"
)

var (
	daVinciFormNames []string
)

var cleanDaVinciFormsCmd = &cobra.Command{
	Use:   "davinci-forms",
	Short: "Clean unwanted demo DaVinci form configuration",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-cleanconfig davinci-forms --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-cleanconfig davinci-forms --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --form-name "Default DaVinci Form" --form-name "Default DaVinci Form 2" --dry-run
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()
		l.Debug().Msgf("Clean Command called for DaVinci forms.")
		l.Debug().Msgf("Dry run setting: %t", dryRun)
		l.Debug().Msgf(`DaVinci Form names: "%s"`, strings.Join(daVinciFormNames, `", "`))

		var err error
		apiClient, err = initApiClient(cmd.Context(), cmd.Version)
		if err != nil {
			return err
		}

		cleanConfig := platform.CleanEnvironmentPlatformDaVinciFormsConfig{
			Environment: clean.CleanEnvironmentConfig{
				Client:        apiClient.API.ManagementAPIClient,
				EnvironmentID: environmentID,
				DryRun:        dryRun,
			},
			BootstrapDaVinciFormNames: daVinciFormNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanDaVinciFormsCmd.PersistentFlags().StringSliceVar(&daVinciFormNames, "form-name", platform.BootstrapDaVinciFormNames, "The list of DaVinci form names to search for to delete.  Case sensitive.")
}
