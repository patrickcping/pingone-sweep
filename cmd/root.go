package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/patrickcping/pingone-sweep/internal/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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

	rootConfigurationParamMapping = map[string]string{
		"region":                "pingone.region",
		"target-environment-id": "pingone.target-environment-id",
		"dry-run":               "dry-run",
		"worker-environment-id": "pingone.worker-environment-id",
		"worker-client-id":      "pingone.worker-client-id",
		"worker-client-secret":  "pingone.worker-client-secret",
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pingone-sweep",
	Short: "pingone-sweep is a CLI to clean demo bootstrap configuration from a PingOne environment.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := initConfig(cmd)
		if err != nil {
			return err
		}

		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if v, ok := rootConfigurationParamMapping[f.Name]; ok && viper.IsSet(v) {
				cmd.Flags().SetAnnotation(f.Name, cobra.BashCompOneRequiredFlag, []string{"false"})
			}
		})

		return nil
	},
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
		cleanAuthenticationPoliciesCmd,
		cleanBrandingThemesCmd,
		cleanDaVinciFormsCmd,
		cleanDirectoryAttributesCmd,
		cleanKeysCmd,
		cleanMfaFido2PoliciesCmd,
		cleanNotificationPoliciesCmd,
		cleanPasswordPoliciesCmd,
		cleanRiskPoliciesCmd,
		cleanVerifyPoliciesCmd,
	)

	// Add config flags
	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", viper.GetString("PINGONE_REGION"), "The region code of the service (NA, EU, AP, CA).")
	rootCmd.MarkPersistentFlagRequired("region")

	// Worker token auth
	rootCmd.PersistentFlags().StringVar(&workerEnvironmentId, "worker-environment-id", viper.GetString("PINGONE_ENVIRONMENT_ID"), "The ID of the PingOne environment that contains the worker token client used to authenticate.")
	rootCmd.PersistentFlags().StringVar(&workerClientId, "worker-client-id", viper.GetString("PINGONE_CLIENT_ID"), "The ID of the worker app (also the client ID) used to authenticate.")
	rootCmd.PersistentFlags().StringVar(&workerClientSecret, "worker-client-secret", viper.GetString("PINGONE_CLIENT_SECRET"), "The client secret of the worker app used to authenticate.")

	rootCmd.MarkFlagsRequiredTogether("worker-environment-id", "worker-client-id", "worker-client-secret")

	// Target environment
	rootCmd.PersistentFlags().StringVar(&environmentID, "target-environment-id", viper.GetString("PINGONE_TARGET_ENVIRONMENT_ID"), "The ID of the target environment to clean.")
	rootCmd.MarkPersistentFlagRequired("target-environment-id")

	// Dry run
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "Run a clean routine but don't delete any configuration - instead issue a warning if configuration were to be deleted.")

	// Do the binds
	for k, v := range rootConfigurationParamMapping {
		viper.BindPFlag(v, rootCmd.PersistentFlags().Lookup(k))
	}
}

func initConfig(cmd *cobra.Command) error {
	l := logger.Get()

	l.Debug().Msgf("Initialising configuration..")

	viper.SetConfigName(".pingone-sweep")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	l.Debug().Msgf("Reading configuration..")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			l.Err(err).Msgf("Error reading configuration file.")
			return err
		}
	}

	viper.SetEnvPrefix("PINGONE")

	viper.AutomaticEnv()

	l.Debug().Msgf("Setting configuration..")

	return nil
}

func initApiClient(ctx context.Context, version string) (*sdk.Client, error) {
	l := logger.Get()

	if apiClient != nil {
		return apiClient, nil
	}

	l.Debug().Msgf("Initialising API client..")

	apiConfig := sdk.Config{
		ClientID:      viper.GetString("pingone.worker-client-id"),
		ClientSecret:  viper.GetString("pingone.worker-client-secret"),
		EnvironmentID: viper.GetString("pingone.worker-environment-id"),
		Region:        viper.GetString("pingone.region"),
	}

	return apiConfig.APIClient(ctx, version)

}
