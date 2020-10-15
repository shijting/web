package inits
import (
	"context"
	"github.com/go-pg/pg/v10"
	_ "github.com/go-pg/pg/v10/orm"
)
var db *pg.DB
func init()  {
	db = pg.Connect(&pg.Options{
		Addr:     "47.103.11.253:5432",
		User:     "postgres",
		Password: "shijinting0510",
		Database: "test",
	})
	ctx := context.Background()

	if err := db.Ping(ctx); err != nil {
		panic(err)
	}
}
func GetDB() *pg.DB {
	return db
}