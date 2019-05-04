package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"upsilon_cities_go/config"

	// needed for postgres driver
	"github.com/lib/pq"
)

// Handler Contains DB related informations
type Handler struct {
	db   *sql.DB
	open bool
}

// New Create a new handler for database, ensure database is created
func New() *Handler {
	handler := new(Handler)
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		config.DB_USER, config.DB_PASSWORD, config.DB_NAME, config.DB_HOST, config.DB_PORT)

	db, _ := sql.Open("postgres", dbinfo)

	errPing := db.Ping()
	if err, ok := errPing.(*pq.Error); ok {
		log.Fatalf("DB: Database failed to be connected: %s", err)
	} else {
		log.Printf("DB: Successfully connected to : %s %s", config.DB_HOST, config.DB_NAME)
	}

	handler.db = db
	handler.open = true
	return handler
}

// Exec executes provided query and check if it's correctly executed or not.
// Abort app if not.
func (dbh *Handler) Exec(query string) (result *sql.Rows) {
	dbh.CheckState()
	log.Printf("DB: About to Exec: %s", query)
	result, err := dbh.db.Query(query)
	errorCheck(query, err)
	return result
}

// Query Just like Exec but uses Postgres formater.
func (dbh *Handler) Query(format string, a ...interface{}) (result *sql.Rows) {
	dbh.CheckState()
	log.Printf("DB: About to Query: %s", format)
	result, err := dbh.db.Query(format, a...)
	errorCheck(format, err)
	return result
}

// CheckState assert that connection to DB is still alive. or break
func (dbh *Handler) CheckState() {
	if !dbh.open {
		log.Fatal("DB: Can't use this connection, it's been closed")
	}
	err := dbh.db.Ping()
	if err != nil {
		log.Fatalf("DB: Can't use this connection, an error occured: %s", err)
	}
}

// Close frees db ressource
func (dbh *Handler) Close() {
	if dbh.open {
		dbh.open = false
		defer dbh.db.Close()
	} else {
		log.Print("DB: Already Closed")
	}
}

// ErrorCheck checks if query result has an error or not
func errorCheck(query string, err error) bool {
	if err != nil {
		log.Printf("DB: Failed to execute query: %s", query)

		// fatal aborts app
		log.Fatalf("DB: Aborting: %s", err)

		return true
	}

	return false
}

func CheckVersion(dbh *Handler) {

	dbh.CheckState()
	log.Printf("DB: About to Query: select * from versions")
	result, err := dbh.db.Query("select applied, file from versions order by applied DESC;")

	// ensure last migration date is way in the past.
	applied_migrations := make(map[string]time.Time)
	if err != nil {
		// version table doesn't exist: create database.
		f, ferr := os.Open(config.DB_SCHEMA)
		if ferr != nil {
			log.Fatalln("DB: No schema file found can't initialize database")
		}
		schema, ferr := ioutil.ReadAll(f)
		if ferr != nil {
			log.Fatalln("DB: Schema found but unable to read it all.")
		}

		_, err := dbh.db.Query(string(schema))

		if err != nil {
			log.Printf("DB: While executing: %s", string(schema))
			log.Fatalf("DB: Unable to apply schema %s ", err)
		}

		// should insert at least something in version ...

		dbh.Query("insert into versions(file) values ('schema.sql');")

		applied_migrations["schema.sql"] = time.Now().UTC()

		// expect schema to be same as all migrations ... ;)
		err = filepath.Walk(config.DB_MIGRATIONS, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatalf("DB: prevent panic by handling failure accessing a path %q: %v\n", config.DB_MIGRATIONS, err)
				return err
			}
			if strings.HasSuffix(info.Name(), ".sql") {

				dbh.Query("insert into versions(file) values ($1);", path)

				applied_migrations[path] = time.Now().UTC()
			}

			return nil
		})
	} else {

		// get lastest applied migration date.
		for result.Next() {
			var applied time.Time
			var file string

			result.Scan(&applied, &file)
			applied_migrations[file] = applied
		}

	}

	// thus we keep order here ;)
	var orderedFiles []string
	log.Printf("DB: Attempting to find migrations in: %s", config.DB_MIGRATIONS)
	err = filepath.Walk(config.DB_MIGRATIONS, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("DB: prevent panic by handling failure accessing a path %q: %v\n", config.DB_MIGRATIONS, err)
			return err
		}
		if strings.HasSuffix(info.Name(), ".sql") {

			migrationFilename := strings.TrimLeft(strings.Replace(path, config.DB_MIGRATIONS, "", 1), "\\")
			dateString := strings.Split(migrationFilename, "_")[0]
			_, err := time.Parse("200601021504", dateString)

			if err != nil {
				log.Fatalf("DB: One one the migration files: %s, has an invalid format. Expected YYYYMMDDHHMM_name.sql", migrationFilename)
				return err
			}
			log.Printf("DB: Read Migration file: %s", migrationFilename)

			orderedFiles = append(orderedFiles, path)

		}
		return nil
	})

	// ensure file are ordered (they should be by date ;)
	sort.Strings(orderedFiles)

	if err != nil {
		return
	}

	for _, k := range orderedFiles {
		if _, found := applied_migrations[k]; !found {
			log.Printf("DB: Applying migration: %s", k)
			f, ferr := os.Open(k)
			if ferr != nil {
				log.Fatalf("DB: Unable to open migration file %s", k)
			}

			migration, ferr := ioutil.ReadAll(f)
			if ferr != nil {
				log.Fatalf("DB: Unable to read migration file %s", k)
			}

			log.Printf("DB: Applying migration: %s", string(migration))
			_, err := dbh.db.Query(string(migration))

			if err != nil {
				log.Fatalf("DB: Unable to apply migration file %s: %s ", k, err)
			}

			dbh.db.Query("insert into versions(file) values ($1);", k)
		}
	}
	log.Printf("DB: DB is up to date ! ")

}
