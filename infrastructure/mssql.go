package infrastructure

import (
	"database/sql"
	"github.com/arzonus/cveapi/interfaces"
	_ "github.com/denisenkom/go-mssqldb"
)

type MSSQLHandler struct {
	Conn   *sql.DB
	URI    string
	Logger Logger
}

func (handler *MSSQLHandler) Execute(statement string) {
	handler.connect()
	_, err := handler.Conn.Exec(statement)
	if err != nil {
		handler.Logger.Info(statement)
		handler.Logger.Info(err)
	}

}

func (handler *MSSQLHandler) Query(statement string) interfaces.Row {
	handler.connect()
	rows, err := handler.Conn.Query(statement)
	if err != nil {
		handler.Logger.Info(err)
		return new(MSSQLRow)
	}
	row := new(MSSQLRow)
	row.Rows = rows
	return row
}

func (handler *MSSQLHandler) connect() {

	if handler.Conn == nil {
		handler.Logger.Info("Connect to MSSQL")
		var err error
		handler.Conn, err = sql.Open("mssql", handler.URI)
		if err != nil {
			handler.Logger.Info(err)
		}
	}
}

type MSSQLRow struct {
	Rows *sql.Rows
}

func (r MSSQLRow) Scan(dest ...interface{}) {
	r.Rows.Scan(dest...)
}

func (r MSSQLRow) Next() bool {
	return r.Rows.Next()
}

func NewMSSQLHandler(uri string) *MSSQLHandler {
	mssqlHandler := new(MSSQLHandler)
	mssqlHandler.URI = uri
	mssqlHandler.connect()
	return mssqlHandler
}
