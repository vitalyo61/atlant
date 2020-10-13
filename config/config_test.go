package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "config.*.yaml")
	assert.NoError(t, err)

	defer os.Remove(tmpfile.Name())

	address := "mongodb://localhost:27017"
	name := "atlant"
	var timeout int64 = 5

	err = os.Setenv("DB_TIMEOUT", "5")
	assert.NoError(t, err)

	yaml := []byte(fmt.Sprintf(`
database:
  address: %s
  database: %s
  timeout: %d
`, address, name, 15))

	_, err = tmpfile.Write(yaml)
	assert.NoError(t, err)

	tmpfile.Close()

	cfg, err := Get(tmpfile.Name())
	assert.NoError(t, err)
	assert.Equal(t, cfg.DB.Address, address)
	assert.Equal(t, cfg.DB.Name, name)
	assert.Equal(t, cfg.DB.Timeout, timeout)
}
