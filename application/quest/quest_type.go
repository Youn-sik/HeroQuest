package quest

type QuestVerification struct {
	Id     string `json:"id"`
	Qid    string `json:"qid"`
	Uid    string `json:"uid"`
	Status string `json:"status"`
	Url    string `json:"url"`
}

type JoinQuestReq struct {
	Qid string `json:"qid"`
	Uid string `json:"uid"`
}

type QuitQuestReq struct {
	Qid    string `json:"qid"`
	Uid    string `json:"uid"`
	Status string `json:"status"`
}

type VerifyQuestReq struct {
	Qid string `json:"qid"`
	Url string `json:"url"`
}

type JudgeQuestReq struct {
	Qid    string `json:"qid"`
	Status string `json:"status"`
}

type GetUsersReq struct {
	Qid string `json:"qid"`
}

type QuestUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
