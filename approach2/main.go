// with-concurrency-handling-choose-seats on random
// 6*20
package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func main(){
	var err error
	Db,err=sql.Open("mysql","root:@/airline");
	if err!=nil{
		panic(err)
	}
	Db.SetConnMaxLifetime(time.Minute * 3)
	Db.SetMaxOpenConns(10)
	Db.SetMaxIdleConns(10)
	defer Db.Query(`DROP TABLE IF EXISTS seats`)
	defer Db.Close()
	var wg sync.WaitGroup
	rows,err:=Db.Query(`SELECT * FROM users`)
	if err!=nil{
		log.Fatalf("Error in fetching data from seats table: %v",err)	
	}
	for rows.Next(){
		var id int
		var name string
		err:=rows.Scan(&id,&name)
		if err!=nil{
			log.Fatalf("Error in scanning rows: %v",err)
		}
		wg.Add(1)

		go func (){
			defer wg.Done()
			allocateSeat(id)
		}()
	}
	wg.Wait()
	printSeats()
	
}	


func allocateSeat(id int) {
	tx,err:=Db.Begin()
	if err!=nil{
		log.Fatalf("Error in beginning transaction: %v",err)
	}
    seat := rand.Intn(120) + 1
    log.Printf("Allocating seat %d to user %d", seat, id)
    _, err1 := tx.Exec(`INSERT INTO bookings (id,user_id,seat) VALUES (NULL,?,?)`, id, seat)
    if err1 != nil {
        log.Printf("Error allocating seat: %v", err)
    }
	err=tx.Commit()
	if err!=nil{
		log.Fatalf("Error in committing transaction: %v",err)
	}
}
func printSeats(){
	rows,err:=Db.Query(`SELECT bookings.id as booking_id,users.id as user_id,users.username,bookings.seat FROM bookings JOIN users ON bookings.user_id=users.id`)
	var seats [6][20]bool
	if err!=nil{
		log.Fatalf("Error in fetching data from bookings table: %v",err)	
	}	
	for rows.Next(){
		var booking_id int
		var user_id int
		var username string
		var seat int
		err:=rows.Scan(&booking_id,&user_id,&username,&seat)
		if err!=nil{
			log.Fatalf("Error in scanning rows: %v",err)
		}
		fmt.Printf("Booking id: %d, User id: %d, Username: %s\n",booking_id,user_id,username)
		r,c:=(seat-1)/20,(seat-1)%20
		seats[r][c]=true
	}
	for i:=0;i<6;i++{
		for j:=0;j<20;j++{
			if seats[i][j]{
				fmt.Print("x")
			}else{
				fmt.Print(".")
			}
			
		
	}
		fmt.Println()
	}
}