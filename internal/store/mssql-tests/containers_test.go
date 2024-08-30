package sqlstore

import (
	"database/sql"
	"embed"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

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
		suite.testDB, err = sql.Open(driver, dsnConnString)
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

func (suite *storeTestSuite) migrateDB() error {

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("sqlserver"); err != nil {
		return errors.Wrap(err, "could not set dialect")
	}

	if err := goose.Up(suite.testDB, "migrations"); err != nil {
		return errors.Wrap(err, "could not run up")
	}
	return nil
}
