package server

import (
	"github.com/guricerin/stop-now-smoking/entity"
)

type userTable struct {
	Id        int64
	Name      string
	AccountId string
	Password  string
}

func toUserTable(u entity.User) (table userTable) {
	table = userTable{
		Id:        u.Id,
		Name:      u.Name,
		AccountId: u.AccountId,
		Password:  u.Password,
	}
	return
}

func toUserEntity(t userTable) (u entity.User) {
	u = entity.User{
		Id:        t.Id,
		Name:      t.Name,
		AccountId: t.AccountId,
		Password:  t.Password,
	}
	return
}

type userStore struct {
	db DbDriver
}

func NewUserStore(db DbDriver) *userStore {
	return &userStore{db: db}
}

func (repo *userStore) Create(u entity.User) (user entity.User, err error) {
	table := toUserTable(u)
	res, err := repo.db.Exec("insert into users (name, account_id, password) values (?, ?, ?)", table.Name, table.AccountId, table.Password)
	if err != nil {
		return
	}
	id64, err := res.LastInsertId()
	if err != nil {
		return
	}
	user, err = repo.RetrieveById(id64)
	return
}

func (repo *userStore) CheckAccountIdExists(u entity.User) bool {
	rows, err := repo.db.Query("select * from users where account_id = ?", u.AccountId)
	if err == nil && rows.Next() {
		defer rows.Close()
		return true
	} else {
		return false
	}
}

func (repo *userStore) RetrieveById(id int64) (u entity.User, err error) {
	table := userTable{}
	err = repo.db.QueryRow("select id, name, account_id, password from users where id = ?", id).
		Scan(&table.Id, &table.Name, &table.AccountId, &table.Password)
	if err != nil {
		return
	}
	u = toUserEntity(table)
	return
}

func (repo *userStore) RetrieveByAccountId(account_id string) (u entity.User, err error) {
	table := userTable{}
	err = repo.db.QueryRow("select id, name, account_id, password from users where account_id = ?", account_id).
		Scan(&table.Id, &table.Name, &table.AccountId, &table.Password)
	if err != nil {
		return
	}
	u = toUserEntity(table)
	return
}

func (repo *userStore) Update(u entity.User) (err error) {
	return
}

func (repo *userStore) DeleteById(id int64) (err error) {
	_, err = repo.db.Exec("delete from users where id = ?", id)
	return
}
