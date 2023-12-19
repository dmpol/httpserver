package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

// For simplification, we're storing the users information as an in-memory map in our code
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Signin(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func Welcome(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the JWT string from the cookie
	tknStr := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the token
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code until this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// (END) The code until this point is the same as the first part of the `Welcome` route

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	// За 30 секунд до окончания токена можно обновиться и получить новый токен
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// immediately clear the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Expires: time.Now(),
	})
}

func main() {
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/logout", Logout)

	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}

//type Product struct {
//	Id          int
//	Name        string
//	Slug        string
//	Description string
//}
//
//var mySigningKey = []byte("my-super-secret-559f")
//
//var products = []Product{
//	{Id: 1, Name: "Hover Shooters", Slug: "hover-shooters",
//		Description: "Shoot your way to the top on 14 different hoverboards"},
//	{Id: 2, Name: "Ocean Explorer", Slug: "ocean-explorer",
//		Description: "Explore the depths of the sea in this one of a kind"},
//	{Id: 3, Name: "Dinosaur Park", Slug: "dinosaur-park",
//		Description: "Go back 65 million years in the past and ride a T-Rex"},
//	{Id: 4, Name: "Cars VR", Slug: "cars-vr",
//		Description: "Get behind the wheel of the fastest cars in the world."},
//	{Id: 5, Name: "Robin Hood", Slug: "robin-hood",
//		Description: "Pick up the bow and arrow and master the art of archery"},
//	{Id: 6, Name: "Real World VR", Slug: "real-world-vr",
//		Description: "Explore the seven wonders of the world in VR"},
//}
//
//var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("API is up and running"))
//})
//
//var ProductsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	payload, _ := json.Marshal(products)
//	w.Header().Set("Content-Type", "application/json")
//	w.Write([]byte(payload))
//})
//
//var AddFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	var product Product
//	vars := mux.Vars(r)
//	slug := vars["slug"]
//
//	for _, p := range products {
//		if p.Slug == slug {
//			product = p
//		}
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	if product.Slug != "" {
//		payload, _ := json.Marshal(product)
//		w.Write(payload)
//	} else {
//		w.Write([]byte("Product Not Found"))
//	}
//})
//
//var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	w.Write([]byte("NotImplemented"))
//})
//
//var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//	token := jwt.New(jwt.SigningMethodHS256)
//
//	token.Claims["admin"] = true
//	token.Claims["name"] = "Ado Kukic"
//	token.Claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
//
//	tokenString, _ := token.SignedString(mySigningKey)
//	w.Write([]byte(tokenString))
//})
//
//var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
//	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
//		return mySigningKey, nil
//	},
//	SigningMethod: jwt.SigningMethodHS256,
//})
//
//func main() {
//	router := mux.NewRouter()
//
//	router.Handle("/", http.FileServer(http.Dir("./views/")))
//	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
//		http.FileServer(http.Dir("./static/"))))
//	router.Handle("/get-token", GetTokenHandler).Methods("GET")
//	router.Handle("/status", StatusHandler).Methods("GET")
//	router.Handle("/products",
//		jwtMiddleware.Handler(ProductsHandler)).Methods("GET")
//	router.Handle("/products/{slug}/feedback",
//		jwtMiddleware.Handler(AddFeedbackHandler)).Methods("POST")
//
//	http.ListenAndServe(":9000", handlers.LoggingHandler(os.Stdout, router))
//}
