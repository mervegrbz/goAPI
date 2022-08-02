package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"fmt"
	"io/ioutil"
	"encoding/json"
)



func getIP(c *gin.Context) {
	c.IndentedJSON(http.StatusOK,getip2())
	// fmt.Print(GetOutboundIP())
	
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

type album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, albums)
}

func main() {
	router := gin.Default()
	fmt.Print(getip2())
	
	router.GET("/albums", getAlbums)
	router.GET("/ip",getIP)
	router.Run(":8080")
}
