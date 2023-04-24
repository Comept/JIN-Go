package main

import (
	"fmt"
	"net/http"

	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type phone struct {
	Id      int    `json:"id"`
	Model   string `json:"model"`
	Company string `json:"company"`
	Price   int    `json:"price"`
}

func main() {
	router := gin.Default()
	router.GET("/albums", getPhone)
	router.GET("/albums/:id", findPhoneByCompany)
	router.POST("/albums", postPhone)

	router.Run("localhost:8080")
}

// обрабатывает get запрос для получения списка всех телефонов в бд
func getPhone(c *gin.Context) {
	var phones []phone = dbQuere("select * from phones")
	fmt.Println(phones)
	c.JSON(http.StatusOK, phones)
}

// обрабатывает post запрос для записи новой модели телефона в бд
func postPhone(c *gin.Context) {
	var newphone phone

	if err := c.BindJSON(&newphone); err != nil {
		return
	}
	query := fmt.Sprintf("insert into Phones (Model, Company, Price) values (%s, %s, %d)", newphone.Model, newphone.Company, newphone.Price)
	dbQuere(query)
	// Add the new album to the slice.
}

// ищет модель телефона по id
func findPhoneByCompany(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// ... handle error
		panic(err)
	}
	query := fmt.Sprintf("select * from Phones WHERE id = %d", id)
	var xj []phone = dbQuere(query)

	fmt.Println(xj)

	c.IndentedJSON(http.StatusOK, xj)
	return
}

// обрабатывает запросы в бд
func dbQuere(SQL_query string) []phone {

	connStr := "user=postgres port=5433 password=q dbname=go sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query(SQL_query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	phones := []phone{}

	for rows.Next() {
		p := phone{}
		err := rows.Scan(&p.Id, &p.Model, &p.Company, &p.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}
		phones = append(phones, p)
	}
	fmt.Println(phones)
	return phones
}
