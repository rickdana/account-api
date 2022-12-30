package main

import (
	"account-api/handler"
	"account-api/model"
	"account-api/repository"
	"account-api/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func intDB() (*gorm.DB, error) {
	dsn := "host=localhost user=dbuser password=4eyesP@ssw0rd dbname=4eyes_poc port=5432 sslmode=disable TimeZone=Europe/Paris"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: logger.Default.LogMode(logger.Info),

	})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	db.AutoMigrate(model.Account{}, model.User{})
	return db, err
}

func main() {
	db, err := intDB()

	//initialize Repo
	accountRepo := repository.NewAccountRepository(db)
	userRepo := repository.NewUserRepository(db)

	//initialize service
	kafkaConfig := service.KafkaConfig{
		Url:       "localhost:9092",
		Topic:     "create",
		Partition: 0,
	}
	kafkaService := service.NewKafkaSender(kafkaConfig)
	accountValidator := service.NewAccountValidator()
	accountService := service.NewAccountService(accountRepo)
	userService := service.NewUserServiceImpl(kafkaService, userRepo)

	if err != nil {
		panic("failed to connect database")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.BasicAuth("4eyes", map[string]string{
		"admin": "admin",
		"user1": "user1",
	}))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger", http.StatusFound)
	})
	r.Mount("/swagger", httpSwagger.WrapHandler)
	r.Mount("/users", handler.NewUsersResource(userService).Routes())
	r.Mount("/accounts", handler.NewAccountsResource(accountService, accountValidator, kafkaService).Routes())
	http.ListenAndServe(":3000", r)
}
