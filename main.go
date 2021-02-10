package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"strconv"
)

type Purchases struct {
	ID string `json:"item_id"`
	Name string `json:"item_name"`
	Quantity string `json:"item_quantity"`
	Rate string `json:"item_rate"`
	Date string `json:"item_purchase_date"`
}

var allPurchases []Purchases

//create database connection function and return the connection
func createConnection() *sql.DB  {
	db,err := sql.Open("mysql", "root:13628@tcp(localhost:3306)/items?charset=utf8")
	if err != nil {
		fmt.Println(`Could not connect to db`)
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func getPurchases(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")
	db := createConnection()
	defer db.Close()
	rows,err := db.Query(`SELECT * FROM purchases`)
	if err != nil {
		panic(err)
	}
	fmt.Println(rows)
	var col1 string
	var col2 string
	var col3 string
	var col4 string
	var col5 string
	allPurchases = nil
	for rows.Next() {
		rows.Scan(&col1,&col2,&col3,&col4,&col5)
		allPurchases = append(allPurchases,Purchases{ID: col1,Name: col2,Quantity: col3,Rate: col4,Date: col5})
	}
	json.NewEncoder(w).Encode(allPurchases)

}

func getPurchase(w http.ResponseWriter,r *http.Request)  {
	w.Header().Set("Content-Type","application/json")
	w.Header().Set("Access-Control-Allow-Origin","*")
	db := createConnection()
	defer db.Close()
	params := mux.Vars(r)
	id,err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the srting into int. %v",err)
	}

	rows := db.QueryRow(`SELECT * FROM purchases WHERE item_id=? limit 1`,id)

	var col1 string
	var col2 string
	var col3 string
	var col4 string
	var col5 string
	allPurchases = nil

	rows.Scan(&col1,&col2,&col3,&col4,&col5)
	//fmt.Println(col1,col2,col3,col4,col5)
	allPurchases = append(allPurchases,Purchases{ID: col1,Name: col2,Quantity: col3,Rate: col4,Date: col5})
	json.NewEncoder(w).Encode(allPurchases)

}

func createPurchase(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Access-Control-Allow-Origin","*")
	var p Purchases
	err := json.NewDecoder(r.Body).Decode(&p)
	fmt.Println("From api : in create purchase : this is post p : ",p)
	fmt.Println("From api : in create purchase : this is error:",err)
	db := createConnection()
	defer db.Close()
	row,err := db.Exec(`INSERT INTO purchases (item_id,item_name,item_quantity,item_rate,
                       item_purchase_date) VALUES ($1,$2,$3,$4,$5)`,p.ID,p.Name,p.Quantity,p.Rate,p.Date)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v",err)
	}
	fmt.Println(row)
	fmt.Printf("Insert a single record %v",p.ID)
}

//update post
func updatePurchase(w http.ResponseWriter, r *http.Request)  {
	//w.Header().Set("Content-Type","application/json")
	//w.Header().Set("Access-Control-Allow-Origin","*")


	var p Purchases
	err := json.NewDecoder(r.Body).Decode(&p)
	fmt.Println("this is post p : ",p)
	fmt.Println("this is error: ",err)

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v",err)
	}
	fmt.Println("this is id : ",id)

	db := createConnection()

	defer db.Close()

	row,err := db.Exec(`UPDATE purchases SET item_id=$1,item_name=$2,item_quantity=$3,item_rate=$4,
					item_purchase_date=$5 WHERE item_id=$6`,p.ID,p.Name,p.Quantity,p.Rate,p.Date,id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v",err)
	}
	fmt.Println(row)
	fmt.Printf("Inserted a single record %v\n", p.ID)
}

//delete
func deletePurchase(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Access-Control-Allow-Origin","*")
	w.Header().Set("Access-Control-Allow-Methods","DELETE")
	w.Header().Add("Access-Control-Allow-Headers","Content-Type")
	w.Header().Add("Content-Type","application/json")


	db := createConnection()

	defer db.Close()

	params := mux.Vars(r)

	id,err := strconv.Atoi(params["id"])
	fmt.Println(id)

	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v",err)
	}

	rows,err := db.Exec(`DELETE FROM purchases WHERE item_id=$1`,id)

	if err != nil {
		panic(err)
	}
	fmt.Println(rows)
}


func main() {
	r := mux.NewRouter()

	r.HandleFunc("/purchases",getPurchases).Methods("GET")
	r.HandleFunc("/purchases",getPurchase).Methods("GET")
	r.HandleFunc("/purchases/{id}",getPurchase).Methods("GET")
	r.HandleFunc("/purchases",createPurchase).Methods("POST")
	r.HandleFunc("/purchases/{id}",updatePurchase).Methods("PUT")
	r.HandleFunc("/purchases/{id}",deletePurchase).Methods("DELETE")

	handler := cors.AllowAll().Handler(r)
	//handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":8080",handler))
}
