package quest

type QuestVerification struct {
	Id     string `json:"id"`
	Qid    string `json:"qid"`
	Uid    string `json:"uid"`
	status string `json:"status"`
	url    string `json:"url"`
}
