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

func query(ctx context.Context, db *sql.DB, stmt string, args ...any) *sql.Rows {
	res, err := db.QueryContext(ctx, stmt, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query %s: %s", stmt, err)
		os.Exit(1)
	}
	return res
}

func queryConn(ctx context.Context, conn *sql.Conn, stmt string, args ...any) *sql.Rows {
	res, err := conn.QueryContext(ctx, stmt, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query %s: %s", stmt, err)
		os.Exit(1)
	}
	return res
}

func execTx(ctx context.Context, tx *sql.Tx, stmt string, args ...any) sql.Result {
	res, err := tx.ExecContext(ctx, stmt, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute statement %s: %s", stmt, err)
		os.Exit(1)
	}
	return res
}

func queryTx(ctx context.Context, tx *sql.Tx, stmt string, args ...any) *sql.Rows {
	res, err := tx.QueryContext(ctx, stmt, args...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute query %s: %s", stmt, err)
		os.Exit(1)
	}
	return res
}

func runCounterExample(dbPath string) {
	db, err := sql.Open("libsql", dbPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", dbPath, err)
		os.Exit(1)
	}
	ctx := context.Background()
	exec(ctx, db, "CREATE TABLE IF NOT EXISTS counter(country TEXT, city TEXT, value INT, PRIMARY KEY(country, city)) WITHOUT ROWID")

	incCounterStatementPositionalArgs := "INSERT INTO counter(country, city, value) VALUES(?, ?, 1) ON CONFLICT DO UPDATE SET value = IFNULL(value, 0) + 1 WHERE country = ? AND city = ?"
	exec(ctx, db, incCounterStatementPositionalArgs, "PL", "WAW", "PL", "WAW")
	exec(ctx, db, incCounterStatementPositionalArgs, "FI", "HEL", "FI", "HEL")
	exec(ctx, db, incCounterStatementPositionalArgs, "FI", "HEL", "FI", "HEL")
	incCounterStatementNamedArgs := "INSERT INTO counter(country, city, value) VALUES(:country, :city, 1) ON CONFLICT DO UPDATE SET value = IFNULL(value, 0) + 1 WHERE country = :country AND city = :city"
	exec(ctx, db, incCounterStatementNamedArgs, sql.Named("country", "PL"), sql.Named("city", "WAW"))
	exec(ctx, db, incCounterStatementNamedArgs, sql.Named("country", "FI"), sql.Named("city", "HEL"))
	incCounterStatementNamedArgs2 := "INSERT INTO counter(country, city, value) VALUES(@country, @city, 1) ON CONFLICT DO UPDATE SET value = IFNULL(value, 0) + 1 WHERE country = @country AND city = @city"
	exec(ctx, db, incCounterStatementNamedArgs2, sql.Named("country", "PL"), sql.Named("city", "WAW"))
	exec(ctx, db, incCounterStatementNamedArgs2, sql.Named("country", "FI"), sql.Named("city", "HEL"))
	incCounterStatementNamedArgs3 := "INSERT INTO counter(country, city, value) VALUES($country, $city, 1) ON CONFLICT DO UPDATE SET value = IFNULL(value, 0) + 1 WHERE country = $country AND city = $city"
	exec(ctx, db, incCounterStatementNamedArgs3, sql.Named("country", "PL"), sql.Named("city", "WAW"))
	exec(ctx, db, incCounterStatementNamedArgs3, sql.Named("country", "FI"), sql.Named("city", "HEL"))

	// try prepared statements
	{
		stmt, err := db.Prepare("UPDATE counter SET value = value + 1 WHERE country = ? AND city = ?")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to prepare statement %s: %s", incCounterStatementPositionalArgs, err)
			os.Exit(1)
		}
		defer stmt.Close()
		_, err = stmt.Exec("FI", "HEL")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to execute prepared statement %s for FI: %s", incCounterStatementPositionalArgs, err)
			os.Exit(1)
		}
		_, err = stmt.Exec("PL", "WAW")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to execute prepared statement %s for PL: %s", incCounterStatementPositionalArgs, err)
			os.Exit(1)
		}
	}
	{
		stmt, err := db.Prepare("SELECT * FROM counter")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to prepare statement %s: %s", "SELECT * FROM counter", err)
			os.Exit(1)
		}
		defer stmt.Close()
		rows, err := stmt.Query()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to execute prepared statement %s: %s", "SELECT * FROM counter", err)
			os.Exit(1)
		}
		for rows.Next() {
			var row struct {
				country string
				city    string
				value   int
			}
			if err := rows.Scan(&row.country, &row.city, &row.value); err != nil {
				fmt.Fprintf(os.Stderr, "failed to scan row: %s", err)
				os.Exit(1)
			}
			fmt.Println(row)
		}
		if err := rows.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "errors from query: %s", err)
			os.Exit(1)
		}
	}

	rows := query(ctx, db, "SELECT * FROM counter")
	for rows.Next() {
		var row struct {
			country string
			city    string
			value   int
		}
		if err := rows.Scan(&row.country, &row.city, &row.value); err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan row: %s", err)
			os.Exit(1)
		}
		fmt.Println(row)
	}
	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "errors from query: %s", err)
		os.Exit(1)
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start a transaction: %s", err)
		os.Exit(1)
	}
	// Defer a rollback in case anything fails.
	defer tx.Rollback()
	rows = queryTx(ctx, tx, `SELECT * FROM counter WHERE (country = "PL" AND city = "WAW") OR (country = "FI" AND city = "HEL")`)
	wawValue := -1
	helValue := -1
	for rows.Next() {
		var row struct {
			country string
			city    string
			value   int
		}
		if err := rows.Scan(&row.country, &row.city, &row.value); err != nil {
			fmt.Fprintf(os.Stderr, "failed to scan row: %s", err)
			os.Exit(1)
		}
		if row.country == "PL" && row.city == "WAW" {
			wawValue = row.value
		}
		if row.country == "FI" && row.city == "HEL" {
			helValue = row.value
		}
	}
	if err := rows.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "errors from query: %s", err)
		os.Exit(1)
	}
	if helValue > wawValue {
		execTx(ctx, tx, `INSERT INTO counter(country, city, value) VALUES("PL", "WAW", ?) ON CONFLICT DO UPDATE SET value = ? WHERE country = "PL" AND city = "WAW"`, helValue, helValue)
	}
	if err = tx.Commit(); err != nil {
		fmt.Fprintf(os.Stderr, "error commiting the transaction: %s", err)
		os.Exit(1)
	}
}

func ConnectToDb() {
	var dbUrl = "file:" + dbLocation
	runCounterExample(dbUrl)
	fmt.Println("connected to db")
}
