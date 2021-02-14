package infra

type InfraCollection interface {
	GetSqlDb() SqlDb
}

type infraCollection struct {
	sqlDb SqlDb
}

func NewInfraCollection(sqlDatabaseName string) InfraCollection {
	return &infraCollection{
		sqlDb: NewSqlDb(sqlDatabaseName),
	}
}

func (i *infraCollection) GetSqlDb() SqlDb {
	return i.sqlDb
}
