package main

import (
	"fmt"
	"github.com/whoisnian/getMyCourses/fetch"
	"github.com/whoisnian/getMyCourses/generate"
	"github.com/whoisnian/getMyCourses/login"
	"io/ioutil"
	"net/http/cookiejar"
	"path/filepath"
	"time"
)

func main() {
	// 选择登录方式
	var choice int
	fmt.Println("1.树维教务系统登录: http://219.216.96.4/eams/loginExt.action")
	fmt.Println("2.东大统一身份认证: https://pass.neu.edu.cn")
	fmt.Printf("\n请选择登录方式（1 或 2）：")
	_, err := fmt.Scanln(&choice)
	fmt.Printf("\n")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 获取帐号和密码
	var username, password string
	fmt.Printf("帐号: ")
	fmt.Scanln(&username)
	fmt.Printf("密码: ")
	fmt.Scanln(&password)

	// 登录
	var cookieJar *cookiejar.Jar
	if choice == 1 {
		cookieJar, err = login.LoginViaSupwisdom(username, password)
	} else {
		cookieJar, err = login.LoginViaTpass(username, password)
	}

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 获取包含课程表的html源码
	html, err := fetch.FetchCourses(cookieJar)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 获取当前教学周
	learnWeek, err := fetch.FetchLearnWeek(cookieJar)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 计算校历第一周周日
	now := time.Now()
	location := time.FixedZone("UTC+8", 8*60*60)
	daySum := int(now.Weekday()) + learnWeek*7 - 7
	schoolStartDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, location).AddDate(0, 0, -daySum)

	fmt.Println("\n当前为第", learnWeek, "教学周。")
	fmt.Println("计算得到本学期开始于", schoolStartDay.Format("2006-01-02"))
	fmt.Println("官方校历 http://www.neu.edu.cn/xl/list.htm")

	// 从html源码生成ics文件内容
	ics, err := generate.GenerateIcs(html, schoolStartDay)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 保存到文件
	err = ioutil.WriteFile("myCourses.ics", []byte(ics), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 提示文件路径
	path, err := filepath.Abs("myCourses.ics")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("\n已保存为：", path)
}
