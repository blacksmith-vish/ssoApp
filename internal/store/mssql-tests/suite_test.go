package sqlstore

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"

	// _ "github.com/denisenkom/go-mssqldb"

	"assmr-chat-mobile-rubricator.cbr.ru/model"
	"assmr-chat-mobile-rubricator.cbr.ru/utils"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
)

var (
	host     = "localhost"
	port     = "5437"
	user     = "sa"
	password = "SuperPassword123"
	dbName   = "db_test"
	driver   = "sqlserver"
)

var (
	dsnConnString = fmt.Sprintf("%s://%s:%s@%s:%s", driver, user, password, host, port)
)

type storeTestSuite struct {
	suite.Suite
	ctx         context.Context
	dbContainer *dockertest.Pool
	resource    *dockertest.Resource
	testDB      *sql.DB
	//sut                 *store.SqlServerMoviesStore
	dbHelper *SqlSupplier
}

func TestMain(t *testing.T) {

	testSuite := new(storeTestSuite)

	suite.Run(t, testSuite)

}

func (suite *storeTestSuite) SetupSuite() {
	var err error

	suite.ctx = context.Background()

	log.Println("Pre create")

	err = createDatabaseContainer(suite)
	suite.closeIfErr(err)

	suite.T().Cleanup(
		func() {
			suite.dbContainer.Purge(suite.resource)
		})

	log.Println("Pre migrate")

	err = suite.migrateDB()
	suite.closeIfErr(err)

	// dbmap = &gorp.DbMap{Db: db,
	// 	TypeConverter: mattermConverter{}, Dialect: gorp.SqlServerDialect{"2016"}, QueryTimeout: connectionTimeout}
	err = utils.TranslationsPreInit()
	suite.closeIfErr(err)

	var (
		driverMSSQL = "mssql"

		maxIdleConns = 20
		maxOpenConns = 300

		queryTimeout                = 30
		connMaxLifetimeMilliseconds = 3600000
	)

	suite.dbHelper = NewSqlSupplier(model.SqlSettings{
		DriverName:                  &driverMSSQL,
		DataSource:                  &dsnConnString,
		MaxIdleConns:                &maxIdleConns,
		MaxOpenConns:                &maxOpenConns,
		Trace:                       false,
		AtRestEncryptKey:            "tisc6gmn7zohwxjc866cz6pxn79n4pmf",
		QueryTimeout:                &queryTimeout,
		ConnMaxLifetimeMilliseconds: &connMaxLifetimeMilliseconds,
	})

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
