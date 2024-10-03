package service

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("Hoosen@Azis")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:@tcp(localhost:3306)/husenDB")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test connection
	if err = db.Ping(); err != nil {
		log.Fatal("Database connection is not alive:", err)
	}
}

// Login Handler
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var storedPassword string
	err = db.QueryRow("SELECT password FROM USERS WHERE username = ?", creds.Username).Scan(&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Buat JWT token
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error creating token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   tokenString,
	})
}

// AuthMiddleware untuk memvalidasi JWT dari cookie
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mengambil token dari cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// Jika cookie tidak ditemukan, kembalikan status Unauthorized
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// Error lain dalam pembacaan cookie
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Mendapatkan token string dari cookie
		tokenStr := cookie.Value
		claims := &Claims{}

		// Parsing dan validasi token JWT
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				// Jika signature JWT invalid, kembalikan status Unauthorized
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			// Error lain dalam parsing JWT
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Jika token tidak valid, kembalikan status Unauthorized
		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Token valid, melanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Insert post into database
	stmt, err := db.Prepare("INSERT INTO USERS (USERNAME, PASSWORD) VALUES (?, ?)")
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(creds.Username, hashedPassword)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	})
}	
