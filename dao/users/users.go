package users

import (
	"github.com/go-pg/pg/v10"
	"github.com/shijting/web/inits/psql"
	"github.com/shijting/web/models/users"
)
// 插入一条用户数据
func Insert(tx *pg.DB, data *users.User) (*users.User, error) {
	if tx == nil {
		tx = psql.GetDB()
	}
	_,err := tx.Model(data).Insert()
	data.Password = ""
	return data, err
}

// 根据xx查询一条用户数据
func GetOneByUsername(tx *pg.DB, username string) (*users.User, error) {
	if tx == nil {
		tx = psql.GetDB()
	}
	result := new(users.User)
	err := tx.Model(result).
		Column("id", "username", "password", "email", "created_date", "updated_date").
		Where("username = ?", username).
		First()

	return result, err
}