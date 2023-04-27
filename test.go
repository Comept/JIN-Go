package main

import (
	"fmt"
	"net/http"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type Phone struct {
	Id      int    `json:"id"`
	Model   string `json:"model"`
	Company string `json:"company"`
	Price   int    `json:"price"`
}

var db *pg.DB

func main() {
	db = pg.Connect(&pg.Options{
		Addr:     ":5433",
		User:     "postgres",
		Password: "q",
		Database: "go",
	})
	defer db.Close()
	router := gin.Default()
	router.GET("/phones", getPhone)
	router.GET("/phones/:id", findPhoneById)
	router.POST("/phones", postPhone)

	router.Run("localhost:8080")
}

// обрабатывает get запрос для получения списка всех телефонов в бд
func getPhone(c *gin.Context) {
	var phones []Phone
	err := db.Model(&phones).Select()
	if err != nil {
		panic(err)
	}
	fmt.Println(phones)
	c.JSON(http.StatusOK, phones)
}

// обрабатывает post запрос для записи новой модели телефона в бд
func postPhone(c *gin.Context) {
	var newphone Phone
	if err := c.BindJSON(&newphone); err != nil {
		return
	}
	_, err := db.Model(&newphone).Insert()
	if err != nil {
		panic(err)
	}
}

// ищет модель телефона по id
func findPhoneById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		// ... handle error
		panic(err)
	}
	phone := &Phone{Id: id}
	err = db.Model(phone).WherePK().Select()
	if err != nil {
		panic(err)
	}
	fmt.Println(phone)

	c.IndentedJSON(http.StatusOK, phone)
	return
}

func createSchema(db *pg.DB) error {
	models := []interface{}{
		(*Phone)(nil),
	}
	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
