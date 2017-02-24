package tpl

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/drud/drud-go/utils/system"
	"github.com/oleiade/reflections"
)

// WordpressConfig encapsulates all the configurations for a Wordpress site.
type WordpressConfig struct {
	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string
	DatabaseHost     string
	DatabaseDriver   string
	DatabasePort     int
	DatabasePrefix   string
	AuthKey          string
	SecureAuthKey    string
	LoggedInKey      string
	NonceKey         string
	AuthSalt         string
	SecureAuthSalt   string
	LoggedInSalt     string
	NonceSalt        string
	Docroot          string
	SiteURL          string
	CoreDir          string
	ContentDir       string
	UploadDir        string
	FileSrc          string
}

// DefaultWordpressConfig produces a WordpressConfig object with defaults.
func DefaultWordpressConfig() *WordpressConfig {
	return &WordpressConfig{
		ContentDir:       "wp-content",
		CoreDir:          "",
		DatabaseName:     "data",
		DatabaseUsername: "root",
		DatabasePassword: "root",
		DatabaseHost:     "127.0.0.1",
		DatabaseDriver:   "mysql",
		DatabasePort:     3306,
		DatabasePrefix:   "wp_",
		SiteURL:          os.Getenv("DEPLOY_URL"),
		FileSrc:          os.Getenv("FILE_SRC"),
	}
}

// WriteAppConfig produces a valid settings.php file from the defined configurations
func (c *WordpressConfig) WriteAppConfig(in *Config) error {
	srcFieldList, err := reflections.Items(in)
	if err != nil {
		return err
	}

	for field, val := range srcFieldList {
		if val != "" {
			has, err := reflections.HasField(c, field)
			if err != nil {
				return err
			}
			if has && val != 0 {
				err = reflections.SetField(c, field, val)
				if err != nil {
					return err
				}
			}
		}
	}

	if c.CoreDir != "" {
		c.CoreDir = SlashIt(in.CoreDir, false)
	}

	if in.UploadDir == "" {
		c.UploadDir = c.ContentDir + "/uploads"
	}

	// If there's no salt ask for some.
	saltKeys := []string{"AuthKey", "SecureAuthKey", "LoggedInKey", "NonceKey", "AuthSalt", "SecureAuthSalt", "LoggedInSalt", "NonceSalt"}
	for _, key := range saltKeys {
		val, err := reflections.GetField(c, key)
		if err != nil {
			return err
		}
		if val == "" {
			err = reflections.SetField(c, key, PassTheSalt())
			if err != nil {
				return err
			}
		}
	}

	tmpl, err := template.New("conf").Funcs(sprig.TxtFuncMap()).Parse(wordpressTemplate)
	if err != nil {
		return err
	}

	filepath := ""
	if in.ConfigPath != "" {
		filepath = SlashIt(in.ConfigPath, true)
	}

	file, err := os.Create(filepath + "wp-config.php")
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, c)
	if err != nil {
		return err
	}

	return nil
}

// PlaceFiles determines where file upload directories should go.
func (c *WordpressConfig) PlaceFiles(move bool) error {
	src := c.FileSrc
	dest := c.UploadDir

	if !system.FileExists(src) {
		log.Fatalf("source path for files does not exist")
	}

	if system.FileExists(dest) {
		log.Printf("destination path exists, removing")
		os.RemoveAll(dest)
	}

	// ensure parent of destination is writable
	os.Chmod(path.Dir(src), 0755)

	if move {
		out, err := system.RunCommand(
			"rsync",
			[]string{"-avz", "--recursive", src + "/", dest},
		)
		if err != nil {
			return fmt.Errorf("%s - %s", err.Error(), string(out))
		}
	} else {
		err := os.Symlink(src, dest)
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteWebConfig updates the web server configuration to support the provided app configurations
// @TODO: need to update rules for other WP concerns, holding off until more firm on approach for this task.
func (c *WordpressConfig) WriteWebConfig() error {
	dest := os.Getenv("NGINX_SITE_CONF")
	root := "root /var/www/html"

	if !system.FileExists(dest) {
		log.Fatalf("target %s does not exist", dest)
	}

	conf, err := ioutil.ReadFile(dest)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(root)
	new := re.ReplaceAllString(string(conf), root+"/"+c.Docroot)

	err = ioutil.WriteFile(dest, []byte(new), 0644)
	if err != nil {
		return err
	}

	return nil
}
