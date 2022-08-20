package models

type Property struct {
	ID            int64    `json:"id"`
	Price         int      `json:"price"`
	IsForSale     bool     `json:"isForSale"`
	NumberOfRooms int      `json:"numberOfRooms"`
	Location      string   `json:"location"`
	Address       string   `json:"address"`
	InternalArea  int      `json:"internalArea"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	PublishDate   string   `json:"publishDate"`
	AuthorID      int      `json:"authorID"`
	Orienter      string   `json:"orienter"`
	PropertyType  string   `json:"propertyType"`
	Images        []string `json:"images"`
	IsPublished   bool     `json:"isPublished"`
}

type User struct {
	ID             int64  `json:"id"`
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

type PropertyWithImage struct {
	Property       Property        `json:"property"`
	PropertyImages []PropertyImage `json:"propertyImages"`
}

type PropertyImage struct {
	ID       int64  `json:"id"`
	FileName string `json:"fileName"`
}

type PropertyImages struct {
	Images []PropertyImage `json:"propertyImages"`
}
type AddPropertyImages struct {
	PropertyID int64    `json:"propertyImage"`
	Images     []string `json:"images"`
}
