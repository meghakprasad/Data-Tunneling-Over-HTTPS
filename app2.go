package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Person struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

func main() {
	// MongoDB setup
	uri := "mongodb://mongoAdmin:admin123@192.168.162.202:27017/?authSource=admin"

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("admin").Collection("people")

	// Connect to RabbitMQ
	conn, err := amqp.Dial("amqp://admin1:admin123@192.168.162.53:5672/")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// Create a channel
	channel, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	defer channel.Close()

	// Declare a queue
	queueName := "people_queue"
	queue, err := channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Start consuming messages from the queue
	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Launch a separate goroutine to handle each message as it arrives
	for i := 0; i < 10; i++ {
		go func() {
			for message := range messages {
				var person Person

				err = json.Unmarshal(message.Body, &person)
				if err != nil {
					log.Printf("Error unmarshalling message: %s", err)
					continue
				}

				log.Printf("Received message: %s", message.Body)

				// Insert the person into MongoDB
				_, err = collection.InsertOne(context.Background(), person)
				if err != nil {
					log.Printf("Error inserting person into MongoDB: %s", err)
				}
			}
		}()
	}

	http.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			people := []Person{}

			cursor, err := collection.Find(context.Background(), bson.M{})
			if err != nil {
				log.Fatal(err)
			}

			defer cursor.Close(context.Background())

			for cursor.Next(context.Background()) {
				var person Person
				err := cursor.Decode(&person)
				if err != nil {
					log.Fatal(err)
				}

				people = append(people, person)
			}

			json.NewEncoder(w).Encode(people)

		default:
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":8090", nil))
}
