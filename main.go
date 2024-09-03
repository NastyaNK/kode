package main

import (
	"log/slog"
	"myproject/pkg/api"
	"myproject/pkg/config"
	"myproject/pkg/repo"
	"myproject/pkg/service"
	"os"
)

func main() {
	// Инициализация логгера
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		//AddSource: true,
		Level: slog.LevelDebug,
	}))
	slog.SetDefault(logger)
	// Конфиг
	appConfig, err := config.LoadApplicationConfig("./application.yaml")
	if err != nil {
		slog.Error("Не удалось загрузить конфиг", slog.String("error", err.Error()))
	}

	// Подключение к базе
	connection, err := repo.NewConnection(&appConfig.Database)
	if err != nil {
		slog.Error("Не удалось подключиться к базе", slog.String("error", err.Error()))
		return
	}

	// Репозитории
	userRepo := repo.NewUserRepository(connection)
	noteRepo := repo.NewNoteRepository(connection)

	// Сервисы
	userService := service.NewUserService(userRepo)
	noteService := service.NewNoteService(noteRepo)

	// API
	kode := api.NewAPI(appConfig, userService, noteService)
	slog.Error("%v", kode.Run())
}
