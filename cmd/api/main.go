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
	bookRepo := repository.NewBookRepository(database.DB)
	bookService := services.NewBookServices(bookRepo)
	bookHandlers := handlers.NewBookHandlers(bookService)
	reservationRepo := repository.NewResservationRepository(database.DB)
	reservationService := services.NewReservationServices(reservationRepo)
	reservationHandlers := handlers.NewReservationHandler(reservationService)

	r := gin.Default()
	r.POST("api/users/register", userHandlers.UserRegister)
	r.POST("api/users/login", userHandlers.UserLogin)
	r.GET("api/users", userHandlers.UserList)
	r.DELETE("api/users/:id", userHandlers.UserDelete)
	r.PUT("api/users/:id", userHandlers.UserUpdate)
	r.GET("api/users/:id", userHandlers.UserGet)

	r.POST("api/books", bookHandlers.BookAdd)
	r.GET("api/books/:id", bookHandlers.BookGet)
	r.DELETE("api/books/:id", bookHandlers.BookDelete)
	r.PUT("api/books/:id", bookHandlers.BookUpdate)
	r.GET("api/books", bookHandlers.BookList)

	r.POST("api/reservation", reservationHandlers.ReservationCreate)
	r.PUT("api/reservation/:id", reservationHandlers.ReservationUpdate)
	r.GET("api/reservation/:id", reservationHandlers.ReservationGetUserByID)

	r.Run(":8080")
}
