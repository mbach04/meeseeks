package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kabukky/httpscerts"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mbach04/meeseeks/handlers"
	"github.com/spf13/viper"
)

func main() {
	log.SetFlags(log.LstdFlags)

	//Read config file and establish some basic default values
	meeseeksConfig, err := readConfig("config", map[string]interface{}{
		"HOSTNAME": "localhost",
		"API_PORT": 8080,
		"DEBUG":    true,
		"TLS_CERT": "cert.pem",
		"TLS_KEY":  "key.pem",
	})

	if err != nil {
		log.Println("Error reading config file:", err)
	}

	// Extract CONFIG info
	debug := meeseeksConfig.GetBool("DEBUG")
	hostname := meeseeksConfig.GetString("HOSTNAME")
	apiPort := meeseeksConfig.GetString("API_PORT")
	certPath := meeseeksConfig.GetString("TLS_CERT")
	keyPath := meeseeksConfig.GetString("TLS_KEY")

	// Dump CONFIG info to Log
	log.Println("DEBUG:", debug)
	log.Println("HOSTNAME:", hostname)
	log.Println("API_PORT:", apiPort)
	log.Println("TLS_CERT:", certPath)
	log.Println("TLS_KEY:", keyPath)

	generateHTTPSCert(hostname, apiPort, certPath, keyPath)
	startEcho(hostname, apiPort, certPath, keyPath)

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

func generateHTTPSCert(hostname string, apiPort string, certPath string, keyPath string) {
	// Check if the cert files are available.
	err := httpscerts.Check(certPath, keyPath)
	// If they are not available, generate new ones.
	if err != nil {
		err = httpscerts.Generate(certPath, keyPath, hostname+":"+apiPort)
		if err != nil {
			log.Fatal("Error: Couldn't create https certs.")
		}
	}
}

func startEcho(hostname string, apiPort string, certPath string, keyPath string) {
	meeseeksEcho := echo.New()
	// Middleware
	meeseeksEcho.Use(middleware.Logger())
	meeseeksEcho.Use(middleware.Recover())

	// Login route
	meeseeksEcho.POST("/login", login)

	// Unauthenticated route for testing
	meeseeksEcho.GET("/test", accessible)

	// Restricted group (all sub routes will require jwt auth)
	meeseeksRestricted := meeseeksEcho.Group("/api/v1")

	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	meeseeksRestricted.Use(middleware.JWTWithConfig(config))

	meeseeksRestricted.GET("", restricted)
	meeseeksRestricted.POST("/ls", handlers.LsCmd)

	meeseeksEcho.Logger.Fatal(meeseeksEcho.StartTLS(hostname+":"+apiPort, certPath, keyPath))
}

//jwtCustomClaims are custom claims extending default ones.
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
		}, "  ") //note the double spaced empty string as 3rd param here for formatting the return
	}

	return echo.ErrUnauthorized
}

//example
func accessible(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"hello": "world",
	})
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
