package sqlite

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func (suite *storeTestSuite) TestPingDB() {

	t := suite.T()

	t.Run("ping DB, should be no error", func(t *testing.T) {
		err := suite.testDB.Ping()
		if err != nil {
			t.Error("can't ping database")
		}
	})

}

type examples struct {
	ID   int64
	Text string
}

func (suite *storeTestSuite) TestSelect() {

	t := suite.T()

	t.Run("select & check result from examples", func(t *testing.T) {
		m := make([]examples, 0)
		err := suite.testDB.Select(
			&m,
			"SELECT * FROM examples",
		)
		if err != nil {
			t.Error("can't select database " + err.Error())
		}

		for i := range m {
			assert.NotEmpty(t, m[i])
		}
	})

}
