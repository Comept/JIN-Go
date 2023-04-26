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
	router.GET("/phones", getPhone)
	router.GET("/phones/:id", findPhoneById)
	router.POST("/phones", postPhone)

	router.Run("localhost:8080")
}

// обрабатывает get запрос для получения списка всех телефонов в бд
func getPhone(c *gin.Context) {
	var phones []phone = takeFromDB("select * from phones")
	fmt.Println(phones)
	c.JSON(http.StatusOK, phones)
}

// обрабатывает post запрос для записи новой модели телефона в бд
func postPhone(c *gin.Context) {
	fmt.Print("213")
	var newphone phone
	fmt.Print("juugugug")
	if err := c.BindJSON(&newphone); err != nil {
		fmt.Print("213")
		return
	}
	fmt.Print(newphone.Company)

	writeToDB(newphone)
}

// ищет модель телефона по id
func findPhoneById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// ... handle error
		panic(err)
	}
	query := fmt.Sprintf("select * from Phones WHERE id = %d", id)
	var xj []phone = takeFromDB(query)

	fmt.Println(xj)

	c.IndentedJSON(http.StatusOK, xj)
	return
}

// обрабатывает запросы в бд
func takeFromDB(SQL_query string) []phone {

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

func writeToDB(newphone phone) {

	connStr := "user=postgres port=5433 password=q dbname=go sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	result, err := db.Exec("insert into Phones (Model, Company, Price) values ($1, $2, $3)",
		newphone.Model, newphone.Company, newphone.Price)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected()) // количество добавленных строк

	return
}
