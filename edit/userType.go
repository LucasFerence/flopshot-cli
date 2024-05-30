package edit

const UserType = "user"

type User struct {
	Id    string `json:"_id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func init() {
	RegisterType("user", User{})
}
