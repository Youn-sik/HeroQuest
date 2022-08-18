package chaincode

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Quest struct {
	Id              string `json:"id"`
	Title           string `json:"title"`
	Content         string `json:"content"`
	Deadline        string `json:"deadline"`
	Creator         string `json:"creator"`
	TokenAmount     int    `json:"token_amount"`
	MaxParticipants int    `json:"max_participants"`
	Status          string `json:"status"`
}

/*
mysql> desc quest;
	+----+--------+--------------------+---------------------+---------+--------------+------------------+--------+
	| id | title  | content            | deadline            | creator | token_amount | max_participants | status |
	+----+--------+--------------------+---------------------+---------+--------------+------------------+--------+
	|  1 | quest1 | this is test quest | 2022-08-28 00:00:00 |       1 |         1000 |               10 | N      |
	+----+--------+--------------------+---------------------+---------+--------------+------------------+--------+
*/

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
6 rows in set (0.00 sec)

mysql> select * from user;
+--------------------------------------+---------+----------+------------+---------------+------+
| id                                   | account | password | name       | token_balance | qid  |
+--------------------------------------+---------+----------+------------+---------------+------+
| c09863c2-1ef8-11ed-84df-9c5c8ed2592b | cho     | 12345    | choyounsik |             0 | NULL |
+--------------------------------------+---------+----------+------------+---------------+------+
1 row in set (0.00 sec)

mysql> desc quest_verification;
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

func getUUID() string {
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Println(err)
	}
	uuidStr := uuid.String()
	return uuidStr
}

// InitLedger adds a base set of quest to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	uuid := getUUID()
	quests := []Quest{
		{Id: uuid, Title: "qtitle_value", Content: "qcontent_value", Deadline: "2022-08-28 00:00:00", Creator: "c09863c2-1ef8-11ed-84df-9c5c8ed2592b", TokenAmount: 1000, MaxParticipants: 10, Status: "N"},
	}

	for _, quest := range quests {
		questJSON, err := json.Marshal(quest)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(quest.Id, questJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

func (s *SmartContract) CreateQuest(ctx contractapi.TransactionContextInterface, title string, content string, deadline string, creator string, tokenAmount int, maxParticipants int, status string) error {
	id := getUUID()
	quest := Quest{
		Id:              id,
		Title:           title,
		Content:         content,
		Deadline:        deadline,
		Creator:         creator,
		TokenAmount:     tokenAmount,
		MaxParticipants: maxParticipants,
		Status:          status,
	}
	questJSON, err := json.Marshal(quest)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, questJSON)
}

// QuestExists returns true when quest with given ID exists in world state
func (s *SmartContract) QuestExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	questJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return questJSON != nil, nil
}

func (s *SmartContract) UpdateQuest(ctx contractapi.TransactionContextInterface, id string, title string, content string, deadline string, creator string, tokenAmount int, maxParticipants int, status string) error {
	exists, err := s.QuestExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	quest := Quest{
		Id:              id,
		Title:           title,
		Content:         content,
		Deadline:        deadline,
		Creator:         creator,
		TokenAmount:     tokenAmount,
		MaxParticipants: maxParticipants,
		Status:          status,
	}
	questJSON, err := json.Marshal(quest)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, questJSON)
}

func (s *SmartContract) GetQuest(ctx contractapi.TransactionContextInterface, id string) (*Quest, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var quest Quest
	err = json.Unmarshal(assetJSON, &quest)
	if err != nil {
		return nil, err
	}

	return &quest, nil
}

func (s *SmartContract) GetAllQuest(ctx contractapi.TransactionContextInterface) ([]*Quest, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var quests []*Quest
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var quest Quest
		err = json.Unmarshal(queryResponse.Value, &quest)
		if err != nil {
			return nil, err
		}
		quests = append(quests, &quest)
	}

	return quests, nil
}

func (s *SmartContract) DeleteQuest(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.QuestExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

/*
// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:             id,
		Color:          color,
		Size:           size,
		Owner:          owner,
		AppraisedValue: appraisedValue,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}

	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, color string, size int, owner string, appraisedValue int) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		ID:              id,
		Title:           title,
		Content:         content,
		Deadline:        deadline,
		Creator:         creator,
		TokenAmount:     tokenAmount,
		MaxParticipants: maxParticipants,
		Status:          status,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, newOwner string) error {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return err
	}

	asset.Owner = newOwner
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, assetJSON)
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}

	return assets, nil
}

*/
