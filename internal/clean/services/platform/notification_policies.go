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
	BootstrapNotificationPolicyNames = []string{
		"Default Notification Policy",
	}
)

type CleanEnvironmentPlatformNotificationPoliciesConfig struct {
	Environment                      clean.CleanEnvironmentConfig
	BootstrapNotificationPolicyNames []string
}

func (c *CleanEnvironmentPlatformNotificationPoliciesConfig) Clean(ctx context.Context) error {
	l := logger.Get()

	l.Debug().Msgf(`Cleaning bootstrap notification policies for environment ID "%s"..`, c.Environment.EnvironmentID)

	if len(c.BootstrapNotificationPolicyNames) == 0 {
		l.Warn().Msgf("No bootstrap notification policy names configured - skipping")
		return nil
	}

	var response *management.EntityArray
	err := sdk.ParseResponse(
		ctx,

		func() (any, *http.Response, error) {
			return c.Environment.Client.NotificationsPoliciesApi.ReadAllNotificationsPolicies(ctx, c.Environment.EnvironmentID).Execute()
		},
		"ReadAllNotificationsPolicies",
		sdk.DefaultCreateReadRetryable,
		&response,
	)

	if err != nil {
		return err
	}

	if response == nil {
		return fmt.Errorf("No notification policies found - the API responded with no data")
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasNotificationsPolicies() {

		l.Debug().Msg("Notification policies found, looping..")
		for _, policy := range embedded.GetNotificationsPolicies() {

			l.Debug().Msgf(`Looping bootstrapped policy names for "%s"..`, policy.GetName())
			for _, defaultPolicyName := range c.BootstrapNotificationPolicyNames {

				if policy.GetName() == defaultPolicyName {
					l.Debug().Msgf(`Found "%s" notification policy`, defaultPolicyName)

					if policy.GetDefault() {
						l.Warn().Msgf(`The "%s" notification policy is set as the default policy - this will not be deleted`, defaultPolicyName)

						break
					}

					if !c.Environment.DryRun {
						err := sdk.ParseResponse(
							ctx,

							func() (any, *http.Response, error) {
								r, err := c.Environment.Client.NotificationsPoliciesApi.DeleteNotificationsPolicy(ctx, c.Environment.EnvironmentID, policy.GetId()).Execute()
								return nil, r, err
							},
							"DeleteNotificationsPolicy",
							sdk.DefaultCreateReadRetryable,
							nil,
						)

						if err != nil {
							return err
						}
						l.Info().Msgf(`Notification policy "%s" deleted`, defaultPolicyName)
					} else {
						l.Warn().Msgf(`Dry run: notification policy "%s" with ID "%s" would be deleted`, defaultPolicyName, policy.GetId())
					}

					break
				}
			}
		}

	} else {
		l.Debug().Msg("No Notification policies found in the target environment")
	}

	return nil
}
