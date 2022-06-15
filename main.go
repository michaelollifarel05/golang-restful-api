package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var username = goDotEnvVariable("DB_USER")
var password = goDotEnvVariable("DB_PASS")
var hostname = goDotEnvVariable("DB_URL")
var dbName = goDotEnvVariable("DB_NAME")

// var dbgorm *gorm.DB

type User struct {
	ID       int64  `gorm:"primary_key;auto_increment;not_null"`
	Name     string `json:"username"`
	Password string `json:"password"`
}

func goDotEnvVariable(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: /homePage")
}

func gorms(p User) {

}

func register(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var person User
	json.Unmarshal(requestBody, &person)
	dbgorm, err1 := gorm.Open("mysql", dsn())
	if err1 != nil {
		panic("failed to connect database")
	}
	dbgorm.Save(&person)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(person)
	fmt.Println("Endpoint Hit: /register")

}

func initialMigrate() {
	var person User
	dbgorm, err1 := gorm.Open("mysql", dsn())
	if err1 != nil {
		panic("failed to connect database")
	}
	dbgorm.AutoMigrate(person)
	fmt.Println("DB Migrated")

}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/register", register).Methods("POST")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	// fmt.Print(dsn())
	initialMigrate()
	// checkDBConn()
	handleRequests()

}
