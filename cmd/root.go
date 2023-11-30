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

const (
	regionParamName      = "region"
	regionParamConfigKey = "pingone.region"

	environmentIDParamName      = "target-environment-id"
	environmentIDParamConfigKey = "pingone.target-environment-id"

	dryRunParamName      = "dry-run"
	dryRunParamConfigKey = "dry-run"

	outputJsonParamName      = "json"
	outputJsonParamConfigKey = "output.json"

	outputNoColorParamName      = "no-color"
	outputNoColorParamConfigKey = "output.no-color"

	workerEnvironmentIDParamName      = "worker-environment-id"
	workerEnvironmentIDParamConfigKey = "pingone.worker-environment-id"

	workerClientIDParamName      = "worker-client-id"
	workerClientIDParamConfigKey = "pingone.worker-client-id"

	workerClientSecretParamName      = "worker-client-secret"
	workerClientSecretParamConfigKey = "pingone.worker-client-secret"
)

var (
	region              string
	workerEnvironmentId string
	workerClientId      string
	workerClientSecret  string
	environmentID       string
	dryRun              bool
	outputJson          bool
	outputNoColor       bool
	apiClient           *sdk.Client

	rootConfigurationParamMapping = map[string]string{
		regionParamName:              regionParamConfigKey,
		environmentIDParamName:       environmentIDParamConfigKey,
		dryRunParamName:              dryRunParamConfigKey,
		outputJsonParamName:          outputJsonParamConfigKey,
		outputNoColorParamName:       outputNoColorParamConfigKey,
		workerEnvironmentIDParamName: workerEnvironmentIDParamConfigKey,
		workerClientIDParamName:      workerClientIDParamConfigKey,
		workerClientSecretParamName:  workerClientSecretParamConfigKey,
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pingone-sweep",
	Short: "pingone-sweep is a CLI to clean demo bootstrap configuration from a PingOne environment.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		err := initConfig()
		if err != nil {
			return err
		}

		cmd.Flags().VisitAll(func(f *pflag.Flag) {
			if v, ok := rootConfigurationParamMapping[f.Name]; ok && viper.IsSet(v) {
				if err = cmd.Flags().SetAnnotation(f.Name, cobra.BashCompOneRequiredFlag, []string{"false"}); err != nil {
					l.Err(err).Msgf("Error setting required status for flag %s", f.Name)
					return
				}
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
	l := logger.Get()

	// General function commands
	rootCmd.AddCommand(
		cleanAuthenticationPoliciesCmd,
		cleanBrandingThemesCmd,
		cleanDaVinciFormsCmd,
		cleanDirectoryAttributesCmd,
		cleanKeysCmd,
		cleanMfaDevicePoliciesCmd,
		cleanMfaFido2PoliciesCmd,
		cleanNotificationPoliciesCmd,
		cleanPasswordPoliciesCmd,
		cleanRiskPoliciesCmd,
		cleanVerifyPoliciesCmd,
	)

	// Add config flags
	rootCmd.PersistentFlags().StringVarP(&region, regionParamName, "r", viper.GetString("PINGONE_REGION"), "The region code of the service (NA, EU, AP, CA).")
	if err := rootCmd.MarkPersistentFlagRequired(regionParamName); err != nil {
		l.Err(err).Msgf("Error marking flag %s as required.", regionParamName)
	}

	// Worker token auth
	rootCmd.PersistentFlags().StringVar(&workerEnvironmentId, workerEnvironmentIDParamName, viper.GetString("PINGONE_ENVIRONMENT_ID"), "The ID of the PingOne environment that contains the worker token client used to authenticate.")
	rootCmd.PersistentFlags().StringVar(&workerClientId, workerClientIDParamName, viper.GetString("PINGONE_CLIENT_ID"), "The ID of the worker app (also the client ID) used to authenticate.")
	rootCmd.PersistentFlags().StringVar(&workerClientSecret, workerClientSecretParamName, viper.GetString("PINGONE_CLIENT_SECRET"), "The client secret of the worker app used to authenticate.")

	rootCmd.MarkFlagsRequiredTogether(workerEnvironmentIDParamName, workerClientIDParamName, workerClientSecretParamName)

	// Target environment
	rootCmd.PersistentFlags().StringVar(&environmentID, environmentIDParamName, viper.GetString("PINGONE_TARGET_ENVIRONMENT_ID"), "The ID of the target environment to clean.")
	if err := rootCmd.MarkPersistentFlagRequired(environmentIDParamName); err != nil {
		l.Err(err).Msgf("Error marking flag %s as required.", environmentIDParamName)
	}

	// Dry run
	rootCmd.PersistentFlags().BoolVar(&dryRun, dryRunParamName, false, "Run a clean routine but don't delete any configuration - instead issue a warning if configuration were to be deleted.")

	// Output format
	rootCmd.PersistentFlags().BoolVar(&outputJson, outputJsonParamName, false, "Output in JSON format.")

	// Output color
	rootCmd.PersistentFlags().BoolVar(&outputNoColor, outputNoColorParamName, false, "Output without color formatting.")

	if err := bindParams(rootConfigurationParamMapping, rootCmd); err != nil {
		l.Err(err).Msgf("Error binding parameters: %s", err)
	}
}

func initConfig() error {
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
		ClientID:      viper.GetString(workerClientIDParamConfigKey),
		ClientSecret:  viper.GetString(workerClientSecretParamConfigKey),
		EnvironmentID: viper.GetString(workerEnvironmentIDParamConfigKey),
		Region:        viper.GetString(regionParamConfigKey),
	}

	return apiConfig.APIClient(ctx, version)

}

func bindParams(paramlist map[string]string, command *cobra.Command) error {
	// Do the binds
	for k, v := range paramlist {
		err := viper.BindPFlag(v, command.PersistentFlags().Lookup(k))
		if err != nil {
			return err
		}
	}

	return nil
}
