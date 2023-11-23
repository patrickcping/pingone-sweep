package platform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/patrickcping/pingone-clean-config/internal/clean"
	"github.com/patrickcping/pingone-clean-config/internal/logger"
	"github.com/patrickcping/pingone-clean-config/internal/sdk"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
)

var (
	BootstrapPasswordPolicyNames = []string{
		"Standard",
		"Basic",
		"Passphrase",
	}
)

type CleanEnvironmentPlatformPasswordPoliciesConfig struct {
	Environment                  clean.CleanEnvironmentConfig
	BootstrapPasswordPolicyNames []string
}

func (c *CleanEnvironmentPlatformPasswordPoliciesConfig) Clean(ctx context.Context) error {
	l := logger.Get()

	debugModule := "Password Policies"

	l.Debug().Msgf(`[%s] Cleaning bootstrap config for environment ID "%s"..`, debugModule, c.Environment.EnvironmentID)

	if len(c.BootstrapPasswordPolicyNames) == 0 {
		l.Warn().Msgf("[%s] No bootstrap names configured - skipping", debugModule)
		return nil
	}

	var response *management.EntityArray
	err := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return c.Environment.Client.PasswordPoliciesApi.ReadAllPasswordPolicies(ctx, c.Environment.EnvironmentID).Execute()
		},
		"ReadAllPasswordPolicies",
		sdk.DefaultCreateReadRetryable,
		&response,
	)

	if err != nil {
		return err
	}

	if response == nil {
		return fmt.Errorf("[%s] No configuration items found - the API responded with no data", debugModule)
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasPasswordPolicies() {

		l.Debug().Msgf("[%s] Configuration items found, looping..", debugModule)
		for _, policy := range embedded.GetPasswordPolicies() {

			l.Debug().Msgf(`[%s] Looping names for "%s"..`, debugModule, policy.GetName())
			for _, defaultPolicyName := range c.BootstrapPasswordPolicyNames {

				if policy.GetName() == defaultPolicyName {
					l.Debug().Msgf(`[%s] Found "%s"`, debugModule, defaultPolicyName)

					if policy.GetDefault() {
						l.Warn().Msgf(`[%s] "%s" is set as the environment default - this configuration will not be deleted`, debugModule, policy.GetName())

						break
					}

					if !c.Environment.DryRun {
						err := sdk.ParseResponse(
							ctx,

							func() (any, *http.Response, error) {
								r, err := c.Environment.Client.PasswordPoliciesApi.DeletePasswordPolicy(ctx, c.Environment.EnvironmentID, policy.GetId()).Execute()
								return nil, r, err
							},
							"DeletePasswordPolicy",
							sdk.DefaultCreateReadRetryable,
							nil,
						)

						if err != nil {
							return err
						}
						l.Info().Msgf(`[%s] "%s" deleted`, debugModule, policy.GetName())
					} else {
						l.Warn().Msgf(`[%s] Dry run: "%s" with ID "%s" would be deleted`, debugModule, policy.GetName(), policy.GetId())
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
