// =============================================================================
//
// Consume messages from Confluent Cloud
// Using Confluent Golang Client for Apache Kafka
//
// =============================================================================

package main

/**
 * Copyright 2020 Confluent Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/go-kafka/app/models"
	"example.com/go-kafka/app/libs"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	// "github.com/brianvoe/gofakeit/v6"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// RecordValue represents the struct of the value in a Kafka message
type User models.User

func main() {

	path, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
	}

	// Initialization
	conf := utils.ReadKafkaConfig(path + "/config/kaka.config")
	topic := flag.String("t", "go-example", "Topic name")

	// Create Consumer instance
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": conf["bootstrap.servers"],
		"sasl.mechanisms":   conf["sasl.mechanisms"],
		"security.protocol": conf["security.protocol"],
		"sasl.username":     conf["sasl.username"],
		"sasl.password":     conf["sasl.password"],
		"group.id":          "go_example_group_1",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		fmt.Printf("Failed to create consumer: %s", err)
		os.Exit(1)
	}

	// Create db connection

	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:root@mongo:27017/go-example?ssl=false&authSource=admin"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// databases, err := client.ListDatabaseNames(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(databases)

	coll := client.Database("go-example").Collection("users")

	// doc := &RecordValue{
	// 	Name: gofakeit.Name(),
	// 	Email: gofakeit.Email(),
	// 	Phone: gofakeit.Phone(),
	// 	Company: gofakeit.Company(),
	// 	JobTitle: gofakeit.JobTitle(),
	// }
	// result, err := coll.InsertOne(context.TODO(), doc)

	// fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)

	// Subscribe to topic
	err = c.SubscribeTopics([]string{*topic}, nil)
	// Set up a channel for handling Ctrl-C, etc
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Process messages
	// totalCount := 0
	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			msg, err := c.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Errors are informational and automatically handled by the consumer
				continue
			}
			recordKey := string(msg.Key)
			recordValue := msg.Value
			data := User{}
			err = json.Unmarshal(recordValue, &data)
			if err != nil {
				fmt.Printf("Failed to decode JSON at offset %d: %v", msg.TopicPartition.Offset, err)
				continue
			}
			
			fmt.Printf("Consumed record with key %s and value %s, and updated total count to %d\n", recordKey, recordValue)

			result, err := coll.InsertOne(context.TODO(), data)

			fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
		}
	}

	fmt.Printf("Closing consumer\n")
	c.Close()

}
