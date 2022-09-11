package db

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"go-rest/models"
	"go-rest/utils"
	"strconv"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "630991"
	dbname   = "property"
)

var DB *sql.DB

func InitDatabase() {
	db, err := sql.Open("postgres", getInfo())
	if err != nil {
		fmt.Println("Couldn't connect to database")
		return
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(200)
	DB = db
}

func SelectLocations() []string {
	var location string
	var locations []string
	rows, err := DB.Query("SELECT DISTINCT location FROM Properties")
	check(err)
	for rows.Next() {
		rows.Scan(&location)
		locations = append(locations, location)
	}
	return locations
}
func rowsToProperties(rows *sql.Rows) []models.PropertyWithImage {
	var properties []models.PropertyWithImage
	var selectProperty models.PropertyWithImage
	for rows.Next() {
		property := models.Property{}
		rows.Scan(&property.ID,
			&property.Price,
			&property.IsForSale,
			&property.NumberOfRooms,
			&property.Location,
			&property.Address,
			&property.InternalArea,
			&property.Title,
			&property.Description,
			&property.PublishDate,
			&property.AuthorID,
			&property.Orienter,
			&property.PropertyType,
			&property.IsPublished)
		var empty []string
		property.Images = empty
		images := GetPropertyImages(property.ID)
		selectProperty = models.PropertyWithImage{
			Property:       property,
			PropertyImages: images,
		}
		properties = append(properties, selectProperty)
	}
	return properties
}
func SelectSavedProperties(userID int64) []models.PropertyWithImage {
	rows, err := DB.Query(`SELECT * FROM Properties WHERE ID IN (SELECT propertyID FROM Saved WHERE userID=$1)`, userID)
	check(err)
	return rowsToProperties(rows)
}
func SelectProperty(location string, max_price int, min_price int, max_rooms int, min_rooms int,
	max_internal_area int, min_internal_area int, is_for_sale bool) []models.PropertyWithImage {
	rows, err := DB.Query(`
	SELECT * FROM Properties WHERE 
	price < $1 AND price > $2 AND
	numberOfRooms < $3 AND numberOfRooms > $4 AND
	internalArea < $5 AND internalArea > $6 AND
	location = $7 AND isPublished = true AND isForSale=$8`, max_price,
		min_price,
		max_rooms,
		min_rooms,
		max_internal_area,
		min_internal_area,
		location,
		is_for_sale)
	check(err)
	return rowsToProperties(rows)
}
func GetUsersProperties(userID int64) []models.PropertyWithImage {
	rows, err := DB.Query("SELECT * FROM Properties WHERE authorID=$1", userID)
	check(err)
	return rowsToProperties(rows)
}

func EditProperty(property *models.Property) {
	stmt, err := DB.Prepare(`UPDATE Properties SET 
	price=$1 , isForSale= $2, 
	numberOfRooms = $3 , location = $4,
	address = $5 , internalArea = $6, 
	title = $7 , description = $8, 
	orienter = $9 , propertyType = $10, isPublished = $11  
	WHERE ID=$12`)
	check(err)
	_, err = stmt.Exec(property.Price,
		property.IsForSale,
		property.NumberOfRooms,
		property.Location,
		property.Address,
		property.InternalArea,
		property.Title,
		property.Description,
		property.Orienter,
		property.PropertyType,
		property.IsPublished,
		property.ID)
	check(err)
}
func DeleteProperty(property *models.Property) {
	stmt, err := DB.Prepare("DELETE FROM Properties WHERE ID=$1")
	check(err)
	stmt.Exec(property.ID)
}
func InsertProperty(property *models.Property) {
	id := GetSequenceValue("propertiesid_sequence")
	property.ID = id
	stmt, err := DB.Prepare(`INSERT INTO Properties(ID, price, isForSale, 
		numberOfRooms, location, address, 
		internalArea,title, description, publishDate, authorID, orienter, propertyType, isPublished)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10, $11, $12, $13, $14)`)
	check(err)
	_, err = stmt.Exec(
		property.ID,
		property.Price,
		property.IsForSale,
		property.NumberOfRooms,
		property.Location,
		property.Address,
		property.InternalArea,
		property.Title,
		property.Description,
		property.PublishDate,
		property.AuthorID,
		property.Orienter,
		property.PropertyType,
		property.IsPublished,
	)
	check(err)
}

func IsUserAuthorOfProperty(userID, propertyID int64) error {
	var id int64
	row := DB.QueryRow("SELECT ID FROM Properties WHERE authorID = $1 AND ID=$2", userID, propertyID)
	err := row.Scan(&id)
	return err
}

func GetPropertyImages(propertyID int64) []models.PropertyImage {
	var images []models.PropertyImage
	var image models.PropertyImage
	rows, err := DB.Query("SELECT ID, fileName FROM PropertyImages WHERE PropertyID = $1", propertyID)
	if err != nil {
		fmt.Println(err.Error())
		return images
	}
	for rows.Next() {
		rows.Scan(&image.ID, &image.FileName)
		images = append(images, image)
	}
	return images
}
func GetSequenceValue(sequenceName string) int64 {
	var id int64
	err := DB.QueryRow("SELECT nextval($1) ", sequenceName).Scan(&id)
	if err != nil {
		fmt.Println(err.Error())
		return 0
	}
	return id
}
func InsertPropertyImage(id, propertyID int64, fileName string) {
	stmt, err := DB.Prepare("INSERT INTO PropertyImages(ID, propertyID, fileName) VALUES($1, $2, $3)")
	check(err)
	_, err = stmt.Exec(id, propertyID, fileName)
	check(err)
}

func InsertUser(user models.User) int64 {
	stmt, err := DB.Prepare(`INSERT INTO Users(email, companyName, isAgent, phoneNumber, password, isMailVerified) 
	VALUES($1,$2,$3,$4,$5,$6)`)
	check(err)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	check(err)
	pass := base64.StdEncoding.EncodeToString([]byte(hashedPassword))
	stmt.Exec(
		user.Email,
		user.CompanyName,
		user.IsAgent,
		user.PhoneNumber,
		pass,
		false,
	)
	id := getUserIDFromMail(user.Email)
	return id
}

func VerifyUser(key string, userID int64) int64 {
	stmt, err := DB.Prepare(`UPDATE Users SET isMailVerified = true 
	WHERE ID IN (SELECT userID FROM TemporaryKeys WHERE key=$1 AND userID = $2)`)
	check(err)
	result, err := stmt.Exec(key, userID)
	check(err)
	id, err := result.RowsAffected()
	check(err)
	return id
}

func InsertTemporaryKey(userID int64, runes string) {
	stmt, err := DB.Prepare("INSERT INTO TemporaryKeys(userID, key) VALUES($1, $2)")
	check(err)
	stmt.Exec(userID, runes)
}

func CanUserLogIn(password string, email string) int64 {
	var is_verified bool
	var pass string
	var _id int64
	e := DB.QueryRow("SELECT ID,password,isMailVerified FROM Users WHERE email = $1", email).Scan(&_id, &pass, &is_verified)
	if e != nil {
		return 0
	}
	hash, err := base64.StdEncoding.DecodeString(pass)
	check(err)
	e = bcrypt.CompareHashAndPassword(hash, []byte(password))
	check(e)
	return _id
}

func getUserIDFromMail(mail string) int64 {
	var id int64
	res := DB.QueryRow(`SELECT ID FROM Users WHERE email=$1`, mail)
	res.Scan(&id)
	return id
}

func getInfo() string {
	val := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return val
}

func check(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}
func SaveProperty(userID int, propertyID int) error {
	stmt, err := DB.Prepare("INSERT INTO Saved(userID, propertyID) VALUES($1,$2)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userID, propertyID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteImage(images *models.PropertyImages) error {
	stmt, err := DB.Prepare("DELETE FROM PropertyImages WHERE propertyID=$1 AND fileName = $2")
	if err != nil {
		return err
	}
	for _, image := range images.Images {
		stmt.Exec(image.ID, image.FileName)
	}
	return nil
}

func ImageInsert(images []string, propertyID int64) error {
	var nextval int64
	for _, image := range images {
		nextval = GetSequenceValue("propertyimages_sequence")
		decoded, err := base64.StdEncoding.DecodeString(image)
		if err == nil {
			filename, err := utils.SaveImageToDisk(strconv.FormatInt(nextval, 16), decoded) // Second bug occurs here
			if err == nil {
				InsertPropertyImage(nextval, propertyID, filename)
			} else {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func SelectProfile(userID int64) models.User {
	var user models.User
	err := DB.QueryRow("SELECT email, companyName, isAgent,phoneNumber FROM Users WHERE ID = $1", userID).Scan(
		&user.Email,
		&user.CompanyName,
		&user.IsAgent,
		&user.PhoneNumber,
	)
	check(err)
	user.ID = userID
	user.Password = ""
	user.IsMailVerified = true
	return user
}

func EditUser(user *models.User) {
	stmt, err := DB.Prepare("UPDATE Users SET companyName=$1, phoneNumber=$2, password= $3")
	check(err)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	check(err)
	pass := base64.StdEncoding.EncodeToString([]byte(hashedPassword))
	_, err = stmt.Exec(user.CompanyName, user.PhoneNumber, pass)
	check(err)
}
