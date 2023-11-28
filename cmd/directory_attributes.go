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
	directoryAttributeNames []string
)

const (
	directoryAttributesCmdName = "directory-attributes"

	directoryAttributeNamesParamName      = "attribute-names"
	directoryAttributeNamesParamConfigKey = "pingone.services.platform.directory-schema.attribute-names"
)

var (
	directoryAttributesConfigurationParamMapping = map[string]string{
		directoryAttributeNamesParamName: directoryAttributeNamesParamConfigKey,
	}
)

var cleanDirectoryAttributesCmd = &cobra.Command{
	Use:   directoryAttributesCmdName,
	Short: "Disable unwanted demo directory schema attributes",
	Long: fmt.Sprintf(`Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s accountId,address,email,externalId,locale,mobilePhone --%s
	
	`, directoryAttributesCmdName, environmentIDParamName, directoryAttributesCmdName, environmentIDParamName, dryRunParamName, directoryAttributesCmdName, environmentIDParamName, directoryAttributeNamesParamName, dryRunParamName),
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool(dryRunParamConfigKey)
		directoryAttributeNames := viper.GetStringSlice(directoryAttributeNamesParamConfigKey)

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
				EnvironmentID: viper.GetString(environmentIDParamConfigKey),
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapAttributeNames: directoryAttributeNames,
			SchemaName:              nil,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanDirectoryAttributesCmd.PersistentFlags().StringSliceVar(&directoryAttributeNames, directoryAttributeNamesParamName, platform.BootstrapDirectoryAttributeNames, "The list of directory attribute names to search for to disable.")

	// Do the binds
	for k, v := range directoryAttributesConfigurationParamMapping {
		viper.BindPFlag(v, cleanDirectoryAttributesCmd.PersistentFlags().Lookup(k))
	}
}
