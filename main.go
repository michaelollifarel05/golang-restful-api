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

type User struct {
	ID       int64  `gorm:"primary_key;auto_increment;not_null"`
	Name     string `json:"username"`
	Password string `json:"password"`
}

func checkDBConn() {
	dbgorm, err1 := gorm.Open("mysql", dsn())
	_ = dbgorm
	if err1 != nil {
		panic("failed to connect database")
	}
	fmt.Println("checked")
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

func register(w http.ResponseWriter, r *http.Request) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	var person User
	json.Unmarshal(requestBody, &person)
	dbgorm, err1 := gorm.Open("mysql", dsn())
	if err1 != nil {
		panic("failed to connect database")
	}

	var check int = checkUser(person)
	if check == 1 {
		dbgorm.Save(&person)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(person)
		fmt.Println("Endpoint Hit: /register")
	} else {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		wrong := map[string]string{"Message": "User already Exist"}

		json.NewEncoder(w).Encode(wrong)
	}
}

func searchUser(w http.ResponseWriter, r *http.Request) {
	dbgorm, err1 := gorm.Open("mysql", dsn())
	if err1 != nil {
		panic("failed to connect database")
	}
	vars := mux.Vars(r)
	key := vars["usname"]
	var persona User
	dbgorm.First(&persona, "name = ?", key)
	// fmt.Print(persona.ID)
	if persona.ID == 0 {
		wrong := map[string]string{"Message": "User already Exist"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(wrong)
	} else {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(persona)
	}

}

func checkUser(person User) (check int) {
	dbgorm, err1 := gorm.Open("mysql", dsn())
	if err1 != nil {
		panic("failed to connect database")
	}
	dbgorm.Where("name = ?", person.Name).Find(&person)
	fmt.Print(person.ID)
	if person.ID == 0 {
		check = 1
	} else {
		check = 0
	}

	return
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
	myRouter.HandleFunc("/search-user/{usname}", searchUser).Methods("GET")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	// fmt.Print(dsn())
	initialMigrate()
	checkDBConn()
	handleRequests()

}
