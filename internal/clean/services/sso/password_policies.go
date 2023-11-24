package sso

import (
	"context"
	"net/http"

	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/logger"
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

	configKey := "Password Policies"

	l.Debug().Msgf(`[%s] Cleaning bootstrap config for environment ID "%s"..`, configKey, c.Environment.EnvironmentID)

	if len(c.BootstrapPasswordPolicyNames) == 0 {
		l.Warn().Msgf("[%s] No bootstrap names configured - skipping", configKey)
		return nil
	}

	ok, err := clean.CheckBillOfMaterials(ctx, configKey, c.Environment, management.ENUMPRODUCTTYPE_ONE_BASE)
	if err != nil {
		return err
	}

	if !ok {
		l.Warn().Msgf("[%s] Bill of materials does not contain applicable service %s - skipping", configKey, management.ENUMPRODUCTTYPE_ONE_BASE)
		return nil
	}

	var response *management.EntityArray
	err = clean.ReadAllConfig(
		ctx,
		configKey,
		c.Environment,
		func() (any, *http.Response, error) {
			return c.Environment.Client.ManagementAPIClient.PasswordPoliciesApi.ReadAllPasswordPolicies(ctx, c.Environment.EnvironmentID).Execute()
		},
		&response,
	)
	if err != nil {
		return err
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasPasswordPolicies() {

		l.Debug().Msgf("[%s] Configuration items found, looping..", configKey)
		for _, policy := range embedded.GetPasswordPolicies() {

			err := clean.TryCleanConfig(
				ctx,
				configKey,
				c.Environment,
				clean.ConfigItem{
					IdentifierToEvaluate: policy.Name,
					Id:                   *policy.Id,
					Default:              policy.Default,
				},
				clean.ConfigItemEval{
					IdentifierListToSearch: c.BootstrapPasswordPolicyNames,
					StartsWithStringMatch:  false,
				},
				func() (any, *http.Response, error) {
					fR, fErr := c.Environment.Client.ManagementAPIClient.PasswordPoliciesApi.DeletePasswordPolicy(ctx, c.Environment.EnvironmentID, policy.GetId()).Execute()
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
