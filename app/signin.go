package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUserID struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `josn:"verified_email"`
	Picture       string `json:"picture"`
}

var googleOauthConfig = oauth2.Config{
	// RedirectURL:  "http://localhost:8080/auth/google/callback",
	RedirectURL:  os.Getenv("DOMAIN_NAME") + "/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET_KEY"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := generateStateOauthCookie(w)

	url := googleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
func generateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

func googleAuthCallback(w http.ResponseWriter, r *http.Request) {
	oauthstate, _ := r.Cookie("oauthstate")

	if r.FormValue("state") != oauthstate.Value {
		// //공격 시도 log
		// log.Printf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, r.FormValue("state"))
		// //해커가 감잡기 힘들게 그냥 기본 경로로, 정보 안 주려고
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		// return

		errMsg := fmt.Sprintf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, r.FormValue("state"))
		log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	data, err := getGoogleUserInfo(r.FormValue("code"))
	if err != nil {
		// log.Println(err.Error())
		// http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		// return

		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Store ID info into Session cookie
	var userInfo GoogleUserID
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := store.Get(r, "session")
	// Set some session values.
	session.Values["id"] = userInfo.ID
	// session.Values[42] = 43 다른거 저장해도 됨.
	// Save it before we write to the response/return from the handler.
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func getGoogleUserInfo(code string) ([]byte, error) {
	//context : thread safe한 저장소
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s", err.Error())
	}

	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	return ioutil.ReadAll(resp.Body)
}
