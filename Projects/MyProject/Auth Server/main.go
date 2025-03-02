package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/argon2"
)

type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterForm struct {
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	CreatedAt time.Time
}

const (
	Connection = "mongodb://localhost:27017"
	Database   = "Employee"
	Collection = "Data"
	TokenExpiry = time.Hour * 24
)

var credentials *mongo.Collection
var JWT_SECRET_KEY string
func init() {
	clientOption := options.Client().ApplyURI(Connection)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	credentials = client.Database(Database).Collection(Collection)
	fmt.Println("The Collection is Ready!!!")

	JWT_SECRET_KEY = os.Getenv("JWT_SECRET_KEY")
}

func main() {
	router := gin.Default()

	router.POST("/Login", Login)
	router.POST("/Register", Register)

	err := router.RunTLS(
		":5050",
		"D:\\Projects\\Employee\\MyProject\\Certificate\\authservice.crt", 
		"D:\\Projects\\Employee\\MyProject\\Certificate\\authservice.key",
	)
	
	if err != nil {
		log.Fatal("Failed to run server with HTTPS: ", err)
	}
}

func GenerateHash(password, salt string) string {
	timeCost := uint32(1)
	memCost := uint32(64 * 1024)
	parallelism := uint8(1)

	hash := argon2.IDKey([]byte(password), []byte(salt), timeCost, memCost, parallelism, 32)
	return hex.EncodeToString(hash)
}

func GenerateJWT(email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(TokenExpiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWT_SECRET_KEY))
}

func Register(c *gin.Context) {
	var newUser RegisterForm
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingUser RegisterForm
	err := credentials.FindOne(context.TODO(), bson.M{"Email": newUser.Email}).Decode(&existingUser)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email is already registered"})
		return
	}

	salt := newUser.Email
	hashedPassword := GenerateHash(newUser.Password, salt)

	userRecord := bson.M{
		"FirstName": newUser.FirstName,
		"LastName":  newUser.LastName,
		"Email":     newUser.Email,
		"Password":  hashedPassword,
		"CreatedAt": time.Now(),
	}

	_, err = credentials.InsertOne(context.TODO(), userRecord)
	if err != nil {
		log.Println("Error inserting user:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Employee registered successfully"})
}

func Login(c *gin.Context) {
	var loginUser Credentials
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var storedUser RegisterForm
	err := credentials.FindOne(context.TODO(), bson.M{"Email": loginUser.Username}).Decode(&storedUser)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	salt := storedUser.Email
	hashedPassword := GenerateHash(loginUser.Password, salt)

	if hashedPassword == storedUser.Password {
		token, err := GenerateJWT(storedUser.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Login successful", "token":   token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
	}
}
