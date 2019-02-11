package domain

type (
	Record struct {
		ID        string
		FirstName string
		LastName  string
		Email     string
		Phone     string
	}
)

func NewMessageFromCSV(id, firstName, lastName, email, phone string) *Record {
	return &Record{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}
}
