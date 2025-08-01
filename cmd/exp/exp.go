package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type pgConfig struct {
	host, port, user, password, dbname, sslmode string
}

func (p pgConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.host, p.port, p.user, p.password, p.dbname, p.sslmode,
	)
}

var config = pgConfig{
	"localhost",
	"5432",
	"baloo",
	"junglebook",
	"lenslocked",
	"disable",
}

func main() {
	db, err := sql.Open(
		"pgx",
		config.String(),
	)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("db connection opened")
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name TEXT,
		email TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		amount INT,
		description TEXT
		);`)
	if err != nil {
		panic(err)
	}
	fmt.Println("Tasbles initialised")
	// name := "Soh Wen Ming"
	// email := "wenming.soh@gmail.com"

	// row := db.QueryRow(`
	// 	INSERT INTO users (name, email)
	// 	VALUES ($1, $2) RETURNING id;`,
	// 	name, email)
	// var id int
	// err = row.Scan(&id)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("User created. id=", id)
	// user_id := 1

	// for i := 0; i < 5; i++ {
	// 	amount := i * 100
	// 	description := "test description"
	// 	result := db.QueryRow(`
	// 	INSERT INTO orders (user_id, amount, description)
	// 	VALUES($1, $2, $3) RETURNING id;`,
	// 		user_id, amount, description,
	// 	)
	// 	var returnedId int
	// 	err = result.Scan(&returnedId)
	// 	if err != nil {
	// 		if err == sql.ErrNoRows {
	// 			fmt.Println("no id was returned!")
	// 		} else {
	// 			panic(err)
	// 		}
	// 	}
	// 	fmt.Println("returned Id: ", returnedId)
	// }
	userId := 1
	type Order struct {
		ID          int
		UserID      int
		Amount      int
		Description string
	}

	// the definition of the struct that we want to fill out
	orders := []Order{}
	// presetting the slice or orders, that will hold the results

	rows, err := db.Query(
		`SELECT id, amount, description from orders 
		WHERE user_id = $1`, userId,
	)
	for rows.Next() {
		/*
			rows .Next() will return true if there is still a next row
			to read, and will return false if there is no more row to read

			For this reason, next has to be called for every Scan() function
			that is called on rows - even before the first
		*/
		order := Order{
			UserID: userId,
		}
		// pre population the UserId, because we're hard coding it here and not
		// getting it from the table

		err = rows.Scan(&order.ID, &order.Amount, &order.Description)
		if err != nil {
			panic(err)
		}
		// for each row, first check if there is an error
		// if no error, then each column will end up being populated to the
		// corresponding pointer to the field within the struct
		orders = append(orders, order)
	}
	err = rows.Err()
	// find out if there was any error during interation
	if err != nil {
		panic(err)
	}
	for _, order := range orders {
		fmt.Println("order result: ", order)
	}

}
