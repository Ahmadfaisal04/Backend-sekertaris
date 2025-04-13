package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func ConnectDB()(db *sql.DB, err error) {
	err = godotenv.Load()
	if err != nil {
		panic(err)
		
	}

	dbName := os.Getenv("DB_NAME")

	mysql := fmt.Sprintf("root:cMBXwyBNKnQGIgeWZZnGNDraEvvbcPen@tcp(caboose.proxy.rlwy.net:40008)/%s", dbName)
	db, err = sql.Open("mysql", mysql)
	if err != nil {
		panic(err)
		
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}