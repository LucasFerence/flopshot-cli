package edit

const UserType = "user"

type User struct {
	Id    string `json:"_id" mapstructure:"_id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (user User) Label() string {
	return user.Email
}

func init() {
	RegisterType(UserType, User{})
}
