package main

import (
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"os"
)

type Response struct {
	Message string `json:"message"`
}

func Handler(request Request) (events.APIGatewayProxyResponse, error) {
	//fmt.Println("Received body: ", request.Body)

	dbConnStr := fmt.Sprintf("%s:%s@tcp(%s:3306)/cities", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_HOST"))
	db, err := sql.Open("mysql", dbConnStr)
	if err != nil {
		fmt.Println(err)
	}
	var city City
	err2 := db.QueryRow("SELECT `id`, `name`, `state_id`, `population` FROM `cities` WHERE state_id = ? AND `lowercase_name` = ?", request.Body.StateId, request.Body.Name).Scan(&city.ID, &city.Name, &city.StateId, &city.Population)
	switch err2 {
		case sql.ErrNoRows:
			return events.APIGatewayProxyResponse{Body: "City not found", StatusCode: 404}, nil
		case nil:
			// fmt.Println(city)
		default:
			panic(err2)
	}
	defer db.Close()

	cityJson, err := json.Marshal(city)
	return events.APIGatewayProxyResponse{Body: string(cityJson), StatusCode: 200}, nil

}

func main() {
	lambda.Start(Handler)
}

type City struct {
	ID 			int `json:"id"`
	Name 		string 	`json:"name"`
	StateId		int `json:"state_id"`
	Population	int `json:"population"`
}

type Request struct {
	Body `json:"body"`
}

type Body struct {
	Name 	string `json:"name"`
	StateId	int `json:"state_id"`
}