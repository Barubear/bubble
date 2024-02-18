package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func main() {

	//create  a database bubble
	//link to bubble
	err := initMySql()
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&Todo{})

	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("./templates/*")
	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	DBgroup := r.Group("v1")
	{
		DBgroup.POST("/todo", postData)
		DBgroup.GET("/todo", getData)
		DBgroup.GET("/todo/:id", getData)
		DBgroup.PUT("/todo/:id", putData)
		DBgroup.DELETE("/todo/:id", deleteData)
	}

	r.Run("localhost:9000")
}

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status bool   `json:"status"`
}

func initMySql() (err error) {
	dsn := "root:980606@tcp(localhost:3306)/bubble"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	return err
}

func postData(c *gin.Context) {
	var todo Todo
	c.BindJSON(&todo)
	if err := DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, todo)
	}
}
func getData(c *gin.Context) {
	var todoList []Todo
	if err := DB.Find(&todoList).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
	} else {
		c.JSON(http.StatusOK, todoList)
	}
}
func putData(c *gin.Context) {
	id, _ := c.Params.Get("id")

	var todo Todo

	err := DB.Where("id=?", id).First(&todo).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})
		return
	}

	c.BindJSON(&todo)

	err = DB.Save(&todo).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})

	} else {
		c.JSON(http.StatusOK, todo)
	}
}
func deleteData(c *gin.Context) {
	id, _ := c.Params.Get("id")

	err := DB.Where("id=?", id).Delete(Todo{}).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err})

	} else {
		c.JSON(http.StatusOK, gin.H{id: "deleted"})
	}

}
