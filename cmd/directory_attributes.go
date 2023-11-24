package cmd

import (
	"strings"

	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/clean/services/platform"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	directoryAttributeNames []string
)

var cleanDirectoryAttributesCmd = &cobra.Command{
	Use:   "directory-attributes",
	Short: "Disable unwanted demo directory schema attributes",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep directory-attributes --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36
	pingone-sweep directory-attributes --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-sweep directory-attributes --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --attribute-names accountId,address,email,externalId,locale,mobilePhone --dry-run
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool("dry-run")
		directoryAttributeNames := viper.GetStringSlice("pingone.services.platform.directory-schema.attribute-names")

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
				EnvironmentID: viper.GetString("pingone.target-environment-id"),
				DryRun:        viper.GetBool("dry-run"),
				Client:        apiClient.API,
			},
			BootstrapAttributeNames: directoryAttributeNames,
			SchemaName:              nil,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanDirectoryAttributesCmd.PersistentFlags().StringSliceVar(&directoryAttributeNames, "attribute-names", platform.BootstrapDirectoryAttributeNames, "The list of directory attribute names to search for to disable.")
	viper.BindPFlag("pingone.services.platform.directory-schema.attribute-names", cleanDirectoryAttributesCmd.PersistentFlags().Lookup("attribute-names"))
}
