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
	"strconv"
	"github.com/go-redis/redis"
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

var psqlcomm = fmt.Sprintf("host= %s port= %d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

var db, err = gorm.Open("postgres", psqlcomm)
//CheckError(err)

func Healthz(w http.ResponseWriter, r *http.Request) {
	log.Info("API Health is OK")
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"alive": true}`)
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
    // psqlcomm := fmt.Sprintf("host= %s port= %d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

    // db, err:= gorm.Open("postgres", psqlcomm)
    // CheckError(err)

    description := r.FormValue("description")
    log.WithFields(log.Fields{"description": description}).Info("Add new TodoItem. Saving to database.")
    todo := &TodoItemModel{Description: description, Completed: false}
    db.Create(&todo)
    result := db.Last(&todo)
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result.Value)
	
	// Set the element in the redis cache when user creates a new entry
 	models.SetInCache(models.REDIS, todo.TaskId, todo)
 	json.NewEncoder(w).Encode(todo)
}


// func CreateTodo(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var todo models.ToDo
// 	json.NewDecoder(r.Body).Decode(&todo)
// 	fmt.Print(todo)
// 	models.DB.Create(&todo)

// 	// Set the element in the redis cache when user creates a new entry
// 	models.SetInCache(models.REDIS, todo.TaskId, todo)
// 	json.NewEncoder(w).Encode(todo)
// }

func UpdateItem(w http.ResponseWriter, r *http.Request) {
	
	
   
	// Get URL parameter from mux
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	
	// Test if the TodoItem exist in DB
	err := GetItemByID(id)
	if err == false {
	    w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": false, "error": "Record Not Found"}`)
	} else {
	    completed, _ := strconv.ParseBool(r.FormValue("completed"))
		log.WithFields(log.Fields{"Id": id, "Completed": completed}).Info("Updating TodoItem")
		todo := &TodoItemModel{}
		db.First(&todo, id)
		todo.Completed = completed
		db.Save(&todo)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"updated": true}`)
	}
}

func DeleteItem(w http.ResponseWriter, r *http.Request) {
	// Get URL parameter from mux
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	// Test if the TodoItem exist in DB
	err := GetItemByID(id)
	if err == false {
	    w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": false, "error": "Record Not Found"}`)
	} else {
		log.WithFields(log.Fields{"Id": id}).Info("Deleting TodoItem")
		todo := &TodoItemModel{}
		db.First(&todo, id)
		db.Delete(&todo)
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"deleted": true}`)
	}
}
func GetItemByID(Id int) bool {
	todo := &TodoItemModel{}
	result := db.First(&todo, Id)
	if result.Error != nil{
		log.Warn("TodoItem not found in database")
		return false
	}
	return true
}

func GetCompletedItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get completed TodoItems")
    completedTodoItems := GetTodoItems(true)
    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(completedTodoItems)
}

func GetIncompleteItems(w http.ResponseWriter, r *http.Request) {
	log.Info("Get Incomplete TodoItems")
	IncompleteTodoItems := GetTodoItems(false)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(IncompleteTodoItems)
}

func GetTodoItems(completed bool) interface{} {
	var todos []TodoItemModel
	TodoItems := db.Where("completed = ?", completed).Find(&todos).Value
	return TodoItems
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
	router.HandleFunc("/todo-completed", GetCompletedItems).Methods("GET")
	router.HandleFunc("/todo-incomplete", GetIncompleteItems).Methods("GET")
	router.HandleFunc("/todo", CreateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", UpdateItem).Methods("POST")
	router.HandleFunc("/todo/{id}", DeleteItem).Methods("DELETE")
	http.ListenAndServe(":8000", router)

}



	func CheckError(err error) {
		if err!=nil {
			panic(err)
		}
	}

//curl -X POST -d "description= type the description" localhost:8000/todo