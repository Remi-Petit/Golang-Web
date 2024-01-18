package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"

	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Todo struct {
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

func connectDB() (*sql.DB, error) {
	connStr := "user=admin dbname=test sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("Connexion à la base de données établie")
	return db, nil
}

func fetchSomeData(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM test")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		// Lire les résultats dans des variables
		err := rows.Scan(&id, &name)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("ID: %d, Nom: %s\n", id, name)
	}
}

func closeDB(db *sql.DB) {
	db.Close()
	fmt.Println("Connexion à la base de données fermée")
}

func main() {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := TodoPageData{
			PageTitle: "My TODO list",
			Todos: []Todo{
				{Title: "Task 1", Done: false},
				{Title: "Task 2", Done: true},
				{Title: "Task 3", Done: true},
				{Title: "Task 4", Done: true},
				{Title: "Task 4", Done: true},
			},
		}
		tmpl.Execute(w, data)
	})

	r := mux.NewRouter()
	r.HandleFunc("/books/{title}/page/{page}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		title := vars["title"]
		page := vars["page"]

		fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)
	})

	http.ListenAndServe(":80", nil)
}
