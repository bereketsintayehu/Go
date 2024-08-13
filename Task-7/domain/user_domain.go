package domain

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole int

const (
	UserR UserRole = iota
	Admin
	Super
)

func (ur UserRole) String() string {
	return [...]string{"User", "Admin", "Super Admin"}[ur]
}

func (ur UserRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(ur.String())
}

func (ur *UserRole) UnmarshalJSON(data []byte) error {
	var roleStr string
	if err := json.Unmarshal(data, &roleStr); err != nil {
		return err
	}

	switch roleStr {
	case "Admin":
		*ur = Admin
	case "User":
		*ur = UserR
	case "Super Admin":
		*ur = Super
	default:
		return errors.New("invalid UserRole")
	}

	return nil
}


type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email    string             `json:"email" bind:"required" validate:"required & unique"`
	Password string             `json:"password" validate:"required & min=6 & max=32"`
	Role     UserRole           `json:"role"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		Alias
		Password string `json:"password,omitempty"`
	}{
		Alias:    (Alias)(*u),
		Password: "",
	})
}

type UserRepository interface{
	EnsureSuperAdmin() error
	FindAll() ([]User, error)
	FindByID(id primitive.ObjectID) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAllAdmin() ([]User, error)
	CreateUser(user User) (*User, error)
	UpdateUser(user User) (*User, error)
	DeleteUser(id primitive.ObjectID) (*User, error)
}

