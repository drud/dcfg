package tpl

import (
	"log"
	"os"

	"github.com/drud/drud-go/utils/pretty"
)

// Config implements the Template Action
type Config struct {
	App              string `yaml:"app"`
	Core             string `yaml:"core"`
	ConfigPath       string `yaml:"configPath"`
	Docroot          string `yaml:"docroot"`
	DatabaseName     string `yaml:"databaseName"`
	DatabaseUsername string `yaml:"databaseUsername"`
	DatabasePassword string `yaml:"databasePassword"`
	DatabaseHost     string `yaml:"databaseHost"`
	DatabaseDriver   string `yaml:"databaseDriver"`
	DatabasePort     int    `yaml:"databasePort"`
	DatabasePrefix   string `yaml:"databasePrefix"`
	IgnoreFiles      bool   `yaml:"ignoreFiles"`
	PublicFiles      string `yaml:"publicFiles"`
	PrivateFiles     string `yaml:"privateFiles"`
	ConfigSyncDir    string `yaml:"configSyncDir"`
	SiteURL          string `yaml:"siteURL"`
	CoreDir          string `yaml:"coreDir"`
	ContentDir       string `yaml:"contentDir"`
	UploadDir        string `yaml:"uploadDir"`
	FileSrc          string `yaml:"fileSrc"`
}

// String prints the Task
func (c Config) String() string {
	return pretty.Prettify(c)
}

// Run creates configurations for an application
func (c *Config) Run() error {
	if !isValidApp(c.App) {
		log.Fatalf("'%s' is not a valid app type", c.App)
	}

	app := TplMap[c.App]

	err := app.WriteAppConfig(c)
	if err != nil {
		return err
	}

	if !c.IgnoreFiles {
		if os.Getenv("DEPLOY_NAME") == "local" {
			app.PlaceFiles(true)
		} else {
			app.PlaceFiles(false)
		}
	}

	if c.Docroot != "" {
		app.WriteWebConfig()
	}

	return nil
}
