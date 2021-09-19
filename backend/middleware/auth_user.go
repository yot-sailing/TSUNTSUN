package middleware

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/yot-sailing/TSUNTSUN/body"
)

type VerifyAccessTokenResponseBody struct {
	scope      string
	client_id  string
	expired_in int
}

func AuthUser(accessToken string) (user body.LINEUser, err error) {
	// アクセストークンの有効性のチェック
	accessTokenStatus, accessTokenResponse := VerifyAccessToken(accessToken)
	if accessTokenStatus != 200 {
		fmt.Println("アクセストークンが有効でありません。")
		fmt.Println("status : " + strconv.Itoa(accessTokenStatus))
		fmt.Println(accessTokenResponse)
		return user, err
	}

	user, err = GetLINEProfile(accessToken)
	if err != nil {
		return user, err
	}
	user := userController.PrepareUser(userID, userName)
	return user
}

func VerifyAccessToken(access_token string) (int, VerifyAccessTokenResponseBody) {
	endpoint := "https://api.line.me/oauth2/v2.1/verify?access_token=" + access_token[7:]
	resp, err := http.Get(endpoint)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	var verifyAccessTokenResponseBody VerifyAccessTokenResponseBody
	err = json.Unmarshal(byteArray, &verifyAccessTokenResponseBody)
	if err != nil {
		fmt.Println(err)
	}

	return resp.StatusCode, verifyAccessTokenResponseBody
}

func GetLINEProfile(access_token string) (body.LINEUser, error) {
	endpoint := "https://api.line.me/v2/profile"
	var line_user_profile body.LINEUser

	req, _ := http.NewRequest("GET", endpoint, nil)
	req.Header.Set("Authorization", access_token)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return line_user_profile, err
	}
	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return line_user_profile, err
	}

	err = json.Unmarshal(byteArray, &line_user_profile)
	if err != nil {
		return line_user_profile, err
	}
	return line_user_profile, nil
}
