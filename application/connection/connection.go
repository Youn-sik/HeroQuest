package connection

import (
	"database/sql"
	"log"
	"questAPP/database"
	"questAPP/sdk"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

var mysql *sql.DB
var contract *gateway.Contract

func init() {
	var err error
	err, mysql = database.NewMysqlConnection()
	if err != nil {
		log.Println("Fail to create connection of mysql", err)
	}
}

func GetMysqlClient() *sql.DB {
	return mysql
}
func GetContractClient() *gateway.Contract {
	err, contract := sdk.GetConnection()
	if err != nil {
		log.Println("Fail to create connection of contract", err)
	}
	return contract
}
