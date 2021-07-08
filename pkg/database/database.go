package database

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
)

type DB struct {
	DBClient *pgxpool.Pool
}

func GetDBConnection(connectionStr string) (*DB, error) {
	db, err := getDBConnection(connectionStr)

	if err != nil {
		return nil, err
	}

	return &DB{
		DBClient: db,
	}, nil
}

func CloseDBConnection(db *DB) {
	db.DBClient.Close()
}

func getDBConnection(connectionStr string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.Connect(context.Background(), connectionStr)

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}

	return conn, nil
}

// Helper function to generate bulk insert query
func (db *DB) GetBulkInsertQuery(SQLString string, numArgsPerRow int, numRows int) string {
	questionMarks := make([]string, 0, numArgsPerRow)
	for i := 0; i < numArgsPerRow; i++ {
		questionMarks = append(questionMarks, "?")
	}
	rowValueSQL := strings.Join(questionMarks, ", ")
	return getBulkInsertSQL(SQLString, rowValueSQL, numRows)
}

func getBulkInsertSQL(SQLString string, rowValueSQL string, numRows int) string {
	// Combine the base SQL string and N value strings
	valueStrings := make([]string, 0, numRows)
	for i := 0; i < numRows; i++ {
		valueStrings = append(valueStrings, "("+rowValueSQL+")")
	}
	allValuesString := strings.Join(valueStrings, ",")
	SQLString = fmt.Sprintf(SQLString, allValuesString)

	// Convert all of the "?" to "$1", "$2", "$3", etc.
	// (which is the way that pgx expects query variables to be)
	numArgs := strings.Count(SQLString, "?")
	SQLString = strings.ReplaceAll(SQLString, "?", "$%v")
	numbers := make([]interface{}, 0, numRows)
	for i := 1; i <= numArgs; i++ {
		numbers = append(numbers, strconv.Itoa(i))
	}
	return fmt.Sprintf(SQLString, numbers...)
}
