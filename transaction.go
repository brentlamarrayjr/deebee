package deebee

type transaction struct{}

func (tx *transaction) Commit() {}

func (tx *transaction) Rollback() {}
