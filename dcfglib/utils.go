package dcfglib

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"crypto/sha256"

	stringutil "github.com/drud/drud-go/utils/strings"
	"github.com/ghodss/yaml"
)

// GetTaskSetList unmarshalls config groups from the config file into structs
func GetTaskSetList(confByte []byte) (TaskSetList, error) {
	var groups TaskSetList
	jbytes, err := yaml.YAMLToJSON(confByte)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(jbytes, &groups)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// maybeChdir changes to a directory if there is one
func maybeChdir(d string) {
	if d != "" {
		err := os.Chdir(d)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// HasVars takes a guess at whether there are vars in the command string
func HasVars(command string) bool {
	if varStart := strings.Index(command, "{{"); varStart > -1 {
		if strings.Index(command[varStart:], "}}") > -1 {
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
