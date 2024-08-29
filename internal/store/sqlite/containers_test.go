package sqlite

import (
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlserver"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pkg/errors"
)

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
