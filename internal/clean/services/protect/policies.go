package protect

import (
	"context"
	"net/http"

	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/patrickcping/pingone-sweep/internal/clean"
	"github.com/patrickcping/pingone-sweep/internal/logger"
)

var (
	BootstrapRiskPolicyNames = []string{
		"Default Risk Policy",
	}
)

type CleanEnvironmentProtectRiskPoliciesConfig struct {
	Environment              clean.CleanEnvironmentConfig
	BootstrapRiskPolicyNames []string
}

func (c *CleanEnvironmentProtectRiskPoliciesConfig) Clean(ctx context.Context) error {
	l := logger.Get()

	configKey := "Risk Policies"

	l.Debug().Msgf(`[%s] Cleaning bootstrap config for environment ID "%s"..`, configKey, c.Environment.EnvironmentID)

	if len(c.BootstrapRiskPolicyNames) == 0 {
		l.Info().Msgf("[%s] No bootstrap names configured - skipping", configKey)
		return nil
	}

	ok, err := clean.CheckBillOfMaterials(ctx, configKey, c.Environment, management.ENUMPRODUCTTYPE_ONE_RISK)
	if err != nil {
		return err
	}

	if !ok {
		l.Info().Msgf("[%s] Bill of materials does not contain applicable service %s - skipping", configKey, management.ENUMPRODUCTTYPE_ONE_RISK)
		return nil
	}

	var response *risk.EntityArray
	err = clean.ReadAllConfig(
		ctx,
		configKey,
		c.Environment,
		func() (any, *http.Response, error) {
			return c.Environment.Client.RiskAPIClient.RiskPoliciesApi.ReadRiskPolicySets(ctx, c.Environment.EnvironmentID).Execute()
		},
		&response,
	)
	if err != nil {
		return err
	}

	if embedded, ok := response.GetEmbeddedOk(); ok && embedded.HasRiskPolicySets() {

		l.Debug().Msgf("[%s] Configuration items found, looping..", configKey)
		for _, policy := range embedded.GetRiskPolicySets() {

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
					IdentifierListToSearch: c.BootstrapRiskPolicyNames,
					StartsWithStringMatch:  false,
				},
				func() (any, *http.Response, error) {
					fR, fErr := c.Environment.Client.RiskAPIClient.RiskPoliciesApi.DeleteRiskPolicySet(ctx, c.Environment.EnvironmentID, policy.GetId()).Execute()
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
