package quest

import (
	"log"

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

// quest.Join(c)
// quest.Quit(c)
// quest.Verify(c)
// quest.Judge(c)
// quest.GetUsers(c)
// quest.GetCreatorVerifyList(c)
// quest.GetParticipantVerifyList(c)

// func Join(c *gin.Context) {
// 	requData := JoinQuestReq{}
// }
