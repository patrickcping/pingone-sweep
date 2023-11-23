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

	l.Debug().Msgf(`Cleaning bootstrap password policies for environment ID "%s"..`, c.Environment.EnvironmentID)

	if len(c.BootstrapPasswordPolicyNames) == 0 {
		l.Warn().Msgf("No bootstrap password policy names configured - skipping")
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
		return fmt.Errorf("No password policies found - the API responded with no data")
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasPasswordPolicies() {

		l.Debug().Msg("Password policies found, looping..")
		for _, policy := range embedded.GetPasswordPolicies() {

			l.Debug().Msgf(`Looping bootstrapped policy names for "%s"..`, policy.GetName())
			for _, defaultThemeName := range c.BootstrapPasswordPolicyNames {

				if policy.GetName() == defaultThemeName {
					l.Debug().Msgf(`Found "%s" password policy`, defaultThemeName)

					if policy.GetDefault() {
						l.Warn().Msgf(`The "%s" password policy is set as the default policy - this will not be deleted`, defaultThemeName)

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
						l.Info().Msgf(`Password policy "%s" deleted`, defaultThemeName)
					} else {
						l.Warn().Msgf(`Dry run: password policy "%s" with ID "%s" would be deleted`, defaultThemeName, policy.GetId())
					}

					break
				}
			}
		}

	} else {
		l.Debug().Msg("No Password policies found in the target environment")
	}

	return nil
}
