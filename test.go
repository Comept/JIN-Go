package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/spf13/viper"
)

type Phone struct {
	Id      int    `json:"id"`
	Model   string `json:"model"`
	Company string `json:"company"`
	Price   int    `json:"price"`
}

var db *pg.DB
var servAddr string

func main() {
	initt()
	defer db.Close()
	router := gin.Default()
	fmt.Print(servAddr)
	router.GET("/phones", getPhone)
	router.GET("/phones/:id", findPhoneById)
	router.POST("/phones", postPhone)

	router.Run(servAddr)
}

// инициализация
func initt() {
	conf := viper.New()
	conf.SetConfigName("conf")
	conf.SetConfigType("env")
	conf.AddConfigPath("./util")
	err := conf.ReadInConfig()
	if err != nil {
		panic(err)
	}
	db = pg.Connect(&pg.Options{
		Addr:     conf.GetString("DBADDR"),
		User:     conf.GetString("DBUSER"),
		Password: conf.GetString("DBPASSWORD"),
		Database: conf.GetString("DB"),
	})
	servAddr = conf.GetString("SERVADDR")
	return
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
