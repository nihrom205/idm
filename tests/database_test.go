package tests

import (
	"github.com/jmoiron/sqlx"
	"github.com/nihrom205/idm/inner/common"
	"github.com/nihrom205/idm/inner/database"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	configName = ".env_test"
)

func TestConnectWithCfgWhenCorrectDSNThenReturnConnect(t *testing.T) {
	assert := assert.New(t)
	_, err := os.Create(configName)
	if err != nil {
		return
	}
	defer func() {
		err = os.Remove(configName)
		if err != nil {
			return
		}
	}()

	err = os.WriteFile(configName, []byte(
		"DB_DRIVER_NAME=postgres\n"+
			"DB_DSN='host=localhost port=5432 user=user password=user dbname=idm_db_test sslmode=disable'\n"+
			"APP_NAME=idm_test\n"+
			"APP_VERSION=0.0.1\n"+
			"LOG_LEVEL=INFO\n"+
			"LOG_DEVELOP_MODE=true\n"+
			"SSL_CERT=test.cert\n"+
			"SSL_KEY=test.key\n"+
			"KEYCLOAK_JWK_URL=keycloak_jwt_url"), 0644)
	if err != nil {
		return
	}

	cfg := common.GetConfig(configName)
	con := database.ConnectDbWithCfg(cfg)

	assert.NotNil(con)
}

func TestConnectWithCfgWhenNotCorrectDSNThenReturn(t *testing.T) {
	assert := assert.New(t)
	_, err := os.Create(configName)
	if err != nil {
		return
	}
	defer func() {
		err = os.Remove(configName)
		if err != nil {
			return
		}
	}()

	err = os.WriteFile(configName, []byte("DB_DRIVER_NAME=postgres\nDB_DSN='host=localhost port=5432 user=user password=user dbname=idm_db_test sslmode=disable'\n"), 0644)
	if err != nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			assert.NotNil(r)
		}
	}()
	cfg := common.GetConfig(configName)
	con := database.ConnectDbWithCfg(cfg)
	defer func(con *sqlx.DB) {
		err = con.Close()
		if err != nil {
			return
		}
	}(con)
}
