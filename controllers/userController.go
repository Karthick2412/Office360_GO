package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"taskupdate/initializers"
	"taskupdate/models"
	"time"

	"gopkg.in/gomail.v2"
	"gorm.io/gorm"

	"github.com/thanhpk/randstr"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type SignupDTO struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var signUpBody SignupDTO
	json.NewDecoder(r.Body).Decode(&signUpBody)

	var usRef models.User
	u1 := uuid.Must(uuid.NewV4())
	usRef.Id = u1.String()
	usRef.Name = signUpBody.Name
	usRef.Email = signUpBody.Email

	hash, err := bcrypt.GenerateFromPassword([]byte(signUpBody.Password), 10)

	if err != nil {
		panic(err)
	}
	usRef.Password = string(hash)
	fmt.Println("User Details is", usRef)
	initializers.ConnectToDb().Create(&usRef)
	json.NewEncoder(w).Encode(usRef)

}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type LoginDTO struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var loginBody LoginDTO
	json.NewDecoder(r.Body).Decode(&loginBody)

	var usRef models.User
	record := initializers.ConnectToDb().Where("email = ?", loginBody.Email).First(&usRef)

	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		//fmt.Println("Hello")
		errorRes := "There is no such user with this mail id"
		json.NewEncoder(w).Encode(errorRes)
		return
	}
	// if usRef.Id == nil {
	// 	fmt.Println("User Not found", usRef)

	// }
	err := bcrypt.CompareHashAndPassword([]byte(usRef.Password), []byte(loginBody.Password))
	if err != nil {
		fmt.Println("Error is", err)
		errorRes := "Username Password mismatch"
		json.NewEncoder(w).Encode(errorRes)
		return

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": usRef.Id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})
	secret := "djaxcompany"
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err2 := token.SignedString([]byte(secret))
	if err2 != nil {
		fmt.Println(err2)
	}

	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   3600 * 24,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	fmt.Println(cookie)
	http.SetCookie(w, &cookie)
	json.NewEncoder(w).Encode(tokenString)
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ForgotEmail struct {
		Email string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&ForgotEmail)
	var userRef models.User
	record := initializers.ConnectToDb().Where("email = ?", ForgotEmail.Email).First(&userRef)
	fmt.Println(record)
	if errors.Is(record.Error, gorm.ErrRecordNotFound) {
		fmt.Println("Hello")
		errorRes := "There is no such user with this mail id"
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	token := randstr.String(6)

	res := "Your OTP IS " + token

	var passwordReset models.PasswordReset
	u1 := uuid.Must(uuid.NewV4())
	passwordReset.RequestId = u1.String()
	passwordReset.RequestedAt = time.Now().Local()
	passwordReset.Status = 0
	passwordReset.UserId = userRef.Id
	passwordReset.OTP = token
	initializers.ConnectToDb().Create(&passwordReset)
	fmt.Println("password Reset Details Are", passwordReset)
	m := gomail.NewMessage()
	m.SetHeader("From", "karthikgm2412@gmail.com")
	m.SetHeader("To", "karthikgm2412@gmail.com")
	m.SetHeader("Subject", "Password Reset OTP")
	m.SetBody("text/html", "Hello <b>"+userRef.Name+"</b><br> Your password reset OTP is "+token+" <br> <a style=\"text-decoration: none;\"href=\"http://localhost:8080/resetpassword\">Clickhere</a> to reset Your password")

	d := gomail.NewDialer("smtp.gmail.com", 587, "karthikgm2412@gmail.com", "azadxudwkzbestsy")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	json.NewEncoder(w).Encode(res)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var resetPasswordDTO struct {
		OTP             string `json:"otp"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	json.NewDecoder(r.Body).Decode(&resetPasswordDTO)
	var passwordReset models.PasswordReset
	record := initializers.ConnectToDb().Where("otp = ?", resetPasswordDTO.OTP).Find(&passwordReset)
	fmt.Println("PasswordReset Record is", passwordReset.RequestId)
	if errors.Is(record.Error, gorm.ErrRecordNotFound) || !passwordReset.RequestedAt.Add(time.Minute*5).Local().After(time.Now().Local()) {
		if errors.Is(record.Error, gorm.ErrRecordNotFound) || passwordReset.RequestId == "" {
			fmt.Println("Hello")
			errorRes := "OTP MISMATCH"
			json.NewEncoder(w).Encode(errorRes)
		} else if !passwordReset.RequestedAt.Add(time.Minute * 5).Local().After(time.Now().Local()) {
			errorRes := "OTP EXPIRED"
			fmt.Println("OTP expired")
			// var k time.Time = time.Now().Add(time.Minute * 5)
			// fmt.Printf("OTP expired %v and Type is %T\n", k, k)
			//fmt.Println("Requested Unix is", passwordReset.RequestedAt.Add(time.Minute*5).Local().UTC().UnixNano(), "Current Time + 5 min is", time.Now().Local().UTC().UnixNano())
			fmt.Println("Requested Unix is", passwordReset.RequestedAt.Add(time.Minute*5).Local(), "Current Time + 5 min is", time.Now().Local())
			// fmt.Println(k.UTC().Unix(), "and req at ", passwordReset.RequestedAt.UTC().Unix())
			// fmt.Printf("Requested Time is %v and Type is %T\n", passwordReset.RequestedAt, passwordReset.RequestedAt)

			json.NewEncoder(w).Encode(errorRes)
		}
		return
	}

	if resetPasswordDTO.Password == resetPasswordDTO.ConfirmPassword {
		var usRef models.User
		initializers.ConnectToDb().Where("Id = ?", passwordReset.UserId).First(&usRef)
		hash, err := bcrypt.GenerateFromPassword([]byte(resetPasswordDTO.Password), 10)

		if err != nil {
			panic(err)
		}
		usRef.Password = string(hash)
		initializers.ConnectToDb().Model(&models.User{}).Where("Id = ?", passwordReset.UserId).Update("password", string(hash))
		fmt.Println("Password Updated Successfully")
		initializers.ConnectToDb().Model(&models.PasswordReset{}).Where("request_id = ?", passwordReset.RequestId).Update("status", 1)
		initializers.ConnectToDb().Where("user_id = ? and status = 0", passwordReset.UserId).Delete(&passwordReset)
		fmt.Println("Reseted Datas Deleted Successfully")
	} else {
		errorRes := "Password and confirm password mismatch"
		fmt.Println("Password and confirm password mismatch")
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	//fmt.Println("O?S Requested Unix is", passwordReset.RequestedAt.Add(time.Minute*5).Local(), "Current Time + 5 min is", time.Now().Local())

	Res := "Password updated successfully"

	json.NewEncoder(w).Encode(Res)

}

func Validate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Println("Working fine")
	json.NewEncoder(w).Encode("You are Logged in")

}
