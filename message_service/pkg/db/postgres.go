package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pavel/message_service/config"
	"regexp"
	"strings"
)

type DB struct {
	QueryBuilder QueryBuilder
	*sqlx.DB
}

func InitPostgres(cfg *config.Config, queryBuilder QueryBuilder) (error, *DB) {
	connection, err := sqlx.Connect(cfg.DB.Driver, fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB.Driver,
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
		cfg.DB.SSLMode),
	)
	db := &DB{
		DB:           connection,
		QueryBuilder: queryBuilder,
	}

	if err != nil {
		return err, nil
	}

	err = db.Ping()
	if err != nil {
		return err, nil
	}

	return nil, db
}

type QueryBuilder interface {
	Select(query string, args ...interface{}) QueryBuilder
	Insert(query string, args ...interface{}) QueryBuilder
	Values(query string, args ...interface{}) QueryBuilder
	AddValue(query string, args ...interface{}) QueryBuilder
	From(query string, args ...interface{}) QueryBuilder
	Join(query string, args ...interface{}) QueryBuilder
	Where(query string, args ...interface{}) QueryBuilder
	AndWhere(query string, args ...interface{}) QueryBuilder
	OrderBy(query string, args ...interface{}) QueryBuilder
	GroupBy(query string, args ...interface{}) QueryBuilder
	Limit(limit uint64) QueryBuilder
	Offset(offset uint64) QueryBuilder
	ToSql() (sql string, args []interface{})
	NewQueryBuilder() QueryBuilder
}

func InitPostgresQueryBuilder() *PostgresQuery {
	return &PostgresQuery{}
}

type PostgresQuery struct {
	query string
	args  []interface{}
}

func (pq *PostgresQuery) NewQueryBuilder() QueryBuilder {
	return &PostgresQuery{}
}

func (pq *PostgresQuery) Limit(limit uint64) QueryBuilder {
	pq.query = fmt.Sprintf("%s LIMIT %d", pq.query, limit)
	return pq
}

func (pq *PostgresQuery) Offset(offset uint64) QueryBuilder {
	pq.query = fmt.Sprintf("%s OFFSET %d", pq.query, offset)
	return pq
}

func (pq *PostgresQuery) Select(query string, args ...interface{}) QueryBuilder {
	pq.addQueryChain("SELECT ", query, args)
	return pq
}

func (pq *PostgresQuery) From(query string, args ...interface{}) QueryBuilder {
	pq.addQueryChain("FROM ", query, args)
	return pq
}

func (pq *PostgresQuery) Insert(query string, args ...interface{}) QueryBuilder {
	pq.addQueryChain("INSERT INTO ", query, args)
	return pq
}

func (pq *PostgresQuery) Values(query string, args ...interface{}) QueryBuilder {
	pq.addQueryChain("VALUES ", query, args)
	return pq
}

func (pq *PostgresQuery) AddValue(query string, args ...interface{}) QueryBuilder {
	pq.addQueryChain(" ", query, args)
	return pq
}

func (pq *PostgresQuery) Join(query string, args ...interface{}) QueryBuilder {
	pq.addQueryChain("JOIN ", query, args)
	return pq
}

func (pq *PostgresQuery) Where(query string, args ...interface{}) QueryBuilder {
	pq.addQueryChain("WHERE ", query, args)
	return pq
}

func (pq *PostgresQuery) AndWhere(query string, args ...interface{}) QueryBuilder {
	pq.addQueryChain("AND ", query, args)
	return pq
}

func (pq *PostgresQuery) OrderBy(query string, args ...interface{}) QueryBuilder {
	pq.addQueryChain("ORDER BY ", query, args)
	return pq
}

func (pq *PostgresQuery) GroupBy(query string, args ...interface{}) QueryBuilder {
	pq.addQueryChain("GROUP BY ", query, args)

	return pq
}

func (pq *PostgresQuery) addQueryChain(operator, query string, args []interface{}) {
	pq.query = pq.query + operator + " " + strings.Trim(query, " ") + " "
	if args != nil {
		pq.addArgument(args)
	}
}

func (pq *PostgresQuery) ToSql() (sql string, args []interface{}) {
	query, arguments := pq.changeArgsSymbolToNumber(), pq.args
	pq.query = ""
	pq.args = nil
	return query, arguments
}

func (pq *PostgresQuery) addArgument(args []interface{}) {
	pq.args = append(pq.args, args...)
}

func (pq *PostgresQuery) changeArgsSymbolToNumber() string {
	counter := 0
	repl := func(match string) string {
		counter++
		return fmt.Sprintf("$%d", counter)
	}
	re := regexp.MustCompile("@")

	return re.ReplaceAllStringFunc(pq.query, repl)
}
