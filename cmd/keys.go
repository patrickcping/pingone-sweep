package cmd

import (
	"strings"

	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/clean/services/platform"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/spf13/cobra"
)

var (
	keyIssuerDNPrefixes []string
	keyCaseSensitive    bool
)

var cleanKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Clean unwanted demo key configuration",
	Long: `Clean away demo configuration and prepare an environment for production-ready configuration.

	Examples:
	
	pingone-sweep keys --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --dry-run
	pingone-sweep keys --target-environment-id 4457a4b7-332e-4e38-9956-09d6e8a19d36 --issuer-dn-prefix "C=US,O=Ping Identity,OU=Ping Identity" --dry-run
	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := logger.Get()
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
				EnvironmentID: environmentID,
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
	cleanKeysCmd.PersistentFlags().StringArrayVar(&keyIssuerDNPrefixes, "issuer-dn-prefix", platform.BootstrapKeyIssuerDNPrefixes, "The list of issuer DN prefixes to search for to delete.")
	cleanKeysCmd.PersistentFlags().BoolVar(&keyCaseSensitive, "case-sensitive", false, "The issuer DN prefix search is case sensitive.")
}
