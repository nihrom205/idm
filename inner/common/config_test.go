package common

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	configName        = ".env_test"
	db_dsn            = "DB_DSN"
	dsn               = "host=test_url port=5432"
	db_driver_name    = "DB_DRIVER_NAME"
	db_driver         = "postgres"
	app_name          = "APP_NAME"
	app_name_value    = "test_idm"
	app_version       = "APP_VERSION"
	app_version_value = "0.0.1"
)

func TestGetConfigWhenNotFoundEnvReturnEnvironmentVariable(t *testing.T) {
	assert := assert.New(t)
	err := os.Setenv(db_dsn, dsn)
	if err != nil {
		return
	}
	err = os.Setenv(db_driver_name, db_driver)
	if err != nil {
		return
	}
	err = os.Setenv(app_name, app_name_value)
	if err != nil {
		return
	}
	err = os.Setenv(app_version, app_version_value)
	if err != nil {
		return
	}
	err = os.Setenv("LOG_LEVEL", "INFO")
	if err != nil {
		return
	}
	err = os.Setenv("LOG_DEVELOP_MODE", "true")
	if err != nil {
		return
	}
	err = os.Setenv("SSL_CERT", "test_cert")
	if err != nil {
		return
	}
	err = os.Setenv("SSL_KEY", "test_key")
	if err != nil {
		return
	}
	err = os.Setenv("REDIS_ADDR", "url")
	if err != nil {
		return
	}
	defer func() {
		err := os.Unsetenv(db_dsn)
		if err != nil {
			return
		}
	}()
	defer func() {
		err := os.Unsetenv(db_driver_name)
		if err != nil {
			return
		}
	}()

	got := GetConfig("fakeFile")

	assert.NotNil(got)
	assert.Equal(got.DSN, dsn)
	assert.Equal(got.DbDriverName, db_driver)
}

func TestGetConfigWhenEnvFileNotValuesReturnEmptyStructure(t *testing.T) {
	assert := assert.New(t)
	err := os.Setenv("VAR_ONE", "VAR_ONE_VALUE")
	if err != nil {
		return
	}
	err = os.Setenv("VAR_TWO", "VAR_TWO_VALUE")
	if err != nil {
		return
	}
	err = os.Setenv("VAR_THREE", "VAR_THREE_VALUE")
	if err != nil {
		return
	}
	err = os.Setenv("VAR_FOUR", "VAR_FOUR_VALUE")
	if err != nil {
		return
	}
	defer func() {
		err = os.Unsetenv("VAR_ONE")
		if err != nil {
			return
		}
	}()
	defer func() {
		err = os.Unsetenv("VAR_TWO")
		if err != nil {
			return
		}
	}()
	defer func() {
		err = os.Unsetenv("VAR_THREE")
		if err != nil {
			return
		}
	}()
	defer func() {
		err = os.Unsetenv("VAR_FOUR")
		if err != nil {
			return
		}
	}()

	_, err = os.Create(configName)
	if err != nil {
		return
	}
	defer func() {
		err := os.Remove(configName)
		if err != nil {
			return
		}
	}()

	err = os.WriteFile(configName, []byte("PARAM"), 0644)
	if err != nil {
		return
	}

	//_ = GetConfig(configName)
	assert.Panics(func() {
		GetConfig(configName)
	})
	//assert.PanicsWithValue(t, func() {
	//
	//	GetConfig(configName)
	//})

	//assert.NotNil(got)
	//assert.Equal(got.DSN, "")
	//assert.Equal(got.DbDriverName, "")
}

func TestGetConfigWhenEnvFileNotValuesAndCorrectVariableReturnStructure(t *testing.T) {
	assert := assert.New(t)
	err := os.Setenv(db_dsn, dsn)
	if err != nil {
		return
	}
	err = os.Setenv(db_driver_name, db_driver)
	if err != nil {
		return
	}
	defer func() {
		err = os.Unsetenv(db_dsn)
		if err != nil {
			return
		}
	}()
	defer func() {
		err = os.Unsetenv(db_driver_name)
		if err != nil {
			return
		}
	}()

	_, err = os.Create(configName)
	if err != nil {
		return
	}
	defer func() {
		err = os.Remove(configName)
		if err != nil {
			return
		}
	}()

	err = os.WriteFile(configName, []byte("PARAM"), 0644)
	if err != nil {
		return
	}

	got := GetConfig(configName)

	assert.NotNil(got)
	assert.Equal(got.DSN, dsn)
	assert.Equal(got.DbDriverName, db_driver)
}

func TestGetConfigWhenEnvFileValuesReturnStructure(t *testing.T) {
	assert := assert.New(t)
	err := os.Setenv("VAR_ONE", "VAR_ONE_VALUE")
	if err != nil {
		return
	}
	err = os.Setenv("VAR_TWO", "VAR_TWO_VALUE")
	if err != nil {
		return
	}
	defer func() {
		err = os.Unsetenv("VAR_ONE")
		if err != nil {
			return
		}
	}()
	defer func() {
		err = os.Unsetenv("VAR_TWO")
		if err != nil {
			return
		}
	}()

	_, err = os.Create(configName)
	if err != nil {
		return
	}
	defer func() {
		err = os.Remove(configName)
		if err != nil {
			return
		}
	}()

	err = os.WriteFile(configName, []byte("DB_DRIVER_NAME=postgres\nDB_DSN='host=test_url port=5432'\n"), 0644)
	if err != nil {
		return
	}

	got := GetConfig(configName)

	assert.NotNil(got)
	assert.Equal(got.DSN, dsn)
	assert.Equal(got.DbDriverName, db_driver)
}

func TestGetConfigWhenEnvFileEqualsValuesReturnStructure(t *testing.T) {
	assert := assert.New(t)
	err := os.Setenv(db_dsn, dsn)
	if err != nil {
		return
	}
	err = os.Setenv(db_driver_name, db_driver)
	if err != nil {
		return
	}
	defer func() {
		err = os.Unsetenv(db_dsn)
		if err != nil {
			return
		}
	}()
	defer func() {
		err = os.Unsetenv(db_driver_name)
		if err != nil {
			return
		}
	}()

	_, err = os.Create(configName)
	if err != nil {
		return
	}
	defer func() {
		err = os.Remove(configName)
		if err != nil {
			return
		}
	}()

	err = os.WriteFile(configName, []byte("DB_DRIVER_NAME=mysql\nDB_DSN='host=localhost port=5432'\n"), 0644)
	if err != nil {
		return
	}

	got := GetConfig(configName)

	assert.NotNil(got)
	assert.Equal(got.DSN, dsn)
	assert.Equal(got.DbDriverName, db_driver)
}
