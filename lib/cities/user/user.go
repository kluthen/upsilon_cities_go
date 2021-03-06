package user

import (
	"regexp"
	"time"
	"upsilon_cities_go/lib/misc/config/system"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int
	Login     string
	Email     string
	Password  string `json:"-"`
	LastLogin time.Time

	// admin

	NeedNewPassword bool
	Enabled         bool
	Admin           bool
}

// courtesy to https://gowebexamples.com/password-hashing/

//New create a new user.
func New() *User {
	usr := new(User)

	usr.NeedNewPassword = false
	usr.Enabled = system.GetBool("user_enabled_by_default", false)
	usr.Admin = system.GetBool("user_admin_by_default", false)
	return usr
}

//CheckPassword password validate regex
func CheckPassword(password string) bool {
	re := regexp.MustCompile("[A-Za-z0-9@#$%^!&+=]{8,}")
	return re.Match([]byte(password))
}

//CheckLogin login validate regex
func CheckLogin(login string) bool {
	re := regexp.MustCompile("[A-Za-z][A-Za-z0-9_-]{3,}")
	return re.Match([]byte(login))
}

//CheckMail mail validate regex
func CheckMail(mail string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.Match([]byte(mail))
}

//HashPassword generate a hash based on nice password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

//CheckPasswordHash will check provided password against hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

//PrettyLastLogin stringify last login.
func (user *User) PrettyLastLogin() string {
	return user.LastLogin.Format(time.RFC3339)
}
