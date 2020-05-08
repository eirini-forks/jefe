package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golangcollege/sessions"
	"github.com/herrjulz/jefe/pkg/database/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kelseyhightower/envconfig"
	"golang.org/x/oauth2"
)

const (
	githubAuthorizeUrl = "https://github.com/login/oauth/authorize"
	githubTokenUrl     = "https://github.com/login/oauth/access_token"
)

type specs struct {
	Port           string `default:"5555"`
	Dsn            string `default:"web:pass@"`
	SessionSecret  string `default:"sessionsecret" split_words:"true"`
	GithubClientID string `split_words:"true" required:"true"`
	GithubSecret   string `split_words:"true" required:"true"`
	GithubOAuthOrg string `split_words:"true" required:"true"`
	AdminUser      string `split_words:"true" default:"admin"`
	AdminPassword  string `split_words:"true" required:"true"`
	TlsEnabled     string `split_words:"true" default:"false"`
}

type BasicAuth struct {
	User     string
	Password string
}

type application struct {
	Envs       mysql.Environments
	OAuthConf  *oauth2.Config
	GHOAuthOrg string
	AuthURL    string
	Session    *sessions.Session
	BasicAuth  BasicAuth
}

func main() {
	var s specs
	err := envconfig.Process("jefe", &s)
	if err != nil {
		log.Fatalf("Can't process environment variables: %w", err)
	}

	dsn := fmt.Sprintf("%s/environmentsdb?parseTime=true", s.Dsn)

	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	envs := mysql.Environments{DB: db}
	oauthCfg := &oauth2.Config{
		ClientID:     s.GithubClientID,
		ClientSecret: s.GithubSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  githubAuthorizeUrl,
			TokenURL: githubTokenUrl,
		},
	}

	session := sessions.New([]byte(s.SessionSecret))
	session.Lifetime = 12 * time.Hour
	// Secure should be used for https connections (check with ingress)
	// session.Secure = true

	basicAuth := BasicAuth{
		User:     s.AdminUser,
		Password: s.AdminPassword,
	}

	a := application{
		Envs:       envs,
		OAuthConf:  oauthCfg,
		Session:    session,
		GHOAuthOrg: "eirini-forks",
		AuthURL: fmt.Sprintf(
			"%s?client_id=%s&scope=read:user+read:org",
			githubAuthorizeUrl,
			s.GithubClientID,
		),
		BasicAuth: basicAuth,
	}

	handler := a.routes()

	var port string
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	} else {
		port = s.Port
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), handler)
}

func openDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
