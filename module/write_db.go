package module

import (
	"loki_client/config"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/sirupsen/logrus"
)

func Conninflux() client.Client {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.ReadConfig("influx.url").(string),      //数据库地址
		Username: config.ReadConfig("influx.user").(string),     //数据库用户名
		Password: config.ReadConfig("influx.password").(string), //数据库密码
	})
	if err != nil {
		logrus.Error("Error creating InfluDB Client: ", err)
	}

	return cli
}

func Writeinflux(cli client.Client, modules string, fields map[string]interface{}) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  config.ReadConfig("influx.db").(string),        //数据库名称
		Precision: config.ReadConfig("influx.precision").(string), //时间精度（很重要，不然循环写入会覆盖之前的数据，influxdb是以时间戳为单位）
	})
	if err != nil {
		logrus.Error("Connection influxdb fail :", err)
	}

	// 写入metrics表
	tags := map[string]string{
		"table": modules,
	}

	server, err := client.NewPoint(modules, tags, fields, time.Now()) //并插入对应字段和tag，如果表不存在自动创建
	if err != nil {
		logrus.Error("Create table fail: ", err)
	}
	bp.AddPoint(server)
	err = cli.Write(bp)
	if err != nil {
		logrus.Error("Inster fields fail: ", err)
	} else {
		requestLogger := logrus.WithFields(logrus.Fields{"table": modules, "fields": fields})
		requestLogger.Info("insert sucess.")
	}
}
