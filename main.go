package main

import(
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type Person_table struct{
	Name string
	Email string
	Phone int
	Password string
}

 func DBConnection() *gorm.DB{
	db, error := gorm.Open("postgres", " user=postgres dbname=gorm password=1234 sslmode=disable")
	if error!=nil{
		fmt.Println("error",error)
	}
	var createtable Person_table
	db.AutoMigrate(&createtable)
	return db
}
  

func main(){
	router := mux.NewRouter()
	router.HandleFunc("/createuser", createuser).Methods("GET")
	router.HandleFunc("/deleteuser", deleteuser).Methods("DELETE")
	router.HandleFunc("/updateuser",updateuser).Methods("PUT")
	router.HandleFunc("/Alluser",alluser).Methods("GET")
	router.HandleFunc("/password",password).Methods("GET")
	router.Handle("/",router)
	log.Fatal(http.ListenAndServe(":8000", router))
}
func createuser(w http.ResponseWriter, r *http.Request){
	var person Person_table
	var result Person_table
	json.NewDecoder(r.Body).Decode(&person)
	db := DBConnection()
	db.Where("Email=?",person.Email).Find(&result)
	if result.Email==""{
		db.Create(&person)
		fmt.Fprint(w,"You are register")
	}else{
		fmt.Fprint(w,"already exist")
	}
	defer db.Close()
}

func deleteuser(w http.ResponseWriter, r *http.Request){
	var delete Person_table
	var resultperson Person_table
	json.NewDecoder(r.Body).Decode(&delete)
	db :=DBConnection()
	db.Where("Email=?",delete.Email).Find(&resultperson)
	if resultperson.Email==delete.Email{
		db.Where("Email=?",delete.Email).Delete(resultperson)
		fmt.Fprint(w,"you are delete that person data")
	}else{
		fmt.Fprint(w,"No such data in person table")
	}
	defer db.Close()
}

func updateuser(w http.ResponseWriter, r *http.Request){
	var update Person_table
	var resultperson Person_table
	json.NewDecoder(r.Body).Decode(&update)
	db :=DBConnection()
	db.Where("Email=?",update.Email).Find(&resultperson)
	resultperson.Name="Rohit Kumar"
	resultperson.Email="rkrrko@gmail.com"
	resultperson.Password="1234"
	resultperson.Phone=12345
	db.Save(&resultperson)
	defer db.Close()
}

func alluser(w http.ResponseWriter, r *http.Request){
	var all []Person_table
	db := DBConnection()
	db.Find(&all)
	userJSON, _ := json.Marshal(all)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
	defer db.Close()

}

func password(w http.ResponseWriter, r *http.Request){
	var hide Person_table
	
	json.NewDecoder(r.Body).Decode(&hide)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(hide.Password), bcrypt.DefaultCost)
    if err != nil {
        panic(err)
	}
	
	
	hide.Password=(string(hashedPassword))
	userJSON, _ := json.Marshal(hide)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJSON)
	
}