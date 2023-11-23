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

	debugModule := "Keys"

	l.Debug().Msgf(`[%s] Cleaning bootstrap config for environment ID "%s"..`, debugModule, c.Environment.EnvironmentID)

	if len(c.BootstrapIssuerDNPrefixes) == 0 {
		l.Warn().Msgf("[%s] No bootstrap names configured - skipping", debugModule)
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
		return fmt.Errorf("[%s] No configuration items found - the API responded with no data", debugModule)
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasKeys() {

		l.Debug().Msgf("[%s] Configuration items found, looping..", debugModule)
		for _, key := range embedded.GetKeys() {

			l.Debug().Msgf(`[%s] Looping names for "%s"..`, debugModule, key.GetName())
			for _, defaultIssuerPrefix := range c.BootstrapIssuerDNPrefixes {

				matchesPrefix := false
				if c.CaseSensitive {
					matchesPrefix = strings.HasPrefix(key.GetIssuerDN(), defaultIssuerPrefix)
				} else {
					matchesPrefix = strings.HasPrefix(strings.ToLower(key.GetIssuerDN()), strings.ToLower(defaultIssuerPrefix))
				}

				if matchesPrefix {
					l.Debug().Msgf(`[%s] Found "%s" (%s) key`, debugModule, key.GetName(), key.GetUsageType())

					if key.GetDefault() {
						l.Warn().Msgf(`[%s] "%s" (%s) is set as the environment default - this configuration will not be deleted`, debugModule, key.GetName(), key.GetUsageType())

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
						l.Info().Msgf(`[%s] "%s" (%s) deleted`, debugModule, key.GetName(), key.GetUsageType())
					} else {
						l.Warn().Msgf(`[%s] Dry run: key "%s" (%s) with ID "%s" would be deleted`, debugModule, key.GetName(), key.GetUsageType(), key.GetId())
					}

					break
				}
			}
		}
		l.Debug().Msgf("[%s] Done", debugModule)

	} else {
		l.Debug().Msgf("[%s] No configuration items found in the target environment", debugModule)
	}

	return nil
}
