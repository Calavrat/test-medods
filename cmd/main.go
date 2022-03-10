package main

import (
	"context"

	"github.com/Calavrat/TestMedods/internal/apiserver"
	"github.com/Calavrat/TestMedods/internal/handler"
	"github.com/Calavrat/TestMedods/internal/service"
	"github.com/Calavrat/TestMedods/internal/store"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := InitConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
		return
	}

	db, err := store.NewMongoDB(context.Background(), &store.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Database: viper.GetString("db.database"),
		AuthDB:   viper.GetString("db.auth-db"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
	})
	if err != nil {
		logrus.Fatalf("error to initialize db: %v", err)
	}

	store := store.NewRepository(db, viper.GetString("db.collection"))
	service := service.NewService(store)
	handler := handler.NewHandler(service)

	srv := new(apiserver.Server)

	// user1 := model.User{
	// 	ID:           "",
	// 	Username:     "Alex",
	// 	PasswordHash: "12345",
	// }

	//db.Collection(viper.GetString("db.collection")).InsertOne(context.Background(), user1)

	logrus.Info("server is listening port:", viper.GetString("port"))

	if err := srv.Start(viper.GetString("port"), handler.InitRoutes()); err != nil {
		logrus.Fatalf("error running server: %s", err.Error())
		return
	}

}

func InitConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
