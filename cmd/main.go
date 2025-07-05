package main

import (
	"context"
	"github.com/nihrom205/idm/inner/common"
	validator2 "github.com/nihrom205/idm/inner/common/validator"
	database2 "github.com/nihrom205/idm/inner/database"
	"github.com/nihrom205/idm/inner/employee"
	"github.com/nihrom205/idm/inner/info"
	"github.com/nihrom205/idm/inner/role"
	"github.com/nihrom205/idm/inner/web"
	"go.uber.org/zap"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	// читаем конфиги
	cfg := common.GetConfig(".env")
	// Создаем логгер
	logger := common.NewLogger(cfg)
	// Отложенный вызов записи сообщений из буфера в лог. Необходимо вызывать перед выходом из приложения
	defer func() { _ = logger.Sync() }()

	server := build(cfg, logger)
	go func() {
		if err := server.App.Listen(":8080"); err != nil {
			logger.Panic("error starting server", zap.Error(err))
		}
	}()

	// Создаем группу для ожидания сигнала завершения работы сервера
	wg := &sync.WaitGroup{}
	wg.Add(1)
	// Запускаем gracefulShutdown в отдельной горутине
	go gracefulShutdown(server, wg, logger)
	// Ожидаем сигнал от горутины gracefulShutdown, что сервер завершил работу
	wg.Wait()
	logger.Info("Graceful shutdown complete.")
}

func gracefulShutdown(server *web.Server, wg *sync.WaitGroup, logger *common.Logger) {
	// Уведомить основную горутину о завершении работы
	defer wg.Done()
	// Создаём контекст, который слушает сигналы прерывания от операционной системы
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	defer stop()
	// Слушаем сигнал прерывания от операционной системы
	<-ctx.Done()
	logger.Info("shutting down gracefully")
	// Контекст используется для информирования веб-сервера о том,
	// что у него есть 5 секунд на выполнение запроса, который он обрабатывает в данный момент
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.App.ShutdownWithContext(ctx); err != nil {
		logger.Error("Server forced to shutdown with error", zap.Error(err))
	}
	logger.Info("Server exiting")
}

func build(cfg common.Config, logger *common.Logger) *web.Server {

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
	employeeController := employee.NewController(server, employeeService, logger)
	employeeController.RegisterRoutes()

	// создаём контроллер role
	roleController := role.NewController(server, roleService, logger)
	roleController.RegisterRoutes()

	// создаём контроллер info
	infoController := info.NewController(server, cfg, db)
	infoController.RegisterRouters()

	return server
}
