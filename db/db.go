package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

const dbLocation = "db/development.db"

func exec(ctx context.Context, db *sql.DB, stmt string, args ...any) sql.Result {
	res, err := db.ExecContext(ctx, stmt, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute statement %s: %s", stmt, err)
		os.Exit(1)
	}
	return res
}

func ConnectToDb() {
	var dbUrl = "file:" + dbLocation
	// dbUrl = "http://127.0.0.1:8080"
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbUrl, err)
		os.Exit(1)
	}

	ctx := context.Background()
	exec(ctx, db, "CREATE TABLE IF NOT EXISTS counter(country TEXT, city TEXT, value INT, PRIMARY KEY(country, city)) WITHOUT ROWID")
	fmt.Println("connected to db")
}
