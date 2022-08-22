package user

import (
	"database/sql"
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
	}

	conn := database.NewMysqlConnection()
	defer conn.Close()

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

	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Body Parsing Error"})
		return
	}

	conn := database.NewMysqlConnection()
	defer conn.Close()

	_, err = conn.Query("update user set account = ?, password = ?, name = ?, token_balance = ?, qid = ?",
		reqData.Account, reqData.Password, reqData.Name, reqData.TokenBalance, reqData.QId)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}

func Delete(c *gin.Context) {
	reqData := DeleteUserReq{}

	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Body Parsing Error"})
		return
	}

	conn := database.NewMysqlConnection()
	defer conn.Close()

	_, err = conn.Query("delete from user where id = ?", reqData.Id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": true})
}

func ListAll(c *gin.Context) {
	var userArr []User

	conn := database.NewMysqlConnection()
	defer conn.Close()

	rows, err := conn.Query("select * from user")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	for rows.Next() {
		user := User{}

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

func Info(c *gin.Context) {
	user := User{}
	reqData := InfoUserReq{}

	err := c.Bind(&reqData)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Body Parsing Error"})
		return
	}

	conn := database.NewMysqlConnection()
	defer conn.Close()

	rows, err := conn.Query("select * from user where account = ?", reqData.Id)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Database Query Error"})
		return
	}

	for rows.Next() {
		err := rows.Scan(&user)
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

	conn := database.NewMysqlConnection()
	defer conn.Close()

	rows, err := conn.Query("select * from user where account = ?", reqData.Id)
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
	c.JSON(http.StatusBadRequest, gin.H{"result": false, "errStr": "Wrong Password"})
}
