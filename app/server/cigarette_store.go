package server

import (
	"time"

	"github.com/guricerin/stop-now-smoking/entity"
)

type cigaretteTable struct {
	Id          int64
	SmokedCount int
	UserId      int64
	CreatedAt   time.Time
}

func toCigaretteTable(c entity.Cigarette) cigaretteTable {
	table := cigaretteTable{
		SmokedCount: c.SmokedCount,
		UserId:      c.UserId,
		CreatedAt:   c.CreatedAt,
	}
	return table
}

func toCigarreteEntity(t cigaretteTable) entity.Cigarette {
	c := entity.Cigarette{
		Id:          t.Id,
		SmokedCount: t.SmokedCount,
		UserId:      t.UserId,
		CreatedAt:   t.CreatedAt,
	}
	return c
}

type cigaretteStore struct {
	db DbDriver
}

func NewCigaretteStore(db DbDriver) *cigaretteStore {
	return &cigaretteStore{db: db}
}

func (store *cigaretteStore) Create(cigarette entity.Cigarette) (c entity.Cigarette, err error) {
	table := toCigaretteTable(cigarette)
	res, err := store.db.Exec("insert into cigarettes (smoked_count, user_id, created_at) values (?, ?, ?)", table.SmokedCount, table.UserId, table.CreatedAt)
	if err != nil {
		return
	}
	id64, err := res.LastInsertId()
	if err != nil {
		return
	}
	c, err = store.RetrieveById(id64)
	return
}

func (store *cigaretteStore) RetrieveById(id int64) (c entity.Cigarette, err error) {
	table := cigaretteTable{}
	err = store.db.QueryRow("select id, smoked_count, user_id, created_at from cigarettes where id = ?", id).
		Scan(&table.Id, &table.SmokedCount, &table.UserId, &table.CreatedAt)
	if err != nil {
		return
	}
	c = toCigarreteEntity(table)
	return
}
