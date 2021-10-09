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
		UserId:    u.Id,
		CreatedAt: time.Now(),
	}
	err = store.db.QueryRow("insert into sessions (uuid, user_id, created_at) values ($1, $2, $3) returning id", table.Uuid, table.UserId, table.CreatedAt).Scan(&table.Id)
	if err != nil {
		return
	}
	sess = toSessionEntity(table)
	return
}

func (store *sessionStore) RetrieveByUuid(uuid string) (sess entity.Session, err error) {
	table := sessionTable{}
	err = store.db.QueryRow("select id, uuid, user_id, created_at from sessions where uuid = $1", uuid).Scan(&table.Id, &table.Uuid, &table.UserId, &table.CreatedAt)
	if err != nil {
		return
	}
	sess = toSessionEntity(table)
	return
}

func (store *sessionStore) DeleteByUuid(uuid string) (err error) {
	_, err = store.db.Exec("delete from sessions where uuid = $1", uuid)
	return
}
