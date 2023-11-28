package cmd

import (
	"fmt"
	"strings"

	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/clean/services/mfa"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	mfaDevicePolicyNames []string
)

const (
	mfaDevicePolicyCmdName = "policy-name"

	mfaDevicePolicyNamesParamName      = "policy-name"
	mfaDevicePolicyNamesParamConfigKey = "pingone.services.mfa.device-policies.names"
)

var (
	mfaDevicePolicyConfigurationParamMapping = map[string]string{
		mfaDevicePolicyNamesParamName: mfaDevicePolicyNamesParamConfigKey,
	}
)

var cleanMfaDevicePoliciesCmd = &cobra.Command{
	Use:   mfaDevicePolicyCmdName,
	Short: "Clean unwanted demo MFA Device policy configuration",
	Long: fmt.Sprintf(`Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s "Default MFA Policy" --%s "Default MFA Policy 2" --%s
	
	`, mfaDevicePolicyCmdName, environmentIDParamName, dryRunParamName, mfaDevicePolicyCmdName, environmentIDParamName, mfaDevicePolicyNamesParamName, mfaDevicePolicyNamesParamName, dryRunParamName),
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool(dryRunParamConfigKey)
		mfaDevicePolicyNames := viper.GetStringSlice(mfaDevicePolicyNamesParamConfigKey)

		l.Debug().Msgf("Clean Command called for MFA Device policies.")
		l.Debug().Msgf("Dry run setting: %t", dryRun)
		l.Debug().Msgf(`MFA Device Policy names: "%s"`, strings.Join(mfaDevicePolicyNames, `", "`))

		var err error
		apiClient, err = initApiClient(cmd.Context(), cmd.Version)
		if err != nil {
			return err
		}

		cleanConfig := mfa.CleanEnvironmentPlatformMFADevicePoliciesConfig{
			Environment: clean.CleanEnvironmentConfig{
				EnvironmentID: viper.GetString(environmentIDParamConfigKey),
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapMFADevicePolicyNames: mfaDevicePolicyNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	l := logger.Get()

	cleanMfaDevicePoliciesCmd.PersistentFlags().StringSliceVar(&mfaDevicePolicyNames, mfaDevicePolicyNamesParamName, mfa.BootstrapMFADevicePolicyNames, "The list of MFA Device policy names to search for to delete.  Case sensitive.")

	if err := bindParams(mfaDevicePolicyConfigurationParamMapping, cleanMfaDevicePoliciesCmd); err != nil {
		l.Err(err).Msgf("Error binding parameters: %s", err)
	}
}
