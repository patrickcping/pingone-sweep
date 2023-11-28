package cmd

import (
	"fmt"
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

const (
	davinciFormsCmdName = "davinci-forms"

	davinciFormNamesParamName      = "form-name"
	davinciFormNamesParamConfigKey = "pingone.services.davinci.forms.names"
)

var (
	davinciFormsConfigurationParamMapping = map[string]string{
		davinciFormNamesParamName: davinciFormNamesParamConfigKey,
	}
)

var cleanDaVinciFormsCmd = &cobra.Command{
	Use:   davinciFormsCmdName,
	Short: "Clean unwanted demo DaVinci form configuration",
	Long: fmt.Sprintf(`Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s "Default DaVinci Form" --%s "Default DaVinci Form 2" --%s
	
	`, davinciFormsCmdName, environmentIDParamName, dryRunParamName, davinciFormsCmdName, environmentIDParamName, davinciFormNamesParamName, davinciFormNamesParamName, dryRunParamName),
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool(dryRunParamConfigKey)
		daVinciFormNames := viper.GetStringSlice(davinciFormNamesParamConfigKey)

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
				EnvironmentID: viper.GetString(environmentIDParamConfigKey),
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapDaVinciFormNames: daVinciFormNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanDaVinciFormsCmd.PersistentFlags().StringSliceVar(&daVinciFormNames, davinciFormNamesParamName, davinci.BootstrapDaVinciFormNames, "The list of DaVinci form names to search for to delete.  Case sensitive.")

	// Do the binds
	for k, v := range davinciFormsConfigurationParamMapping {
		viper.BindPFlag(v, cleanDaVinciFormsCmd.PersistentFlags().Lookup(k))
	}
}
