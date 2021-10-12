package server

import "github.com/guricerin/stop-now-smoking/entity"

type followTable struct {
	Id           int64
	SrcAccountId string
	DstAccountId string
}

func toFollowEntity(t followTable) entity.Follow {
	return entity.Follow{
		Id:           t.Id,
		SrcAccountId: t.SrcAccountId,
		DstAccountId: t.DstAccountId,
	}
}

type followStore struct {
	db DbDriver
}

func NewFollowStore(db DbDriver) *followStore {
	return &followStore{db: db}
}

func (s *followStore) Create(srcAccountId, dstAccountId string) (err error) {
	_, err = s.db.Exec("insert into follows (src_account_id, dst_account_id) values ($1, $2)", srcAccountId, dstAccountId)
	return
}

func (s *followStore) IsFollowing(srcAccountId, dstAccountId string) (bool, error) {
	rows, err := s.db.Query("select * from follows where src_account_id = $1 and dst_account_id = $2", srcAccountId, dstAccountId)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), nil
}

func (s *followStore) RetrieveFollows(srcAccountId string) ([]entity.Follow, error) {
	rows, err := s.db.Query("select id, src_account_id, dst_account_id from follows where src_account_id = $1", srcAccountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]entity.Follow, 0)
	for rows.Next() {
		t := followTable{}
		err := rows.Scan(&t.Id, &t.SrcAccountId, &t.DstAccountId)
		if err != nil {
			return nil, err
		}
		res = append(res, toFollowEntity(t))
	}
	return res, nil
}

func (s *followStore) RetrieveFollowers(dstAccountId string) ([]entity.Follow, error) {
	rows, err := s.db.Query("select id, src_account_id, dst_account_id from follows where dst_account_id = $1", dstAccountId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := make([]entity.Follow, 0)
	for rows.Next() {
		t := followTable{}
		err := rows.Scan(&t.Id, &t.SrcAccountId, &t.DstAccountId)
		if err != nil {
			return nil, err
		}
		res = append(res, toFollowEntity(t))
	}
	return res, nil
}

func (s *followStore) Delete(srcAccountId, dstAccountId string) (err error) {
	_, err = s.db.Exec("delete from follows where src_account_id = $1 and dst_account_id = $2", srcAccountId, dstAccountId)
	return
}

func (s *followStore) DeleteAllByAccountId(u entity.User) (err error) {
	table := toUserTable(u)
	_, err = s.db.Exec("delete from follows where src_account_id = $1 or dst_account_id = $2", table.AccountId, table.AccountId)
	return
}
