package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"
)

type OauthUser struct {
	ID           uint64
	UserID       uint64
	Provider     string // Like Facebook, Google, Twitter etc.
	UUID         string
	Name         string
	Email        string
	AccessToken  string
	ExpiresAt    time.Time
	TokenType    string
	RefreshToken string
	ProfilePic   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func OauthUserFBStoreDetails(token *oauth2.Token) *OauthUser {
	// Get the details of the user
	resp, _ := http.Get("https://graph.facebook.com/me?access_token=" +
		url.QueryEscape(token.AccessToken) +
		"&fields=id,name,email")

	defer resp.Body.Close()
	me := map[string]interface{}{}

	err := json.NewDecoder(resp.Body).Decode(&me)
	if err != nil {
		fmt.Println("json decode error", "error", err)
	}

	uuid := me["id"].(string)
	email := me["email"].(string)
	name := me["name"].(string)
	provider := "Facebook"

	// Get Profile pic
	// TODO
	profilePic := getFBProfilePic(uuid, token.AccessToken)

	oauthUser := OauthUser{}
	DB.Where("provider = ? and email = ?", provider, email).First(&oauthUser)

	if oauthUser.ID != 0 {
		// Update the oauth user
		DB.Model(&oauthUser).Updates(OauthUser{Name: name, ProfilePic: profilePic,
			AccessToken: token.AccessToken,
			ExpiresAt:   token.Expiry, TokenType: token.TokenType,
			RefreshToken: token.RefreshToken,
			UpdatedAt:    time.Now()})
	} else {
		// Create new entry
		oauthUser = OauthUser{Provider: provider, Email: email, UUID: uuid,
			UserID: 0, // There is no user yet
			Name:   name, ProfilePic: profilePic,
			AccessToken: token.AccessToken,
			ExpiresAt:   token.Expiry, TokenType: token.TokenType,
			RefreshToken: token.RefreshToken,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now()}
		DB.Create(&oauthUser)
	}
	return &oauthUser
}

func (oauthUser *OauthUser) UpdateUserID(user *User) {
	oauthUser.UserID = user.ID
	DB.Save(&oauthUser)
}

func getFBProfilePic(uuid string, accessToken string) string {
	profilePicApi := "https://graph.facebook.com/" + uuid +
		"?fields=picture.width(720).height(720)&access_token=" + accessToken

	resp, _ := http.Get(profilePicApi)
	defer resp.Body.Close()
	picResp := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&picResp); err != nil {
		fmt.Println("json decode error", err)
	}
	return picResp["picture"].(map[string]interface{})["data"].(map[string]interface{})["url"].(string)
}
