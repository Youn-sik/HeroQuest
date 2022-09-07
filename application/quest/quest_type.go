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
}

type QuitQuestReq struct {
	Qid    string `json:"qid"`
}

type VerifyQuestReq struct {
	Qid string `json:"qid"`
	Url string `json:"url"`
}

type JudgeQuestReq struct {
	Qid    string `json:"qid"`
	Uid    string `json:"uid"`
	Status string `json:"status"`
}

type GetUsersReq struct {
	Qid string `json:"qid"`
}

type QuestUser struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetQuestInfoReq struct {
	Qid string `json:"qid"`
}

type CreateQuestReq struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Deadline    string `json:"deadline"`
	TokenAmount string `json:"token_amount"`
}

type ModifyQuestReq struct {
	Qid         string `json:"qid,omitempty"`
	Title       string `json:"title,omitempty"`
	Content     string `json:"content,omitempty"`
	Deadline    string `json:"deadline,omitempty"`
	TokenAmount string `json:"token_amount,omitempty"`
}

type DeleteQuestReq struct {
	Qid string `json:"qid"`
}
