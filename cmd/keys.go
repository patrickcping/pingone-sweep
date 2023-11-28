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
	keyIssuerDNPrefixes []string
	keyCaseSensitive    bool
)

const (
	keysCmdName = "keys"

	keysIssuerDNPrefixesParamName      = "issuer-dn-prefix"
	keysIssuerDNPrefixesParamConfigKey = "pingone.services.platform.keys.issuer-dn-prefixes"

	keysCaseSensitiveParamName      = "case-sensitive"
	keysCaseSensitiveParamConfigKey = "pingone.services.platform.keys.case-sensitive"
)

var (
	keysConfigurationParamMapping = map[string]string{
		keysIssuerDNPrefixesParamName: keysIssuerDNPrefixesParamConfigKey,
		keysCaseSensitiveParamName:    keysCaseSensitiveParamConfigKey,
	}
)

var cleanKeysCmd = &cobra.Command{
	Use:   keysCmdName,
	Short: "Clean unwanted demo key configuration",
	Long: fmt.Sprintf(`Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s
	pingone-sweep %s --%s 4457a4b7-332e-4e38-9956-09d6e8a19d36 --%s "C=US,O=Ping Identity,OU=Ping Identity" --%s
	
	`, keysCmdName, environmentIDParamName, dryRunParamName, keysCmdName, environmentIDParamName, keysIssuerDNPrefixesParamName, dryRunParamName),
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()

		dryRun := viper.GetBool(dryRunParamConfigKey)
		keyIssuerDNPrefixes := viper.GetStringSlice(keysIssuerDNPrefixesParamConfigKey)
		keyCaseSensitive := viper.GetBool(keysCaseSensitiveParamConfigKey)

		l.Debug().Msgf("Clean Command called for keys.")
		l.Debug().Msgf("Dry run setting: %t", dryRun)
		l.Debug().Msgf(`Issuer DN prefixes: "%s"`, strings.Join(keyIssuerDNPrefixes, `", "`))

		var err error
		apiClient, err = initApiClient(cmd.Context(), cmd.Version)
		if err != nil {
			return err
		}

		cleanConfig := platform.CleanEnvironmentPlatformKeysConfig{
			Environment: clean.CleanEnvironmentConfig{
				EnvironmentID: viper.GetString(environmentIDParamConfigKey),
				DryRun:        dryRun,
				Client:        apiClient.API,
			},
			BootstrapIssuerDNPrefixes: keyIssuerDNPrefixes,
			CaseSensitive:             keyCaseSensitive,
		}

		return cleanConfig.Clean(cmd.Context())
	},
}

func init() {
	cleanKeysCmd.PersistentFlags().StringArrayVar(&keyIssuerDNPrefixes, keysIssuerDNPrefixesParamName, platform.BootstrapKeyIssuerDNPrefixes, "The list of issuer DN prefixes to search for to delete.")
	cleanKeysCmd.PersistentFlags().BoolVar(&keyCaseSensitive, keysCaseSensitiveParamName, false, "The issuer DN prefix search is case sensitive.")

	// Do the binds
	for k, v := range keysConfigurationParamMapping {
		viper.BindPFlag(v, cleanKeysCmd.PersistentFlags().Lookup(k))
	}
}
