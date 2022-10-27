package quest

type Quest struct {
	Id              string       `json:"id"`
	Title           string       `json:"title"`
	Content         string       `json:"content"`
	Deadline        string       `json:"deadline"`
	Creator         string       `json:"creator"`
	TokenAmount     int          `json:"token_amount"`
	MaxParticipants int          `json:"max_participants"`
	Status          string       `json:"status"`
	Participant     Participant  `json:"participant,omitempty"`
	Verification    Verification `json:"verification,omitempty"`
}

type Participant map[string]string

type Verification map[string]VerificationData

type VerificationData struct {
	Uid    string `json:"uid,omitempty"`
	Status string `json:"status,omitempty"`
	Url    string `json:"url,omitempty"`
}

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
	Creator string `json:"creator"`
	TokenAmount string `json:"tokenAmount"`
	MaxParticipants string `json:"maxParticipants"`
	Status string `json:"status"`
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
