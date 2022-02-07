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

// type Employee struct {
// 	gorm.Model
// 	Name string  `gorm:"type:varchar(255);" json:"title"`
// 	EmpId  int  `gorm:"type:text" json:"body"`

type Employee struct {
	gorm.Model
	Name string `gorm:"column:Name"`
	EmpId  int  `gorm:"column:EmpId"` 
}

func main() {
	psqlcomm := fmt.Sprintf("host= %s port= %d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err:= gorm.Open("postgres", psqlcomm)
	CheckError(err)

    log.Println("DB Connection established...")

	defer db.Close()

	// InsertStmt := `insert into "Employee"("Name", "EmpId") values('Rohit', 21)`
	// _, e := db.Exec(InsertStmt)
	// CheckError(e)

	// InsertDynStmt := `insert into "Employee"("Name", "EmpId") values($1, $2)`
	// _, e = db.Exec(InsertDynStmt, "krish", 03)
	// CheckError(e)
	employee1 := &Employee{Name: "Karen", EmpId: 56}
    err = db.Create(&employee1).Error
    CheckError(err)

}

 

	func CheckError(err error) {
		if err!=nil {
			panic(err)
		}
	}
