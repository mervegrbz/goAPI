package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
	"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)
var envs,err = godotenv.Read(".env")

var (
	mysqlHost    = envs["HOST"]
	mysqlUser    = envs["USER"]    
	mysqlPass    = envs["PASSWORD"]
	mysqlDatabase= envs["DB"]
)
var result = fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", mysqlUser, 
                          mysqlPass, mysqlHost, mysqlDatabase)

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

type User struct {
	Id  int        `json:"id"`
	Name string    `json:"name"`
	Surname string `json:"surname"`
}


func getIP(c *gin.Context) {
	c.IndentedJSON(http.StatusOK,getip2())
}
type IP struct {
    Query string 
}
type Ip struct {
	IP string `json:"ip"`
}
func getip2() string {
    req, err := http.Get("http://ip-api.com/json/")
    if err != nil {
        return err.Error()
    }
    defer req.Body.Close()
    body, err := ioutil.ReadAll(req.Body)
    if err != nil {
        return err.Error()
    }
    var ip IP
    json.Unmarshal(body, &ip)
    return ip.Query
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
	
	c.IndentedJSON(http.StatusOK, albums)
}
func createTable(){
	db, err := sql.Open("mysql", result)
    if err != nil {
        panic(err.Error())
    }

	_, err2 := db.Exec("CREATE TABLE IF NOT EXISTS mytable (id BIGINT NOT NULL AUTO_INCREMENT PRIMARY KEY, some_text TEXT NOT NULL)")
	
    if err2 != nil {
        panic(err2)
    }

}
func createUser(name string,surname string, db sql.DB ) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second))
	defer cancel()
	trx, err := db.BeginTx(ctx, nil)

	if err != nil {
		trx.Rollback()
		fmt.Printf("error occured while creating transaction: %v", err)
	}

	result, err := trx.Exec("INSERT INTO mytable (name, surname) VALUES (?, ?)", name, surname)
	if err != nil {
		fmt.Printf("error occured while inserting user with name %v: %v", name, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("error occured while getting last insert id for %v: %v", name, err)
	}
	fmt.Print(id)
	err = trx.Commit()
	if err != nil {
		fmt.Printf("error occured while commiting transaction")
	}

}
func FetchUsers () (userlist []User, err error ) {
	users := []User{}
	db, err := sql.Open("mysql", result)
    if err != nil {
        panic(err.Error())
		return nil,err

    }
	
	rows, err := db.Query("SELECT id,name,surname FROM mytable where surname IS NOT NULL")
	if err != nil {
		fmt.Printf("error fetching users: %v", err)
		return nil, err
	}

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Name, &u.Surname); err != nil {
			fmt.Printf("error scanning the result set: %v", err)
			return nil, err
		}
		users = append(users, u)
	}
	
	return users, nil

}
func getUsers(c *gin.Context) {
	userlist, _  := FetchUsers()
	
	c.IndentedJSON(http.StatusOK, userlist)
}


func main() {
	// var users =[]User{
	// 	{id:0, name: "Blue Train", surname: "John Coltrane"},
	// 	{id:0, name: "Jeru", surname: "Gerry Mulligan"},
	// 	{id:0, name: "Sarah Vaughan and Clifford Brown", surname: "Sarah Vaughan"},
	// }
	
	// for _, elem := range users {
	//     createUser(elem.name, elem.surname, *db)
	// }
    db, err := sql.Open("mysql", result)
    if err != nil {
        panic(err.Error())
    }
	
	
	defer db.Close()
	fmt.Println(mysqlDatabase)
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/ip",getIP)
	router.GET("/users",getUsers)
	router.Run(":8080")
}
