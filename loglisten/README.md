## 日志监控使用方法

e := errorLog.New(errorLog.Option{})

go e.Start()

Option字段描述
```json
{
  ServerName     string //服务名称
  Path           string //日志路径
  Email          string //接收者邮箱，多个用","进行分割
  WechatPhone    string //企业微信手机号，多个用","进行分割
  TotalThreshold []int  //总日志数预警
  CycleThreshold int    //周期新增日志数预警
  CycleTime      int    //周期时间（分钟）
  HourNum        int    //每小时预警次数
  DayNum         int    //每天预警次数
}
```
option例子
```json
listen, _ := errorLog.New(errorLog.Option{
		ServerName:     "服务名称",
		Path:           "./logs/log-error.%s",
		Email:          "446xxxxx0@qq.com",
		WechatPhone:    "1355xxx1566",
		TotalThreshold: []int{10, 20},
		CycleThreshold: 3,
		CycleTime:      1,
		HourNum:        2,
		DayNum:         5,
	})
go listen.Start()
```
