package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/nihrom205/idm/inner/common"
	validator2 "github.com/nihrom205/idm/inner/common/validator"
	"github.com/nihrom205/idm/inner/database"
	"github.com/nihrom205/idm/inner/employee"
	"github.com/nihrom205/idm/inner/info"
	"github.com/nihrom205/idm/inner/role"
	"github.com/nihrom205/idm/inner/web"
)

func main() {
	cfg := common.GetConfig(".env")
	db := database.ConnectDbWithCfg(cfg)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("error closing db: %v", err)
		}
	}()

	server := build(db)
	if err := server.App.Listen(":8080"); err != nil {
		panic(fmt.Sprintf("error starting server: %v", err))
	}
}

func build(db *sqlx.DB) *web.Server {
	// читаем конфиги
	cfg := common.GetConfig(".env")

	// создаём веб-сервер
	server := web.NewServer()

	// создаём репозиторий
	employeeRepo := employee.NewEmployeeRepository(db)
	roleRepo := role.NewRoleRepository(db)

	// создаём валидатор
	vld := validator2.NewValidator()

	// создаём сервис
	employeeService := employee.NewService(employeeRepo, vld)
	roleService := role.NewService(roleRepo, vld)

	// создаём контроллер employee
	employeeController := employee.NewController(server, employeeService)
	employeeController.RegisterRoutes()

	// создаём контроллер role
	roleController := role.NewController(server, roleService)
	roleController.RegisterRoutes()

	// создаём контроллер info
	infoController := info.NewController(server, cfg, db)
	infoController.RegisterRouters()

	return server
}
