package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/avast/retry-go"
	"github.com/majd/ipatool/v2/pkg/appstore"
)

type AccountInfo struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Success bool   `json:"success"`
}

type AuthCodeCallback func() (string, error)

func Login(email, password, authCode string, authCodeCallback AuthCodeCallback) error {
	var lastErr error

	return retry.Do(func() error {
		currentAuthCode := authCode
		fmt.Println(email, password, currentAuthCode)
		fmt.Println("In retry.Do, lastErr:", lastErr)
		if errors.Is(lastErr, appstore.ErrAuthCodeRequired) && authCodeCallback != nil {
			fmt.Println("ErrAuthCodeRequired calling authCodeCallback")
			var err error
			currentAuthCode, err = authCodeCallback()
			if err != nil {
				return err
			}
		}

		_, err := dependencies.AppStore.Login(appstore.LoginInput{
			Email:    email,
			Password: password,
			AuthCode: currentAuthCode,
		})
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
			return errors.Is(err, appstore.ErrAuthCodeRequired) && authCodeCallback != nil
		}),
	)
}

func GetAccountInfo() (*AccountInfo, error) {
	output, err := dependencies.AppStore.AccountInfo()
	if err != nil {
		return nil, err
	}

	return &AccountInfo{
		Name:    output.Account.Name,
		Email:   output.Account.Email,
		Success: true,
	}, nil
}

func RevokeCredentials() error {
	return dependencies.AppStore.Revoke()
}

// nolint:wrapcheck
func infoCmd() (AccountInfo, error) {
	output, err := dependencies.AppStore.AccountInfo()
	if err != nil {
		return AccountInfo{
			Success: false,
		}, err
	}

	return AccountInfo{
		Name:    output.Account.Name,
		Email:   output.Account.Email,
		Success: true,
	}, nil
}

// nolint:wrapcheck
func revokeCmd() bool {
	err := dependencies.AppStore.Revoke()
	if err != nil {
		println(err)
		return false
	}

	return true
}
