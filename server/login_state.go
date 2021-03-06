package server

type LoginState int

const (
	Guest               = iota // ゲストユーザ
	LoginButNotRsrcUser        // ログインユーザとリソースページのユーザが一致しない
	LoginAndRsrcUser           // 一致する
	RsrcNotFound               // 該当するリソースユーザは存在しない
)
