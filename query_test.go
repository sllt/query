package query

import (
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

var DB *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})

	DB = db
}

type User struct {
	ID   int `gorm:"primaryKey"`
	Sex  string
	Age  int
	Name string
}

func TestQueryAnd(t *testing.T) {
	u := &User{
		Sex:  "male",
		Age:  18,
		Name: "john",
	}
	sql := DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Scopes(Where(u)).Find(&User{})
	})
	assert.Equal(t, sql,
		"SELECT * FROM `users` WHERE `sex` = \"male\" AND `age` = 18 AND `name` = \"john\"",
		"should be equal",
	)
}

func TestQueryDefault(t *testing.T) {
	u := &User{
		Age:  18,
		Name: "john",
	}
	sql := DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Scopes(Where(u)).Find(&User{})
	})
	assert.Equal(t, sql,
		"SELECT * FROM `users` WHERE `age` = 18 AND `name` = \"john\"",
		"should be equal",
	)
}

type Pagination struct {
	Page int `op:"page"`
	Size int `op:"size"`
}

type UserRequest struct {
	Pagination
	Name string
	Age  int
}

func TestStructAnonymous(t *testing.T) {
	req := &UserRequest{
		Pagination: Pagination{
			Page: 3,
			Size: 20,
		},
		Name: "john",
		Age:  18,
	}

	sql := DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Scopes(Where(req)).Find(&User{})
	})
	assert.Equal(t, sql,
		"SELECT * FROM `users` WHERE `name` = \"john\" AND `age` = 18 LIMIT 20 OFFSET 40",
		"should be equal",
	)
}
