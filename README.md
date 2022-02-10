## Loki API Middleware
###  介绍
#### 当高并发日志场景，Loki 聚合查询展示至Grafana性能严重不足，虽然可以通过分布式loki解决query高密度计算问题，但是需要大量机器完成 且没有必要在这个方面花费大量金钱，所以通过此中间件可以通过loki query api查询日志小间隔范围，比如cron 5分钟查询timelimit 5分钟范围，这样小范围查询loki性能则影响较小，并把响应的metrics存入influxdb通过influxdb对接grafana展示大范围聚合查询；此中间件相当于做了metrics过滤 利用influxdb减少loki性能开销（loki 2.0 目前处理 高并发日志聚哈查询还不够优秀）

### Yaml Config explain:  
timer:
  -second         定时器执行秒间隔
loki
  -url            loki query url填写
  -timelimit      loki query 时间范围
  -hosts          loki host 查询标识，字段对应