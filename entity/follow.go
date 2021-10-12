package entity

type Follow struct {
	Id int64
	// フォロー元ユーザID
	SrcAccountId string
	// フォロー先ユーザID
	DstAccountId string
}
