package main

import (
    //"database/sql"
    "fmt"
    "io"
    "net/http"
	"github.com/gorilla/mux"
    log "github.com/sirupsen/logrus"
    _"github.com/lib/pq" 
    "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" 
	
	"encoding/json"
)

const (
	host=     "localhost"
	port=     5432
	user=     "postgres"
	password= "iamtir3d"
	dbname=   "first_db"
)


type TodoItemModel struct{
	 Id int `gorm:"primary_key"`
	 Description string
	 Completed bool
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
    psqlcomm := fmt.Sprintf("host= %s port= %d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

    db, err:= gorm.Open("postgres", psqlcomm)
    CheckError(err)

    description := r.FormValue("description")
    log.WithFields(log.Fields{"description": description}).Info("Add new TodoItem. Saving to database.")
    todo := &TodoItemModel{Description: description, Completed: false}
    db.Create(&todo)
    result := db.Last(&todo)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result.Value)
	
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetReportCaller(true)
}


func main() {
	psqlcomm := fmt.Sprintf("host= %s port= %d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err:= gorm.Open("postgres", psqlcomm)
	CheckError(err)

    log.Println("DB Connection established...")

	defer db.Close()

	 //db.Debug().DropTableIfExists(&TodoItemModel{})
    db.Debug().AutoMigrate(&TodoItemModel{})

	todo1 := &TodoItemModel{Description: "Buy bread", Completed: false}
    err = db.Create(&todo1).Error
    CheckError(err)



    log.Info("Starting Todolist API server")
	router := mux.NewRouter()
	router.HandleFunc("/healthz", Healthz).Methods("GET")
	router.HandleFunc("/todo", CreateItem).Methods("POST")
	http.ListenAndServe(":8000", router)

}



	func CheckError(err error) {
		if err!=nil {
			panic(err)
		}
	}

//curl -X POST -d "description= type the description" localhost:8000/todo