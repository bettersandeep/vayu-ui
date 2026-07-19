package optimization

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatabaseFilterRoundTrip(t *testing.T) {
	require.Equal(t, "", buildDatabaseFilter(nil))
	require.Equal(t, "", buildDatabaseFilter([]string{" ", ""}))
	require.Equal(t, "^(sales_db)$", buildDatabaseFilter([]string{"sales_db"}))
	require.Equal(t, "^(sales_db|olake_prod)$", buildDatabaseFilter([]string{"sales_db", "olake_prod"}))

	require.Nil(t, parseDatabaseFilter(""))
	require.Nil(t, parseDatabaseFilter("^(?!legacy_db$).*$"), "hand-written filters do not round-trip, and must not be mangled")
	require.Equal(t, []string{"sales_db", "olake_prod"}, parseDatabaseFilter("^(sales_db|olake_prod)$"))
	require.Equal(t, []string{"sales_db"}, parseDatabaseFilter(buildDatabaseFilter([]string{"sales_db"})))
}
