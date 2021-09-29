package server

import "github.com/guricerin/stop-now-smoking/entity"

type followTable struct {
	Id        int64
	SrcUserId int64
	DstUserid int64
}

type followStore struct {
	db DbDriver
}

func NewFollowStore(db DbDriver) *followStore {
	return &followStore{db: db}
}

func (s *followStore) Create(srcAccountId, dstAccountId string) (err error) {
	_, err = s.db.Exec("insert into follows (src_account_id, dst_account_id) values (?, ?)", srcAccountId, dstAccountId)
	return
}

func (s *followStore) DeleteAllByAccountId(u entity.User) (err error) {
	table := toUserTable(u)
	_, err = s.db.Exec("delete from follows where src_account_id = ? or dst_account_id = ?", table.AccountId, table.AccountId)
	return
}
