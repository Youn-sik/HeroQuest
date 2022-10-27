package user

import (
	"database/sql"
	"log"
	"net/http"
	con "questAPP/connection"
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
mysql> desc user;
+---------------+--------------+------+-----+---------+-------+
| Field         | Type         | Null | Key | Default | Extra |
+---------------+--------------+------+-----+---------+-------+
| id            | varchar(255) | NO   | PRI | NULL    |       |
| account       | varchar(255) | NO   |     | NULL    |       |
| password      | varchar(255) | NO   |     | NULL    |       |
| name          | varchar(255) | NO   |     | NULL    |       |
| token_balance | int          | YES  |     | 0       |       |
| qid           | varchar(255) | YES  |     | NULL    |       |
+---------------+--------------+------+-----+---------+-------+
6 rows in set (0.02 sec)
*/

func Create(c *gin.Context) {
	reqData := CreateUserReq{}

	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Body Parsing Error"})
		return
	} else if reqData.Account == "" || reqData.Password == "" || reqData.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Wrong Body Form"})
		return
	}

	conn := con.GetMysqlClient()

	rows, err := conn.Query("select id from user where account = ?", reqData.Account)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}
	for rows.Next() {
		var uid string
		err = rows.Scan(&uid)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Parsing Error"})
			return
		}
		if uid != "" {
			c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Account is already exsist"})
			return
		}
	}

	_, err = conn.Query("insert into user (id, account, password, name) "+
		"value (?,?,?,?)", getUUID(), reqData.Account, reqData.Password, reqData.Name)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}

func Modify(c *gin.Context) {
	reqData := ModifyUserReq{}
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
	} else if reqData.Account == "" || reqData.Name == "" || reqData.Password == "" || reqData.QId == "" || reqData.TokenBalance == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Wrong Body Form"})
		return
	}

	conn := con.GetMysqlClient()

	_, err = conn.Query("update user set account = ?, password = ?, name = ?, token_balance = ?, qid = ? where id = ?",
		reqData.Account, reqData.Password, reqData.Name, reqData.TokenBalance, reqData.QId, uid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}

func Delete(c *gin.Context) {
	result, errStr, uid := middleware.GetIdFromToken(c.GetHeader("Authorization"))
	if !result {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": errStr})
		return
	}

	conn := con.GetMysqlClient()

	_, err := conn.Query("delete from user where id = ?", uid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}

func ListAll(c *gin.Context) {
	var userArr []User

	conn := con.GetMysqlClient()

	rows, err := conn.Query("select * from user")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	for rows.Next() {
		user := User{}

		err := rows.Scan(&user.Id, &user.Account, &user.Password, &user.Name, &user.TokenBalance, &user.QId)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Parsing Error"})
			return
		}
		userArr = append(userArr, user)
	}

	c.JSON(http.StatusOK, gin.H{"result": true, "userArr": userArr})
}

func Info(c *gin.Context) {
	user := User{}
	result, errStr, uid := middleware.GetIdFromToken(c.GetHeader("Authorization"))
	if !result {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": errStr})
		return
	}

	conn := con.GetMysqlClient()

	rows, err := conn.Query("select * from user where id = ?", uid)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Account, &user.Password, &user.Name, &user.TokenBalance, &user.QId)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Parsing Error"})
			return
		}
	}

	if user.Id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "No User Data in Database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true, "userInfo": user})
}

func Login(c *gin.Context) {
	var qid sql.NullString
	user := User{}
	reqData := LoginReq{}

	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Body Parsing Error"})
		return
	}

	conn := con.GetMysqlClient()

	rows, err := conn.Query("select * from user where account = ?", reqData.Account)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Query Parsing Error"})
		return
	}

	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Account, &user.Password, &user.Name, &user.TokenBalance, &qid)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Parsing Error"})
			return
		}
	}

	if user.Id != "" && (user.Password == reqData.Password) {
		token := middleware.TokenBuild(user.Id)

		if token == "false" {
			c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Error Occurred With Creating Token"})
			return
		}
		if qid.Valid {
			user.QId = ""
		} else {
			user.QId = qid.String
		}

		c.JSON(http.StatusOK, gin.H{"result": true, "userInfo": user, "token": token})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "The Account Does Not Exist or The Password Is Incorrect"})
}
