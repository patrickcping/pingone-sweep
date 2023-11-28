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
	mfaFido2PolicyNames []string
)

const (
	mfaFido2PoliciesCmdName = "mfa-fido2-policies"

	mfaFido2PolicyNamesParamName      = "policy-name"
	mfaFido2PolicyNamesParamConfigKey = "pingone.services.mfa.fido2-policies.names"
)

var (
	mfaFido2PolicyConfigurationParamMapping = map[string]string{
		mfaFido2PolicyNamesParamName: mfaFido2PolicyNamesParamConfigKey,
	}
)

var cleanMfaFido2PoliciesCmd = &cobra.Command{
	Use:   mfaFido2PoliciesCmdName,
	Short: "Clean unwanted demo MFA FIDO2 policy configuration",
	Long: fmt.Sprintf(`Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s "Passkeys" --%s "Security Keys" --%s
	
	`, mfaFido2PoliciesCmdName, environmentIDParamName, dryRunParamName, mfaFido2PoliciesCmdName, environmentIDParamName, mfaFido2PolicyNamesParamName, mfaFido2PolicyNamesParamName, dryRunParamName),
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool(dryRunParamConfigKey)
		mfaFido2PolicyNames := viper.GetStringSlice(mfaFido2PolicyNamesParamConfigKey)

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
				EnvironmentID: viper.GetString(environmentIDParamConfigKey),
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapMFAFIDO2PolicyNames: mfaFido2PolicyNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanMfaFido2PoliciesCmd.PersistentFlags().StringSliceVar(&mfaFido2PolicyNames, mfaFido2PolicyNamesParamName, mfa.BootstrapMFAFIDO2PolicyNames, "The list of MFA FIDO2 policy names to search for to delete.  Case sensitive.")

	// Do the binds
	for k, v := range mfaFido2PolicyConfigurationParamMapping {
		viper.BindPFlag(v, cleanMfaFido2PoliciesCmd.PersistentFlags().Lookup(k))
	}
}
