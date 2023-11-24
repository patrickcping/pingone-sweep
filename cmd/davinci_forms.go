package cmd

import (
	"strings"

	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/clean/services/davinci"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	daVinciFormNames []string
)

var cleanDaVinciFormsCmd = &cobra.Command{
	Use:   "davinci-forms",
	Short: "Clean unwanted demo DaVinci form configuration",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep davinci-forms --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-sweep davinci-forms --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --form-name "Default DaVinci Form" --form-name "Default DaVinci Form 2" --dry-run
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool("dry-run")
		daVinciFormNames := viper.GetStringSlice("pingone.services.davinci.forms.names")

		l.Debug().Msgf("Clean Command called for DaVinci forms.")
		l.Debug().Msgf("Dry run setting: %t", dryRun)
		l.Debug().Msgf(`DaVinci Form names: "%s"`, strings.Join(daVinciFormNames, `", "`))

		var err error
		apiClient, err = initApiClient(cmd.Context(), cmd.Version)
		if err != nil {
			return err
		}

		cleanConfig := davinci.CleanEnvironmentDaVinciFormsConfig{
			Environment: clean.CleanEnvironmentConfig{
				EnvironmentID: viper.GetString("pingone.target-environment-id"),
				DryRun:        viper.GetBool("dry-run"),
				Client:        apiClient.API,
			},
			BootstrapDaVinciFormNames: daVinciFormNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanDaVinciFormsCmd.PersistentFlags().StringSliceVar(&daVinciFormNames, "form-name", davinci.BootstrapDaVinciFormNames, "The list of DaVinci form names to search for to delete.  Case sensitive.")
	viper.BindPFlag("pingone.services.davinci.forms.names", cleanDaVinciFormsCmd.PersistentFlags().Lookup("form-name"))
}
