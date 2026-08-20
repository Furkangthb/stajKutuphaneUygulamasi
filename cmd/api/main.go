package main

import (
	"os"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/database"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/middleware"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/services"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/handlers"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	database.VeriTabaniBaglanma()
	database.RedisConnect()

	userRepo := repository.NewUserRepository(database.DB)
	userService := services.NewUserServices(userRepo)
	userHandlers := handlers.NewUserHandlers(userService)

	bookRepo := repository.NewBookRepository(database.DB)
	bookService := services.NewBookServices(bookRepo)
	bookHandlers := handlers.NewBookHandlers(bookService)

	reservationRepo := repository.NewResservationRepository(database.DB)
	reservationService := services.NewReservationServices(reservationRepo)
	reservationHandlers := handlers.NewReservationHandler(reservationService)

	authRepo := repository.NewAuthRepository(database.RedisClient)
	authService := services.NewAuthServices(authRepo, userRepo, os.Getenv("JWT_SECRET"))
	authHandler := handlers.NewAuthHandlers(authService)

	r := gin.Default()

	r.POST("/api/users/register", userHandlers.UserRegister)
	r.POST("/api/auth/login", authHandler.Login)
	r.POST("/api/auth/logout", authHandler.Logout)
	r.GET("/api/users", middleware.RequireAuth(authService), middleware.RequireRole("admin"), userHandlers.UserList)
	r.DELETE("/api/users/:id", middleware.RequireAuth(authService), middleware.RequireRole("admin"), userHandlers.UserDelete)
	r.PUT("/api/users/:id", middleware.RequireAuth(authService), userHandlers.UserUpdate)
	r.GET("/api/users/:id", middleware.RequireAuth(authService), userHandlers.UserGet)

	r.POST("/api/books", middleware.RequireAuth(authService), middleware.RequireRole("admin"), bookHandlers.BookAdd)
	r.GET("/api/books/:id", bookHandlers.BookGet)
	r.DELETE("/api/books/:id", middleware.RequireAuth(authService), middleware.RequireRole("admin"), bookHandlers.BookDelete)
	r.PUT("/api/books/:id", middleware.RequireAuth(authService), middleware.RequireRole("admin"), bookHandlers.BookUpdate)
	r.GET("/api/books", middleware.RequireAuth(authService), bookHandlers.BookList)

	r.POST("/api/reservation", middleware.RequireAuth(authService), reservationHandlers.ReservationCreate)
	r.PUT("/api/reservation/:id", middleware.RequireAuth(authService), reservationHandlers.ReservationUpdate)
	r.GET("/api/reservation/:id", middleware.RequireAuth(authService), reservationHandlers.ReservationGetUserByID)

	r.Run(":8080")
}
