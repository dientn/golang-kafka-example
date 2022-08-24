# 🚀 Getting started with Go example

First of, clone this repo to your local machine, then use following command to install required dependencies:

```
go mod download
```

### `build`

```
go build app/*.go

MAC M1 

go build -tags dynamic app/*.go
```

### `Development`

#### `producer`
```
go run app/producer.go

MAC M1 

go run -tags dynamic app/consumer.go
```

#### `consumer`
```
go run app/consumer.go

MAC M1 

go run -tags dynamic app/consumer.go
```


### `Docker`

```
docker-compose -p go-example up -d --build
```
