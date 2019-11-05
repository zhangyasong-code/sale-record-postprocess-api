package customer

import (
	"github.com/go-xorm/xorm"
)

func InitDB(db *xorm.Engine) error {
	return db.Sync(
		new(PostMileage),
	)
}
