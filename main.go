package main

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"log"

	"os"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type Params struct {
	Server   string `json:"server"`
	Port     string `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
	Query    string `json:"query"`
}

func readJSON(filename string) Params {
	file, _ := ioutil.ReadFile(filename)
	params := Params{}

	if err := json.Unmarshal(file, &params); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return params
}

func postgres(arg1 string, arg2 string) string {

	request := readJSON(arg1)

	result := ""

	connectionString := "postgres://" + request.User + ":" + request.Password + "@" + request.Server + ":" + request.Port + "/" + request.Database + "?sslmode=disable"

	db, err := sqlx.Connect("pgx", connectionString)
	if err != nil {
		fmt.Println(err)
	}

	a := []map[string]interface{}{}

	rows, err := db.Queryx(request.Query)
	for rows.Next() {
		results := make(map[string]interface{})
		err = rows.MapScan(results)
		if err != nil {
			log.Fatalln(err)
		}
		a = append(a, results)
	}

	b, _ := json.Marshal(a)

	result = string(b)

	return result
}

func main() {
	argsWithoutProg := os.Args[1:]
	fmt.Println(len(argsWithoutProg))

	arg1 := os.Args[1]
	arg2 := os.Args[2]

	str := postgres(arg1, arg2)

	f, err := os.Create(arg2)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err2 := f.WriteString(str)
	if err2 != nil {
		log.Fatal(err2)
	}

}
