package cmd

import (
	"errors"
	"time"

	"github.com/avast/retry-go"
	"github.com/majd/ipatool/v2/pkg/appstore"
)

// nolint:wrapcheck
func purchaseCmd() error {
	var bundleID string

	var lastErr error
	var acc appstore.Account

	return retry.Do(func() error {
		infoResult, err := dependencies.AppStore.AccountInfo()
		if err != nil {
			return err
		}

		acc = infoResult.Account

		if errors.Is(lastErr, appstore.ErrPasswordTokenExpired) {
			loginResult, err := dependencies.AppStore.Login(appstore.LoginInput{Email: acc.Email, Password: acc.Password})
			if err != nil {
				return err
			}

			acc = loginResult.Account
		}

		lookupResult, err := dependencies.AppStore.Lookup(appstore.LookupInput{Account: acc, BundleID: bundleID})
		if err != nil {
			return err
		}

		err = dependencies.AppStore.Purchase(appstore.PurchaseInput{Account: acc, App: lookupResult.App})
		if err != nil {
			return err
		}

		return nil
	},
		retry.LastErrorOnly(true),
		retry.DelayType(retry.FixedDelay),
		retry.Delay(time.Millisecond),
		retry.Attempts(2),
		retry.RetryIf(func(err error) bool {
			lastErr = err

			return errors.Is(err, appstore.ErrPasswordTokenExpired)
		}),
	)

}
