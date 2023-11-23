package platform

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/patrickcping/pingone-clean-config/internal/clean"
	"github.com/patrickcping/pingone-clean-config/internal/logger"
	"github.com/patrickcping/pingone-clean-config/internal/sdk"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

var (
	BootstrapKeyIssuerDNPrefixes = []string{
		"C=US,O=Ping Identity,OU=Ping Identity",
	}
)

type CleanEnvironmentPlatformKeysConfig struct {
	Environment               clean.CleanEnvironmentConfig
	BootstrapIssuerDNPrefixes []string
	CaseSensitive             bool
}

func (c *CleanEnvironmentPlatformKeysConfig) Clean(ctx context.Context) error {
	l := logger.Get()

	l.Debug().Msgf(`Cleaning bootstrap keys for environment ID "%s"..`, c.Environment.EnvironmentID)

	if len(c.BootstrapIssuerDNPrefixes) == 0 {
		l.Warn().Msgf("No bootstrap key names configured - skipping")
		return nil
	}

	var response *management.EntityArray
	err := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return c.Environment.Client.CertificateManagementApi.GetKeys(ctx, c.Environment.EnvironmentID).Execute()
		},
		"GetKeys",
		sdk.DefaultCreateReadRetryable,
		&response,
	)

	if err != nil {
		return err
	}

	if response == nil {
		return fmt.Errorf("No keys found - the API responded with no data")
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasKeys() {

		l.Debug().Msg("Keys found, looping..")
		for _, key := range embedded.GetKeys() {

			l.Debug().Msgf(`Looping bootstrapped key names for "%s"..`, key.GetName())
			for _, defaultIssuerPrefix := range c.BootstrapIssuerDNPrefixes {

				matchesPrefix := false
				if c.CaseSensitive {
					matchesPrefix = strings.HasPrefix(key.GetIssuerDN(), defaultIssuerPrefix)
				} else {
					matchesPrefix = strings.HasPrefix(strings.ToLower(key.GetIssuerDN()), strings.ToLower(defaultIssuerPrefix))
				}

				if matchesPrefix {
					l.Debug().Msgf(`Found "%s" (%s) key`, key.GetName(), key.GetUsageType())

					if key.GetDefault() {
						l.Warn().Msgf(`The "%s" (%s) key is set as the default key - this will not be deleted`, key.GetName(), key.GetUsageType())

						break
					}

					if !c.Environment.DryRun {
						err := sdk.ParseResponse(
							ctx,

							func() (any, *http.Response, error) {
								r, err := c.Environment.Client.CertificateManagementApi.DeleteKey(ctx, c.Environment.EnvironmentID, key.GetId()).Execute()
								return nil, r, err
							},
							"DeleteKey",
							sdk.DefaultCreateReadRetryable,
							nil,
						)

						if err != nil {
							return err
						}
						l.Info().Msgf(`Key "%s" (%s) deleted`, key.GetName(), key.GetUsageType())
					} else {
						l.Warn().Msgf(`Dry run: key "%s" (%s) with ID "%s" would be deleted`, key.GetName(), key.GetUsageType(), key.GetId())
					}

					break
				}
			}
		}
		l.Debug().Msg("Keys done")

	} else {
		l.Debug().Msg("No Keys found in the target environment")
	}

	return nil
}
