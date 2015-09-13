package main
import (
	"database/sql"
	"log"
	"net/http"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db := NewDB()
	log.Println("Listeneing on :3000")
	http.ListenAndServe(":3000", ShowBooks(db))
}

func NewDB() *sql.DB {
	db, err := sql.Open("sqlite3", "example.sqlite")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("create table if not exists books(title text, author text)")
	if err != nil {
		panic(err)
	}

	return db
}

func ShowBooks(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		var title, author string
		err := db.QueryRow("select title, author from books").Scan(&title, &author)
		if err != nil {
			InsertBooks(db)
			panic(err)
		}

		fmt.Fprintf(rw, "The first book is '%s' by '%s'", title, author)
	})
}

func InsertBooks(db *sql.DB) {
	for i := 1; i < 3; i ++ {
		sql := fmt.Sprintf("insert into books values ('book%d', 'author%d')", i, i)
		fmt.Print(sql)
		_, err := db.Exec(sql)
		if err != nil {
			panic(err)
		}
	}
}