package main

import (
	_ "avito2/docs"
	"avito2/pkg/handler"
	"avito2/pkg/repository"
	"avito2/pkg/repository/postgres"
	"avito2/pkg/service"
	_scheduler "avito2/pkg/service/scheduler"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

// @title			Avito Intership Task
// @version			1.0
// @description		Avito Segments
// @contact.name	Gvozdoslav
// @host			localhost:8080
func main() {
	logrus.SetFormatter(new(logrus.TextFormatter))
	logrus.SetLevel(logrus.DebugLevel)

	if err := InitConfig(); err != nil {
		logrus.Fatalf("error while initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading environment variables: %s", err.Error())
	}

	db, err := postgres.NewPostgresDB(postgres.DbConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.name"),
	})
	if err != nil {
		logrus.Fatalf("error initializing db: %s", err.Error())
	}

	repos := repository.NewRepositories(db)
	services := service.NewServices(repos)
	handlers := handler.NewHandler(services)

	scheduler := _scheduler.NewScheduler(&repos.UserSegmentRepository)
	scheduler.ScheduleDeletion()

	server := new(Server)
	err = server.Run(viper.GetString("port"), handlers.InitRoutes())
	if err != nil {
		logrus.Fatalf("error while running http server: %s", err.Error())
	}
}

func InitConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
