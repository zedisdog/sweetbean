package seed

import "gorm.io/gorm"

type GormSeedFunc func(db *gorm.DB) error

var seeders = make([]GormSeedFunc, 0)

func Add(f ...GormSeedFunc) {
	seeders = append(seeders, f...)
}

func Seed(db *gorm.DB) (err error) {
	for _, seeder := range seeders {
		err = seeder(db)
		if err != nil {
			return
		}
	}

	return
}
