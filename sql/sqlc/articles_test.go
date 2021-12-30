package sqlc_db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListArticles(t *testing.T) {
	_, err := testQueries.ListArticles(context.Background())
	require.NoError(t, err)
}
