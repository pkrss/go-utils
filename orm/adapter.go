package orm

type OrmAdapterInterface interface {
	RegModel(model BaseModelInterface)

	ExecSql(sql string, val ...interface{}) error
	QueryOneBySql(outputRecord interface{}, sql string, val ...interface{}) error
	QueryBySql(outputRecords interface{}, sql string, val ...interface{}) error

	SqlInArg(arg interface{}) interface{}
	SqlReturnSql() string
	SqlLimitStyle() string
}

var DefaultOrmAdapter OrmAdapterInterface
