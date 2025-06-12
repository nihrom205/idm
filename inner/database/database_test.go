package database

import (
	"github.com/nihrom205/idm/inner/common"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	configName = ".env_test"
)

func TestConnectWithCfgWhenCorrectDSNThenReturnConnect(t *testing.T) {
	assert := assert.New(t)
	os.Create(configName)
	defer os.Remove(configName)

	os.WriteFile(configName, []byte("DB_DRIVER_NAME=postgres\nDB_DSN='host=localhost port=5432 user=user password=user dbname=idm_db sslmode=disable'\n"), 0644)

	cfg := common.GetConfig(configName)
	con := ConnectDbWithCfg(cfg)

	assert.NotNil(con)
}

func TestConnectWithCfgWhenNotCorrectDSNThenReturn(t *testing.T) {
	assert := assert.New(t)
	os.Create(configName)
	defer os.Remove(configName)

	os.WriteFile(configName, []byte("DB_DRIVER_NAME=postgres\nDB_DSN='host=localhost port=54321 user=user password=user dbname=idm_db sslmode=disable'\n"), 0644)

	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(r)
		}
	}()
	cfg := common.GetConfig(configName)
	con := ConnectDbWithCfg(cfg)
	defer con.Close()
}
