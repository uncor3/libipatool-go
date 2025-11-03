package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/99designs/keyring"
	cookiejar "github.com/juju/persistent-cookiejar"
	"github.com/majd/ipatool/v2/pkg/http"
	"github.com/majd/ipatool/v2/pkg/keychain"
	"github.com/majd/ipatool/v2/pkg/util"
	"github.com/majd/ipatool/v2/pkg/util/machine"
	"github.com/majd/ipatool/v2/pkg/util/operatingsystem"
)

var KeychainPassphrase = ""

func newCookieJar(machine machine.Machine) http.CookieJar {
	return util.Must(cookiejar.New(&cookiejar.Options{
		Filename: filepath.Join(machine.HomeDirectory(), ConfigDirectoryName, CookieJarFileName),
	}))
}

func newKeychain(machine machine.Machine, enabledBackends []string) keychain.Keychain {

	for _, backend := range enabledBackends {
		println("Enabled keychain backend: %s\n", backend)
	}

	allowedBackends := []keyring.BackendType{}
	for _, backend := range enabledBackends {
		switch backend {
		case "keychain":
			allowedBackends = append(allowedBackends, keyring.KeychainBackend)
		case "secret-service":
			allowedBackends = append(allowedBackends, keyring.SecretServiceBackend)
		case "wincred":
			allowedBackends = append(allowedBackends, keyring.WinCredBackend)
		case "file":
			allowedBackends = append(allowedBackends, keyring.FileBackend)
		}
	}

	if len(allowedBackends) == 0 {
		allowedBackends = []keyring.BackendType{
			keyring.FileBackend,
		}
	}

	ring := util.Must(keyring.Open(keyring.Config{
		AllowedBackends: allowedBackends,
		ServiceName:     KeychainServiceName,
		FileDir:         filepath.Join(machine.HomeDirectory(), ConfigDirectoryName),
		FilePasswordFunc: func(s string) (string, error) {
			return KeychainPassphrase, nil
		},
	}))

	return keychain.New(keychain.Args{Keyring: ring})
}

func createConfigDirectory(os operatingsystem.OperatingSystem, machine machine.Machine) error {
	configDirectoryPath := filepath.Join(machine.HomeDirectory(), ConfigDirectoryName)
	_, err := os.Stat(configDirectoryPath)

	if err != nil && os.IsNotExist(err) {
		err = os.MkdirAll(configDirectoryPath, 0700)
		if err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("could not read metadata: %w", err)
	}

	return nil
}
