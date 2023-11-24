package cmd

import (
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

var cleanMfaDevicePoliciesCmd = &cobra.Command{
	Use:   "mfa-device-policies",
	Short: "Clean unwanted demo MFA Device policy configuration",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep mfa-device-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-sweep mfa-device-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --policy-name "Default MFA Policy" --policy-name "Default MFA Policy 2" --dry-run
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool("dry-run")
		mfaDevicePolicyNames := viper.GetStringSlice("pingone.services.mfa.device-policies.names")

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
				EnvironmentID: viper.GetString("pingone.target-environment-id"),
				DryRun:        viper.GetBool("dry-run"),
				Client:        apiClient.API,
			},
			BootstrapMFADevicePolicyNames: mfaDevicePolicyNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanMfaDevicePoliciesCmd.PersistentFlags().StringSliceVar(&mfaDevicePolicyNames, "policy-name", mfa.BootstrapMFADevicePolicyNames, "The list of MFA Device policy names to search for to delete.  Case sensitive.")
	viper.BindPFlag("pingone.services.mfa.device-policies.names", cleanMfaFido2PoliciesCmd.PersistentFlags().Lookup("policy-name"))
}
