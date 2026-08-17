package main

import (
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/database"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/services"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/handlers"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	database.VeriTabaniBaglanma()
	userRepo := repository.NewUserRepository(database.DB)
	userService := services.NewUserServices(userRepo)
	userHandlers := handlers.NewUserHandlers(userService)

	r := gin.Default()
	r.POST("api/users", userHandlers.UserRegister)

	r.Run(":8080")
}
