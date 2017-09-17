package deebee

//Provider is an interface implemented by various DB providers
type Provider interface {
	Query(stmt *Statement) ([]*Row, error)
	Close() error
	Ping() error
	HasTable(i interface{}) bool
	CreateTable(i interface{}, foreignKeys ...ForeignKey) (stmt *Statement, err error)
	Migrate(i interface{}) error
	Select(i interface{}, excluding ...string) (*Statement, error)
	Insert(i interface{}) (*Statement, error)
	Update(i interface{}, excluding ...string) (*Statement, error)
	Delete(i interface{}) (*Statement, error)
}
