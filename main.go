package main

import (
	"fmt"
	"getMyCourses/fetch"
	"getMyCourses/generate"
	"getMyCourses/login"
)

func main() {
	// 获取用户名和密码
	var username, password string

	fmt.Printf("username: ")
	fmt.Scanln(&username)
	fmt.Printf("password: ")
	fmt.Scanln(&password)

	// 登录
	cookieJar, err := login.LoginViaSupwisdom(username, password)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 获取包含课程表的html源代码
	html, err := fetch.FetchCourses(cookieJar)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 从html代码生成ics文件内容
	ics, err := generate.GenerateIcs(html)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 保存文件
	fmt.Println(len(ics))
}
