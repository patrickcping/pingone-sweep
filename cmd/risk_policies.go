package cmd

import (
	"fmt"
	"strings"

	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/clean/services/protect"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	riskPolicyNames []string
)

const (
	riskPoliciesCmdName = "risk-policies"

	riskPolicyNamesParamName      = "policy-name"
	riskPolicyNamesParamConfigKey = "pingone.services.protect.risk-policies.names"
)

var (
	riskPolicyConfigurationParamMapping = map[string]string{
		riskPolicyNamesParamName: riskPolicyNamesParamConfigKey,
	}
)

var cleanRiskPoliciesCmd = &cobra.Command{
	Use:   riskPoliciesCmdName,
	Short: "Clean unwanted demo Risk policy configuration",
	Long: fmt.Sprintf(`Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s "Default Risk Policy" --%s "Default Risk Policy 2" --%s
	
	`, riskPoliciesCmdName, environmentIDParamName, dryRunParamName, riskPoliciesCmdName, environmentIDParamName, riskPolicyNamesParamName, riskPolicyNamesParamName, dryRunParamName),
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool(dryRunParamConfigKey)
		riskPolicyNames := viper.GetStringSlice(riskPolicyNamesParamConfigKey)

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
				EnvironmentID: viper.GetString(environmentIDParamConfigKey),
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapRiskPolicyNames: riskPolicyNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	l := logger.Get()

	cleanRiskPoliciesCmd.PersistentFlags().StringSliceVar(&riskPolicyNames, riskPolicyNamesParamName, protect.BootstrapRiskPolicyNames, "The list of Risk policy names to search for to delete.  Case sensitive.")

	if err := bindParams(riskPolicyConfigurationParamMapping, cleanRiskPoliciesCmd); err != nil {
		l.Err(err).Msgf("Error binding parameters: %s", err)
	}
}
