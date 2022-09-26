package models

import (
	"time"

	"golang.org/x/oauth2"
)

type User struct {
	ID        uint64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func UserUpdateOrCreate(token *oauth2.Token) *User {
	// Store details in OauthUser table
	oauthUser := OauthUserFBStoreDetails(token)

	user := User{}
	DB.Where("email = ?", oauthUser.Email).First(&user)
	if user.ID == 0 {
		user := User{Name: oauthUser.Name, Email: oauthUser.Email,
			CreatedAt: time.Now(), UpdatedAt: time.Now()}
		DB.Create(&user)
	}
	// Update the user id in OauthUser
	oauthUser.UpdateUserID(&user)
	return (&user)
}

func UserGetByID(userID uint64) *User {
	user := User{}
	DB.Where("id = ?", userID).First(&user)
	return (&user)
}
