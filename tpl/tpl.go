package tpl

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/drud/drud-go/utils/stringutil"
)

// Tpl is the interface that each plugin must implement
type Tpl interface {
	WriteAppConfig(in *Config) error
	PlaceFiles(move bool) error
	WriteWebConfig() error
}

// TplMap is used to retrieve the correct plugin
var TplMap = map[string]Tpl{
	"drupal":    DefaultDrupalConfig(),
	"wordpress": DefaultWordpressConfig(),
}

// isValidApp determines if a given app matches one of the defined plugins.
func isValidApp(app string) bool {
	for valApp := range TplMap {
		if app == valApp {
			return true
		}
	}
	return false
}

// PassTheSalt generates a hash salt
func PassTheSalt() string {
	salt := sha256.New()
	random := stringutil.RandomString(20)
	salt.Write([]byte(random))

	return hex.EncodeToString(salt.Sum(nil))
}

// SlashIt ensures you have a preceding or trailing slash on a string
func SlashIt(val string, trailing bool) string {
	if trailing && !strings.HasSuffix(val, "/") {
		return val + "/"
	} else if !trailing && !strings.HasPrefix(val, "/") {
		return "/" + val
	}
	return val
}
