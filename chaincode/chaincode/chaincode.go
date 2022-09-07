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

/*
추가 내용
참여자 []string
검증 []struct {
	uid stirng
	status string
	url string
}
*/

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
	// uuid := getUUID()
	var tmpVerificationData VerificationData
	tmpMapParticipant := make(map[string]string)
	tmpMapVerification := make(map[string]VerificationData)

	uuid := "TEST-uuid"
	tmpVerificationData.Uid = "TEST"
	tmpVerificationData.Status = "TEST"
	tmpVerificationData.Url = "TEST"
	tmpMapParticipant["TEST"] = "TEST"
	tmpMapVerification["TEST"] = tmpVerificationData
	quests := []Quest{
		{
			Id: uuid, 
			Title: "Genesis_Quest", 
			Content: "Genesis_Content", 
			Deadline: "2100-12-31 00:00:00", 
			Creator: "c09863c2-1ef8-11ed-84df-9c5c8ed2592b", 
			TokenAmount: 10000, 
			MaxParticipants: 10, 
			Status: "W",
			Participant: tmpMapParticipant,
			Verification: tmpMapVerification,
		},
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

func (s *SmartContract) UpdateQuest(ctx contractapi.TransactionContextInterface,
	id string, title string, content string, deadline string, creator string, tokenAmount int, maxParticipants int, status string, partiParticipant Participant, verification Verification) error {
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
		Participant:     partiParticipant,
		Verification:    verification,
	}
	questJSON, err := json.Marshal(quest)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, questJSON)
}

func (s *SmartContract) UpdateQuestInfo(ctx contractapi.TransactionContextInterface,
	id string, title string, content string, deadline string, creator string, tokenAmount int, maxParticipants int, status string) error {

	questJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if questJSON == nil {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	var quest Quest
	err = json.Unmarshal(questJSON, &quest)
	if err != nil {
		return err
	}

	var titleArg string
	var contentArg string
	var deadlineArg string
	var creatorArg string
	var tokenAmountArg int
	var maxParticipantsArg int
	var statusArg string

	if title == "" {
		titleArg = quest.Title
	} else {
		titleArg = title
	}

	if content == "" {
		contentArg = quest.Content
	} else {
		contentArg = content
	}

	if deadline == "" {
		deadlineArg = quest.Deadline
	} else {
		deadlineArg = deadline
	}

	if creator == "" {
		creatorArg = quest.Creator
	} else {
		creatorArg = creator
	}

	if tokenAmount == 0 {
		tokenAmountArg = quest.TokenAmount
	} else {
		tokenAmountArg = tokenAmount
	}

	if maxParticipants == 0 {
		maxParticipantsArg = quest.MaxParticipants
	} else {
		maxParticipantsArg = maxParticipants
	}

	if status == "" {
		statusArg = quest.Status
	} else {
		statusArg = status
	}

	// overwriting original asset with new asset
	questModified := Quest{
		Id:              id,
		Title:           titleArg,
		Content:         contentArg,
		Deadline:        deadlineArg,
		Creator:         creatorArg,
		TokenAmount:     tokenAmountArg,
		MaxParticipants: maxParticipantsArg,
		Status:          statusArg,
		Participant:     quest.Participant,
		Verification:    quest.Verification,
	}

	questJSON, err = json.Marshal(questModified)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, questJSON)
}

func (s *SmartContract) AddParticipantQuest(ctx contractapi.TransactionContextInterface, id string, participant string) error {
	exists, err := s.QuestExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the quest %s does not exist", id)
	}

	questJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if questJSON == nil {
		return fmt.Errorf("the quest %s does not exist", id)
	}

	var quest Quest
	err = json.Unmarshal(questJSON, &quest)
	if err != nil {
		return err
	}

	// participant 추가
	// 현재 퀘스트에 이미 등록되어있는지 확인
	if quest.Participant[participant] != "" {
		return fmt.Errorf("%s is already in quest", id)
	}

	quest.Participant[participant] = participant

	questJSON, err = json.Marshal(quest)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, questJSON)
}

func (s *SmartContract) QuitParticipantQuest(ctx contractapi.TransactionContextInterface, id string, participant string) error {
	exists, err := s.QuestExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	questJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if questJSON == nil {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	var quest Quest
	err = json.Unmarshal(questJSON, &quest)
	if err != nil {
		return err
	}

	// participant 제거
	if quest.Participant[participant] == "" {
		return fmt.Errorf("%s is not in quest", id)
	}

	delete(quest.Participant, participant)

	questJSON, err = json.Marshal(quest)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, questJSON)
}

func (s *SmartContract) AddVerificationQuest(ctx contractapi.TransactionContextInterface, id string, uid string, url string) error {
	exists, err := s.QuestExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the quest %s does not exist", id)
	}
	// 해당 퀘스트 존재 여부 및 해당 유저가 participant 인지 확인 필요에 관한 에러 처리 필요.

	questJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if questJSON == nil {
		return fmt.Errorf("the quest %s does not exist", id)
	}

	var quest Quest
	err = json.Unmarshal(questJSON, &quest)
	if err != nil {
		return err
	}

	// verification 추가
	// 현재 퀘스트에 이미 등록되어있는지 확인
	if quest.Verification[uid].Uid != "" {
		return fmt.Errorf("%s is already in quest verification", id)
	}

	verification := VerificationData{}
	verification.Uid = uid
	verification.Status = "W"
	verification.Url = url

	quest.Verification[uid] = verification

	questJSON, err = json.Marshal(quest)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, questJSON)
}

func (s *SmartContract) JudgeVerificationQuest(ctx contractapi.TransactionContextInterface, id string, uid string, status string) error {
	exists, err := s.QuestExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the quest %s does not exist", id)
	}

	questJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if questJSON == nil {
		return fmt.Errorf("the quest %s does not exist", id)
	}

	var quest Quest
	err = json.Unmarshal(questJSON, &quest)
	if err != nil {
		return err
	}

	// verification 추가
	// 현재 퀘스트에 이미 등록되어있는지 확인
	if quest.Verification[uid].Uid == "" {
		return fmt.Errorf("%s is not in quest verification", uid)
	}

	verification := VerificationData{}
	verification.Uid = quest.Verification[uid].Uid
	verification.Status = status
	verification.Url = quest.Verification[uid].Url

	quest.Verification[uid] = verification

	questJSON, err = json.Marshal(quest)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(id, questJSON)
}

func (s *SmartContract) GetQuest(ctx contractapi.TransactionContextInterface, id string) (*Quest, error) {
	questJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if questJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var quest Quest
	err = json.Unmarshal(questJSON, &quest)
	if err != nil {
		return nil, err
	}

	return &quest, nil
}

func (s *SmartContract) GetCreatorQuest(ctx contractapi.TransactionContextInterface, uid string) ([]*Quest, error) {
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

		// creator == uid 인 퀘스트만 append 하기
		if quest.Creator == uid {
			quests = append(quests, &quest)
		}
	}

	return quests, nil
}

func (s *SmartContract) GetParticipantQuest(ctx contractapi.TransactionContextInterface, uid string) ([]*Quest, error) {
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

		// participant[uid] 이 존재하는 경우만 append 하기
		if quest.Participant[uid] != "" {
			quests = append(quests, &quest)
		}
	}

	return quests, nil
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

func (s *SmartContract) DeleteQuest(ctx contractapi.TransactionContextInterface, uid string, qid string) error {
	questJSON, err := ctx.GetStub().GetState(qid)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if questJSON == nil {
		return fmt.Errorf("the asset %s does not exist", qid)
	}

	var quest Quest
	err = json.Unmarshal(questJSON, &quest)
	if err != nil {
		return err
	}

	if quest.Creator != uid {
		return fmt.Errorf("this asset is not belongs to %s", uid)
	}

	return ctx.GetStub().DelState(qid)
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
