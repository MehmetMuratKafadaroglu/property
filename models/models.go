package models

type Property struct {
	ID            int    `json:"id"`
	Price         int    `json:"price"`
	IsForSale     bool   `json:"isForSale"`
	NumberOfRooms int    `json:"numberOfRooms"`
	Location      string `json:"location"`
	Address       string `json:"address"`
	InternalArea  int    `json:"internalArea"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	PublishDate   string `json:"publishDate"`
	AuthorID      int    `json:"authorID"`
	Orienter      string `json:"orienter"`
	PropertyType  string `json:"propertyType"`
	//	Images        []string `json:"image"`
}

type User struct {
	ID             int    `json:"id"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	IsMailVerified bool   `json:"isMailVerified"`
	CompanyName    string `json:"companyName"`
	IsAgent        bool   `json:"isAgent"`
	PhoneNumber    int    `json:"phoneNumber"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
