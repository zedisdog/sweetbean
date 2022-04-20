package entity

type Entity interface {
	//Assimilate copy field from [from] to self
	Assimilate(from interface{}) error
	IsDirty() bool
}
