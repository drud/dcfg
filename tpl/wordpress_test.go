package tpl

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/drud/drud-go/utils/system"
	"github.com/stretchr/testify/assert"
)

var wpconf = "wp-config.php"
var wapp = "wordpress"

func TestWordPressWriteAppConfig(t *testing.T) {
	assert := assert.New(t)

	in := new(Config)
	in.App = wapp
	in.IgnoreFiles = true
	wp := DefaultWordpressConfig()

	err := wp.WriteAppConfig(in)
	assert.NoError(err)

	content, err := ioutil.ReadFile(wpconf)
	assert.NoError(err)
	assert.Contains(string(content), "define( 'WP_CONTENT_DIR', dirname( __FILE__ ) . '/wp-content' );")
	os.Remove(wpconf)
}

func TestWordPressPlaceFiles(t *testing.T) {
	assert := assert.New(t)

	src := path.Join(os.TempDir(), "file_src")
	dest := "wp-content"
	os.Setenv("FILE_SRC", src)
	os.Mkdir(dest, 0755)
	os.MkdirAll(src, 0755)
	os.Create(path.Join(src, "testfile"))

	in := new(Config)
	in.App = wapp
	wp := DefaultWordpressConfig()
	wp.WriteAppConfig(in)
	os.Remove(wpconf)

	err := wp.PlaceFiles(false)
	assert.NoError(err)
	assert.True(system.FileExists(path.Join(dest, "uploads")))
	assert.True(system.FileExists(path.Join(dest, "uploads", "testfile")))
	link, err := os.Readlink(path.Join(dest, "uploads"))
	assert.NoError(err)
	assert.Contains(link, "file_src")
	os.Remove(wpconf)
	os.Remove(dest)

	err = wp.PlaceFiles(true)
	assert.NoError(err)
	assert.True(system.FileExists(path.Join(dest, "uploads")))
	assert.True(system.FileExists(path.Join(dest, "uploads", "testfile")))
	os.Remove(wpconf)
	os.RemoveAll(src)
	os.RemoveAll(dest)
	os.Unsetenv("FILE_SRC")
}

func TestWordPressWebConfig(t *testing.T) {
	assert := assert.New(t)

	in := new(Config)
	in.App = wapp
	in.IgnoreFiles = true
	in.Docroot = "potato"
	wp := DefaultWordpressConfig()
	wp.WriteAppConfig(in)
	os.Remove(wpconf)
	webConf := "test.conf"
	ioutil.WriteFile(webConf, []byte(testConf), os.FileMode(0644))
	os.Setenv("NGINX_SITE_CONF", webConf)

	err := wp.WriteWebConfig()
	assert.NoError(err)
	result, err := ioutil.ReadFile(webConf)
	assert.NoError(err)
	assert.Contains(string(result), "root /var/www/html/potato;")
	os.Remove(webConf)
	os.Remove(wpconf)
}
