package server

import (
	"fmt"
	"time"

	"github.com/guricerin/stop-now-smoking/entity"
)

type cigaretteTable struct {
	Id          int64
	SmokedCount uint
	UserId      uint64
	CreatedAt   time.Time
}

func toCigaretteTable(c entity.Cigarette) cigaretteTable {
	table := cigaretteTable{
		SmokedCount: c.SmokedCount,
		UserId:      uint64(c.UserId),
		CreatedAt:   c.CreatedAt,
	}
	return table
}

func toCigarreteEntity(t cigaretteTable) entity.Cigarette {
	c := entity.Cigarette{
		Id:          t.Id,
		SmokedCount: t.SmokedCount,
		UserId:      int64(t.UserId),
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

func (store *cigaretteStore) RetrieveAllByUserId(id int64) ([]entity.Cigarette, error) {
	rows, err := store.db.Query("select id, smoked_count, user_id, created_at from cigarettes where user_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cigarettes []entity.Cigarette
	for rows.Next() {
		var t cigaretteTable
		err := rows.Scan(&t.Id, &t.SmokedCount, &t.UserId, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		c := toCigarreteEntity(t)
		cigarettes = append(cigarettes, c)
	}
	return cigarettes, nil
}

const layout = "2006-01-02"

func (store *cigaretteStore) toDateFormat(start, end time.Time) (string, string) {
	startStr := start.Format(layout)
	endStr := end.Format(layout)
	startStr = fmt.Sprintf("%s 00:00:00", startStr)
	endStr = fmt.Sprintf("%s 23:59:59", endStr)
	return startStr, endStr
}

func (store *cigaretteStore) RetrieveAllByUserIdAndBetweenDate(id int64, start, end time.Time) ([]entity.Cigarette, error) {
	startStr, endStr := store.toDateFormat(start, end)
	rows, err := store.db.Query("select id, smoked_count, user_id, created_at from cigarettes where user_id = ? and created_at between ? and ?", id, startStr, endStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cigarettes []entity.Cigarette
	for rows.Next() {
		var t cigaretteTable
		err := rows.Scan(&t.Id, &t.SmokedCount, &t.UserId, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		c := toCigarreteEntity(t)
		cigarettes = append(cigarettes, c)
	}
	return cigarettes, nil
}

func (store *cigaretteStore) ExistByUserIdAndDate(cig entity.Cigarette) (bool, error) {
	table := toCigaretteTable(cig)
	startStr, endStr := store.toDateFormat(table.CreatedAt, table.CreatedAt)
	rows, err := store.db.Query("select * from cigarettes where user_id = ? and created_at between ? and ?", table.UserId, startStr, endStr)
	if err != nil {
		return false, err
	} else {
		defer rows.Close()
		return rows.Next(), nil
	}
}

func (store *cigaretteStore) UpdateByUserIdAndDate(cig entity.Cigarette) (err error) {
	table := toCigaretteTable(cig)
	startStr, endStr := store.toDateFormat(table.CreatedAt, table.CreatedAt)
	_, err = store.db.Exec("update cigarettes set smoked_count = ? where user_id = ? and created_at between ? and ?", table.SmokedCount, table.UserId, startStr, endStr)
	return
}
