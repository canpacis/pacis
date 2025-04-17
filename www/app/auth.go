package app

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/canpacis/pacis/pages"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauthConfig *oauth2.Config
)

func InitAuth() {
	oauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("OAuthCallbackURL"),
		ClientID:     os.Getenv("GoogleOAuthClientID"),
		ClientSecret: os.Getenv("GoogleOAuthClientSecret"),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}
}

func randstate() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type User struct {
	UserID  string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (u User) ID() string {
	return u.UserID
}

//pacis:authentication
func AuthHandler(r *http.Request) (*User, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil
		}
		return nil, err
	}

	client := oauthConfig.Client(r.Context(), &oauth2.Token{AccessToken: cookie.Value})
	// TODO: Potentially cache this response
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	user := new(User)
	if err := json.Unmarshal(data, user); err != nil {
		return nil, err
	}

	return user, nil
}

//pacis:page path=/auth/login
func LoginPage(ctx *pages.PageContext) I {
	state := randstate()
	url := oauthConfig.AuthCodeURL(state)

	ctx.SetCookie(&http.Cookie{
		Name:     "auth_state",
		Value:    state,
		Path:     "/auth",
		Secure:   true,
		HttpOnly: true,
		// Give the state cookie 5 minutes to expire
		Expires:  time.Now().Add(time.Minute * 5),
		SameSite: http.SameSiteNoneMode,
	})

	return Div(
		Class("container flex-1 flex items-center justify-center flex-col gap-4"),

		H1(Class("text-3xl font-semibold"), Text("Welcome to Pacis")),
		Button(
			ButtonSizeLg,
			Href(url),
			Replace(A),

			GoogleIcon(),
			Text("Login with Google"),
		),
	)
}

//pacis:page path=/auth/logout
func LogoutPage(ctx *pages.PageContext) I {
	// Remove the cookie
	ctx.SetCookie(&http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   -1,
	})
	return ctx.Redirect("/")
}

//pacis:page path=/auth/callback
func AuthCallbackPage(ctx *pages.PageContext) I {
	r := ctx.Request()
	state := r.FormValue("state")

	cookie, err := ctx.GetCookie("auth_state")
	if err != nil {
		return ctx.Redirect("/?error=invalid_state")
	}
	if state != cookie.Value {
		return ctx.Redirect("/?error=invalid_state")
	}

	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return ctx.Redirect("/?error=exchange_fail")
	}

	ctx.SetCookie(&http.Cookie{
		Name:     "auth_token",
		Value:    token.AccessToken,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return ctx.Redirect("/")
}
