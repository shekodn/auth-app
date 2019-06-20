package model

import (
  "fmt"
  "log"
  "os"

  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
  "github.com/joho/godotenv"
)

const hashCost = 8
var db *gorm.DB

type User struct {
  gorm.Model
  Username string `json:"username"`
  Password string `json:"password"`
}

func init() {
  e := godotenv.Load()

  if e !=nil {
    log.Println(e)
  }

  username := os.Getenv("POSTGRES_USER")
  password := os.Getenv("POSTGRES_PASSWORD")
  dbName := os.Getenv("POSTGRES_DB")
  dbHost := os.Getenv("DB_HOST")

  dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)

  conn, err := gorm.Open("postgres", dbUri)

  if err != nil {
    log.Print("Error:", err, "\n")
  }

  db = conn
  db.Debug().AutoMigrate(&User{}) // DB migration # TODO
}

func GetDB() *gorm.DB {
  return db
}
