package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// TODO: Load google oauth config from database on load
var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8000/auth/google/callback",
	ClientID:     "<PUT CLIENT ID HERE>",
	ClientSecret: "<PUT SECRET HERE>",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

//TODO: Load JWT Secret key from database on load, prevent people from keeping this crappy default
var jwtSecret []byte = []byte("thisisahorriblesecretkey")

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type UserData struct {
	Id          string
	Email       string
	Name        string
	Given_name  string
	Family_name string
}

func OauthGoogleLogin(w http.ResponseWriter, r *http.Request) {

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)

	/*
	   AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
	   validate that it matches the the state query parameter on your redirect callback.
	*/
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func OauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Read oauthState from Cookie
	oauthState, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthState.Value {
		log.Println("invalid oauth google state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	data, err := getUserDataFromGoogle(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var userData UserData
	json.Unmarshal(data, &userData)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    userData.Id,
		"email": userData.Email,
	})
	tokenString, error := token.SignedString(jwtSecret)
	if error != nil {
		fmt.Println(error)
	}

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "authtoken", Value: tokenString, Expires: expiration, Path: "/", SameSite: http.SameSiteLaxMode}
	http.SetCookie(w, &cookie)

	fmt.Fprintf(w, "Token: %s\n", data)

}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(code string) ([]byte, error) {
	// Use code to get token and get user info from Google.

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

type User struct {
	Id    string
	Email string
}

func AddUserContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		t, _ := r.Cookie("authtoken")
		if t != nil {
			token, _ := jwt.Parse(t.Value, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There was an error")
				}
				return jwtSecret, nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				user := User{
					Id:    claims["id"].(string),
					Email: claims["email"].(string),
				}

				ctx := context.WithValue(r.Context(), "user", &user)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				ctx := context.WithValue(r.Context(), "user", &User{})
				next.ServeHTTP(w, r.WithContext(ctx))
			}

		} else {
			ctx := context.WithValue(r.Context(), "user", &User{})
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
