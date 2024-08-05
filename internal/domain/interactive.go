package domain

type Interactive struct {
	Biz   string `json:"biz"`
	BizId string `json:"biz_id"`

	ReadCnt    int64 `json:"read_cnt"`
	LikedCnt   int64 `json:"liked_cnt"`
	CollectCnt int64 `json:"collect_cnt"`

	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}

type Self struct {
	Liked     bool `json:"liked"`
	Collected bool `json:"collected"`
}

type Collection struct {
	Name  string
	Uid   int64
	Items []Resource
}
