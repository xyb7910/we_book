package domain

type Interactive struct {
	Biz   string `json:"biz"`
	BizId int64  `json:"biz_id"`

	ReadCnt    int64 `json:"read_cnt"`
	LikedCnt   int64 `json:"liked_cnt"`
	CollectCnt int64 `json:"collect_cnt"`

	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}
