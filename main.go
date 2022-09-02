package main

import (
	"fmt"
	"go-rest/db"
	"go-rest/models"
	"go-rest/settings"
	"go-rest/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func checkErr(e error) {
	if e != nil {
		fmt.Println(e.Error())
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

func addImages(ctx *gin.Context) {
	var images models.AddPropertyImages
	err := ctx.BindJSON(&images)
	checkErr(err)
	err = checkUserFromClaims(ctx, images.PropertyID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": 1})
	} else {
		db.ImageInsert(images.Images, images.PropertyID)
		ctx.JSON(200, gin.H{"error": 0})
	}
}
func deleteImages(ctx *gin.Context) {
	var images models.PropertyImages
	err := ctx.BindJSON(&images)
	checkErr(err)
	err = checkUserFromClaims(ctx, images.Images[0].ID)
	if err != nil {
		ctx.JSON(400, gin.H{"error": 1})
	} else {
		err = db.DeleteImage(&images)
		if err != nil {
			fmt.Println(err.Error())
			ctx.JSON(400, gin.H{"error": 2})
		} else {
			ctx.JSON(200, gin.H{"error": 0})
		}
	}

}
func insert_properties(context *gin.Context) {
	var newProperty models.Property
	err := context.BindJSON(&newProperty)
	checkErr(err)
	db.InsertProperty(&newProperty)
	err = db.ImageInsert(newProperty.Images, newProperty.ID)
	checkErr(err)
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
	utils.EmailQueue.Push(newUser.Email, url)
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

func save(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("userID"))
	checkErr(err)
	propertyID, err := strconv.Atoi(ctx.Param("propertyID"))
	checkErr(err)

	err = db.SaveProperty(userID, propertyID)
	checkErr(err)
	ctx.JSON(200, gin.H{"error": 0})
}

func editDeleteProperties(ctx *gin.Context, fn func(*models.Property)) {
	var newProperty models.Property
	err := ctx.BindJSON(&newProperty)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = checkUserFromClaims(ctx, newProperty.ID)
	if err == nil {
		fn(&newProperty)
	} else {
		ctx.JSON(400, gin.H{"error": 1}) //User is not the Author of the Property
	}
	ctx.JSON(200, gin.H{"error": 0})
}

func profile(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	claims, err := utils.GetClaims(authHeader)
	checkErr(err)
	user := db.SelectProfile(claims.ID)
	ctx.IndentedJSON(200, user)
}

func editUser(ctx *gin.Context) {
	var user models.User
	id, err := getUserIDFromClaims(ctx)
	checkErr(err)
	err = ctx.BindJSON(&user)
	checkErr(err)
	if user.ID != id {
		ctx.JSON(400, gin.H{"error": 1})
	} else {
		db.EditUser(&user)
		ctx.IndentedJSON(200, gin.H{"error": 0})
	}
}

func selectUsersProperties(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("userID"))
	checkErr(err)
	properties := db.GetUsersProperties(int64(userID))
	ctx.IndentedJSON(http.StatusCreated, properties)
}

func selectLocations(ctx *gin.Context) {
	ctx.IndentedJSON(200, gin.H{"error": 0, "locations": db.SelectLocations()})
}
func getUserIDFromClaims(ctx *gin.Context) (int64, error) {
	authHeader := ctx.GetHeader("Authorization")
	claims, err := utils.GetClaims(authHeader)
	if err != nil {
		return 0, err
	}
	return claims.ID, nil
}
func checkUserFromClaims(ctx *gin.Context, propertyID int64) error {
	id, err := getUserIDFromClaims(ctx)
	checkErr(err)
	err = db.IsUserAuthorOfProperty(id, propertyID)
	return err
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
	go utils.MailSender()
	router := gin.Default()
	router.Static("/assets", "./assets")

	private := router.Group("/private")
	private.Use(Auth())
	public := router.Group("/public")

	public.GET("/verify/:userID/:key", verify_user)
	public.GET("/properties/:location/:max_price/:min_price/:max_rooms/:min_rooms/:max_internal_area/:min_internal_area", get_properties)
	public.POST("/user/", insert_user)
	public.POST("/login/", login)
	public.GET("/locations/", selectLocations)

	private.POST("/add/properties/", insert_properties)
	private.POST("/edit/properties/", func(ctx *gin.Context) { editDeleteProperties(ctx, db.EditProperty) })
	private.POST("/delete/properties/", func(ctx *gin.Context) { editDeleteProperties(ctx, db.DeleteProperty) })
	private.POST("/add/images/", addImages)       //Untested
	private.POST("/delete/images/", deleteImages) //Untested
	private.POST("/edit/profile/", editUser)
	private.GET("/profile/", profile)
	private.GET("/properties/:userID", selectUsersProperties)
	private.GET("/save/:userID/:propertyID", save)
	router.Run(settings.ServerName)
}
