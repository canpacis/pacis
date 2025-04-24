package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/canpacis/pacis/pages"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	cachedb     *redis.Client
	oauthConfig *oauth2.Config
)

type CacheStorage struct {
	db *redis.Client
}

func (cs *CacheStorage) Get(key string, val any) error {
	return cs.db.Get(context.Background(), key).Scan(val)
}

func (cs *CacheStorage) Set(key string, val any) error {
	return cs.db.Set(context.Background(), key, val, time.Hour).Err()
}

func Init() {
	oauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("OAuthCallbackURL"),
		ClientID:     os.Getenv("GoogleOAuthClientID"),
		ClientSecret: os.Getenv("GoogleOAuthClientSecret"),
		Scopes:       []string{"email", "profile"},
		Endpoint:     google.Endpoint,
	}

	cachedb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RedisAddress"),
		Username: os.Getenv("RedisUsername"),
		Password: os.Getenv("RedisPassword"),
		DB:       0,
	})
}

func randstate() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

type User struct {
	UserID   string `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Name     string `json:"name,omitempty"`
	Picture  string `json:"picture,omitempty"`
	LoggedIn bool   `json:"logged_in"`
}

func (u User) ID() string {
	return u.UserID
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

//pacis:middleware label=authentication
func AuthHandler(r *http.Request) (*User, error) {
	if oauthConfig == nil {
		return nil, errors.New("no oauth2 config")
	}
	if cachedb == nil {
		return nil, errors.New("no cachedb")
	}

	cookie, err := r.Cookie("auth_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, nil
		}
		return nil, err
	}

	client := oauthConfig.Client(r.Context(), &oauth2.Token{AccessToken: cookie.Value})

	user := new(User)
	err = cachedb.Get(r.Context(), cookie.Value).Scan(user)
	if err == nil {
		return user, nil
	}

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, user); err != nil {
		return nil, err
	}

	user.LoggedIn = true
	err = cachedb.Set(r.Context(), cookie.Value, user, time.Hour).Err()
	if err != nil {
		slog.Error("failed to cache user data", "error", err)
	}

	return user, nil
}

type LoginPage struct {
	User *User `context:"user"`
}

//pacis:page path=/auth/login middlewares=auth
func (p *LoginPage) Page(ctx *pages.Context) I {
	// ctx.SetTitle("Login | Pacis")

	if p.User != nil {
		return pages.Redirect(ctx, "/")
	}

	state := randstate()
	url := oauthConfig.AuthCodeURL(state)

	return Frag(
		pages.Cookie(
			ctx,

			&http.Cookie{
				Name:     "auth_state",
				Value:    state,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				// Give the state cookie 5 minutes to expire
				Expires:  time.Now().Add(time.Minute * 5),
				SameSite: http.SameSiteNoneMode,
			},
		),

		Div(
			Class("container flex-1 flex items-center justify-center flex-col gap-4"),

			H1(Class("text-3xl font-semibold"), Text("Welcome to Pacis")),
			Button(
				ButtonSizeLg,
				Href(url),
				Replace(A),
				Class("!rounded-full"),

				GoogleIcon(),
				Text("Login with Google"),
			),
		),
	)
}

//pacis:page path=/auth/logout middlewares=auth
func LogoutPage(ctx *pages.Context) I {
	return Frag(
		pages.Cookie(
			ctx,

			&http.Cookie{
				Name:     "auth_token",
				Value:    "",
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteNoneMode,
				MaxAge:   -1,
			},
		),
		pages.Redirect(ctx, "/"),
	)
}

type AuthCallbackPage struct {
	Code        string `query:"code"`
	QueryState  string `query:"state"`
	CookieState string `cookie:"auth_state"`
}

//pacis:page path=/auth/callback middlewares=auth
func (p *AuthCallbackPage) Page(ctx *pages.Context) I {
	fmt.Println(p.Code)
	fmt.Println(p.QueryState)
	fmt.Println(p.CookieState)
	if p.QueryState != p.CookieState {
		return pages.Error(ctx, http.StatusBadRequest, &AppError{Code: InvalidAuthStateError, Err: errors.New("app error")})
	}

	token, err := oauthConfig.Exchange(ctx, p.Code)
	if err != nil {
		return pages.Error(ctx, http.StatusBadRequest, &AppError{Code: AuthExchangeError, Err: errors.New("app error")})
	}

	return Frag(
		pages.Cookie(
			ctx,

			&http.Cookie{
				Name:     "auth_token",
				Value:    token.AccessToken,
				Path:     "/",
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			},
		),
		pages.Redirect(ctx, "/"),
	)
}
