package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/clean/services/verify"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	verifyPolicyNames []string
)

const (
	verifyPoliciesCmdName = "verify-policies"

	verifyPolicyNamesParamName      = "policy-name"
	verifyPolicyNamesParamConfigKey = "pingone.services.verify.policies.names"
)

var (
	verifyPolicyConfigurationParamMapping = map[string]string{
		verifyPolicyNamesParamName: verifyPolicyNamesParamConfigKey,
	}
)

var cleanVerifyPoliciesCmd = &cobra.Command{
	Use:   verifyPoliciesCmdName,
	Short: "Clean unwanted demo Verify policy configuration",
	Long: fmt.Sprintf(`Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s "Default Verify Policy" --%s "Default Verify Policy 2" --%s
	
	`, verifyPoliciesCmdName, environmentIDParamName, dryRunParamName, verifyPoliciesCmdName, environmentIDParamName, verifyPolicyNamesParamName, verifyPolicyNamesParamName, dryRunParamName),
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.SetOut(os.Stdout)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool(dryRunParamConfigKey)
		verifyPolicyNames := viper.GetStringSlice(verifyPolicyNamesParamConfigKey)

		l.Debug().Msgf("Clean Command called for Verify policies.")
		l.Debug().Msgf("Dry run setting: %t", dryRun)
		l.Debug().Msgf(`Verify Policy names: "%s"`, strings.Join(verifyPolicyNames, `", "`))

		var err error
		apiClient, err = initApiClient(cmd.Context(), cmd.Version)
		if err != nil {
			return err
		}

		cleanConfig := verify.CleanEnvironmentVerifyPoliciesConfig{
			Environment: clean.CleanEnvironmentConfig{
				EnvironmentID: viper.GetString(environmentIDParamConfigKey),
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapVerifyPolicyNames: verifyPolicyNames,
		}

		cmd.Printf("Tests!!")

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanVerifyPoliciesCmd.PersistentFlags().StringSliceVar(&verifyPolicyNames, verifyPolicyNamesParamName, verify.BootstrapVerifyPolicyNames, "The list of Verify policy names to search for to delete.  Case sensitive.")

	// Do the binds
	for k, v := range verifyPolicyConfigurationParamMapping {
		viper.BindPFlag(v, cleanVerifyPoliciesCmd.PersistentFlags().Lookup(k))
	}
}
