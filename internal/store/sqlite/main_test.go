package sqlite

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "db_test"
	port     = "5437"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=30"
)

var testDB *sqlx.DB

// var testRepo Repo

func TestMain(m *testing.M) {
	// connect to docker; fail if docker not running
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker; is it running? %s", err)
	}

	opts := dockertest.RunOptions{
		Name:       "postgres-test",
		Repository: "postgres",
		Tag:        "14.5", // same as docker compose
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sqlx.Connect("postgres", fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to database: %s", err)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	code := m.Run()

	// testRepo = &repo{db: testDB}

	os.Exit(code)
}

func createTables() error {
	dbSQL, err := os.ReadFile("./testdata/db.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = testDB.Exec(string(dbSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	time.Sleep(3 * time.Second)

	tableSQL, err := os.ReadFile("./testdata/tables.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("can't ping database")
	}
}

type examples struct {
	ID   int64
	Text string
}

func Test_Select(t *testing.T) {
	m := make([]examples, 0)
	err := testDB.Select(
		&m,
		"SELECT * FROM examples",
	)
	if err != nil {
		t.Error("can't select database" + err.Error())
	}

	for i := range m {
		assert.NotEmpty(t, m[i])
	}
}
