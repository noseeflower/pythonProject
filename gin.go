package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

type PostParams struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required,limitAge"`
	Sex  bool   `json:"sex" binding:"required"`
}

var limitAge validator.Func = func(fl validator.FieldLevel) bool {
	age, ok := fl.Field().Interface().(int)
	if ok {
		maxAge := 18
		if age <= maxAge {
			return false
		}
	}
	return true
}

func main() {
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("limitAge", limitAge)
	}

	v1 := r.Group("v1")

	v1.GET("/test", GetContext)
	v1.POST("/test", PostContext)
	v1.POST("/bindjson", BindJson)

	r.POST("/Upload", Upload)
	r.POST("/Uploads", Uploads)
	r.GET("/Download", Download)
	err := r.Run(":9000")
	if err != nil {
		fmt.Println("RunErr:", err)
	}
}

func middel() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("路由方法前")
		c.Next()
		fmt.Println("路由方法后")
	}
}

func middeltwo() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("路由方法前2")
		c.Next()
		fmt.Println("路由方法后2")
	}
}

func BindJson(c *gin.Context) {
	var p PostParams
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(200, gin.H{
			"msg":  "导入json错误",
			"data": gin.H{},
		})
	} else {
		c.JSON(200, gin.H{
			"msg":  "导入json成功",
			"data": p,
		})
	}
}

func PostContext(c *gin.Context) {
	user := c.DefaultPostForm("user", "fafa")
	pwd := c.PostForm("pwd")
	c.JSON(200, gin.H{
		"success": "post",
		"user":    user,
		"pwd":     pwd,
	})
}

func GetContext(c *gin.Context) {
	fmt.Println("Get路由")
	c.JSON(http.StatusOK, gin.H{
		"success": "花老师真tm帅",
	})
}

func Upload(c *gin.Context) {
	file, _ := c.FormFile("file")
	c.SaveUploadedFile(file, "./"+file.Filename)
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Filename))
	c.File("./" + file.Filename)
}

func Uploads(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["file"]
	for _, file := range files {
		log.Println(file.Filename)
		c.SaveUploadedFile(file, "./"+file.Filename)
	}
}

func Download(c *gin.Context) {
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "适当说脏话.jpg"))
	c.File("./" + "适当说脏话.jpg")
}
