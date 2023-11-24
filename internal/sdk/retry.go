package sdk

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/patrickcping/pingone-go-sdk-v2/authorize"
	"github.com/patrickcping/pingone-go-sdk-v2/credentials"
	"github.com/patrickcping/pingone-go-sdk-v2/management"
	"github.com/patrickcping/pingone-go-sdk-v2/mfa"
	"github.com/patrickcping/pingone-go-sdk-v2/pingone/model"
	"github.com/patrickcping/pingone-go-sdk-v2/risk"
	"github.com/patrickcping/pingone-go-sdk-v2/verify"
	"github.com/patrickcping/pingone-sweep/internal/logger"
)

type Retryable func(context.Context, *http.Response, *model.P1Error) bool

var (
	DefaultRetryable = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool { return false }

	DefaultCreateReadRetryable = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

		l := logger.Get()

		if p1error != nil {
			var err error

			// Permissions may not have propagated by this point
			if m, err := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
				l.Warn().Msg("Insufficient PingOne privileges detected")
				return true
			}
			if err != nil {
				l.Warn().Msg("Cannot match error string for retry")
				return false
			}

		}

		return false
	}

	RoleAssignmentRetryable = func(ctx context.Context, r *http.Response, p1error *model.P1Error) bool {

		l := logger.Get()

		if p1error != nil {
			var err error

			// Permissions may not have propagated by this point (1)
			if m, err := regexp.MatchString("^The actor attempting to perform the request is not authorized.", p1error.GetMessage()); err == nil && m {
				l.Warn().Msg("Insufficient PingOne privileges detected")
				return true
			}
			if err != nil {
				l.Warn().Msg("Cannot match error string for retry")
				return false
			}

			// Permissions may not have propagated by this point (2)
			if details, ok := p1error.GetDetailsOk(); ok && details != nil && len(details) > 0 {
				if m, err := regexp.MatchString("^Must have role at the same or broader scope", details[0].GetMessage()); err == nil && m {
					l.Warn().Msg("Insufficient PingOne privileges detected")
					return true
				}
				if err != nil {
					l.Warn().Msg("Cannot match error string for retry")
					return false
				}
			}

		}

		return false
	}
)

func RetryWrapper(ctx context.Context, timeout time.Duration, f SDKInterfaceFunc, isRetryable Retryable) (interface{}, *http.Response, error) {

	l := logger.Get()

	var resp interface{}
	var r *http.Response

	err := retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		var err error

		resp, r, err = f()

		if err != nil || r.StatusCode >= 300 {

			var errorBody []byte
			var errR error

			switch t := err.(type) {
			case *authorize.GenericOpenAPIError:
				errorBody = t.Body()
				err, errR = model.RemarshalGenericOpenAPIErrorObj(t)
			case *credentials.GenericOpenAPIError:
				errorBody = t.Body()
				err, errR = model.RemarshalGenericOpenAPIErrorObj(t)
			case *management.GenericOpenAPIError:
				errorBody = t.Body()
				err, errR = model.RemarshalGenericOpenAPIErrorObj(t)
			case *mfa.GenericOpenAPIError:
				errorBody = t.Body()
				err, errR = model.RemarshalGenericOpenAPIErrorObj(t)
			case *risk.GenericOpenAPIError:
				errorBody = t.Body()
				err, errR = model.RemarshalGenericOpenAPIErrorObj(t)
			case *verify.GenericOpenAPIError:
				errorBody = t.Body()
				err, errR = model.RemarshalGenericOpenAPIErrorObj(t)
			case *url.Error:
				l.Warn().Msgf("Detected HTTP error %s", t.Err.Error())
			default:
				l.Warn().Msgf("Detected unknown error (retry) %+v", t)
			}
			if errR != nil {
				l.Error().Msgf("Cannot remarshal type - %s", errR)
				return retry.NonRetryableError(err)
			}

			var errorModel *model.P1Error
			if len(errorBody) > 0 {
				err1 := json.Unmarshal(errorBody, &errorModel)
				if err1 != nil {
					l.Error().Msgf("Cannot remarshal service error - %s", err1)
					retry.NonRetryableError(err)
				}
			}

			if ((errorModel != nil && errorModel.Id != nil) || r != nil) && (isRetryable(ctx, r, errorModel) || DefaultRetryable(ctx, r, errorModel)) {
				l.Debug().Msgf("Retrying ... ")
				return retry.RetryableError(err)
			}

			return retry.NonRetryableError(err)

		}
		return nil
	})

	if err != nil {
		return nil, r, err
	}

	return resp, r, nil
}
