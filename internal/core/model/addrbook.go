package model

type ContactPhoneType string

const (
	ContactPhoneTypeMobile ContactPhoneType = "MOBILE"
	ContactPhoneTypeHome   ContactPhoneType = "HOME"
	ContactPhoneTypeWork   ContactPhoneType = "WORK"
)

type Contact struct {
	ID        string
	FirstName string
	LastName  string
	Phones    []*ContactPhone
}

type ContactPhone struct {
	PhoneType   ContactPhoneType
	PhoneNumber string
}

type ContactToSave struct {
	FirstName string
	LastName  string
	Phones    []*ContactPhoneToSave
}

type ContactPhoneToSave struct {
	PhoneType   ContactPhoneType
	PhoneNumber string
}
