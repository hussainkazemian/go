package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("hello world")

	if os.Getenv("ENV") != "production" {
		// Load the .env file if not in production
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file:", err)
		}
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MONGODB ATLAS")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,PATCH,DELETE,OPTIONS",
	}))

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	if os.Getenv("ENV") == "production" {
		app.Static("/", "./client/dist")
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))

}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	// Query params: search, status (all|completed|active), sortBy (createdAt|body|completed), order (asc|desc)
	search := strings.TrimSpace(c.Query("search", ""))
	status := strings.ToLower(c.Query("status", "all"))
	sortBy := strings.ToLower(c.Query("sortBy", "createdAt"))
	order := strings.ToLower(c.Query("order", "desc"))

	filter := bson.M{}
	if status == "completed" {
		filter["completed"] = true
	} else if status == "active" {
		filter["completed"] = false
	}

	if search != "" {
		// Case-insensitive regex search on body
		filter["body"] = bson.M{"$regex": search, "$options": "i"}
	}

	// Sorting
	sortField := "createdAt"
	switch sortBy {
	case "body":
		sortField = "body"
	case "completed":
		sortField = "completed"
	default:
		sortField = "createdAt"
	}
	sortOrder := -1 // desc
	if order == "asc" {
		sortOrder = 1
	}

	findOptions := options.Find().SetSort(bson.D{{Key: sortField, Value: sortOrder}})

	cursor, err := collection.Find(context.Background(), filter, findOptions)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		todos = append(todos, todo)
	}
	if err := cursor.Err(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	type createTodoDTO struct {
		Body string `json:"body"`
	}
	var payload createTodoDTO
	if err := c.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid payload")
	}

	body := strings.TrimSpace(payload.Body)
	if body == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Todo body cannot be empty"})
	}
	if len(body) > 200 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Todo body must be 200 characters or less"})
	}

	// Optional: prevent duplicate bodies (case-insensitive)
	dupFilter := bson.M{"body": bson.M{"$regex": fmt.Sprintf("^%s$", regexp.QuoteMeta(body)), "$options": "i"}}
	count, err := collection.CountDocuments(context.Background(), dupFilter)
	if err == nil && count > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "A todo with the same text already exists"})
	}

	todo := &Todo{
		ID:        primitive.NilObjectID,
		Completed: false,
		Body:      body,
		CreatedAt: time.Now().UTC(),
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)
	return c.Status(fiber.StatusCreated).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	// Accept optional payload { completed: boolean }
	var payload struct {
		Completed *bool `json:"completed"`
	}
	_ = c.BodyParser(&payload)

	set := bson.M{"updatedAt": time.Now().UTC()}
	if payload.Completed == nil {
		// default to mark as completed
		set["completed"] = true
	} else {
		set["completed"] = *payload.Completed
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": set}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"success": true})
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectID}
	_, err = collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"success": true})
}
