package clean

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone"
	"github.com/patrickcping/pingone-sweep/internal/logger"
	"github.com/patrickcping/pingone-sweep/internal/sdk"
)

type CleanEnvironmentConfig struct {
	EnvironmentID string
	DryRun        bool
	Client        *pingone.Client
}

type ConfigItem struct {
	IdentifierToEvaluate string
	Id                   string
	Default              *bool
	Enabled              *bool
}

type ConfigItemEval struct {
	IdentifierListToSearch []string
	StartsWithStringMatch  bool
	CaseSensitive          *bool
}

func BillOfMaterialsHasService(ctx context.Context, configKey string, env CleanEnvironmentConfig, productType management.EnumProductType) (bool, error) {
	l := logger.Get()

	var response *management.BillOfMaterials
	err := sdk.ParseResponse(
		ctx,
		func() (any, *http.Response, error) {
			return env.Client.ManagementAPIClient.BillOfMaterialsBOMApi.ReadOneBillOfMaterials(ctx, env.EnvironmentID).Execute()
		},
		fmt.Sprintf("[%s]-CHECKBOM", configKey),
		sdk.DefaultCreateReadRetryable,
		&response,
	)

	if err != nil {
		return false, err
	}

	if response == nil {
		return false, fmt.Errorf("[%s] No bill of materials found - the API responded with no data", configKey)
	}

	for _, product := range response.GetProducts() {
		if product.GetType() == productType {
			l.Debug().Msgf(`[%s] Found product "%s" in the bill of materials`, configKey, string(productType))
			return true, nil
		}
	}

	l.Debug().Msgf(`[%s] Product "%s" cannot be found in the target environment's bill of materials`, configKey, string(productType))
	return false, nil
}

func ReadAllConfig(ctx context.Context, configKey string, env CleanEnvironmentConfig, readAllSdkFunction sdk.SDKInterfaceFunc, targetObject any) error {

	err := sdk.ParseResponse(
		ctx,
		readAllSdkFunction,
		fmt.Sprintf("[%s]-READALL", configKey),
		sdk.DefaultCreateReadRetryable,
		targetObject,
	)

	if err != nil {
		return err
	}

	if targetObject == nil {
		return fmt.Errorf("[%s] No configuration items found - the API responded with no data", configKey)
	}

	return nil
}

func TryCleanConfig(ctx context.Context, configKey string, env CleanEnvironmentConfig, configItem ConfigItem, configItemEval ConfigItemEval, deleteSdkFunction sdk.SDKInterfaceFunc, disableSdkFunction sdk.SDKInterfaceFunc) error {
	l := logger.Get()

	if disableSdkFunction == nil && deleteSdkFunction == nil {
		return fmt.Errorf("[%s] No SDK functions provided", configKey)
	}

	var debugAction string
	var sdkActionFunc sdk.SDKInterfaceFunc
	if deleteSdkFunction != nil {
		debugAction = "DELETE"
		sdkActionFunc = deleteSdkFunction
	}
	if disableSdkFunction != nil {
		debugAction = "DISABLE"
		sdkActionFunc = disableSdkFunction
	}

	l.Debug().Msgf(`[%s] Looping configured list of identifiers for "%s" for action %s..`, configKey, configItem.IdentifierToEvaluate, debugAction)
	for _, identifierToSearch := range configItemEval.IdentifierListToSearch {

		var eqExprResult bool
		if configItemEval.CaseSensitive != nil && *configItemEval.CaseSensitive {
			eqExprResult = configItem.IdentifierToEvaluate == identifierToSearch
		} else {
			eqExprResult = strings.EqualFold(configItem.IdentifierToEvaluate, identifierToSearch)
		}

		if eqExprResult {
			l.Debug().Msgf(`[%s] Found "%s"`, configKey, identifierToSearch)

			if configItem.Default != nil && *configItem.Default {
				l.Warn().Msgf(`[%s] "%s" is set as the environment default - no action will be taken`, configKey, configItem.IdentifierToEvaluate)

				break
			}

			if configItem.Enabled != nil && !*configItem.Enabled && disableSdkFunction != nil {
				l.Info().Msgf(`[%s] "%s" is already disabled - no action taken`, configKey, configItem.IdentifierToEvaluate)

				break
			}

			if !env.DryRun {

				err := sdk.ParseResponse(
					ctx,
					sdkActionFunc,
					fmt.Sprintf("[%s]-%s", configKey, debugAction),
					sdk.DefaultCreateReadRetryable,
					nil,
				)

				if err != nil {
					return err
				}
				l.Info().Msgf(`[%s] %s action completed for "%s"`, configKey, debugAction, configItem.IdentifierToEvaluate)
			} else {
				l.Warn().Msgf(`[%s] Dry run: %s action "%s" with ID "%s"`, configKey, debugAction, configItem.IdentifierToEvaluate, configItem.Id)
			}

			break
		}
	}

	return nil
}
