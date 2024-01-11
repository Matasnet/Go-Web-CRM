package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Client struct {
	ID        int
	Name      string
	Surname   string
	Email     string
	Phone     string
	CreatedAt time.Time
}

func addClient(db *sql.DB, client Client) error {
	client.CreatedAt = time.Now().UTC()
  formattedTime := client.CreatedAt.Format("2006-01-02 15:04")

	query := `INSERT INTO clients (name, surname, email, phone, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, client.Name, client.Surname, client.Email, client.Phone, formattedTime)
	return err
}

func getClients(db *sql.DB) ([]Client, error) {
	query := `SELECT id, name, surname, email, phone, created_at FROM clients`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	clients := []Client{}
	for rows.Next() {
		var id int
		var name, surname, email, phone string
		var createdAt time.Time
		err := rows.Scan(&id, &name, &surname, &email, &phone, &createdAt)
		if err != nil {
			return nil, err
		}
		client := Client{
			ID:        id,
			Name:      name,
			Surname:   surname,
			Email:     email,
			Phone:     phone,
			CreatedAt: createdAt,
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func deleteClient(db *sql.DB, clientID int) error {
	query := `DELETE FROM clients WHERE id = ?`
	_, err := db.Exec(query, clientID)
	return err
}

func filterClients(clients []Client, filter string, filterColumn string) []Client {
	if filter == "" {
		return clients
	}

	var filteredClients []Client

	for _, client := range clients {
		switch filterColumn {
		case "ID":
			if strconv.Itoa(client.ID) == filter {
				filteredClients = append(filteredClients, client)
			}
		case "Name":
			if strings.Contains(client.Name, filter) {
				filteredClients = append(filteredClients, client)
			}
		case "Surname":
			if strings.Contains(client.Surname, filter) {
				filteredClients = append(filteredClients, client)
			}
		case "Email":
			if strings.Contains(client.Email, filter) {
				filteredClients = append(filteredClients, client)
			}
		case "Phone":
			if strings.Contains(client.Phone, filter) {
				filteredClients = append(filteredClients, client)
			}
		case "CreatedAt":
			if strings.Contains(client.CreatedAt.String(), filter) {
				filteredClients = append(filteredClients, client)
			}
		}
	}

	return filteredClients
}

func main() {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if db == nil {
		fmt.Println("Nie można nawiązać połączenia z bazą danych.")
		return
	}
	fmt.Println("Połączono z bazą danych.")
  fmt.Println("Twój kod został uruchomiony wejdź na http://localhost:8080")
  fmt.Println("MatasNET")


  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    html := `
  <!DOCTYPE html>
  <html lang="pl">
  <head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Strona Główna</title>
  </head>
  <body>
    <h1>Witaj na Stronie Głównej prostego CRM!</h1>
  
    <p>Możesz dodać klienta bądź przejrzeć listę tych, którzy są już dodani.</p>

    <p>Wybierz jedną z opcji:</p>
    
    <ul>
      <li><a href="/add-form">Dodaj klienta</a></li>
      <li><a href="/list">Lista klientów</a></li>
    </ul>
  </body>
  </html>
  `
    w.Write([]byte(html))
  })
  

	http.HandleFunc("/add-form", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html lang="pl">
<head>
  <title>Dodaj klienta</title>
</head>
<body>
  <h1>Dodaj klienta</h1>

  <form action="/add" method="post">
    <label for="name">Imię:</label>
    <input type="text" id="name" placeholder="John"  name="name" required><br>

    <label for="surname">Nazwisko:</label>
    <input type="text" id="surname" placeholder="Doe"  name="surname" required><br>

    <label for="email">Adres e-mail:</label>
    <input type="email" placeholder="john.doe@example.com" id="email" name="email" required><br>

    <label for="phone">Numer telefonu:</label>
    <input type="text" id="phone" placeholder="123456789" name="phone" required><br>

    <button type="submit">Dodaj klienta</button>
  </form>
</body>
</html>
`
		w.Write([]byte(html))
	})

	http.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Metoda nieobsługiwana", http.StatusMethodNotAllowed)
			return
		}

		name := r.FormValue("name")
		surname := r.FormValue("surname")
		email := r.FormValue("email")
		phone := r.FormValue("phone")

		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    phoneRegex := regexp.MustCompile(`^\d{9,20}$`)
		if !emailRegex.MatchString(email) || !phoneRegex.MatchString(phone) {
			http.Error(w, "Nieprawidłowy adres e-mail lub numer telefonu", http.StatusBadRequest)
			return
		}

		client := Client{
			Name:    name,
			Surname: surname,
			Email:   email,
			Phone:   phone,
		}

		err := addClient(db, client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/list", http.StatusFound)
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		clients, err := getClients(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		filter := r.FormValue("filter")
		filterColumn := r.FormValue("filterColumn")

		filteredClients := filterClients(clients, filter, filterColumn)

		tmpl, err := template.New("list").Parse(`
<!DOCTYPE html>
<html lang="pl">
<head>
  <title>Lista klientów</title>
</head>
<body>
  <h1>Lista klientów</h1>

  <form action="/list" method="get">
    <label for="filterColumn">Filtruj według kolumny:</label>
    <select id="filterColumn" name="filterColumn">
      <option value="ID">ID</option>
      <option value="Name">Imię</option>
      <option value="Surname">Nazwisko</option>
      <option value="Email">Adres e-mail</option>
      <option value="Phone">Numer telefonu</option>
      <option value="CreatedAt">Data utworzenia</option>
    </select>
    <input type="text" id="filter" name="filter" placeholder="Wpisz wartość do filtrowania">
    <button type="submit">Filtruj</button>
  </form>

  <table>
    <thead>
      <tr>
        <th>ID</th>
        <th>Imię</th>
        <th>Nazwisko</th>
        <th>Adres e-mail</th>
        <th>Numer telefonu</th>
        <th>Data utworzenia</th>
        <th>Akcje</th>
      </tr>
    </thead>
    <tbody>
      {{ range . }}
        <tr>
          <td>{{ .ID }}</td>
          <td>{{ .Name }}</td>
          <td>{{ .Surname }}</td>
          <td>{{ .Email }}</td>
          <td>{{ .Phone }}</td>
          <td>{{ .CreatedAt }}</td>
          <td>
            <form action="/delete" method="post">
              <input type="hidden" name="id" value="{{ .ID }}">
              <button type="submit">Usuń</button>
            </form>
          </td>
        </tr>
      {{ end }}
    </tbody>
  </table>

</body>
</html>
`)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, filteredClients)
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Metoda nieobsługiwana", http.StatusMethodNotAllowed)
			return
		}

		clientID := r.FormValue("id")
		if clientID == "" {
			http.Error(w, "Brak identyfikatora klienta", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(clientID)
		if err != nil {
			http.Error(w, "Nieprawidłowy identyfikator klienta", http.StatusBadRequest)
			return
		}

		err = deleteClient(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/list", http.StatusFound)
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}