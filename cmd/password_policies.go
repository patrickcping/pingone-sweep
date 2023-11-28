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
	passwordPolicyNames []string
)

const (
	passwordPoliciesCmdName = "password-policies"

	passwordPolicyNamesParamName      = "policy-name"
	passwordPolicyNamesParamConfigKey = "pingone.services.sso.password-policies.names"
)

var (
	passwordPolicyConfigurationParamMapping = map[string]string{
		passwordPolicyNamesParamName: passwordPolicyNamesParamConfigKey,
	}
)

var cleanPasswordPoliciesCmd = &cobra.Command{
	Use:   passwordPoliciesCmdName,
	Short: "Clean unwanted demo password policy configuration",
	Long: fmt.Sprintf(`Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s "Standard" --%s "Basic" --%s
	
	`, passwordPoliciesCmdName, environmentIDParamName, dryRunParamName, passwordPoliciesCmdName, environmentIDParamName, passwordPolicyNamesParamName, passwordPolicyNamesParamName, dryRunParamName),
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool(dryRunParamConfigKey)
		passwordPolicyNames := viper.GetStringSlice(passwordPolicyNamesParamConfigKey)

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
				EnvironmentID: viper.GetString(environmentIDParamConfigKey),
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapPasswordPolicyNames: passwordPolicyNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanPasswordPoliciesCmd.PersistentFlags().StringSliceVar(&passwordPolicyNames, passwordPolicyNamesParamName, sso.BootstrapPasswordPolicyNames, "The list of password policy names to search for to delete.  Case sensitive.")

	// Do the binds
	for k, v := range passwordPolicyConfigurationParamMapping {
		viper.BindPFlag(v, cleanPasswordPoliciesCmd.PersistentFlags().Lookup(k))
	}
}
