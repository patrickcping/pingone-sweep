package platform

import (
	"context"
	"net/http"

	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/logger"
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

	configKey := "Keys"

	l.Debug().Msgf(`[%s] Cleaning bootstrap config for environment ID "%s"..`, configKey, c.Environment.EnvironmentID)

	if len(c.BootstrapIssuerDNPrefixes) == 0 {
		l.Info().Msgf("[%s] No bootstrap names configured - skipping", configKey)
		return nil
	}

	var response *management.EntityArray
	err := clean.ReadAllConfig(
		ctx,
		configKey,
		c.Environment,
		func() (any, *http.Response, error) {
			return c.Environment.Client.ManagementAPIClient.CertificateManagementApi.GetKeys(ctx, c.Environment.EnvironmentID).Execute()
		},
		&response,
	)
	if err != nil {
		return err
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasKeys() {

		l.Debug().Msgf("[%s] Configuration items found, looping..", configKey)
		for _, key := range embedded.GetKeys() {

			err := clean.TryCleanConfig(
				ctx,
				configKey,
				c.Environment,
				clean.ConfigItem{
					IdentifierToEvaluate: key.Name,
					Id:                   *key.Id,
					Default:              key.Default,
				},
				clean.ConfigItemEval{
					IdentifierListToSearch: c.BootstrapIssuerDNPrefixes,
					StartsWithStringMatch:  true,
					CaseSensitive:          &c.CaseSensitive,
				},
				func() (any, *http.Response, error) {
					fR, fErr := c.Environment.Client.ManagementAPIClient.CertificateManagementApi.DeleteKey(ctx, c.Environment.EnvironmentID, key.GetId()).Execute()
					return nil, fR, fErr
				},
				nil,
			)

			if err != nil {
				return err
			}

		}
		l.Debug().Msgf("[%s] Done", configKey)

	} else {
		l.Debug().Msgf("[%s] No configuration items found in the target environment", configKey)
	}

	return nil
}
