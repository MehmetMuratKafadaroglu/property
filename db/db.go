package db

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"go-rest/models"

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
func SelectProperty(location string, max_price int, min_price int, max_rooms int, min_rooms int,
	max_internal_area int, min_internal_area int) []models.Property {

	var properties []models.Property

	rows, err := DB.Query(`
	SELECT * FROM Properties WHERE 
	price < $1 AND price > $2 AND
	numberOfRooms < $3 AND numberOfRooms > $4 AND
	internalArea < $5 AND internalArea > $6 AND
	location = $7`, max_price,
		min_price,
		max_rooms,
		min_rooms,
		max_internal_area,
		min_internal_area,
		location)
	check(err)
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
			&property.PropertyType)
		properties = append(properties, property)
	}
	return properties
}

func InsertProperty(property *models.Property) {
	stmt, err := DB.Prepare(`INSERT INTO Properties(price, isForSale, 
		numberOfRooms, location, address, 
		internalArea,title, description, publishDate, authorID, orienter, propertyType)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10, $11, $12)
		`)
	check(err)
	_, err = stmt.Exec(
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
		property.PropertyType)
	check(err)
}

func InsertUser(user models.User) int64 {
	insertUserHelper(user)
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

func insertUserHelper(user models.User) {
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
