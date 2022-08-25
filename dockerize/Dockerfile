# Create another stage called "dev" that is based off of our "base" stage (so we have golang available to us)
FROM amd64/golang:1.16-alpine
# FROM amd64/node:16-alpine

# Install the air binary so we get live code-reloading when we save files
RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin

# Run the air command in the directory where our code will live
RUN apk add build-base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

# RUN go get github.com/brianvoe/gofakeit/v6
# RUN go get go.mongodb.org/mongo-driver
# RUN go get github.com/confluentinc/confluent-kafka-go/kafka

COPY * ./

# COPY ./dockerize/consumer-entrypoint.sh ./consumer-entrypoint.sh
# COPY ./dockerize/producer-entrypoint.sh ./producer-entrypoint.sh

RUN ls -la .
# RUN go build -tags dynamic  app/*.go