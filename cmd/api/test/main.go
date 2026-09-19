package main

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CheckConnection() {
	url := "postgres://postgres:postgres@localhost:5432/myapp?sslmode=disable"
	ctx := context.Background()
	connect, err := pgx.Connect(ctx, url)
	if err != nil {
		panic(err)
	}
	err = connect.Ping(ctx)
	if err != nil {
		panic(err)
	}
	err = connect.Close(ctx)
	if err != nil {
		panic(err)
	}
}

func main() {
	CheckConnection()
}
