package module

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/thedevsaddam/gojsonq/v2"
)

func Loki_api(url string, logql string) *interface{} {
	req, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(logql))
	if err != nil {
		requestLogger := logrus.WithFields(logrus.Fields{"err": err, "url": url, "logql": logql})
		requestLogger.Fatal("Client conn error !")
	}

	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		requestLogger := logrus.WithFields(logrus.Fields{"err": err, "url": url, "logql": logql})
		requestLogger.Error("Client body parse error !")
	}

	info := gojsonq.New().FromString(string(body)).Find("data.result")
	return &info
}
