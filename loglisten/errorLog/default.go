/*
createTime: 2022/11/29
*/
package errorLog

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/robfig/cron"
	util2 "gitlab.droi.cn/utils/serverlisten.git/util"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type ErrorClass struct {
	serverName     string
	path           string
	email          string
	weChatPhone    string
	totalThreshold []int
	cycleThreshold int
	cycleTime      int
	hourNum        int
	dayNum         int
	preLogLine     int
	nowTimeDay     string
	nowTimeHour    int
	sendCycleNum   int
	sendTotalNum   map[int]int
}

type Option struct {
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

func New(op Option) (*ErrorClass, error) {
	totalThreshold := []int{100, 200}
	cycleThreshold := 100
	cycleTime := 3
	hourNum := 2
	dayNum := 2

	if op.Path == "" {
		return nil, errors.New("日志路径不能为空")
	}
	if op.Email == "" && op.WechatPhone == "" {
		return nil, errors.New("接收者不能为空")
	}
	if len(op.TotalThreshold) > 0 {
		totalThreshold = op.TotalThreshold
	}
	if op.CycleThreshold > 0 {
		cycleThreshold = op.CycleThreshold
	}
	if op.CycleTime > 0 {
		cycleTime = op.CycleTime
	}
	if op.HourNum > 0 {
		hourNum = op.HourNum
	}
	if op.DayNum > 0 {
		dayNum = op.DayNum
	}

	return &ErrorClass{
		serverName:     op.ServerName,
		path:           op.Path,
		email:          op.Email,
		weChatPhone:    op.WechatPhone,
		totalThreshold: totalThreshold,
		cycleThreshold: cycleThreshold,
		cycleTime:      cycleTime,
		hourNum:        hourNum,
		dayNum:         dayNum,
	}, nil
}

func (e *ErrorClass) Start() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	c := cron.New()
	cycle := fmt.Sprintf("0 */%d * * * *", e.cycleTime)
	err := c.AddFunc(cycle, e.logListen)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	c.Start()
	defer c.Stop()
	select {}
}

func (e *ErrorClass) logListen() {
	errNum, info := e.errorLogListen()
	if errNum-e.preLogLine >= e.cycleThreshold {
		e.cycleHandler(errNum-e.preLogLine, info)
	}
	e.preLogLine = errNum

	for _, v := range e.totalThreshold {
		if errNum > v {
			e.totalHandler(v, errNum, info)
		}
	}
}

func (e *ErrorClass) errorLogListen() (int, string) {
	filename := fmt.Sprintf(e.path, time.Now().Format("20060102"))
	fileByte, _ := ioutil.ReadFile(filename)

	count := bytes.Count(fileByte, []byte(`"level":"error"`))
	if count > 0 {
		return count, readFile(filename)
	}
	return 0, ""
}

func readFile(file_name string) (info string) {
	file, err := os.Open(file_name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	var lineText string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText = scanner.Text()
		//fmt.Print(lineText)
	}

	return string(lineText)
}

func (e *ErrorClass) cycleHandler(line int, info string) {
	nowDay := time.Now().Format("20060102")
	nowHour := time.Now().Hour()
	if e.nowTimeDay != nowDay {
		e.nowTimeDay = nowDay
		e.nowTimeHour = nowHour
		e.sendCycleNum = 0
		e.sendTotalNum = map[int]int{}
	}

	if e.nowTimeHour != nowHour {
		e.nowTimeHour = nowHour
		e.sendCycleNum = 0
	}

	if e.sendCycleNum >= e.hourNum {
		return
	}

	//发送邮件
	var localIps []string
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIps = append(localIps, ipnet.IP.String())
			}
		}
	}
	sendSuccess := false
	if len(strings.Split(e.email, ",")) > 0 {
		subject := fmt.Sprintf("周期错误日志预警--服务名称：%s-%s,服务器错误预警", localIps[0], e.serverName)
		body := fmt.Sprintf("错误日志在%d分钟内增加%d行，请及时处理。最后一条错误日志为：%s，发送时间：%s", e.cycleTime, line, info, time.Now().Format("2006-01-02 15:04:05"))
		err := util2.SendMail(strings.Split(e.email, ","), subject, body)
		if err != nil {
			return
		}
		sendSuccess = true
	}

	//企业微信通知
	if len(strings.Split(e.weChatPhone, ",")) > 0 {
		body := fmt.Sprintf("周期错误日志预警--服务名称：%s-%s,服务器错误预警;错误日志在%d分钟内增加%d行，请及时处理。最后一条错误日志为：%s，发送时间：%s", e.serverName, localIps[0], e.cycleTime, line, info, time.Now().Format("2006-01-02 15:04:05"))
		code := util2.SendMsg(body, strings.Split(e.weChatPhone, ","))
		if code != 200 {
			return
		}
		sendSuccess = true
	}

	if sendSuccess {
		e.sendCycleNum += 1
	}
}

func (e *ErrorClass) totalHandler(threshold, line int, info string) {
	nowDay := time.Now().Format("20060102")
	nowHour := time.Now().Hour()
	if e.nowTimeDay != nowDay {
		e.nowTimeDay = nowDay
		e.nowTimeHour = nowHour
		e.sendCycleNum = 0
		e.sendTotalNum = map[int]int{}
	}

	if e.nowTimeHour != nowHour {
		e.nowTimeHour = nowHour
		e.sendCycleNum = 0
	}

	sendNum, ok := e.sendTotalNum[threshold]
	if !ok {
		e.sendTotalNum[threshold] = 0
	}
	if sendNum >= e.dayNum {
		return
	}

	var localIps []string
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIps = append(localIps, ipnet.IP.String())
			}
		}
	}

	//发送邮件
	sendSuccess := false
	if len(strings.Split(e.email, ",")) > 0 {
		subject := fmt.Sprintf("总错误日志预警--服务名称：%s-%s,服务器错误预警", e.serverName, localIps[0])
		body := fmt.Sprintf("服务器总错误日志达到%d行，请及时处理。最后一条错误日志为：%s，发送时间：%s", line, info, time.Now().Format("2006-01-02 15:04:05"))
		err := util2.SendMail(strings.Split(e.email, ","), subject, body)
		if err != nil {
			fmt.Println("邮件发送失败")
			return
		}
		sendSuccess = true
	}

	//企业微信通知
	if len(strings.Split(e.weChatPhone, ",")) > 0 {
		body := fmt.Sprintf("总错误日志预警--服务名称：%s-%s,服务器错误预警;服务器总错误日志达到%d行，请及时处理。最后一条错误日志为：%s，发送时间：%s", e.serverName, localIps[0], line, info, time.Now().Format("2006-01-02 15:04:05"))
		code := util2.SendMsg(body, strings.Split(e.weChatPhone, ","))
		if code != 200 {
			return
		}
		sendSuccess = true
	}

	if sendSuccess {
		e.sendTotalNum[threshold] += 1
	}
}
