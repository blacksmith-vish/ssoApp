package sqlite

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/golang-migrate/migrate/v4"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
)

const (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "db_test"
	port     = "5437"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=30"

	migrationPath = "./migrations"
)

var (
	dsnConnString = fmt.Sprintf(dsn, host, port, user, password, dbName)
)

type storeTestSuite struct {
	suite.Suite
	ctx         context.Context
	dbContainer *dockertest.Pool
	resource    *dockertest.Resource
	testDB      *sqlx.DB
	//sut                 *store.SqlServerMoviesStore
	//dbHelper            *databaseHelper
}

func TestMain(t *testing.T) {
	suite.Run(t, new(storeTestSuite))
}

func (suite *storeTestSuite) SetupSuite() {
	var err error

	suite.ctx = context.Background()

	err = createDatabaseContainer(suite)
	if err != nil {
		suite.dbContainer.Purge(suite.resource)
		log.Fatal(err)
	}

	// err = suite.createScheme()
	// if err != nil {
	// 	suite.dbContainer.Purge(suite.resource)
	// 	log.Fatal(err)
	// }

	err = suite.migrateDB()
	if err != nil {
		suite.dbContainer.Purge(suite.resource)
		log.Fatal(err)
	}

}

func (suite *storeTestSuite) TearDownSuite() {
	if err := suite.dbContainer.Purge(suite.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func createDatabaseContainer(suite *storeTestSuite) error {

	pool, err := dockertest.NewPool("")
	if err != nil {
		return errors.Wrap(err, "could not connect to docker")
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

	suite.resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		return errors.Wrap(err, "could not start resource")
	}

	if err := pool.Retry(func() error {
		var err error
		suite.testDB, err = sqlx.Connect("postgres", dsnConnString)
		if err != nil {
			log.Println("Error:", err)
			return err
		}
		return suite.testDB.Ping()
	}); err != nil {
		_ = pool.Purge(suite.resource)
		return errors.Wrap(err, "could not connect to database")
	}

	suite.dbContainer = pool

	return nil
}

func (suite *storeTestSuite) createScheme() error {

	_, err := suite.testDB.Exec("CREATE DATABASE " + dbName)
	if err != nil {
		return err
	}

	// tableSQL, err := os.ReadFile("./testdata/tables.sql")
	// if err != nil {
	// 	return err
	// }

	// _, err = suite.testDB.Exec(string(tableSQL))
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (suite *storeTestSuite) migrateDB() error {

	driver, err := postgres.WithInstance(suite.testDB.DB, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "unable to create db instance")
	}

	migrationSource := "file://" + migrationPath

	SrcDriver, err := source.Open(migrationSource)
	if err != nil {
		return errors.Wrap(err, "unable to get source")
	}

	m, err := migrate.NewWithInstance("migration_embeded_sql_files", SrcDriver, "psql_db", driver)
	if err != nil {
		return errors.Wrap(err, "failed to init migrate")
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}

		return errors.Wrap(err, "failed to run migrate")
	}

	return nil
}
