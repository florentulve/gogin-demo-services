package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/jinzhu/gorm"

	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Posts struct {
	ID        int8      `gorm:"type:serial;primary_key;NOT NULL" json:"id"`
	Title     string    `gorm:"type:serial;NOT NULL" json:"title"`
	Body      string    `gorm:"type:serial;NOT NULL" json:"body"`
	Published bool      `gorm:"type:serial;NOT NULL" json:"published"`
	UUID      uuid.UUID `gorm:"type:uuid;NOT NULL" json:"uuid"`
}

//BCryptHashRequest password hash with bcrypt request
type BCryptHashRequest struct {
	Password string `form:"password" json:"password"`
	Hash     string `form:"Hash" json:"Hash"`
}

func setupRouter(db *gorm.DB) *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.POST("bcrypt/hash", func(c *gin.Context) {
		var data BCryptHashRequest
		if c.BindJSON(&data) == nil {
			hashb, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "hash failed"})
			}
			//hash := base64.StdEncoding.EncodeToString(hashb)
			c.JSON(http.StatusOK, gin.H{"password": string(hashb)})
		}
	})

	r.POST("bcrypt/verify", func(c *gin.Context) {
		var data BCryptHashRequest
		if c.BindJSON(&data) == nil {
			err := bcrypt.CompareHashAndPassword([]byte(data.Hash), []byte(data.Password))
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": err})
			} else {
				c.JSON(http.StatusOK, gin.H{"status": "Ok"})
			}

		}
	})

	r.GET("posts", func(c *gin.Context) {
		var posts []Posts
		db.Find(&posts)

		c.JSON(http.StatusOK, gin.H{"results": posts})

	})

	return r
}

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Service Port")
	flag.Parse()

	db, err := gorm.Open("postgres", "host=localhost port=5432 user=DB password=DB dbname=DB sslmode=disable")
	defer db.Close()
	if err != nil {
		panic(err)
	}
	//db.SetConnMaxLifetime(25)
	r := setupRouter(db)
	// Listen and Server in 0.0.0.0:8080
	r.Run(fmt.Sprintf(":%d", port))
}
