package cmd

import (
	"strings"

	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/clean/services/mfa"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/spf13/cobra"
)

var (
	mfaFido2PolicyNames []string
)

var cleanMfaFido2PoliciesCmd = &cobra.Command{
	Use:   "mfa-fido2-policies",
	Short: "Clean unwanted demo MFA FIDO2 policy configuration",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep mfa-fido2-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-sweep mfa-fido2-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --policy-name "Passkeys" --policy-name "Security Keys" --dry-run
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()
		l.Debug().Msgf("Clean Command called for MFA FIDO2 policies.")
		l.Debug().Msgf("Dry run setting: %t", dryRun)
		l.Debug().Msgf(`MFA FIDO2 Policy names: "%s"`, strings.Join(mfaFido2PolicyNames, `", "`))

		var err error
		apiClient, err = initApiClient(cmd.Context(), cmd.Version)
		if err != nil {
			return err
		}

		cleanConfig := mfa.CleanEnvironmentPlatformMFAFIDO2PoliciesConfig{
			Environment: clean.CleanEnvironmentConfig{
				EnvironmentID: environmentID,
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapMFAFIDO2PolicyNames: mfaFido2PolicyNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanMfaFido2PoliciesCmd.PersistentFlags().StringSliceVar(&mfaFido2PolicyNames, "policy-name", mfa.BootstrapMFAFIDO2PolicyNames, "The list of MFA FIDO2 policy names to search for to delete.  Case sensitive.")
}
