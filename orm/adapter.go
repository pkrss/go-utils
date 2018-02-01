package orm

type OrmAdapterInterface interface {
	RegModel(m BaseModelInterface)
	FindOneByCond(m BaseModelInterface, cond string, val []interface{}, selCols []string) error
	FindOneBySql(m BaseModelInterface, sql string, val ...interface{}) error
	UpdateByCond(m BaseModelInterface, cond string, val []interface{}, selCols []string) error
	DeleteByCond(m BaseModelInterface, cond string, val ...interface{}) error
	Insert(m BaseModelInterface, selCols ...string) error
	LimitSqlStyle() string
	QueryOneBySql(recordPointer interface{}, sql string, val ...interface{}) error
	QueryBySql(recordPointer interface{}, sql string, val ...interface{}) error
	InArg(arg interface{}) interface{}
}

var DefaultOrmAdapter OrmAdapterInterface
