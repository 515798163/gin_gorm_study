package main

import (
	"fmt"
	"gin-class/middleware"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID    uint
	Name  string
	Email string
	Age   uint8
}

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm_class?charset=utf8mb4&parseTime=True&loc=Local"
	db, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&User{})
	r := gin.Default() // 携带基础中间件启动

	// 使用中间件
	r.Use(middleware.RecoveryMiddleware())

	// 模拟一个panic
	r.GET("/panic", func(c *gin.Context) {
		panic("模拟一个panic")
	})

	r.POST("/user", func(c *gin.Context) {
		var user User
		_ = c.BindJSON(&user)
		fmt.Println(user)
		db.Create(&user)
	})

	r.GET("/user/:age", func(c *gin.Context) {
		age := c.Param("age")
		var users []User

		db.Where("age>?", age).Find(&users).Limit(3)
		c.JSON(200, gin.H{
			"users": users,
		})
	})

	r.GET("/user", func(c *gin.Context) {
		var users []User
		db.Scopes(Paginate(c)).Find(&users)
		c.JSON(200, gin.H{
			"users": users,
		})
	})

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("mustBig", mustBig)
	}

	r.POST("/testBind", func(c *gin.Context) {
		var p PostParams
		err := c.ShouldBindJSON(&p)

		if err != nil {
			fmt.Println(err.Error())
			c.JSON(200, gin.H{
				"msg":  "报错",
				"data": gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"msg":  "成功",
				"data": p,
			})
		}
	})
	r.Run(":1010")
	////r.Run(":1010") // 监听并在 0.0.0.0:8080 上启动服务
	//r.GET("/path/:id", func(c *gin.Context) {
	//	id := c.Param("id")
	//	user := c.DefaultQuery("user", "runze")
	//	pwd := c.Query("pwd")
	//	c.JSON(200, gin.H{
	//		"id":   id,
	//		"user": user,
	//		"pwd":  pwd,
	//	})
	//})
	//r.POST("/path", func(c *gin.Context) {
	//	user := c.DefaultQuery("user", "runze")
	//	pwd := c.PostForm("pwd")
	//	c.JSON(200, gin.H{
	//		"user": user,
	//		"pwd":  pwd,
	//	})
	//})
	//
	//r.DELETE("/path/:id", func(c *gin.Context) {
	//	id := c.Param("id")
	//	user := c.DefaultQuery("user", "runze")
	//	pwd := c.Query("pwd")
	//	c.JSON(200, gin.H{
	//		"id":   id,
	//		"user": user,
	//		"pwd":  pwd,
	//	})
	//})
	//
	//r.PUT("/path", func(c *gin.Context) {
	//	user := c.DefaultQuery("user", "runze")
	//	pwd := c.PostForm("pwd")
	//	c.JSON(200, gin.H{
	//		"user": user,
	//		"pwd":  pwd,
	//	})
	//})
	//
	//r.Run(":1010") // 监听并在 0.0.0.0:1010 上启动服务
}

type PostParams struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required,mustBig"`
	Sex  bool   `json:"sex" binding:"required"`
}

func mustBig(fl validator.FieldLevel) bool {
	//fmt.Println(fl.Field().Interface().(int))
	if fl.Field().Interface().(int) <= 18 {
		return false
	}
	return true
}
