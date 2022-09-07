package main

import (
	"database/sql"
	"log"

	con "questAPP/connection"
	m "questAPP/middleware"
	"questAPP/quest"
	"questAPP/user"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

var mysql *sql.DB
var contract *gateway.Contract

func setupRouter() *gin.Engine {

	router := gin.Default()
	router.Use(m.CORSMiddleware())
	/*
		- 사용자 생성 함수
		- 사용자 수정 함수
		- 사용자 삭제 함수
		- 사용자 리스트 조회 함수
		- 사용자 정보 조회 함수
	*/

	// 인증이 필요할 경우 해당 그룹 사용
	authRouter := router.Group("/auth")
	authRouter.Use(m.TokenAuthenticate)

	userRouter := router.Group("/user")
	userRouter.POST("/create", func(c *gin.Context) {
		user.Create(c)
	})
	userRouter.POST("/modify", func(c *gin.Context) {
		user.Modify(c)
	})
	userRouter.POST("/delete", func(c *gin.Context) {
		user.Delete(c)
	})
	userRouter.GET("/listall", func(c *gin.Context) {
		user.ListAll(c)
	})
	userRouter.GET("/info", func(c *gin.Context) {
		user.Info(c)
	})
	userRouter.POST("/login", func(c *gin.Context) {
		user.Login(c)
	})

	questRouter := router.Group("/quest")
	questRouter.Use(m.TokenAuthenticate)
	// questInfo, questCreate, questModify, questDelete 함수 구현 되어있음.
	questRouter.GET("/list", func(c *gin.Context) {
		quest.GetQuestList(c)
	})
	questRouter.POST("/join", func(c *gin.Context) {
		quest.Join(c)
	})
	questRouter.POST("/quit", func(c *gin.Context) {
		quest.Quit(c)
	})
	questRouter.POST("/verify", func(c *gin.Context) {
		quest.Verify(c)
	}) // 퀘스트 검증 요청 (퀘스트 완료 요청)
	questRouter.POST("/verify/judge", func(c *gin.Context) {
		quest.Judge(c)
	}) // 퀘스트 검증
	questRouter.GET("/users", func(c *gin.Context) {
		quest.GetUsers(c)
	}) // 퀘스트에 참여한 유저 리스트
	questRouter.GET("/verify/list/creator", func(c *gin.Context) {
		quest.GetCreatorVerifyList(c)
	}) // 생성자: 본인이 검증할(mysql) or 검증한 퀘스트(block-chain) 리스트
	questRouter.GET("/verify/list/participant", func(c *gin.Context) {
		quest.GetParticipantVerifyList(c)
	}) // 참여자: 본인이 검증받을(mysql) or 검증받은(block-chain) 퀘스트 리스트

	/*
		2. 퀘스트에 참여한 사용자 리스트 조회 함수
		3. 유저 퀘스트 수락 및 진행
			- 퀘스트 참여 함수
			- 퀘스트 포기 함수
		4. 퀘스트 완료 검증
			- 퀘스트 검증 생성 함수
			- 퀘스트 검증 수락(거절)
			- 퀘스트 검증 리스트 조회 함수
			- 퀘스트 검증 결과 조회 함수
		5. 퀘스트 보상
			- 토큰 추가 함수
			- 토큰 차감 함수
	*/

	return router
}

func main() {
	port := ":3000"

	mysql = con.GetMysqlClient()
	defer mysql.Close()

	router := setupRouter()
	log.Println("[SERVER] => Backend Admin application is listening on port " + port)
	router.Run(port)
}
