package internal

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func init() {
	http.DefaultClient = &http.Client{
		Timeout: time.Second * 3,
	}
}

// Dispatch 分发请求
func Dispatch(w http.ResponseWriter, r *http.Request) {
	if err := dispatch(r); err != nil {
		Logger().WithError(err).Error("dispatch im callback")
		w.Write([]byte("error"))
		return
	}
	w.Write([]byte("success"))
}

func dispatch(r *http.Request) error {
	config := GetConfig()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	envs, err := ConvertCallbackCommandToEnvs(config, body)
	Logger().WithField("message", string(body)).Info("callback message")
	if err != nil {
		return err
	}

	var g errgroup.Group
	for _, env := range envs {
		env := env
		g.Go(func() error {
			Logger().WithFields(logrus.Fields{
				"EnvID":    env.ID,
				"Callback": env.URL,
				"Length":   env.Length,
			}).Info("env found")
			result, err := request(env.URL+config.URI, string(body))
			if err != nil {
				return err
			} else if string(result) == "error" {
				return errors.New("reserve request")
			}
			return nil
		})
	}
	return g.Wait()
}

// request 转发请求
func request(url string, content string) ([]byte, error) {
	resp, err := http.Post(url, "application/json", strings.NewReader(content))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
