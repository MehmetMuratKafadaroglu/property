package main

import (
	"fmt"
	"go-rest/db"
	"go-rest/models"
	"go-rest/settings"
	"go-rest/utils"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}
func get_properties(context *gin.Context) {
	location := context.Param("location")
	max_price, err := strconv.Atoi(context.Param("max_price"))
	checkErr(err)
	min_price, err := strconv.Atoi(context.Param("min_price"))
	checkErr(err)
	max_rooms, err := strconv.Atoi(context.Param("max_rooms"))
	checkErr(err)
	min_rooms, err := strconv.Atoi(context.Param("min_rooms"))
	checkErr(err)
	max_internal_area, err := strconv.Atoi(context.Param("max_internal_area"))
	checkErr(err)
	min_internal_area, err := strconv.Atoi(context.Param("min_internal_area"))
	checkErr(err)

	properties := db.SelectProperty(location, max_price, min_price, max_rooms, min_rooms, max_internal_area, min_internal_area)
	context.IndentedJSON(http.StatusOK, properties)
}
func insert_properties(context *gin.Context) {
	var newProperty models.Property
	if err := context.BindJSON(&newProperty); err != nil {
		fmt.Println(err.Error())
		return
	}
	db.InsertProperty(&newProperty)
	context.IndentedJSON(http.StatusCreated, newProperty)
}

func insert_user(context *gin.Context) {
	var newUser models.User
	if err := context.BindJSON(&newUser); err != nil {
		return
	}
	id := db.InsertUser(newUser)
	runes := utils.RandStringRunes(128)
	url := settings.ServerUrl + "/verify/" + strconv.FormatInt(id, 10) + "/" + runes
	utils.SendMail(newUser.Email, url)
	db.InsertTemporaryKey(id, runes)
	context.IndentedJSON(http.StatusCreated, newUser)
}

func verify_user(context *gin.Context) {
	userID, err := strconv.Atoi(context.Param("userID"))
	checkErr(err)
	key := context.Param("key")
	db.VerifyUser(key, int64(userID))
}

func login(context *gin.Context) {
	var user models.LoginUser
	if err := context.BindJSON(&user); err != nil {
		fmt.Println(err.Error())
		return
	}
	_id := db.CanUserLogIn(user.Password, user.Email)
	if _id == 0 {
		context.JSON(http.StatusOK, gin.H{"token": "", "error": 1})
		return
	}
	token := utils.GenerateToken(user.Email, _id)
	context.JSON(http.StatusOK, gin.H{"token": token, "error": 0})
}

func logout(context *gin.Context) {
	fmt.Println("Logout")
	context.JSON(200, gin.H{"error": 0})
}

func save(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("userID"))
	checkErr(err)
	propertyID, err := strconv.Atoi(ctx.Param("propertyID"))
	checkErr(err)

	err = db.SaveProperty(userID, propertyID)
	checkErr(err)
	ctx.JSON(200, gin.H{"error": 0})
}

func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		token, err := utils.VerifyToken(authHeader)
		if err != nil {
			context.AbortWithStatus(400)
		}
		context.Header("Token", token)
		context.Next()
	}
}
func main() {
	utils.Init()
	db.InitDatabase()
	router := gin.Default()
	router.Static("/assets", "./assets")
	router.Use(sessions.Sessions("session", cookie.NewStore([]byte(settings.Secret))))

	public := router.Group("/public")
	public.GET("/properties/:location/:max_price/:min_price/:max_rooms/:min_rooms/:max_internal_area/:min_internal_area", get_properties)
	public.POST("/user/", insert_user)
	public.GET("/verify/:userID/:key", verify_user)
	public.POST("/login/", login)

	private := router.Group("/private")
	private.Use(Auth())
	private.POST("/add/properties/", insert_properties)
	private.GET("/logout/", logout)
	private.GET("/save/:userID/:propertyID", save)
	router.Run(settings.ServerName)
}
