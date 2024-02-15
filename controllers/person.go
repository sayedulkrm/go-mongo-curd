package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sayedulkrm/go-mongo-curd/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
	Id          string  `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName   string  `json:"firstname" bson:"firstname"`
	LastName    string  `json:"lastname" bson:"lastname"`
	Age         string  `json:"age" bson:"age"`
	PhoneNumber string  `json:"phonenumber" bson:"phonenumber"`
	Email       string  `json:"email" bson:"email"`
	Address     Address `json:"address" bson:"address"`
}

type Address struct {
	AddressLine1 string `json:"addressline1" bson:"addressline1"`
	AddressLine2 string `json:"addressline2" bson:"addressline2"`
	City         string `json:"city" bson:"city"`
	State        string `json:"state" bson:"state"`
	Country      string `json:"country" bson:"country"`
	ZipCode      string `json:"zipcode" bson:"zipcode"`
}

type UpdatePersonRequest struct {
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Age         string `json:"age"`
	PhoneNumber string `json:"phonenumber"`
	Email       string `json:"email"`
	Address     Address
}

//  Created Person

func CreatePerson(ct *gin.Context) {

	var person Person
	if err := ct.BindJSON(&person); err != nil {
		ct.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unsupported Format"})
		return
	}

	fmt.Println(person)

	SavePersonToDB(ct, person)

}

// Save Person to database

func SavePersonToDB(ct *gin.Context, PersonRecord Person) {

	db := config.New()

	collection := db.Client.Database("gocurd").Collection("AllPersonForDemo")

	response, insertError := collection.InsertOne(ct, PersonRecord)

	println(response)

	if insertError != nil {
		ct.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert person"})
		return

	}

	// Set the _id field of the person object with the generated ObjectID
	PersonRecord.Id = response.InsertedID.(primitive.ObjectID).Hex()

	// Respond with the saved person object including the generated _id
	ct.IndentedJSON(http.StatusCreated, gin.H{"success": true, "message": "Person is created successfully", "person": PersonRecord})
}

// Get person

func GetPerson(ct *gin.Context) {

	id := ct.Param("id")

	fmt.Println(id)

	db := config.New()
	collection := db.Client.Database("gocurd").Collection("AllPersonForDemo")

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// Handle error (e.g., invalid ID format)
		ct.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	filter := bson.D{primitive.E{Key: "_id", Value: idPrimitive}}

	// Find the person document
	var person Person

	dbError := collection.FindOne(ct, filter).Decode(&person)

	fmt.Println("Heyyyyyy", dbError)

	if dbError != nil {
		ct.IndentedJSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	ct.IndentedJSON(http.StatusOK, gin.H{"success": true, "message": "Person found successfully", "person": person})

}

// Get All the persons

func GetAllPerson(ct *gin.Context) {
	// Use the existing database connection
	db := config.New()
	collection := db.Client.Database("gocurd").Collection("AllPersonForDemo")

	// Define an empty slice to store the results
	var persons []Person

	// Find all documents in the collection
	cursor, err := collection.Find(ct, bson.M{})
	if err != nil {
		ct.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve persons"})
		return
	}
	defer cursor.Close(ct)

	// Iterate over the cursor and decode each document into the slice
	for cursor.Next(ct) {
		var person Person
		if err := cursor.Decode(&person); err != nil {
			ct.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode person"})
			return
		}
		persons = append(persons, person)
	}

	// Check if any error occurred during iteration
	if err := cursor.Err(); err != nil {
		ct.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve persons"})
		return
	}

	// Respond with the retrieved persons
	ct.IndentedJSON(http.StatusOK, gin.H{"success": true, "message": "Persons retrieved successfully", "persons": persons})
}

//  Delete User

func DeletePerson(ct *gin.Context) {
	id := ct.Param("id")

	fmt.Println(id)

	db := config.New()
	collection := db.Client.Database("gocurd").Collection("AllPersonForDemo")

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// Handle error (e.g., invalid ID format)
		ct.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	filter := bson.D{primitive.E{Key: "_id", Value: idPrimitive}}

	// Perform the delete operation
	result, err := collection.DeleteOne(ct, filter)
	if err != nil {
		ct.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete person"})
		return
	}

	if result.DeletedCount == 0 {
		// No document was deleted, return a "Person not found" error
		ct.IndentedJSON(http.StatusNotFound, gin.H{"error": "Person not found. Unable to delete."})
		return
	}

	ct.IndentedJSON(http.StatusAccepted, gin.H{"success": true, "message": "Person deleted successfully"})
}

// Update User Profile

func UpdatePerson(ct *gin.Context) {
	id := ct.Param("id")

	fmt.Println(id)

	db := config.New()
	collection := db.Client.Database("gocurd").Collection("AllPersonForDemo")

	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		// Handle error (e.g., invalid ID format)
		ct.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	filter := bson.D{primitive.E{Key: "_id", Value: idPrimitive}}

	fmt.Println("Heyyy An filter", filter)

	var person Person
	if err := ct.BindJSON(&person); err != nil {
		ct.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Unsupported Format"})
		return
	}

	// Define update operation using $set operator
	update := bson.M{"$set": bson.M{
		"firstname":   person.FirstName,
		"lastname":    person.LastName,
		"age":         person.Age,
		"phonenumber": person.PhoneNumber,
		"email":       person.Email,
		"address": bson.M{
			"addressline1": person.Address.AddressLine1,
			"addressline2": person.Address.AddressLine2,
			"city":         person.Address.City,
			"state":        person.Address.State,
			"country":      person.Address.Country,
			"zipcode":      person.Address.ZipCode,
		},
		// Add more fields to update as needed
	}}

	fmt.Println("Heyyy An Updated", update)

	_, dbError := collection.UpdateOne(ct, filter, update)

	if dbError != nil {
		ct.IndentedJSON(http.StatusFailedDependency, gin.H{"message": "Unable to update the person"})
		return
	}

	ct.IndentedJSON(http.StatusAccepted, gin.H{"message": "Person updated successfully"})
}
