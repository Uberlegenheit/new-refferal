package dao

import (
	"fmt"
)

const (
	QueryUpdateMetric = "UPDATE"
	QueryInsertMetric = "INSERT"
	QuerySelectMetric = "SELECT"
	QueryDeleteMetric = "DELETE"
	QueryBuildMetric  = "BUILD"
	QueryGetMetric    = "GET"
	QuerySetMetric    = "SET"
	QuerySetExMetric  = "SETEX"
)

var ErrAlreadyExistsDesc = fmt.Sprintf("duplicate key value")
var ErrDuplicate = fmt.Errorf("DB duplicate error")
var ErrNoRows = fmt.Errorf("DB no rows in result set")

type (
	DAOTx interface {
		CommitTx() error
		RollbackTx() error
	}

	ServiceDAO interface {
		BeginTx() (DAOTx, error)
	}
)
