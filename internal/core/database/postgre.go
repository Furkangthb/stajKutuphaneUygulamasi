package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func VeriTabaniBaglanma() {

	veriTabaniUrl := os.Getenv("POSTGRES_URL")

	if veriTabaniUrl == "" {
		log.Fatal("Url bulunamadı")
	}
	var err error
	DB, err = sql.Open("pqx", veriTabaniUrl)
	if err != nil {
		log.Fatal("Veri Tabanı Başlatılamadı")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err = DB.PingContext(ctx); err != nil {
		log.Fatal("Ping atılamadı")
	}
	fmt.Print("Bağlantı başarılı")
}
