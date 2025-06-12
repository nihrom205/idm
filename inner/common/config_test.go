package common

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	configName     = ".env_test"
	db_dsn         = "DB_DSN"
	dsn            = "host=test_url port=5432"
	db_driver_name = "DB_DRIVER_NAME"
	db_driver      = "postgres"
)

func TestGetConfigWhenNotFoundEnvReturnEnvironmentVariable(t *testing.T) {
	assert := assert.New(t)
	os.Setenv(db_dsn, dsn)
	os.Setenv(db_driver_name, db_driver)
	defer os.Unsetenv(db_dsn)
	defer os.Unsetenv(db_driver_name)

	got := GetConfig("fakeFile")

	assert.NotNil(got)
	assert.Equal(got.DSN, dsn)
	assert.Equal(got.DbDriverName, db_driver)
}

func TestGetConfigWhenEnvFileNotValuesReturnEmptyStructure(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("VAR_ONE", "VAR_ONE_VALUE")
	os.Setenv("VAR_TWO", "VAR_TWO_VALUE")
	defer os.Unsetenv("VAR_ONE")
	defer os.Unsetenv("VAR_TWO")

	os.Create(configName)
	defer os.Remove(configName)

	os.WriteFile(configName, []byte("PARAM"), 0644)

	got := GetConfig(configName)

	assert.NotNil(got)
	assert.Equal(got.DSN, "")
	assert.Equal(got.DbDriverName, "")
}

func TestGetConfigWhenEnvFileNotValuesAndCorrectVariableReturnStructure(t *testing.T) {
	assert := assert.New(t)
	os.Setenv(db_dsn, dsn)
	os.Setenv(db_driver_name, db_driver)
	defer os.Unsetenv(db_dsn)
	defer os.Unsetenv(db_driver_name)

	os.Create(configName)
	defer os.Remove(configName)

	os.WriteFile(configName, []byte("PARAM"), 0644)

	got := GetConfig(configName)

	assert.NotNil(got)
	assert.Equal(got.DSN, dsn)
	assert.Equal(got.DbDriverName, db_driver)
}

func TestGetConfigWhenEnvFileValuesReturnStructure(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("VAR_ONE", "VAR_ONE_VALUE")
	os.Setenv("VAR_TWO", "VAR_TWO_VALUE")
	defer os.Unsetenv("VAR_ONE")
	defer os.Unsetenv("VAR_TWO")

	os.Create(configName)
	defer os.Remove(configName)

	os.WriteFile(configName, []byte("DB_DRIVER_NAME=postgres\nDB_DSN='host=test_url port=5432'\n"), 0644)

	got := GetConfig(configName)

	assert.NotNil(got)
	assert.Equal(got.DSN, dsn)
	assert.Equal(got.DbDriverName, db_driver)
}

func TestGetConfigWhenEnvFileEqualsValuesReturnStructure(t *testing.T) {
	assert := assert.New(t)
	os.Setenv(db_dsn, dsn)
	os.Setenv(db_driver_name, db_driver)
	defer os.Unsetenv(db_dsn)
	defer os.Unsetenv(db_driver_name)

	os.Create(configName)
	defer os.Remove(configName)

	os.WriteFile(configName, []byte("DB_DRIVER_NAME=mysql\nDB_DSN='host=localhost port=5432'\n"), 0644)

	got := GetConfig(configName)

	assert.NotNil(got)
	assert.Equal(got.DSN, dsn)
	assert.Equal(got.DbDriverName, db_driver)
}
