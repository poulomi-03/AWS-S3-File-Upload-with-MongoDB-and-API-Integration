package main

import (
	"context"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	s3Session *s3.S3
	bucket    string
	mongoURI  string
	dbName    string
	collName  string
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get environment variables
	bucket = os.Getenv("AWS_BUCKET")
	mongoURI = os.Getenv("MONGODB_CONN_URI")
	dbName = os.Getenv("MONGODB_DB_NAME")
	collName = os.Getenv("COLLECTION_NAME")

	// Log errors if any of the critical environment variables are missing
	if bucket == "" {
		log.Fatal("AWS_BUCKET is not set in the environment variables")
	}
	if mongoURI == "" {
		log.Fatal("MONGODB_CONN_URI is not set in the environment variables")
	}
	if dbName == "" {
		log.Fatal("MONGODB_DB_NAME is not set in the environment variables")
	}
	if collName == "" {
		log.Fatal("COLLECTION_NAME is not set in the environment variables")
	}

	// Initialize AWS S3 session
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY"),
			os.Getenv("AWS_SECRET_KEY"),
			"",
		),
	})
	if err != nil {
		log.Fatalf("Failed to initialize AWS session: %v", err)
	}
	s3Session = s3.New(awsSession)

	// Log successful AWS and MongoDB connections
	log.Println("Connected to AWS S3 and MongoDB successfully")
}

// uploadToS3 uploads a file to AWS S3 and returns the file's URL
func uploadToS3(file multipart.File, fileName string) (string, error) {
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}
	file.Seek(0, 0) // Reset file pointer to the beginning

	_, err = s3Session.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fileName),
		Body:        file,
		ContentType: aws.String(http.DetectContentType(buffer)),
		ACL:         aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}

	fileURL := "https://" + bucket + ".s3.amazonaws.com/" + fileName
	return fileURL, nil
}

// connectMongo connects to MongoDB and returns a collection handle
func connectMongo() (*mongo.Collection, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}
	return client.Database(dbName).Collection(collName), nil
}

// postSubmit handles POST requests to save form data
func postSubmit(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	file, header, err := c.Request.FormFile("picture")
	if err != nil {
		log.Printf("Error while uploading file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file upload"})
		return
	}

	// Generate a unique file name
	fileName := time.Now().Format("20060102150405") + "-" + header.Filename
	fileURL, err := uploadToS3(file, fileName)
	if err != nil {
		log.Printf("Error uploading file to S3: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to S3"})
		return
	}

	collection, err := connectMongo()
	if err != nil {
		log.Printf("Error connecting to MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}

	// Create the document to insert into MongoDB
	document := bson.M{
		"name":       name,
		"email":      email,
		"picture":    fileURL,
		"created_at": time.Now(),
	}
	_, err = collection.InsertOne(context.TODO(), document)
	if err != nil {
		log.Printf("Error saving data to MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save data to MongoDB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Form submitted successfully"})
}

// fetchPosts handles GET requests to fetch all posts from MongoDB
func fetchPosts(c *gin.Context) {
	collection, err := connectMongo()
	if err != nil {
		log.Printf("Error connecting to MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB"})
		return
	}

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Printf("Error fetching data from MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from MongoDB"})
		return
	}
	defer cursor.Close(context.TODO())

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Printf("Error parsing data from MongoDB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse data from MongoDB"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func main() {
	r := gin.Default()

	// Enable CORS for specific origins
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Define routes
	r.POST("/admin/post-submit", postSubmit)
	r.GET("/admin/posts", fetchPosts)

	// Start the server
	log.Println("Server is running on http://localhost:8080")
	r.Run(":8080")
}
