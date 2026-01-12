package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:15329554862ph@tcp(127.0.0.1:3306)/sports_management?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, image_url, title FROM carousel_images")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Carousel Images in DB:")
	count := 0
	for rows.Next() {
		var id int
		var url, title string
		var titleNull sql.NullString
		if err := rows.Scan(&id, &url, &titleNull); err != nil {
			log.Fatal(err)
		}
		if titleNull.Valid {
			title = titleNull.String
		}
		fmt.Printf("ID: %d, URL: %s, Title: %s\n", id, url, title)
		count++
	}
	fmt.Printf("Total count: %d\n", count)
}
