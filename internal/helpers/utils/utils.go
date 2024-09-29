package utils

import (
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaulEtcdTimeout = 5
)

var (
	releaseVersion = "development"
)

// GetReleaseVersion returns the release version of the application.
func GetReleaseVersion() string {
	return releaseVersion
}

// ReadEnvVar reads the value of the environment variable with the given name.
//   - If the environment variable is not set, the default value is returned.
func ReadEnvVar(_envName, _defaultVal string) string {
	if val := os.Getenv(_envName); val != "" {
		return val
	}

	return _defaultVal
}

// ErrorMessage returns the error message as a string.
//   - If the error is nil, an empty string is returned.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}

func GetEntryPoints() []string {
	envVal := os.Getenv("ETCD_NODES")

	return strings.Split(envVal, ",")
}

func GetUsername() string {
	return os.Getenv("ETCD_ADMIN_USER")
}

func GetUIUsername() string {
	return os.Getenv("UI_USER")
}

func GetUIPassword() string {
	return os.Getenv("UI_PASSWORD")
}

func GetPassword() string {
	return os.Getenv("ETCD_ADMIN_PASSWORD")
}

func GetSSLCAFile() string {
	return os.Getenv("SSL_CA")
}

func GetSSLCertFile() string {
	return os.Getenv("SSL_CERT")
}

func GetSSLKeyFile() string {
	return os.Getenv("SSL_KEY")
}

func TLSSkipVerify() bool {
	skipVerify := os.Getenv("TLS_SKIP_VERIFY")

	return skipVerify == "1" || skipVerify == "true"
}

func GetOpTimeout() time.Duration {
	timeoutVal := ReadEnvVar("ETCD_TIMEOUT", "5")
	parsedTimeout, parseErr := strconv.ParseInt(timeoutVal, 10, 8)

	if parseErr != nil {
		log.Printf("could not parse: %q is not a proper timeout value\n", timeoutVal)

		parsedTimeout = defaulEtcdTimeout
	}

	if parsedTimeout < 1 {
		parsedTimeout = defaulEtcdTimeout
	}

	return time.Second * time.Duration(parsedTimeout)
}

func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func Base64Decode(input string) string {
	out, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return ""
	}

	return string(out)
}
