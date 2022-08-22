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
func Join(c *gin.Context) {
	reqData := JoinQuestReq{}

	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Body Parsing Error"})
		return
	}

	conn := database.NewMysqlConnection()
	defer conn.Close()

	_, err = conn.Query("update user set qid = ? where uid = ?", reqData.Qid, reqData.Uid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}

func Quit(c *gin.Context) {
	reqData := QuitQuestReq{}

	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Body Parsing Error"})
		return
	}

	conn := database.NewMysqlConnection()
	defer conn.Close()

	_, err = conn.Query("update user set qid = ? where uid = ? and qid = ?", nil, reqData.Uid, reqData.Qid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	rows, err := conn.Query("select * from quest_verification where uid = ? and qid = ?", reqData.Uid, reqData.Qid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}
	if rows != nil {
		conn.Query("update quest_verification set status = ? where uid = ? qid = ?", reqData.Status, reqData.Uid, reqData.Qid)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}

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
	c.JSON(http.StatusBadRequest, gin.H{"result": true})
}

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

func GetCreatorVerifyList(c *gin.Context) {

}
func GetParticipantVerifyList(c *gin.Context) {
	// quest := QuestVerification{}
	// reqData :=
}

/*
mysql> desc quest_verification
+--------+--------------+------+-----+---------+-------+
| Field  | Type         | Null | Key | Default | Extra |
+--------+--------------+------+-----+---------+-------+
| id     | varchar(255) | NO   | PRI | NULL    |       |
| qid    | varchar(255) | NO   |     | NULL    |       |
| uid    | varchar(255) | NO   |     | NULL    |       |
| status | varchar(10)  | NO   |     | NULL    |       |
| url    | varchar(255) | NO   |     | NULL    |       |
+--------+--------------+------+-----+---------+-------+
*/
