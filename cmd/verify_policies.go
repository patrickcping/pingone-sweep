package cmd

import (
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

var cleanVerifyPoliciesCmd = &cobra.Command{
	Use:   "verify-policies",
	Short: "Clean unwanted demo Verify policy configuration",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep verify-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-sweep verify-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --policy-name "Default Verify Policy" --policy-name "Default Verify Policy 2" --dry-run
	
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		cmd.SetOut(os.Stdout)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool("dry-run")
		verifyPolicyNames := viper.GetStringSlice("pingone.services.verify.policies.names")

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
				EnvironmentID: viper.GetString("pingone.target-environment-id"),
				DryRun:        viper.GetBool("dry-run"),
				Client:        apiClient.API,
			},
			BootstrapVerifyPolicyNames: verifyPolicyNames,
		}

		cmd.Printf("Tests!!")

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanVerifyPoliciesCmd.PersistentFlags().StringSliceVar(&verifyPolicyNames, "policy-name", verify.BootstrapVerifyPolicyNames, "The list of Verify policy names to search for to delete.  Case sensitive.")
	viper.BindPFlag("pingone.services.verify.policies.names", cleanVerifyPoliciesCmd.PersistentFlags().Lookup("policy-name"))
}
