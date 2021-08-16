package user

import (
	"chain-demo/middleware/db"
	"fmt"
	"time"
)

type User struct {
	Id        uint64    `json:"id,omitempty"`
	Nickname  string    `form:"nickname" json:"nickname,omitempty"`
	Password  string    `form:"password" json:"-"`
	Gender    int64     `json:"gender,omitempty"`
	Birthday  time.Time `json:"birthday,omitempty"`
	CreatedAt time.Time `gorm:"column:created_time" json:"created_time,omitempty"`
	UpdatedAt time.Time `gorm:"column:updated_time" json:"updated_time,omitempty"`
}

//TableName 为User绑定表名
func (u User) TableName() string {
	return "user"
}

// GetById will populate a user object from a database model with
// a matching id.
func (u *User) GetById() *User {
	db.GetDB().Take(u)

	fmt.Println(u.Nickname)
	return u
}
