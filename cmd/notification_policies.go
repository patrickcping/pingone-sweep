package cmd

import (
	"strings"

	"github.com/patrickcping/pingone-clean-config/internal/clean"
	"github.com/patrickcping/pingone-clean-config/internal/clean/services/platform"
	"github.com/patrickcping/pingone-clean-config/internal/logger"
	"github.com/spf13/cobra"
)

var (
	notificationPolicyNames []string
)

var cleanNotificationPoliciesCmd = &cobra.Command{
	Use:   "notification-policies",
	Short: "Clean unwanted demo notification policy configuration",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-cleanconfig notification-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-cleanconfig notification-policies --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --policy-name "Default Notification Policy" --policy-name "Default Notification Policy 2" --dry-run
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()
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
				Client:        apiClient.API.ManagementAPIClient,
				EnvironmentID: environmentID,
				DryRun:        dryRun,
			},
			BootstrapNotificationPolicyNames: notificationPolicyNames,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanNotificationPoliciesCmd.PersistentFlags().StringSliceVar(&notificationPolicyNames, "policy-name", platform.BootstrapNotificationPolicyNames, "The list of notification policy names to search for to delete.  Case sensitive.")
}
