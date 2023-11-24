package cmd

import (
	"strings"

	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/clean/services/protect"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/spf13/cobra"
)

var (
	riskPolicyNames []string
)

var cleanRiskPoliciesCmd = &cobra.Command{
	Use:   "risk-policies",
	Short: "Clean unwanted demo Risk policy configuration",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep risk-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-sweep risk-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --policy-name "Default Risk Policy" --policy-name "Default Risk Policy 2" --dry-run
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()
		l.Debug().Msgf("Clean Command called for Risk policies.")
		l.Debug().Msgf("Dry run setting: %t", dryRun)
		l.Debug().Msgf(`Risk Policy names: "%s"`, strings.Join(riskPolicyNames, `", "`))

		var err error
		apiClient, err = initApiClient(cmd.Context(), cmd.Version)
		if err != nil {
			return err
		}

		cleanConfig := protect.CleanEnvironmentProtectRiskPoliciesConfig{
			Environment: clean.CleanEnvironmentConfig{
				EnvironmentID: environmentID,
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapRiskPolicyNames: riskPolicyNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanRiskPoliciesCmd.PersistentFlags().StringSliceVar(&riskPolicyNames, "policy-name", protect.BootstrapRiskPolicyNames, "The list of Risk policy names to search for to delete.  Case sensitive.")
}
