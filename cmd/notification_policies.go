package cmd

import (
	"fmt"
	"strings"

	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/clean/services/platform"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	notificationPolicyNames []string
)

const (
	notificationPoliciesCmdName = "notification-policies"

	notificationPolicyNamesParamName      = "policy-name"
	notificationPolicyNamesParamConfigKey = "pingone.services.platform.notification-policies.names"
)

var (
	notificationPolicyConfigurationParamMapping = map[string]string{
		notificationPolicyNamesParamName: notificationPolicyNamesParamConfigKey,
	}
)

var cleanNotificationPoliciesCmd = &cobra.Command{
	Use:   notificationPoliciesCmdName,
	Short: "Clean unwanted demo notification policy configuration",
	Long: fmt.Sprintf(`Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s "Default Notification Policy" --%s "Default Notification Policy 2" --%s
	
	`, notificationPoliciesCmdName, environmentIDParamName, dryRunParamName, notificationPoliciesCmdName, environmentIDParamName, notificationPolicyNamesParamName, notificationPolicyNamesParamName, dryRunParamName),
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool(dryRunParamConfigKey)
		notificationPolicyNames := viper.GetStringSlice(notificationPolicyNamesParamConfigKey)

		l.Debug().Msgf("Clean Command called for notification policies.")
		l.Debug().Msgf("Dry run setting: %t", dryRun)
		l.Debug().Msgf(`Notification Policy names: "%s"`, strings.Join(notificationPolicyNames, `", "`))

		var err error
		apiClient, err = initApiClient(cmd.Context(), cmd.Version)
		if err != nil {
			return err
		}

		cleanConfig := platform.CleanEnvironmentPlatformNotificationPoliciesConfig{
			Environment: clean.CleanEnvironmentConfig{
				EnvironmentID: viper.GetString(environmentIDParamConfigKey),
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapNotificationPolicyNames: notificationPolicyNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	l := logger.Get()

	cleanNotificationPoliciesCmd.PersistentFlags().StringSliceVar(&notificationPolicyNames, notificationPolicyNamesParamName, platform.BootstrapNotificationPolicyNames, "The list of notification policy names to search for to delete.  Case sensitive.")

	if err := bindParams(notificationPolicyConfigurationParamMapping, cleanNotificationPoliciesCmd); err != nil {
		l.Err(err).Msgf("Error binding parameters: %s", err)
	}
}
