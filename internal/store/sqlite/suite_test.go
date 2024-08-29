package sqlite

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlserver"

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
	user     = "sa"
	password = "Lnxpass1"
	dbName   = "db_test"
	port     = "5437"

	driver = "sqlserver"

	migrationPath = "./migrations_sqlserver"
)

var (
	dsnConnString = fmt.Sprintf("%s://%s:%s@%s:%s", driver, user, password, host, port)
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
		_ = suite.dbContainer.Purge(suite.resource)
		log.Fatal(err)
	}

	// err = suite.createScheme()
	// if err != nil {
	// 	_ = suite.dbContainer.Purge(suite.resource)
	// 	log.Fatal(err)
	// }

	err = suite.migrateDB()
	if err != nil {
		_ = suite.dbContainer.Purge(suite.resource)
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
		Name:       "mssql-test",
		Repository: "mcr.microsoft.com/mssql/server",
		Tag:        "2019-latest",
		Env: []string{
			"ACCEPT_EULA=Y",
			"MSSQL_SA_PASSWORD=" + password,
		},
		ExposedPorts: []string{"1433"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"1433/tcp": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	suite.resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(suite.resource)
		return errors.Wrap(err, "could not start resource")
	}

	pool.MaxWait = time.Second * 60

	if err := pool.Retry(func() error {
		var err error
		suite.testDB, err = sqlx.Connect(driver, dsnConnString)
		if err != nil {
			if err.Error() != "EOF" {
				log.Println("Error:", err)
			}
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

	driver, err := sqlserver.WithInstance(suite.testDB.DB, &sqlserver.Config{})
	if err != nil {
		return errors.Wrap(err, "unable to create db instance")
	}

	migrationSource := "file://" + migrationPath

	SrcDriver, err := source.Open(migrationSource)
	if err != nil {
		return errors.Wrap(err, "unable to get source")
	}

	m, err := migrate.NewWithInstance("migration_embeded_sql_files", SrcDriver, dbName, driver)
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
