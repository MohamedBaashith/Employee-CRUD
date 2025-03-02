package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Employee struct {
	Name     string
	Position string
	Contact  int64
}

const Connection = "mongodb://localhost:27017"
const Database = "Employee"
const Collection = "Details"

var group *mongo.Collection

func init() {
	clientOption := options.Client().ApplyURI(Connection)

	connect, err := mongo.Connect(context.TODO(), clientOption)
	Error(err)

	group = connect.Database(Database).Collection(Collection)

	fmt.Println("The Collection is Ready!!!")
}

func main() {
	router := gin.Default()

	router.SetTrustedProxies(nil)

	router.GET("/FetchTheDetails", ReadAllDetails)
	router.POST("/AddTheDetails", CreateEmployeeDetails)
	router.PUT("/UpdateTheDetails/:id", UpdateEmployeeDetails)
	router.DELETE("/DeleteTheParticularDetails/:id", DeleteOneDetail)
	router.DELETE("/DeleteAllTheDetails", DeleteAllDetails)

	err := router.RunTLS(
		":6000",
		"D:\\Projects\\Employee\\MyProject\\Certificate\\crudservice.crt", 
		"D:\\Projects\\Employee\\MyProject\\Certificate\\crudservice.key",
	)
	
	if err != nil {
		log.Fatal("Failed to run server with HTTPS: ", err)
	}
}

func insertRecord(emp Employee) {
	inserted, err := group.InsertOne(context.TODO(), emp)
	Error(err)
	fmt.Println("The Employee detail is inserted and it's Id is", inserted.InsertedID)
}

func updateRecord(empId string, updateFields map[string]interface{}) {
	id, err := primitive.ObjectIDFromHex(empId)
	Error(err)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updateFields}
	result, err := group.UpdateOne(context.TODO(), filter, update)
	Error(err)
	fmt.Println("The Employee detail is Updated and it's modified Id is", result.ModifiedCount)
}

func deleteOneRecord(empId string) {
	id, err := primitive.ObjectIDFromHex(empId)
	Error(err)
	filter := bson.M{"_id": id}

	deletedResult, err := group.DeleteOne(context.TODO(), filter)
	Error(err)
	fmt.Println("The Employee detail is Deleted with the delete count:", deletedResult.DeletedCount)
}

func deleteAllRecord() {
	filter := bson.D{{}}
	deletedResult, err := group.DeleteMany(context.TODO(), filter)
	Error(err)
	fmt.Println("All the Employee details are deleted with delete count:", deletedResult.DeletedCount)
}

func readAllRecord() []bson.M {
	filter := bson.D{{}}
	cursor, err := group.Find(context.TODO(), filter)
	Error(err)
	var details []bson.M
	for cursor.Next(context.TODO()) {
		var detail bson.M
		err := cursor.Decode(&detail)
		Error(err)
		details = append(details, detail)
	}
	defer cursor.Close(context.TODO())
	return details
}

func ReadAllDetails(c *gin.Context) {
	allmovies := readAllRecord()
	c.JSON(http.StatusOK, allmovies)
}

func CreateEmployeeDetails(c *gin.Context) {
	var detail Employee
	if err := c.ShouldBindJSON(&detail); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	insertRecord(detail)
	c.JSON(http.StatusOK, detail)
}

func UpdateEmployeeDetails(c *gin.Context) {
	empId := c.Param("id")

	var updateFields map[string]interface{}

	if err := c.ShouldBindJSON(&updateFields); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updateRecord(empId, updateFields)

	c.JSON(http.StatusOK, gin.H{
		"message": "Details Updated Successfully!!!",
		"id":      empId,
	})
}

func DeleteOneDetail(c *gin.Context) {
	empId := c.Param("id")
	deleteOneRecord(empId)
	c.JSON(http.StatusOK, gin.H{"message": "Employee Detail Deleted!!!", "id": empId})
}

func DeleteAllDetails(c *gin.Context) {
	deleteAllRecord()
	c.JSON(http.StatusOK, gin.H{"message": "All Employee Details are Deleted!!!"})
}

func Error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
