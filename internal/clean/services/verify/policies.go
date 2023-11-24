package verify

import (
	"context"
	"net/http"

	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/logger"
)

var (
	BootstrapVerifyPolicyNames = []string{
		"Default Verify Policy",
	}
)

type CleanEnvironmentVerifyPoliciesConfig struct {
	Environment                clean.CleanEnvironmentConfig
	BootstrapVerifyPolicyNames []string
}

func (c *CleanEnvironmentVerifyPoliciesConfig) Clean(ctx context.Context) error {
	l := logger.Get()

	configKey := "Verify Policies"

	l.Debug().Msgf(`[%s] Cleaning bootstrap config for environment ID "%s"..`, configKey, c.Environment.EnvironmentID)

	if len(c.BootstrapVerifyPolicyNames) == 0 {
		l.Info().Msgf("[%s] No bootstrap names configured - skipping", configKey)
		return nil
	}

	ok, err := clean.CheckBillOfMaterials(ctx, configKey, c.Environment, management.ENUMPRODUCTTYPE_ONE_VERIFY)
	if err != nil {
		return err
	}

	if !ok {
		l.Info().Msgf("[%s] Bill of materials does not contain applicable service %s - skipping", configKey, management.ENUMPRODUCTTYPE_ONE_VERIFY)
		return nil
	}

	var response *verify.EntityArray
	err = clean.ReadAllConfig(
		ctx,
		configKey,
		c.Environment,
		func() (any, *http.Response, error) {
			return c.Environment.Client.VerifyAPIClient.VerifyPoliciesApi.ReadAllVerifyPolicies(ctx, c.Environment.EnvironmentID).Execute()
		},
		&response,
	)
	if err != nil {
		return err
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasVerifyPolicies() {

		l.Debug().Msgf("[%s] Configuration items found, looping..", configKey)
		for _, policy := range embedded.GetVerifyPolicies() {

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
					IdentifierListToSearch: c.BootstrapVerifyPolicyNames,
					StartsWithStringMatch:  false,
				},
				func() (any, *http.Response, error) {
					fR, fErr := c.Environment.Client.VerifyAPIClient.VerifyPoliciesApi.DeleteVerifyPolicy(ctx, c.Environment.EnvironmentID, policy.GetId()).Execute()
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
