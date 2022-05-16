package entity

type Entity interface {
	//Assimilate copy field from [from] to self
	Assimilate(from interface{}) error
	Incarnation(to interface{}) error
}

type DirtyState struct {
	Dirty bool `json:"-" gorm:"-"`
}
