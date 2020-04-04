package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"nebulosa-studio/quicktest/status"
	"time"

	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/yaml.v2"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	targets, err := read("test.yml")
	if err != nil {
		fmt.Println("can not file yaml file")
		return
	}

	for _, i := range targets {
		switch i.Type {
		case "redis":
			redisTest(i)
		case "mysql":
			mysqlTest(i)
		case "mongodb":
			mongodbTest(i)
		default:
			fmt.Println("unknown type:", i.Type, i)
		}
	}
}

// Target of testing
type Target struct {
	Name     string
	Type     string
	Host     string
	Port     string
	Username string
	Password string
}

func read(file string) ([]Target, error) {
	m := []Target{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return m, err
	}
	err = yaml.Unmarshal(data, &m)
	return m, err
}

func redisTest(t Target) {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", t.Host, t.Port),
	})
	r, err := client.Ping().Result()
	if err != nil {
		status.Print("redis", status.Error, err.Error())
		return
	}

	status.Print("redis", status.Success, r)
}

func mysqlTest(t Target) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/", t.Username, t.Password, t.Host))
	if err != nil {
		status.Print("mysql", status.Error, err.Error())
		return
	}

	err = db.Ping()
	if err != nil {
		status.Print("mysql", status.Error, err.Error())
		return
	}

	status.Print("mysql", status.Success, "ping")
}

func mongodbTest(t Target) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", t.Host, t.Port)))
	if err != nil {
		status.Print("mongodb", status.Error, err.Error())
		return
	}
	defer client.Disconnect(ctx)

	ctx, _ = context.WithTimeout(context.Background(), 2*time.Second)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		status.Print("mongodb", status.Error, err.Error())
		return
	}

	status.Print("mongodb", status.Success, "pong")
}
