package utildb_test

import (
	"fmt"
	"testing"

	"blog/util/utildb"

	"github.com/stretchr/testify/require"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
)

func TestSchema(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	mock.ExpectExec("DROP SCHEMA IF EXISTS foo CASCADE").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("CREATE SCHEMA foo").WillReturnResult(sqlmock.NewResult(0, 1))

	err = utildb.CreateSchema(db, "foo")
	require.NoError(t, err)

	mock.ExpectExec("DROP SCHEMA IF EXISTS").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec("CREATE SCHEMA").WillReturnResult(sqlmock.NewResult(0, 1))

	schema, err := utildb.CreateRandomSchema(db)
	require.NoError(t, err)
	require.NotEmpty(t, schema)

	mock.ExpectExec("DROP SCHEMA foo CASCADE").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(fmt.Sprintf("DROP SCHEMA %s CASCADE", schema)).WillReturnResult(sqlmock.NewResult(0, 1))

	err = utildb.DropSchema(db, "foo")
	require.NoError(t, err)
	err = utildb.DropSchema(db, schema)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}
