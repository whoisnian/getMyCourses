# getMyCourses
从NEU新版树维教务系统获取自己的课程表，并生成可导入Google日历的.ics文件。

## 使用方法
* `go run getMyCourses.go`后按提示登录树维教务系统，会在同级目录生成`myCourses.ics`。
* 在Google日历设置页面创建新日历，选择`myCourses.ics`文件导入到新创建的日历，提示`已导入 x 个活动，共 x 个。`后表示导入成功。

## 参考信息
* 获取课程表的具体请求过程：[NEU 新版教务处课程表.md](https://gist.github.com/whoisnian/32b832bd55978fefa042d7c76f9d76c3)
* iCalendar格式介绍：[维基百科：ICalendar](https://en.wikipedia.org/wiki/ICalendar)
* Google日历帮助：[创建或编辑ICAL文件](https://support.google.com/calendar/answer/37118#format_ical)

## 注意
* 登录使用帐号密码为树维教务系统的帐号密码，登录页面为[http://219.216.96.4/eams/loginExt.action](http://219.216.96.4/eams/loginExt.action)，不是[学校统一身份认证平台](https://pass.neu.edu.cn/tpass/login)。  
* 生成.ics文件过程中会在命令行输出识别到的课程，请检查无误后再进行导入。
