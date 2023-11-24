package cmd

import (
	"strings"

	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/clean/services/verify"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/spf13/cobra"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()
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
				EnvironmentID: environmentID,
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapVerifyPolicyNames: verifyPolicyNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanVerifyPoliciesCmd.PersistentFlags().StringSliceVar(&verifyPolicyNames, "policy-name", verify.BootstrapVerifyPolicyNames, "The list of Verify policy names to search for to delete.  Case sensitive.")
}
