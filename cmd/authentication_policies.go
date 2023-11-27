package cmd

import (
	"strings"

	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/clean/services/sso"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	authenticationPolicyNames []string
)

var cleanAuthenticationPoliciesCmd = &cobra.Command{
	Use:   "authentication-policies",
	Short: "Clean unwanted demo sign-on (authentication) policy configuration",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep authentication-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-sweep authentication-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --policy-name "Single_Factor" --policy-name "Multi_Factor" --dry-run
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool("dry-run")
		authenticationPolicyNames := viper.GetStringSlice("pingone.services.platform.authentication-policies.names")

		l.Debug().Msgf("Clean Command called for sign-on (authentication) policies.")
		l.Debug().Msgf("Dry run setting: %t", dryRun)
		l.Debug().Msgf(`sign-on (authentication) Policy names: "%s"`, strings.Join(authenticationPolicyNames, `", "`))

		var err error
		apiClient, err = initApiClient(cmd.Context(), cmd.Version)
		if err != nil {
			return err
		}

		cleanConfig := sso.CleanEnvironmentAuthenticationPoliciesConfig{
			Environment: clean.CleanEnvironmentConfig{
				EnvironmentID: viper.GetString("pingone.target-environment-id"),
				DryRun:        viper.GetBool("dry-run"),
				Client:        apiClient.API,
			},
			BootstrapAuthenticationPolicyNames: authenticationPolicyNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanAuthenticationPoliciesCmd.PersistentFlags().StringSliceVar(&authenticationPolicyNames, "policy-name", sso.BootstrapAuthenticationPolicyNames, "The list of sign-on (authentication) policy names to search for to delete.  Case sensitive.")
	viper.BindPFlag("pingone.services.sso.authentication-policies.names", cleanAuthenticationPoliciesCmd.PersistentFlags().Lookup("policy-name"))
}
