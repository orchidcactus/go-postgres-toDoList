package main

import (
    //"database/sql"
    "fmt"

    _"github.com/lib/pq" 
    "github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" 
	"log"
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

}

 

	func CheckError(err error) {
		if err!=nil {
			panic(err)
		}
	}
