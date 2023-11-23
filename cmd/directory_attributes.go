package cmd

import (
	"strings"

	"github.com/patrickcping/pingone-clean-config/internal/clean"
	"github.com/patrickcping/pingone-clean-config/internal/clean/services/platform"
	"github.com/patrickcping/pingone-clean-config/internal/logger"
	"github.com/spf13/cobra"
)

var (
	directoryAttributeNames []string
)

var cleanDirectoryAttributesCmd = &cobra.Command{
	Use:   "directory-attributes",
	Short: "Disable unwanted demo directory schema attributes",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-cleanconfig directory-attributes --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36
	pingone-cleanconfig directory-attributes --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-cleanconfig directory-attributes --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --attribute-names accountId,address,email,externalId,locale,mobilePhone --dry-run
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()
		l.Debug().Msgf("Clean Command called for directory attributes.")
		l.Debug().Msgf("Dry run setting: %t", dryRun)
		l.Debug().Msgf(`Attribute names: "%s"`, strings.Join(directoryAttributeNames, `", "`))

		var err error
		apiClient, err = initApiClient(cmd.Context(), cmd.Version)
		if err != nil {
			return err
		}

		cleanConfig := platform.CleanEnvironmentPlatformDirectoryAttributeConfig{
			Environment: clean.CleanEnvironmentConfig{
				Client:        apiClient.API.ManagementAPIClient,
				EnvironmentID: environmentID,
				DryRun:        dryRun,
			},
			BootstrapAttributeNames: directoryAttributeNames,
			SchemaName:              nil,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanDirectoryAttributesCmd.PersistentFlags().StringSliceVar(&directoryAttributeNames, "attribute-names", platform.BootstrapDirectoryAttributeNames, "The list of directory attribute names to search for to disable.")
}
