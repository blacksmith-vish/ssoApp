package sqlstore

import (
	"testing"

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
	ID   int
	Text string
}

func (suite *storeTestSuite) TestSelect() {

	t := suite.T()

	t.Run("select & check result from examples", func(t *testing.T) {

		rows, err := suite.testDB.Query("SELECT * FROM examples")
		if err != nil {
			t.Error("can't select database " + err.Error())
		}

		m := make([]examples, 0)

		for rows.Next() {
			var id int
			var text string

			err := rows.Scan(&id, &text)
			if err != nil {
				t.Error("can't select database " + err.Error())
			}
			m = append(m, examples{
				ID:   id,
				Text: text,
			})
		}

		for i := range m {
			assert.NotEmpty(t, m[i])
		}
	})

}
