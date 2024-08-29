package sqlite

import (
	"context"
	"fmt"
	"log"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
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
	suite.closeIfErr(err)

	err = suite.migrateDB()
	suite.closeIfErr(err)
}

func (suite *storeTestSuite) TearDownSuite() {
	if err := suite.dbContainer.Purge(suite.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func (suite *storeTestSuite) closeIfErr(err error) {
	if err != nil {
		if err := suite.dbContainer.Purge(suite.resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
		log.Fatal(err)
	}
}
