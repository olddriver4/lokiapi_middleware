package main

import (
	"fmt"
	"loki_client/config"
	"loki_client/module"
	"strconv"
	"time"
)

func main() {
	url := config.ReadConfig("loki.url").(string)
	hosts := config.ReadConfig("loki.hosts").([]interface{})
	timelimit := config.ReadConfig("loki.timelimit").(string)

	ticker := time.NewTicker(time.Duration(config.ReadConfig("timer.second").(int)) * time.Second)
	for range ticker.C { //定时器运行
		go func() { //goruntine 并行执行
			for _, host := range hosts {
				logqls := make(map[string]string)
				host := host.(string)
				//请求loki query入库
				logqls["access"] = fmt.Sprintf("query=topk(10, sum by (remote_addr, geoip_country_name, geoip_city_name) (count_over_time({host=~\"%s\"} | json | __error__=\"\" [%s])))", host, timelimit)
				logqls["pages"] = fmt.Sprintf("query=topk(10, sum by (request_uri) (count_over_time({host=~\"%s\"} !~ `\\.ico|\\.svg|\\.css|\\.png|\\.txt|\\.js|\\.xml` | json | status = 200 and request_uri != \"\" | __error__=\"\" [%s])))", host, timelimit)
				logqls["referers"] = fmt.Sprintf("query=topk(10, sum by (http_referer) (count_over_time({host=~\"%s\"} | json |  http_referer != \"\" and http_referer !~ \".*?$host.*?\" and http_referer !~ \".*?\\\\*\\\\*\\\\*.*?\" | __error__=\"\" [%s])))", host, timelimit)
				logqls["useragents"] = fmt.Sprintf("query=topk(10, sum by (http_user_agent) (count_over_time({host=~\"%s\"} | json |  __error__=\"\" [%s])))", host, timelimit)
				logqls["total_request"] = fmt.Sprintf("query=sum by(host) (count_over_time({host=~\"%s\"}[%s]))", host, timelimit)
				logqls["realtime_visitors"] = fmt.Sprintf("query=count(sum by (remote_addr) (count_over_time({host=~\"%s\"} | json |  __error__=\"\" [%s])))", host, timelimit)
				logqls["request_status_code"] = fmt.Sprintf("query=sum by (status) (count_over_time({host=~\"%s\"} | json |  __error__=\"\" [%s]))", host, timelimit)
				logqls["geoipcode"] = fmt.Sprintf("query=sum by (geoip_country_code) (count_over_time({host=~\"%s\"} | json | __error__=\"\" [%s]))", host, timelimit)
				logqls["request_time"] = fmt.Sprintf("query=avg_over_time({host=~\"%s\"} | json | unwrap request_time |  __error__=\"\"  [%s]) by (host)", host, timelimit)

				for modules, logql := range logqls {
					metrics := *module.Loki_api(url, logql)
					if modules == "access" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}
							value_int := map[string]int{}
							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "metric" {
									vv1 := v1.(map[string]interface{})
									value["geoip_city_name"] = vv1["geoip_city_name"].(string)
									value["geoip_country_name"] = vv1["geoip_country_name"].(string)
									value["remote_addr"] = vv1["remote_addr"].(string)
								}

								if k1 == "value" {
									vv1 := v1.([]interface{})
									value["request_ip_sum"] = vv1[1].(string)
									value_int["request_ip_sum"], _ = strconv.Atoi(vv1[1].(string))
								}
							}

							// 写入metrics表
							tags := map[string]string{
								"host":               host,
								"geoip_city_name":    value["geoip_city_name"],
								"geoip_country_name": value["geoip_country_name"],
								"remote_addr":        value["remote_addr"],
							}

							fields := map[string]interface{}{
								"request_ip_sum": value_int["request_ip_sum"],
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, tags, fields)
							conn.Close()
						}
					} else if modules == "pages" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}
							value_int := map[string]int{}
							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "metric" {
									vv1 := v1.(map[string]interface{})
									value["request_uri"] = vv1["request_uri"].(string)
								}

								if k1 == "value" {
									vv1 := v1.([]interface{})
									value_int["request_uri_sum"], _ = strconv.Atoi(vv1[1].(string))
								}
							}

							// 写入metrics表
							tags := map[string]string{
								"host":        host,
								"request_uri": value["request_uri"],
							}
							fields := map[string]interface{}{
								"request_uri_sum": value_int["request_uri_sum"],
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, tags, fields)
							conn.Close()
						}
					} else if modules == "referers" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}
							value_int := map[string]int{}
							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "metric" {
									vv1 := v1.(map[string]interface{})
									value["http_referer"] = vv1["http_referer"].(string)
								}

								if k1 == "value" {
									vv1 := v1.([]interface{})
									value_int["http_referer_sum"], _ = strconv.Atoi(vv1[1].(string))
								}
							}

							// 写入metrics表
							tags := map[string]string{
								"host":         host,
								"http_referer": value["http_referer"],
							}
							fields := map[string]interface{}{
								"http_referer_sum": value_int["http_referer_sum"],
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, tags, fields)
							conn.Close()
						}
					} else if modules == "useragents" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}
							value_int := map[string]int{}
							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "metric" {
									vv1 := v1.(map[string]interface{})
									value["http_user_agent"] = vv1["http_user_agent"].(string)
								}

								if k1 == "value" {
									vv1 := v1.([]interface{})
									value_int["http_user_agent_sum"], _ = strconv.Atoi(vv1[1].(string))
								}
							}
							// 写入metrics表
							tags := map[string]string{
								"host":            host,
								"http_user_agent": value["http_user_agent"],
							}
							fields := map[string]interface{}{
								"http_user_agent_sum": value_int["http_user_agent_sum"],
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, tags, fields)
							conn.Close()
						}
					} else if modules == "total_request" {
						for _, v := range metrics.([]interface{}) {
							value_int := map[string]int{}
							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "value" {
									vv1 := v1.([]interface{})
									value_int["request_sum"], _ = strconv.Atoi(vv1[1].(string))
								}
							}
							// 写入metrics表
							tags := map[string]string{
								"host": host,
							}
							fields := map[string]interface{}{
								"request_sum": value_int["request_sum"],
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, tags, fields)
							conn.Close()
						}
					} else if modules == "realtime_visitors" {
						for _, v := range metrics.([]interface{}) {
							value_int := map[string]int{}
							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "value" {
									vv1 := v1.([]interface{})
									value_int["realtime_visitors_sum"], _ = strconv.Atoi(vv1[1].(string))
								}
							}
							// 写入metrics表
							tags := map[string]string{
								"host": host,
							}
							fields := map[string]interface{}{
								"realtime_visitors_sum": value_int["realtime_visitors_sum"],
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, tags, fields)
							conn.Close()
						}
					} else if modules == "request_status_code" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}
							value_int := map[string]int{}

							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "metric" {
									vv1 := v1.(map[string]interface{})
									value["status"] = vv1["status"].(string)
								}

								if k1 == "value" {
									vv1 := v1.([]interface{})
									value_int["status_sum"], _ = strconv.Atoi(vv1[1].(string))
								}
							}
							// 写入metrics表
							tags := map[string]string{
								"host":   host,
								"status": value["status"],
							}
							fields := map[string]interface{}{
								"status_sum": value_int["status_sum"],
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, tags, fields)
							conn.Close()
						}
					} else if modules == "geoipcode" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}
							value_int := map[string]int{}

							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "metric" {
									vv1 := v1.(map[string]interface{})
									value["geoip_country_code"] = vv1["geoip_country_code"].(string)
								}

								if k1 == "value" {
									vv1 := v1.([]interface{})
									value_int["geoip_country_code_num"], _ = strconv.Atoi(vv1[1].(string))
								}
							}
							// 写入metrics表
							tags := map[string]string{
								"host":               host,
								"geoip_country_code": value["geoip_country_code"],
							}
							fields := map[string]interface{}{
								"geoip_country_code_num": value_int["geoip_country_code_num"],
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, tags, fields)
							conn.Close()
						}
					} else if modules == "request_time" {
						for _, v := range metrics.([]interface{}) {
							value_int := map[string]float64{}

							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "value" {
									vv1 := v1.([]interface{})
									value_int["request_time"], _ = strconv.ParseFloat(vv1[1].(string), 64)
								}
							}
							// 写入metrics表
							tags := map[string]string{
								"host": host,
							}
							fields := map[string]interface{}{
								"request_time": value_int["request_time"],
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, tags, fields)
							conn.Close()
						}
					}
				}
			}
		}()
	}
}
