package quest

/*
3. 유저 퀘스트 수락 및 진행
    - 퀘스트 참여 함수
    - 퀘스트 포기 함수
4. 퀘스트 완료 검증
    - 퀘스트 검증 생성 함수
    - 퀘스트 검증 조회 함수
5. 퀘스트 보상
    - 토큰 추가 함수
    - 토큰 차감 함수
*/

import (
	"log"
	"net/http"
	"questAPP/database"
	"questAPP/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func getUUID() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	uuidStr := uuid.String()
	return uuidStr
}

/*
quest.Join(c)
quest.Quit(c)
quest.Verify(c)
quest.Judge(c)
quest.GetUsers(c)
quest.GetCreatorVerifyList(c)
quest.GetParticipantVerifyList(c)
*/

// ChainCode SDK
func Join(c *gin.Context) {
	reqData := JoinQuestReq{}
	result, errStr, uid := middleware.GetIdFromToken(c.GetHeader("Authorization"))
	if !result {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": errStr})
		return
	}

	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Body Parsing Error"})
		return
	}

	conn := database.NewMysqlConnection()
	defer conn.Close()

	_, err = conn.Query("update user set qid = ? where uid = ?", reqData.Qid, uid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	// SDK에 qid 값으로 quest 내용 select 후 참여자 정보 추가
	// AddParticipantQuest(qid, uid) 함수 사용

	c.JSON(http.StatusOK, gin.H{"result": true})
}

// ChainCode SDK
func Quit(c *gin.Context) {
	reqData := QuitQuestReq{}
	result, errStr, uid := middleware.GetIdFromToken(c.GetHeader("Authorization"))
	if !result {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": errStr})
		return
	}

	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Body Parsing Error"})
		return
	}

	conn := database.NewMysqlConnection()
	defer conn.Close()

	_, err = conn.Query("update user set qid = ? where uid = ? and qid = ?", nil, uid, reqData.Qid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	rows, err := conn.Query("select * from quest_verification where uid = ? and qid = ?", uid, reqData.Qid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}
	if rows != nil {
		conn.Query("update quest_verification set status = ? where uid = ? qid = ?", reqData.Status, uid, reqData.Qid)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
			return
		}
	}

	// SDK에 qid 값으로 quest 내용 select 후 참여자 정보 삭제
	// QuitParticipantQuest(qid, uid) 함수 사용

	c.JSON(http.StatusOK, gin.H{"result": true})
}

// ChainCode SDK
func Verify(c *gin.Context) {
	reqData := VerifyQuestReq{}
	result, errStr, uid := middleware.GetIdFromToken(c.GetHeader("Authorization"))
	if !result {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": errStr})
		return
	}

	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Body Parsing Error"})
		return
	}
	if reqData.Qid == "" || reqData.Url == "" {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Body Value Was Missing"})
		return
	}

	conn := database.NewMysqlConnection()
	defer conn.Close()

	_, err = conn.Query("insert into quest_verification values(?,?,?,?,?)", getUUID(), reqData.Qid, uid, "W", reqData.Url)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	// SDK에 qid 값으로 quest 내용 select 후 검증자 정보 추가 (url 및 uid 값 등록)
	// AddVerificationQuest(qid, uid, url) 함수 사용

	c.JSON(http.StatusBadRequest, gin.H{"result": true})
}

// ChainCode SDK
func Judge(c *gin.Context) {
	reqData := JudgeQuestReq{}
	result, errStr, uid := middleware.GetIdFromToken(c.GetHeader("Authorization"))
	if !result {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": errStr})
		return
	}

	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Body Parsing Error"})
		return
	}

	conn := database.NewMysqlConnection()
	defer conn.Close()

	_, err = conn.Query("update quest_verification set status = ? where uid = ? qid = ?", reqData.Status, uid, reqData.Qid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	// SDK에 qid 값으로 quest 내용 select 후 검증자 정보 추가 (status 업데이트)
	// JudgeVerificationQuest(uid, status) 함수 사용

	c.JSON(http.StatusBadRequest, gin.H{"result": true})
}

func GetUsers(c *gin.Context) {
	qid := c.Query("qid")
	var userArr []QuestUser

	conn := database.NewMysqlConnection()
	defer conn.Close()

	rows, err := conn.Query("select * from user where qid = ?", qid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	for rows.Next() {
		user := QuestUser{}

		err := rows.Scan(&user)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Parsing Error"})
			return
		}
		userArr = append(userArr, user)
	}

	c.JSON(http.StatusOK, gin.H{"result": true, "userArr": userArr})
}

// ChainCode SDK
func GetCreatorVerifyList(c *gin.Context) {
	result, errStr, uid := middleware.GetIdFromToken(c.GetHeader("Authorization"))
	if !result {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": errStr})
		return
	}

	// SDK 에 uid 값으로 퀘스트 리스트 요청
	// GetCreatorQuest(uid) 함수 사용
}

// ChainCode SDK
func GetParticipantVerifyList(c *gin.Context) {
	// quest := QuestVerification{}
	// reqData :=

	// SDK 에 uid 값으로 퀘스트 리스트 요청
	// GetParticipantQuest(uid) 함수 사용
}

func GetQuestList(c *gin.Context) {
	var questArr []Quest

	// SDK에 퀘스트 리스트 요청
	// GetAllQuest() 함수 사용
}

func GetQuestInfo(c *gin.Context) {
	reqData := GetQuestInfoReq{}

	// SDK qid로 퀘스트 정보 요청
	// GetQuest() 함수 사용
}

func CreateQuest(c *gin.Context) {
	reqData := CreateQuestReq{}
	result, errStr, uid := middleware.GetIdFromToken(c.GetHeader("Authorization"))
	if !result {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": errStr})
		return
	}

	// SDK uid및 quest property들로 퀘스트 생성 요청
	// CreateQuest(title, content, deadline, uid(creator), TokenAmount) 함수 사용
}

func ModifyQuest(c *gin.Context) {
	reqData := ModifyQuestReq{}
	result, errStr, uid := middleware.GetIdFromToken(c.GetHeader("Authorization"))
	if !result {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": errStr})
		return
	}

	// SDK uid및 quest property들로 퀘스트 수정 요청
	// UpdateQuestInfo(qid, title, content, deadline, tokenAmount)
}

func DeleteQuest(c *gin.Context) {
	reqData := DeleteQuestReq{}
	result, errStr, uid := middleware.GetIdFromToken(c.GetHeader("Authorization"))
	if !result {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": errStr})
		return
	}

	// SDK qid와 uid로 퀘스트 삭제 요청
	// DeleteQuest(uid, qid)
}
