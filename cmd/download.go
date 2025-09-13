package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/majd/ipatool/v2/pkg/appstore"
)

func DownloadApp(bundleID, outputPath, externalVersionID string, acquireLicense bool, progressCallback appstore.ProgressCallback) error {
	if bundleID == "" {
		return errors.New("bundle identifier must be specified")
	}

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
		// app := appstore.App{BundleID:appID}

		lookupResult, err := dependencies.AppStore.Lookup(appstore.LookupInput{Account: acc, BundleID: bundleID})
		if err != nil {
			return err
		}

		app := lookupResult.App
		fmt.Println(app.Name)
		fmt.Println(outputPath)

		if errors.Is(lastErr, appstore.ErrLicenseRequired) && acquireLicense {
			err := dependencies.AppStore.Purchase(appstore.PurchaseInput{Account: acc, App: app})
			if err != nil {
				return err
			}
		}

		// TODO: progress
		out, err := dependencies.AppStore.Download(appstore.DownloadInput{
			Account: acc, App: app, OutputPath: outputPath, Progress: progressCallback, ExternalVersionID: externalVersionID})
		if err != nil {
			fmt.Println("Error in dependencies.AppStore.Download", err)
			return err
		}

		err = dependencies.AppStore.ReplicateSinf(appstore.ReplicateSinfInput{Sinfs: out.Sinfs, PackagePath: out.DestinationPath})
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

			if errors.Is(err, appstore.ErrPasswordTokenExpired) {
				return true
			}

			if errors.Is(err, appstore.ErrLicenseRequired) && acquireLicense {
				return true
			}

			return false
		}),
	)
}
