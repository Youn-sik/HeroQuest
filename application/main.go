package main

import (
	"log"

	"questAPP/user"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, DELETE, POST")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()

	/*
		- 사용자 생성 함수
		- 사용자 수정 함수
		- 사용자 삭제 함수
		- 사용자 리스트 조회 함수
		- 사용자 정보 조회 함수
	*/
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

	// questRouter := router.Group("/quest")
	// questRouter.POST("/join", func(c *gin.Context) {
	// 	quest.Join(c)
	// })
	// questRouter.POST("/quit", func(c *gin.Context) {
	// 	quest.Quit(c)
	// })
	// questRouter.POST("/verify", func(c *gin.Context) {
	// 	quest.Verify(c)
	// })
	// questRouter.POST("/verify/judge", func(c *gin.Context) {
	// 	quest.Judge(c)
	// })
	// questRouter.GET("/users", func(c *gin.Context) {
	// 	quest.GetUsers(c)
	// })
	// questRouter.GET("/verify/list/creator", func(c *gin.Context) {
	// 	quest.GetCreatorVerifyList(c)
	// }) // 생성자: 본인이 검증할(mysql) or 검증한 퀘스트(block-chain) 리스트
	// questRouter.GET("/verify/list/participant", func(c *gin.Context) {
	// 	quest.GetParticipantVerifyList(c)
	// }) // 참여자: 본인이 검증받을(mysql) or 검증받은(block-chain) 퀘스트 리스트

	/*
		퀘스트에 참여한 사용자 리스트 조회 함수
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

	router := setupRouter()
	log.Println("[SERVER] => Backend Admin application is listening on port " + port)
	router.Run(port)
}
