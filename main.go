package main

import (
	"fmt"
	"loki_client/config"
	"loki_client/module"
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
					fmt.Println(logql)
					if modules == "access" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}
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
								}
							}

							fields := map[string]interface{}{
								"geoip_city_name":    value["geoip_city_name"],
								"geoip_country_name": value["geoip_country_name"],
								"remote_addr":        value["remote_addr"],
								"request_ip_sum":     value["request_ip_sum"],
								"host":               host,
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, fields)
							conn.Close()
						}
					} else if modules == "pages" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}
							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "metric" {
									vv1 := v1.(map[string]interface{})
									value["request_uri"] = vv1["request_uri"].(string)
								}

								if k1 == "value" {
									vv1 := v1.([]interface{})
									value["request_uri_sum"] = vv1[1].(string)
								}
							}
							fields := map[string]interface{}{
								"request_uri":     value["request_uri"],
								"request_uri_sum": value["request_uri_sum"],
								"host":            host,
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, fields)
							conn.Close()
						}
					} else if modules == "referers" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}
							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "metric" {
									vv1 := v1.(map[string]interface{})
									value["http_referer"] = vv1["http_referer"].(string)
								}

								if k1 == "value" {
									vv1 := v1.([]interface{})
									value["http_referer_sum"] = vv1[1].(string)
								}
							}
							fields := map[string]interface{}{
								"http_referer":     value["http_referer"],
								"http_referer_sum": value["http_referer_sum"],
								"host":             host,
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, fields)
							conn.Close()
						}
					} else if modules == "useragents" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}
							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "metric" {
									vv1 := v1.(map[string]interface{})
									value["http_user_agent"] = vv1["http_user_agent"].(string)
								}

								if k1 == "value" {
									vv1 := v1.([]interface{})
									value["http_user_agent_sum"] = vv1[1].(string)
								}
							}
							fields := map[string]interface{}{
								"http_user_agent":     value["http_user_agent"],
								"http_user_agent_sum": value["http_user_agent_sum"],
								"host":                host,
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, fields)
							conn.Close()
						}
					} else if modules == "total_request" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}
							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "value" {
									vv1 := v1.([]interface{})
									value["request_sum"] = vv1[1].(string)
								}
							}
							fields := map[string]interface{}{
								"request_sum": value["request_sum"],
								"host":        host,
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, fields)
							conn.Close()
						}
					} else if modules == "realtime_visitors" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}
							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "value" {
									vv1 := v1.([]interface{})
									value["realtime_visitors_sum"] = vv1[1].(string)
								}
							}
							fields := map[string]interface{}{
								"realtime_visitors_sum": value["realtime_visitors_sum"],
								"host":                  host,
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, fields)
							conn.Close()
						}
					} else if modules == "request_status_code" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}

							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "metric" {
									vv1 := v1.(map[string]interface{})
									value["status"] = vv1["status"].(string)
								}

								if k1 == "value" {
									vv1 := v1.([]interface{})
									value["status_sum"] = vv1[1].(string)
								}
							}
							fields := map[string]interface{}{
								"status":     value["status"],
								"status_sum": value["status_sum"],
								"host":       host,
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, fields)
							conn.Close()
						}
					} else if modules == "geoipcode" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}

							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "metric" {
									vv1 := v1.(map[string]interface{})
									value["geoip_country_code"] = vv1["geoip_country_code"].(string)
								}

								if k1 == "value" {
									vv1 := v1.([]interface{})
									value["geoip_country_code_num"] = vv1[1].(string)
								}
							}
							fields := map[string]interface{}{
								"geoip_country_code":     value["geoip_country_code"],
								"geoip_country_code_num": value["geoip_country_code_num"],
								"host":                   host,
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, fields)
							conn.Close()
						}
					} else if modules == "request_time" {
						for _, v := range metrics.([]interface{}) {
							value := map[string]string{}

							for k1, v1 := range v.(map[string]interface{}) {
								if k1 == "metric" {
									vv1 := v1.(map[string]interface{})
									value["host"] = vv1["host"].(string)
								}

								if k1 == "value" {
									vv1 := v1.([]interface{})
									value["request_time"] = vv1[1].(string)
								}
							}
							fields := map[string]interface{}{
								"host":         value["host"],
								"request_time": value["request_time"],
							}
							conn := module.Conninflux()
							module.Writeinflux(conn, modules, fields)
							conn.Close()
						}
					}
				}
			}
		}()
	}
}
