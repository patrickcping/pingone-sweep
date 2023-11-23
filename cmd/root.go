package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/patrickcping/pingone-clean-config/internal/logger"
	"github.com/patrickcping/pingone-clean-config/internal/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"

	// goreleaser can also pass the specific commit if you want
	commit string = ""
)

var (
	region              string
	workerEnvironmentId string
	workerClientId      string
	workerClientSecret  string
	environmentID       string
	dryRun              bool
	apiClient           *sdk.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pingone-cleanconfig",
	Short:   "pingone-cleanconfig is a CLI to clean demo bootstrap configuration from a PingOne environment.",
	Version: fmt.Sprintf("%s-%s", version, commit),
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()
		l.Debug().Msgf("Clean Command called for all services.")

		var err error
		apiClient, err = initApiClient(cmd.Context(), cmd.Version)
		if err != nil {
			return err
		}

		commands := cmd.Commands()

		for _, command := range commands {
			if command.Name() != "completion" && command.Name() != "help" {
				l.Debug().Msgf("Running command: %s", command.Name())
				err := command.RunE(cmd, args)
				if err != nil {
					return err
				}
			}
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// General function commands
	rootCmd.AddCommand(
		cleanBrandingThemesCmd,
		cleanDaVinciFormsCmd,
		cleanDirectoryAttributesCmd,
		cleanPasswordPoliciesCmd,
		cleanKeysCmd,
	)

	// Add config flags
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", viper.GetString("PINGONE_REGION"), "The region code of the service (NA, EU, AP, CA).")
	rootCmd.MarkFlagRequired("region")

	// Worker token auth
	rootCmd.PersistentFlags().StringVarP(&workerEnvironmentId, "worker-environment-id", "", viper.GetString("PINGONE_ENVIRONMENT_ID"), "The ID of the PingOne environment that contains the worker token client used to authenticate.")
	rootCmd.PersistentFlags().StringVarP(&workerClientId, "worker-client-id", "", viper.GetString("PINGONE_CLIENT_ID"), "The ID of the worker app (also the client ID) used to authenticate.")
	rootCmd.PersistentFlags().StringVarP(&workerClientSecret, "worker-client-secret", "", viper.GetString("PINGONE_CLIENT_SECRET"), "The client secret of the worker app used to authenticate.")

	rootCmd.MarkFlagsRequiredTogether("worker-environment-id", "worker-client-id", "worker-client-secret")

	// Target environment
	rootCmd.PersistentFlags().StringVarP(&environmentID, "target-environment-id", "", viper.GetString("PINGONE_TARGET_ENVIRONMENT_ID"), "The ID of the target environment to clean.")
	rootCmd.MarkFlagRequired("target-environment-id")

	// Dry run
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Run a clean routine but don't delete any configuration - instead issue a warning if configuration were to be deleted.")

}

func initApiClient(ctx context.Context, version string) (*sdk.Client, error) {
	l := logger.Get()

	if apiClient != nil {
		return apiClient, nil
	}

	l.Debug().Msgf("Initialising API client..")

	apiConfig := sdk.Config{
		ClientID:      workerClientId,
		ClientSecret:  workerClientSecret,
		EnvironmentID: workerEnvironmentId,
		Region:        region,
	}

	return apiConfig.APIClient(ctx, version)

}
