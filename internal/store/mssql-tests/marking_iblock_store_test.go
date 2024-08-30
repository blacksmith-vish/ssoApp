package sqlstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func (suite *storeTestSuite) TestMarkingIblock() {

	t := suite.T()

	t.Run("MarkingIblock GET, should be error", getMarkingIblock(suite))

}

func getMarkingIblock(suite *storeTestSuite) func(t *testing.T) {

	return func(t *testing.T) {

		response := <-suite.dbHelper.MarkingIblock().Get("sdfsvd", false)

		assert.NotNil(t, response.Err)

		assert.Equal(t, "store.sql_marking_iblock.get.existing.app_error", response.Err.Id)

		assert.Nil(t, response.Data)

	}
}
