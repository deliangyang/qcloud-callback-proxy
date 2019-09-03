package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/deliangyang/qcloud-callback-proxy/internal"
	"github.com/sirupsen/logrus"
)

var (
	config = internal.Config{}
)

type Info struct {
	ToAccount string `json:"To_Account"`
}

type State struct {
	CallbackCommand string `json:"CallbackCommand"`
	Info            Info   `json:"Info"`
}

func TestPost(t *testing.T) {
	internal.Parse("../configs/proxy.toml")
	config = internal.GetConfig()

	im := &State{
		CallbackCommand: "Group.CallbackAfterNewMemberJoin",
	}
	if message, err := json.Marshal(im); err == nil {
		request(string(message), t)
	} else {
		t.Fail()
	}

	im = &State{
		CallbackCommand: "State.StateChange",
		Info: Info{
			ToAccount: "20190000",
		},
	}
	if message, err := json.Marshal(im); err == nil {
		request(string(message), t)
	} else {
		t.Fail()
	}
}

func request(content string, t *testing.T) {
	url := fmt.Sprintf("http://127.0.0.1%s%s", config.Port, config.URI)
	resp, err := http.Post(url, "json", strings.NewReader(content))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	logrus.Info("body:", string(body))
	if string(body) != "success" {
		t.Fail()
	}
}
