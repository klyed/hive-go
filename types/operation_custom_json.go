package types

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"fmt"
	"strings"

	"github.com/klyed/hive-go/encoding/transaction"
	"github.com/pkg/errors"
)

var (
	TypeFollow         = "follow"
	TypeReblog         = "reblog"
	TypeLogin          = "login"
	TypePrivateMessage = "sm_"
	TypeSplinterLand   = "SplinterLandOperation"
	TypeHiveEngine     = "ssc-mainnet-hive"
	TypeHiveSmartChain = "hsc"
)

// FC_REFLECT( hive::chain::custom_json_operation,
//             (required_auths)
//             (required_posting_auths)
//             (id)
//             (json) )

var customJSONDataObjects = map[string]interface{}{
	TypeFollow:         &FollowOperation{},
	TypeReblog:         &ReblogOperation{},
	TypeLogin:          &LoginOperation{},
	TypePrivateMessage: &PrivateMessageOperation{},
	TypeSplinterLand:   &SplinterLandOperation{},
	TypeHiveEngine:     &HiveEngineOperation{},
	TypeHiveSmartChain: &HiveSmartChainOperation{},
}

//CustomJSONOperation represents custom_json operation data.
type CustomJSONOperation struct {
	RequiredAuths        []string `json:"required_auths"`
	RequiredPostingAuths []string `json:"required_posting_auths"`
	ID                   string   `json:"id"`
	JSON                 string   `json:"json"`
}

//FollowOperation the structure for the operation CustomJSONOperation.
type FollowOperation struct {
	Follower  string   `json:"follower"`
	Following string   `json:"following"`
	What      []string `json:"what"`
}

//ReblogOperation the structure for the operation CustomJSONOperation.
type ReblogOperation struct {
	Account  string `json:"account"`
	Author   string `json:"author"`
	Permlink string `json:"permlink"`
}

//LoginOperation the structure for the operation CustomJSONOperation.
type LoginOperation struct {
	Account string `json:"account"`
}

//PrivateMessageOperation the structure for the operation CustomJSONOperation.
type PrivateMessageOperation struct {
	From             string `json:"from"`
	To               string `json:"to"`
	FromMemoKey      string `json:"from_memo_key"`
	ToMemoKey        string `json:"to_memo_key"`
	SentTime         uint64 `json:"sent_time"`
	Checksum         uint32 `json:"checksum"`
	EncryptedMessage string `json:"encrypted_message"`
}

//SplinterLandOperation the structure for the operation CustomJSONOperation.
type SplinterLandOperation struct {
	Follower  string   `json:"follower"`
	Following string   `json:"following"`
	What      []string `json:"what"`
}

//"{"contractName":"tokens","contractAction":"transfer","contractPayload":{"symbol":"DHEDGE","to":"lvr-docudrama","quantity":"0.00011323","memo":"Your daily ARCHON / ARCHONM GP drip for DHEDGE based on 0.002455895067637721 GP with 0.0024559 % share"}}"
//HiveEngineOperation the structure for the operation CustomJSONOperation.
type HiveEngineOperation struct {
	Contract  string   `json:"contractName"`
	Action    string   `json:"contractAction"`
	Payload   []string `json:"contractPalyload"`
}

//HiveSmartChainOperation the structure for the operation CustomJSONOperation.
type HiveSmartChainOperation struct {
	Contract  string   `json:"contractName"`
	Action    string   `json:"contractAction"`
	Payload   []string `json:"contractPalyload"`
}

//Type function that defines the type of operation.
func (op *CustomJSONOperation) Type() OpType {
	return TypeCustomJSON
}

//Data returns the operation data.
func (op *CustomJSONOperation) Data() interface{} {
	return op
}

//UnmarshalData unpacking the JSON parameter in the CustomJSONOperation type.
func (op *CustomJSONOperation) UnmarshalData() (interface{}, error) {
	// Get the corresponding data object template.

	template, ok := customJSONDataObjects[op.ID]
	if !ok {
		// In case there is no corresponding template, return nil.
		return nil, nil
	}

	// Clone the template.
	opData := reflect.New(reflect.Indirect(reflect.ValueOf(template)).Type()).Interface()

	// Prepare the whole operation tuple.
	var bodyReader io.Reader
	if op.JSON[0] == '[' {
		rawTuple := make([]json.RawMessage, 2)
		if err := json.NewDecoder(strings.NewReader(op.JSON)).Decode(&rawTuple); err != nil {
			return nil, errors.Wrapf(err, "failed to unmarshal CustomJSONOperation.JSON: \n%v", op.JSON)
		}
		if len(rawTuple) < 2 || rawTuple[1] == nil {
			return nil, errors.Errorf("invalid CustomJSONOperation.JSON: \n%v", op.JSON)
		}
		bodyReader = bytes.NewReader([]byte(rawTuple[1]))
	} else {
		bodyReader = strings.NewReader(op.JSON)
	}

	// Unmarshal into the new object instance.
	if err := json.NewDecoder(bodyReader).Decode(opData); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal CustomJSONOperation.JSON: \n%v", op.JSON)
	}

	return opData, nil
}

//MarshalTransaction is a function of converting type CustomJSONOperation to bytes.
func (op *CustomJSONOperation) MarshalTransaction(encoder *transaction.Encoder) error {
	enc := transaction.NewRollingEncoder(encoder)
	enc.EncodeUVarint(uint64(TypeCustomJSON.Code()))
	enc.EncodeArrString(op.RequiredAuths)
	enc.EncodeArrString(op.RequiredPostingAuths)
	enc.Encode(op.ID)
	enc.Encode(op.JSON)
	return enc.Err()
}

//MarshalCustomJSON generate a row from the structure fields.
func MarshalCustomJSON(v interface{}) (string, error) {
	var tmp []interface{}
	//res1 := strings.ToLower(str1)
	typeInterface := reflect.TypeOf(v).Name()
	smcheck :=  typeInterface[0:1]
	switch typeInterface {
	case "FollowOperation":
		tmp = append(tmp, TypeFollow)
	case "ReblogOperation":
		tmp = append(tmp, TypeReblog)
	case "LoginOperation":
		tmp = append(tmp, TypeLogin)
	case "PrivateMessageOperation":
		tmp = append(tmp, TypePrivateMessage)
	case "HiveEngineOperation":
		tmp = append(tmp, TypeHiveEngine)
	case "HiveSmartChainOperation":
		tmp = append(tmp,TypeHiveSmartChain)
	default:
		return "", errors.New("Unknown type")
	}
	if smcheck == "sm" {
		switch smcheck {
		case "SplinterLandOperation":
		  tmp = append(tmp, TypeSplinterLand)
		}
	}

	tmp = append(tmp, v)

	b, err := json.Marshal(tmp)
	if err != nil {
		return "", err
	}

	return string(b), nil //strings.Replace(string(b), "\"", "\\\"", -1), nil
}
/*
package main

import (
	"fmt"
	"net/url"
	"strings"
)

var hostnameWhitelist = map[string]struct{}{"*sm_*": struct{}{}, "sm_": struct{}{}}

func SplinterLandsTypePeek() {

	url, _ := url.Parse("sm_find_match")

	fmt.Println("custom_json id: " + url.Hostname())

	split := strings.SplitAfterN(url.Hostname(), ".", 2)
	split[0] = "*"
	hostName := strings.Join(split, ".")

	fmt.Println("Converted Hostname : ",hostName)

	if _, ok := hostnameWhitelist[hostName]; ok {
		fmt.Println("valid domain, allow access")
	} else {
		fmt.Println("NOT valid domain")
	}
}
*/
