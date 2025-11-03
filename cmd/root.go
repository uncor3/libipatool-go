package cmd

import (
	"github.com/majd/ipatool/v2/pkg/appstore"
	ipatoolhttp "github.com/majd/ipatool/v2/pkg/http"
	"github.com/majd/ipatool/v2/pkg/keychain"
	"github.com/majd/ipatool/v2/pkg/util/machine"
	"github.com/majd/ipatool/v2/pkg/util/operatingsystem"
)

var dependencies = Dependencies{}

type Dependencies struct {
	OS        operatingsystem.OperatingSystem
	Machine   machine.Machine
	CookieJar ipatoolhttp.CookieJar
	Keychain  keychain.Keychain
	AppStore  appstore.AppStore
}

func Initialize() error {
	verbose := false
	nonInteractive := true
	format := OutputFormatJSON

	initDependencies(verbose, nonInteractive, format)
	return nil
}

// TODO: implement
func Cleanup() {
	dependencies = Dependencies{}
}

func initDependencies(verbose, nonInteractive bool, format OutputFormat) {
	dependencies.OS = operatingsystem.New()
	dependencies.Machine = machine.New(machine.Args{OS: dependencies.OS})
	dependencies.CookieJar = newCookieJar(dependencies.Machine)
	dependencies.Keychain = newKeychain(dependencies.Machine)
	dependencies.AppStore = appstore.NewAppStore(appstore.Args{
		CookieJar:       dependencies.CookieJar,
		OperatingSystem: dependencies.OS,
		Keychain:        dependencies.Keychain,
		Machine:         dependencies.Machine,
	})

	createConfigDirectory(dependencies.OS, dependencies.Machine)
}
