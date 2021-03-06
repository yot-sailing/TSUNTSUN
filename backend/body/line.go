package body

// LINEAPIでアクセストークンを用いて取得できるユーザー情報
type LINEUser struct {
	UserID        string
	DisplayName   string
	PictureUrl    string
	StatusMessage string
}

// LINEAPIでIDトークンを用いて取得できるユーザー情報
type VerifyResponseBody struct {
	Iss     string
	Sub     string
	Aud     string
	Exp     int
	Iat     int
	Nonce   string
	Amr     []string
	Name    string
	Picture string
	Email   string
}

type RevokeRequestBody struct {
	ClientID      string
	ClientSercret string
	AccessToken   string
}
