package cmd

import (
	"fmt"
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

const (
	authenticationPoliciesCmdName = "authentication-policies"

	authenticationPolicyNamesParamName      = "policy-name"
	authenticationPolicyNamesParamConfigKey = "pingone.services.sso.authentication-policies.names"
)

var (
	authenticationPolicyConfigurationParamMapping = map[string]string{
		authenticationPolicyNamesParamName: authenticationPolicyNamesParamConfigKey,
	}
)

var cleanAuthenticationPoliciesCmd = &cobra.Command{
	Use:   authenticationPoliciesCmdName,
	Short: "Clean unwanted demo sign-on (authentication) policy configuration",
	Long: fmt.Sprintf(`Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s "Single_Factor" --%s "Multi_Factor" --%s
	
	`, authenticationPoliciesCmdName, environmentIDParamName, dryRunParamName, authenticationPoliciesCmdName, environmentIDParamName, authenticationPolicyNamesParamName, authenticationPolicyNamesParamName, dryRunParamName),
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool(dryRunParamConfigKey)
		authenticationPolicyNames := viper.GetStringSlice(authenticationPolicyNamesParamConfigKey)

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
				EnvironmentID: viper.GetString(environmentIDParamConfigKey),
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapAuthenticationPolicyNames: authenticationPolicyNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanAuthenticationPoliciesCmd.PersistentFlags().StringSliceVar(&authenticationPolicyNames, authenticationPolicyNamesParamName, sso.BootstrapAuthenticationPolicyNames, "The list of sign-on (authentication) policy names to search for to delete.  Case sensitive.")

	// Do the binds
	for k, v := range authenticationPolicyConfigurationParamMapping {
		viper.BindPFlag(v, cleanAuthenticationPoliciesCmd.PersistentFlags().Lookup(k))
	}
}
