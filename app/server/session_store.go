package server

import (
	"time"

	"github.com/guricerin/stop-now-smoking/entity"
)

type sessionTable struct {
	Id        int64
	Uuid      string
	UserId    int64
	CreatedAt time.Time
}

func toSessionTable(s entity.Session) (table sessionTable) {
	table = sessionTable{
		Id:        int64(s.Id),
		Uuid:      s.Uuid,
		UserId:    int64(s.UserId),
		CreatedAt: s.CreatedAt,
	}
	return
}

func toSessionEntity(t sessionTable) (s entity.Session) {
	s = entity.Session{
		Id:        t.Id,
		Uuid:      t.Uuid,
		UserId:    t.UserId,
		CreatedAt: t.CreatedAt,
	}
	return
}

type sessionStore struct {
	db DbDriver
}

func NewSessionStore(db DbDriver) *sessionStore {
	return &sessionStore{db: db}
}

func (store *sessionStore) Create(u entity.User) (sess entity.Session, err error) {
	uuid, err := entity.CreateUuid()
	if err != nil {
		return
	}
	table := sessionTable{
		Uuid:      uuid,
		UserId:    int64(u.Id),
		CreatedAt: time.Now(),
	}
	res, err := store.db.Exec("insert into sessions (uuid, user_id, created_at) values (?, ?, ?)", table.Uuid, table.UserId, table.CreatedAt)
	if err != nil {
		return
	}
	id64, err := res.LastInsertId()
	if err != nil {
		return
	}
	table.Id = id64
	sess = toSessionEntity(table)
	return
}

func (store *sessionStore) RetrieveByUuid(uuid string) (sess entity.Session, err error) {
	table := sessionTable{}
	err = store.db.QueryRow("select id, uuid, user_id, created_at from user_sessions where uuid = ?", uuid).Scan(&table.Id, &table.Uuid, &table.UserId, &table.CreatedAt)
	if err != nil {
		return
	}
	sess = toSessionEntity(table)
	return
}
