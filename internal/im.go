package internal

import (
	"fmt"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// IM callback request's body
type IM struct {
	CallbackCommand string `json:"CallbackCommand"`
	FromAccount     string `json:"From_Account"`
}

// https://cloud.tencent.com/document/product/269/1523
var (
	StateStateChange                   = "State.StateChange"
	SnsCallbackFriendAdd               = "Sns.CallbackFriendAdd"
	SnsCallbackFriendDelete            = "Sns.CallbackFriendDelete"
	SnsCallbackBlackListAdd            = "Sns.CallbackBlackListAdd"
	SnsCallbackBlackListDelete         = "Sns.CallbackBlackListDelete"
	C2CCallbackBeforeSendMsg           = "C2C.CallbackBeforeSendMsg"
	C2CCallbackAfterSendMsg            = "C2C.CallbackAfterSendMsg"
	GroupCallbackBeforeCreateGroup     = "Group.CallbackBeforeCreateGroup"
	GroupCallbackAfterCreateGroup      = "Group.CallbackAfterCreateGroup"
	GroupCallbackBeforeApplyJoinGroup  = "Group.CallbackBeforeApplyJoinGroup"
	GroupCallbackBeforeInviteJoinGroup = "Group.CallbackBeforeInviteJoinGroup"
	GroupCallbackAfterNewMemberJoin    = "Group.CallbackAfterNewMemberJoin"
	GroupCallbackAfterMemberExit       = "Group.CallbackAfterMemberExit"
	GroupCallbackBeforeSendMsg         = "Group.CallbackBeforeSendMsg"
	GroupCallbackAfterSendMsg          = "Group.CallbackAfterSendMsg"
	GroupCallbackAfterGroupFull        = "Group.CallbackAfterGroupFull"
	GroupCallbackAfterGroupDestroyed   = "Group.CallbackAfterGroupDestroyed"
	GroupCallbackAfterGroupInfoChanged = "Group.CallbackAfterGroupInfoChanged"
)

const defaultDispatchToAll = "all"

// Command 解析
type Command struct {
	Commands []string
	Parse    func(message []byte) string
}

var events = map[string]Command{
	"State": {
		Commands: []string{
			StateStateChange,
		},
		Parse: func(message []byte) string {
			return jsoniter.Get(message, "Info", "To_Account").ToString()
		},
	},
	"Sns": {
		Commands: []string{
			SnsCallbackFriendAdd,
			SnsCallbackFriendDelete,
			SnsCallbackBlackListAdd,
			SnsCallbackBlackListDelete,
		},
		Parse: func(message []byte) string {
			return jsoniter.Get(message, "PairList", 0).Get("From_Account").ToString()
		},
	},
	"C2C": {
		Commands: []string{
			C2CCallbackBeforeSendMsg,
			C2CCallbackAfterSendMsg,
		},
		Parse: func(message []byte) string {
			return jsoniter.Get(message, "From_Account").ToString()
		},
	},
	"Group": {
		Commands: []string{
			GroupCallbackBeforeCreateGroup,
			GroupCallbackAfterCreateGroup,
			GroupCallbackBeforeApplyJoinGroup,
			GroupCallbackBeforeInviteJoinGroup,
			GroupCallbackAfterNewMemberJoin,
			GroupCallbackAfterMemberExit,
			GroupCallbackBeforeSendMsg,
			GroupCallbackAfterSendMsg,
			GroupCallbackAfterGroupFull,
			GroupCallbackAfterGroupDestroyed,
			GroupCallbackAfterGroupInfoChanged,
		},
		Parse: func(message []byte) string {
			Account := jsoniter.Get(message, "Operator_Account").ToString()
			if len(Account) > 0 && strings.ToLower(Account) != "admin" {
				return Account
			}
			Account = jsoniter.Get(message, "Owner_Account").ToString()
			if len(Account) > 0 && strings.ToLower(Account) != "admin" {
				return Account
			}
			// 只有groupID
			return defaultDispatchToAll
		},
	},
}

// ConvertCallbackCommandToEnvs 获取envs
func ConvertCallbackCommandToEnvs(config Config, message []byte) (ENVS, error) {
	event := jsoniter.Get(message, "CallbackCommand").ToString()
	Logger().WithField("event", event).Info("parse callback command")

	items := strings.Split(event, ".")
	if v, ok := events[items[0]]; ok {
		for _, e := range v.Commands {
			if e == event {
				result := v.Parse(message)
				if result == defaultDispatchToAll {
					return config.ENVS, nil
				}
				envs, err := config.FindEnvs(result)
				if err != nil {
					return nil, err
				}
				return envs, nil
			}
		}
	}

	return nil, fmt.Errorf("callback command event not found")
}
