package main

import (
	"context"
	"fmt"
	"github.com/nihrom205/idm/inner/common"
	validator2 "github.com/nihrom205/idm/inner/common/validator"
	database2 "github.com/nihrom205/idm/inner/database"
	"github.com/nihrom205/idm/inner/employee"
	"github.com/nihrom205/idm/inner/info"
	"github.com/nihrom205/idm/inner/role"
	"github.com/nihrom205/idm/inner/web"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	server := build()
	go func() {
		if err := server.App.Listen(":8080"); err != nil {
			panic(fmt.Sprintf("error starting server: %v", err))
		}
	}()

	// Создаем группу для ожидания сигнала завершения работы сервера
	wg := &sync.WaitGroup{}
	wg.Add(1)
	// Запускаем gracefulShutdown в отдельной горутине
	go gracefulShutdown(server, wg)
	// Ожидаем сигнал от горутины gracefulShutdown, что сервер завершил работу
	wg.Wait()
	fmt.Println("Graceful shutdown complete.")
}

func gracefulShutdown(server *web.Server, wg *sync.WaitGroup) {
	// Уведомить основную горутину о завершении работы
	defer wg.Done()
	// Создаём контекст, который слушает сигналы прерывания от операционной системы
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer stop()
	// Слушаем сигнал прерывания от операционной системы
	<-ctx.Done()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")
	// Контекст используется для информирования веб-сервера о том,
	// что у него есть 5 секунд на выполнение запроса, который он обрабатывает в данный момент
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.App.ShutdownWithContext(ctx); err != nil {
		fmt.Printf("Server forced to shutdown with error: %v\n", err)
	}
	fmt.Println("Server exiting")
}

func build() *web.Server {
	// читаем конфиги
	cfg := common.GetConfig(".env")

	// Создаём подключение к базе данных
	db := database2.ConnectDbWithCfg(cfg)

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
