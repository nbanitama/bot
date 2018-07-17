package postgresql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Connection struct {
	db *sql.DB
}

func NewConnection(host string, dbname string, user string, password string) (*Connection, error) {
	postgresCon := "postgres://" + user + ":" + password + "@" + host + "/" + dbname + "?sslmode=disable"

	db, err := sql.Open("postgres", postgresCon)
	if err != nil {
		return nil, err
	}
	return &Connection{
		db: db,
	}, nil
}

func CloseConnection(db *sql.DB) error {
	return db.Close()
}

func (c *Connection) ExecuteQuery(query string) (*sql.Rows, error) {
	stmt, err := c.db.Prepare(query)

	if err != nil {
		log.Fatal(err)
		return nil, fmt.Errorf("[PopulateData] Error occurred when preparing statement")
	}

	return stmt.Query()
}

func (c *Connection) ExecuteQueryInt(query string, result *int) error {
	return c.db.QueryRow(query).Scan(result)
}

func (c *Connection) Execute(query string, data ...string) (sql.Result, error) {
	stmt, err := c.db.Prepare(query)
	if err != nil {
		log.Println("Preparing")
		log.Println(err)
		return nil, err
	}
	return stmt.Exec(data)
}

func (c *Connection) GetConnection() *sql.DB {
	return c.db
}
