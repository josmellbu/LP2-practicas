package main

import (
	//"model"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Person struct {
    gorm.Model
    Name string
    Age  string
}

type User struct {
	ID uint64            `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	//A sample use
	var user = User{
		ID:             1,
		Username: "username",
		Password: "password",
	}

func main() {
	dsn := "docker:docker@tcp(db:3306)/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Person{})

	//s := Person{Name: "Sean", Age: 50}
	//s.Name = "Sean"

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hola GIN",
		})
	})

	r.GET("/login", Login)

	r.GET("/persons/:id", func(c *gin.Context) {
		id := c.Param("id")
		var d Person
		if err := db.First(&d, id).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
  		//db.First(&d, id)
		c.JSON(http.StatusOK, &d)
	})

	r.GET("/persons/", func(c *gin.Context) {
		var lis []Person
  		db.Find(&lis)
		c.JSON(http.StatusOK, lis)
	})

	r.POST("/persons/", func(c *gin.Context) {
		//d := Person{Name: c.PostForm("name"), Age: c.PostForm("age")}
		//db.Create(&d)
		//c.JSON(200, gin.H{
		//	"name": d.Name,
		//	"age": d.Age,
		//})
		//var d Person
		d := Person{Name: c.PostForm("name"), Age: c.PostForm("age")}
		/*if err := c.BindJSON(&d); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}*/
		db.Create(&d)
		c.JSON(http.StatusOK, &d)
	})

	r.PUT("/persons/:id", func(c *gin.Context) {
		id := c.Param("id")
		var d Person
		if err := db.First(&d, id).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		d.Name = c.PostForm("name")
		d.Age = c.PostForm("age")
  		db.Save(&d)
		c.JSON(http.StatusOK, &d)
	})

	r.DELETE("/persons/:id", func(c *gin.Context){
		id := c.Param("id")
		var d Person
		if err := db.Where("id = ?", id).First(&d).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		db.Unscoped().Delete(&d)
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}




func Login(c *gin.Context) {
	var u User
	if err := c.ShouldBindJSON(&u); err != nil {
	   c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
	   return
	}
	//compare the user from the request, with the one we defined:
	if user.Username != u.Username || user.Password != u.Password {
	   c.JSON(http.StatusUnauthorized, "Please provide valid login details")
	   return
	}
	token, err := CreateToken(user.ID)
	if err != nil {
	   c.JSON(http.StatusUnprocessableEntity, err.Error())
	   return
	}
	c.JSON(http.StatusOK, token)
  }

  func CreateToken(userid uint64) (string, error) {
	var err error
	//Creating Access Token
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
	   return "", err
	}
	return token, nil
  }