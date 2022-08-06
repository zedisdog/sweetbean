package contract

import "github.com/zedisdog/sweetbean/database"

type CanUseTx[Tx any, Self any] interface {
	WithTx(Tx) Self
}

type CanMakeTx[Tx any] interface {
	Transaction(func(Tx) error) error
}

//Transaction 事务接口
type Transaction[Tx any, Self any] interface {
	CanUseTx[Tx, Self]
	CanMakeTx[Tx]
}

//Repo repo基础方法接口
type Repo[Model any] interface {
	Create(*Model) error
	Update(*Model) error
	Delete(...database.Condition) error
	First(...database.Condition) (Model, error)
}

type CanFind[Model any] interface {
	Find(...database.Condition) ([]Model, error)
}

type CanPage[Model any] interface {
	Page(offset, limit int, conditions ...database.Condition) (list []Model, total int, err error)
}

type Assert interface {
	Exists(...database.Condition) (bool, error)
}

type CanCount interface {
	Count(...database.Condition) (count int, err error)
}
