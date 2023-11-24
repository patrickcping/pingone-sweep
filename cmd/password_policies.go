package cmd

import (
	"strings"

	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/clean/services/sso"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/spf13/cobra"
)

var (
	passwordPolicyNames []string
)

var cleanPasswordPoliciesCmd = &cobra.Command{
	Use:   "password-policies",
	Short: "Clean unwanted demo password policy configuration",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep password-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-sweep password-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --policy-name "Standard" --policy-name "Basic" --dry-run
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()
		l.Debug().Msgf("Clean Command called for password policies.")
		l.Debug().Msgf("Dry run setting: %t", dryRun)
		l.Debug().Msgf(`Password Policy names: "%s"`, strings.Join(passwordPolicyNames, `", "`))

		var err error
		apiClient, err = initApiClient(cmd.Context(), cmd.Version)
		if err != nil {
			return err
		}

		cleanConfig := sso.CleanEnvironmentPlatformPasswordPoliciesConfig{
			Environment: clean.CleanEnvironmentConfig{
				EnvironmentID: environmentID,
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapPasswordPolicyNames: passwordPolicyNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanPasswordPoliciesCmd.PersistentFlags().StringSliceVar(&passwordPolicyNames, "policy-name", sso.BootstrapPasswordPolicyNames, "The list of password policy names to search for to delete.  Case sensitive.")
}
