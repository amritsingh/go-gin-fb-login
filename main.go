package main

import (
	"fb_login/helpers"
	"fb_login/middlewares"
	"fb_login/models"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
)

var FACEBOOK = &oauth2.Config{
	// Get your own ClienID and ClientSecret
	ClientID:     "",       // Copy from Facebook Developer page
	ClientSecret: "", 	// Copy from Facebook Developer page
	Scopes:       []string{},
	Endpoint:     facebook.Endpoint,
	RedirectURL:  "http://localhost:9000/facebook/auth",
}

func Auth(c *gin.Context) {
	code := c.Query("code")
	token, err := FACEBOOK.Exchange(oauth2.NoContext, code)

	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/")
	}

	// Create or update user
	user := models.UserUpdateOrCreate(token)

	// Create Login Session
	helpers.SessionSet(c, user.ID)

	c.Redirect(http.StatusMovedPermanently, "/")
}

func Logout(c *gin.Context) {
	// Clear the session
	helpers.SessionClear(c)
	c.Redirect(http.StatusMovedPermanently, "/")
}

func GetUserFromSession(c *gin.Context) *models.User {
	userID := c.GetUint64("user_id")
	if userID > 0 {
		return models.UserGetByID(userID)
	} else {
		return nil
	}
}

func Home(c *gin.Context) {
	// Get user from the session
	user := GetUserFromSession(c)

	authUrl := FACEBOOK.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.HTML(http.StatusOK, "views/index.html", gin.H{
		"authUrl": authUrl,
		"user":    user,
	})
}

func main() {
	r := gin.Default()

	// Init session store
	store := memstore.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("gin_fb_login", store))

	r.Use(gin.Logger())

	models.ConnectDatabase()

	r.Static("/vendor", "./static/vendor")

	r.LoadHTMLGlob("templates/**/*")

	// Add middleware
	r.Use(middlewares.AuthenticateUser())

	r.GET("/", Home)

	r.GET("/facebook/auth", Auth)
	r.POST("/logout", Logout)

	log.Println("Server started!")
	r.Run("0.0.0.0:9000") // Port changed to 9000
}
