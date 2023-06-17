package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"

    "github.com/streadway/amqp"
)

type Person struct {
    Name  string `json:"name,omitempty"`
    Email string `json:"email,omitempty"`
}

func main() {
    // RabbitMQ setup
    conn, err := amqp.Dial("amqp://admin1:admin123@localhost:5672/")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatal(err)
    }
    defer ch.Close()

    q, err := ch.QueueDeclare(
        "people_queue", // name
        false,          // durable
        false,          // delete when unused
        false,          // exclusive
        false,          // no-wait
        nil,            // arguments
    )
    if err != nil {
        log.Fatal(err)
    }

    // HTTP server setup
    http.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            // Return a message indicating that the route is not implemented
            fmt.Fprintln(w, "GET /people is not implemented")

        case "POST":
            var person Person
            err := json.NewDecoder(r.Body).Decode(&person)
            if err != nil {
                log.Fatal(err)
            }

            // Publish the person to the RabbitMQ queue
            body, err := json.Marshal(person)
            if err != nil {
                log.Fatal(err)
            }

            err = ch.Publish(
                "",     // exchange
                q.Name, // routing key
                false,  // mandatory
                false,  // immediate
                amqp.Publishing{
                    ContentType: "application/json",
                    Body:        body,
                },
            )
            if err != nil {
                log.Fatal(err)
            }

fmt.Fprintln(w, "Person added successfully")

        default:
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        }
    })

    log.Fatal(http.ListenAndServe(":8090", nil))
}
