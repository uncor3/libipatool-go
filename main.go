package main

/*
#include <stdlib.h>
typedef char* (*AuthCodeCallbackFunc)();
typedef void (*ProgressCallbackFunc)(long long current, long long total, void* userData);

static char* call_callback(AuthCodeCallbackFunc callback) {
    return callback();
}

static void call_progress_callback(ProgressCallbackFunc callback, long long current, long long total, void* userData) {
    if (callback) {
        callback(current, total, userData);
    }
}
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"strings"
	"unsafe"

	"github.com/majd/ipatool/v2/cmd"
)

var initialized bool

//export IpaToolInitialize
func IpaToolInitialize(backends *C.char) C.int {
	if initialized {
		return 0
	}

	goBackends := C.GoString(backends)
	enabledBackends := []string{}
	if goBackends != "" {
		enabledBackends = append(enabledBackends, splitAndTrim(goBackends, ",")...)
	}

	err := cmd.Initialize(enabledBackends)
	if err != nil {
		return -1
	}

	initialized = true
	return 0
}

//export IpaToolCleanup
func IpaToolCleanup() {
	if initialized {
		cmd.Cleanup()
		initialized = false
	}
}

//export IpaToolSearch
func IpaToolSearch(term *C.char, limit C.longlong) *C.char {
	if !initialized {
		return C.CString("Error: IpaTool not initialized")
	}

	goTerm := C.GoString(term)
	result, err := cmd.Search(goTerm, int64(limit))
	if err != nil {
		return C.CString(fmt.Sprintf("Error: %s", err.Error()))
	}

	jsonResult, _ := json.Marshal(result)
	return C.CString(string(jsonResult))
}

//export IpaToolLogin
func IpaToolLogin(email *C.char, password *C.char, authCode *C.char) C.int {
	if !initialized {
		return -1
	}

	goEmail := C.GoString(email)
	goPassword := C.GoString(password)
	goAuthCode := C.GoString(authCode)

	err := cmd.Login(goEmail, goPassword, goAuthCode, nil)
	if err != nil {
		return -1
	}

	return 0
}

//export IpaToolLoginWithCallback
func IpaToolLoginWithCallback(email *C.char, password *C.char, callback C.AuthCodeCallbackFunc) C.int {
	if !initialized {
		return -1
	}

	goEmail := C.GoString(email)
	goPassword := C.GoString(password)

	authCodeCallback := func() (string, error) {
		result := C.call_callback(callback)
		if result == nil {
			return "", fmt.Errorf("callback returned null")
		}
		return C.GoString(result), nil
	}

	err := cmd.Login(goEmail, goPassword, "", authCodeCallback)
	fmt.Println(err)
	if err != nil {
		return -1
	}

	return 0
}

//export IpaToolGetAccountInfo
func IpaToolGetAccountInfo() *C.char {
	if !initialized {
		return C.CString("Error: IpaTool not initialized")
	}

	result, err := cmd.GetAccountInfo()
	// TODO: handle error
	// returning string on error iDescriptor expect JSON
	if err != nil {
		return C.CString(fmt.Sprintf("Error: %s", err.Error()))
	}

	jsonResult, _ := json.Marshal(result)
	return C.CString(string(jsonResult))
}

//export IpaToolRevokeCredentials
func IpaToolRevokeCredentials() C.int {
	if !initialized {
		return -1
	}

	err := cmd.RevokeCredentials()
	if err != nil {
		return -1
	}

	return 0
}

//export IpaToolDownloadApp
func IpaToolDownloadApp(bundleID *C.char, outputPath *C.char, externalVersionID *C.char, acquireLicense C.int, onProgress C.ProgressCallbackFunc, userData unsafe.Pointer) C.int {
	if !initialized {
		return -1
	}

	goBundleID := C.GoString(bundleID)
	goOutputPath := C.GoString(outputPath)
	goExternalVersionID := C.GoString(externalVersionID)
	// goAcquireLicense := acquireLicense != 0

	var progressCallback func(int64, int64)
	if onProgress != nil {
		progressCallback = func(current, total int64) {
			C.call_progress_callback(onProgress, C.longlong(current), C.longlong(total), userData)
		}
	}

	err := cmd.DownloadApp(goBundleID, goOutputPath, goExternalVersionID, true, progressCallback)
	fmt.Println(err)
	if err != nil {
		return -1
	}

	return 0
}

//export SetKeyChainPassphrase
func SetKeyChainPassphrase(passphrase *C.char) {
	cmd.KeychainPassphrase = C.GoString(passphrase)
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

func main() {}
