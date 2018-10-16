package main

import (
	"github.com/kabukky/httpscerts"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mbach04/meeseeks/handlers"
)

// jwtCustomClaims are custom claims extending default ones.
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
	jwt.StandardClaims
}

//TODO: Make this piece pluggable for future flexibility
func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "admin" && password == "secret" {

		// Set custom claims
		claims := &jwtCustomClaims{
			"admin",
			true,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSONPretty(http.StatusOK, echo.Map{
			"token": t,
		}, "  ")//note the double spaced empty string as 3rd param here for formatting the return
	}

	return echo.ErrUnauthorized
}

//example
func accessible(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"hello": "world",
	})
}

//example
func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}

func main() {
	log.SetFlags(log.LstdFlags)
	// router := mux.NewRouter().StrictSlash(true)


	v1, err := readConfig("meeseeks", map[string]interface{}{
		"api_port": 8080,
		"debug":    true,
		"tls_cert": "cert.pem",
		"tls_key":  "key.pem",
	})
	if err != nil {
		log.Println("Error reading config file:", err)
	}

	// Store CONFIG info
	debug := v1.GetBool("debug")
	apiPort := v1.GetString("API_PORT")
	cert := v1.GetString("tls_cert")
	key := v1.GetString("tls_key")

	// Dump CONFIG info to Log
	log.Println("DEBUG:", debug)
	log.Println("API_PORT:", apiPort)
	log.Println("CERT", cert)
	log.Println("KEY", key)

	// Check if the cert files are available.
	err = httpscerts.Check(cert, key)
	// If they are not available, generate new ones.
	if err != nil {
		err = httpscerts.Generate(cert, key, "localhost:"+apiPort)
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Login route
	e.POST("/login", login)

	// Unauthenticated route for testing
	e.GET("/test", accessible)

	// Restricted group (all sub routes will require jwt auth)
	r := e.Group("/api/v1")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	r.Use(middleware.JWTWithConfig(config))
	r.GET("", restricted)

	r.POST("/ls", handlers.LsCmd)

	e.Logger.Fatal(e.StartTLS(":"+apiPort, cert, key))
}

func readConfig(filename string, defaults map[string]interface{}) (*viper.Viper, error) {
	v := viper.New()
	for key, value := range defaults {
		v.SetDefault(key, value)
	}
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err
}
