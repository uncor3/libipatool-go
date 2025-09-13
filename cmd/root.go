package cmd

import (
	"log"

	"github.com/majd/ipatool/v2/pkg/appstore"
	ipatoolhttp "github.com/majd/ipatool/v2/pkg/http"
	"github.com/majd/ipatool/v2/pkg/keychain"
	"github.com/majd/ipatool/v2/pkg/util/machine"
	"github.com/majd/ipatool/v2/pkg/util/operatingsystem"
)

var version = "dev"
var dependencies = Dependencies{}
var keychainPassphrase string

type Dependencies struct {
	Logger    log.Logger
	OS        operatingsystem.OperatingSystem
	Machine   machine.Machine
	CookieJar ipatoolhttp.CookieJar
	Keychain  keychain.Keychain
	AppStore  appstore.AppStore
}

// Initialize initializes the library dependencies
func Initialize() error {
	verbose := false
	nonInteractive := true
	format := OutputFormatJSON

	initDependencies(verbose, nonInteractive, format)
	return nil
}

// Cleanup cleans up library resources
func Cleanup() {
	// Clean up any resources if needed
	dependencies = Dependencies{}
}

func initDependencies(verbose, nonInteractive bool, format OutputFormat) {
	// dependencies.Logger = newLogger(format, verbose)
	dependencies.OS = operatingsystem.New()
	dependencies.Machine = machine.New(machine.Args{OS: dependencies.OS})
	dependencies.CookieJar = newCookieJar(dependencies.Machine)
	dependencies.Keychain = newKeychain(dependencies.Machine, nil, !nonInteractive)
	dependencies.AppStore = appstore.NewAppStore(appstore.Args{
		CookieJar:       dependencies.CookieJar,
		OperatingSystem: dependencies.OS,
		Keychain:        dependencies.Keychain,
		Machine:         dependencies.Machine,
	})

	createConfigDirectory(dependencies.OS, dependencies.Machine)
}
