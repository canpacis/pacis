package app

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/authorizerdev/authorizer-go"
	"github.com/canpacis/pacis/pages"
	. "github.com/canpacis/pacis/ui/components"
	. "github.com/canpacis/pacis/ui/html"
	"github.com/redis/go-redis/v9"
)

var (
	cachedb *redis.Client
	auth    *authorizer.AuthorizerClient
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

func Init() error {
	// oauthConfig = &oauth2.Config{
	// 	RedirectURL:  os.Getenv("OAUTH_CALLBACK_URL"),
	// 	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	// 	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	// 	Scopes:       []string{"email", "profile"},
	// 	Endpoint:     google.Endpoint,
	// }

	cachedb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	var err error
	auth, err = authorizer.NewAuthorizerClient(
		os.Getenv("AUTHORIZER_ID"),
		os.Getenv("AUTHORIZER_URL"),
		os.Getenv("OAUTH_CALLBACK_URL"),
		map[string]string{},
	)
	if err != nil {
		return fmt.Errorf("failed to instantiate authorizer: %w", err)
	}
	return nil
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
	// if oauthConfig == nil {
	// 	return nil, errors.New("no oauth2 config")
	// }
	// if cachedb == nil {
	// 	return nil, errors.New("no cachedb")
	// }

	// cookie, err := r.Cookie("auth_token")
	// if err != nil {
	// 	if errors.Is(err, http.ErrNoCookie) {
	// 		return nil, nil
	// 	}
	// 	return nil, err
	// }

	// client := oauthConfig.Client(r.Context(), &oauth2.Token{AccessToken: cookie.Value})

	// user := new(User)
	// err = cachedb.Get(r.Context(), cookie.Value).Scan(user)
	// if err == nil {
	// 	return user, nil
	// }

	// resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	// if err != nil {
	// 	return nil, err
	// }
	// defer resp.Body.Close()

	// data, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, err
	// }

	// if err := json.Unmarshal(data, user); err != nil {
	// 	return nil, err
	// }

	// user.LoggedIn = true
	// err = cachedb.Set(r.Context(), cookie.Value, user, time.Hour).Err()
	// if err != nil {
	// 	slog.Error("failed to cache user data", "error", err)
	// }

	return nil, nil
}

type LoginPage struct {
	User *User `context:"user"`
}

//pacis:page path=/auth/login middlewares=auth
func (p *LoginPage) Page(ctx *pages.Context) I {
	// if p.User != nil {
	// 	return pages.Redirect(ctx, "/")
	// }

	// state := randstate()
	// url := oauthConfig.AuthCodeURL(state)

	// pages.SetCookie(
	// 	ctx,

	// 	&http.Cookie{
	// 		Name:     "auth_state",
	// 		Value:    state,
	// 		Path:     "/",
	// 		Secure:   true,
	// 		HttpOnly: true,
	// 		// Give the state cookie 5 minutes to expire
	// 		Expires:  time.Now().Add(time.Minute * 5),
	// 		SameSite: http.SameSiteNoneMode,
	// 	},
	// )

	return Div(
		Class("container flex-1 flex items-center justify-center flex-col gap-4"),

		H1(Class("text-3xl font-semibold"), Text("Welcome to Pacis")),
		Form(
			Input(
				Placeholder("Email"),
				Placeholder("Password"),
			),
		),
		// Button(
		// 	ButtonSizeLg,
		// 	Href(url),
		// 	Replace(A),
		// 	Class("!rounded-full"),

		// 	GoogleIcon(),
		// 	Text("Login with Google"),
		// ),
	)
}

type Signup struct {
	User     *User  `context:"user"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

//pacis:page path=/auth/signup middlewares=auth
func (s *Signup) Page(ctx *pages.Context) I {
	if s.User != nil {
		return pages.Redirect(ctx, "/")
	}

	return Div(
		Class("container flex-1 flex justify-center items-center"),

		Div(
			Class("max-w-80 flex flex-col gap-4"),

			H1(Class("font-semibold"), Text("Sign up to Pacis")),
			Form(
				Action("/auth/signup"),
				Method("POST"),
				Class("flex flex-col gap-2 w-full"),

				Input(
					Type("email"),
					Attr("required"),
					Placeholder("Email"),
					Name("email"),
				),
				Input(
					Type("password"),
					Attr("required"),
					Placeholder("Password"),
					Name("password"),
				),
				Button(
					Type("submit"),

					Text("Sign Up"),
				),
			),
		),
	)
}

//pacis:action path=/auth/signup method=post middlewares=auth
func (s *Signup) Action(ctx *pages.Context) I {
	token, err := auth.SignUp(&authorizer.SignUpInput{
		Email:           &s.Email,
		Password:        s.Password,
		ConfirmPassword: s.Password,
	})
	fmt.Println(token, err)

	// return pages.Redirect(ctx, "/")
	return Frag()
}

//pacis:page path=/auth/logout middlewares=auth
func LogoutPage(ctx *pages.Context) I {
	pages.SetCookie(
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
	)

	return pages.Redirect(ctx, "/")
}

// type AuthCallbackPage struct {
// 	Code        string `query:"code"`
// 	QueryState  string `query:"state"`
// 	CookieState string `cookie:"auth_state"`
// }

// //pacis:page path=/auth/callback middlewares=auth
// func (p *AuthCallbackPage) Page(ctx *pages.Context) I {
// 	if p.QueryState != p.CookieState {
// 		return pages.Error(ctx, NewAppError(InvalidAuthStateError, ErrGenericAppError, http.StatusBadRequest))
// 	}

// 	token, err := oauthConfig.Exchange(ctx, p.Code)
// 	if err != nil {
// 		return pages.Error(ctx, NewAppError(AuthExchangeError, ErrGenericAppError, http.StatusBadRequest))
// 	}

// 	pages.SetCookie(
// 		ctx,

// 		&http.Cookie{
// 			Name:     "auth_token",
// 			Value:    token.AccessToken,
// 			Path:     "/",
// 			Secure:   true,
// 			HttpOnly: true,
// 			SameSite: http.SameSiteLaxMode,
// 		},
// 	)
// 	return pages.Redirect(ctx, "/")
// }
