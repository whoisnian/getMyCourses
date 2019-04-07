package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 课程持续时间，周几第几节
type CourseTime struct {
	dayOfTheWeek int
	timeOfTheDay int
}

// 课程信息
type Course struct {
	courseID    string
	courseName  string
	roomID      string
	roomName    string
	weeks       string
	courseTimes []CourseTime
}

// 作息时间表，上课时间
var classStartTime = []string{
	"083000",
	"093000",
	"104000",
	"114000",
	"140000",
	"150000",
	"161000",
	"171000",
	"183000",
	"193000",
	"203000",
	"213000",
}

// 作息时间表，下课时间
var classEndTime = []string{
	"092000",
	"102000",
	"113000",
	"123000",
	"145000",
	"155000",
	"170000",
	"180000",
	"192000",
	"202000",
	"212000",
	"222000",
}

// ics文件用到的星期几简称
var dayOfWeek = []string{
	"MO",
	"TU",
	"WE",
	"TH",
	"FR",
	"SA",
	"SU",
}

var USERNAME, PASSWORD string
var myCourses []Course

func main() {
	// 获取用户名和密码
	fmt.Printf("username: ")
	fmt.Scanln(&USERNAME)
	fmt.Printf("password: ")
	fmt.Scanln(&PASSWORD)

	// Cookie自动维护
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("ERROR_0: ", err.Error())
		return
	}
	var client http.Client
	client.Jar = cookieJar

	// 第一次请求
	req, err := http.NewRequest(http.MethodGet, "http://219.216.96.4/eams/loginExt.action", nil)
	if err != nil {
		fmt.Println("ERROR_1: ", err.Error())
		return
	}

	resp1, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR_2: ", err.Error())
		return
	}
	defer resp1.Body.Close()

	content, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		fmt.Println("ERROR_3: ", err.Error())
		return
	}

	temp := string(content)
	if !strings.Contains(temp, "CryptoJS.SHA1(") {
		fmt.Println("ERROR_4: GET Failed")
		return
	}

	// 对密码进行SHA1哈希
	temp = temp[strings.Index(temp, "CryptoJS.SHA1(")+15 : strings.Index(temp, "CryptoJS.SHA1(")+52]
	PASSWORD = temp + PASSWORD
	bytes := sha1.Sum([]byte(PASSWORD))
	PASSWORD = hex.EncodeToString(bytes[:])

	fmt.Printf("\n登录中。。。\n")
	time.Sleep(1 * time.Second)
	// 第二次请求
	formValues := make(url.Values)
	formValues.Set("username", USERNAME)
	formValues.Set("password", PASSWORD)
	formValues.Set("session_locale", "zh_CN")
	req, err = http.NewRequest(http.MethodPost, "http://219.216.96.4/eams/loginExt.action", strings.NewReader(formValues.Encode()))
	if err != nil {
		fmt.Println("ERROR_5: ", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0")
	resp2, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR_6: ", err.Error())
		return
	}
	defer resp2.Body.Close()

	content, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		fmt.Println("ERROR_7: ", err.Error())
		return
	}

	temp = string(content)
	if !strings.Contains(temp, "personal-name") {
		fmt.Println("ERROR_8: LOGIN Failed")
		return
	}

	temp = temp[strings.Index(temp, "class=\"personal-name\">")+23 : strings.Index(temp, "class=\"personal-name\">")+42]
	fmt.Println("Login as " + temp)

	fmt.Printf("\n加载课表中。。。\n")
	time.Sleep(1 * time.Second)
	// 第三次请求
	req, err = http.NewRequest(http.MethodGet, "http://219.216.96.4/eams/courseTableForStd.action", nil)
	if err != nil {
		fmt.Println("ERROR_9: ", err.Error())
		return
	}

	resp3, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR_10: ", err.Error())
		return
	}
	defer resp3.Body.Close()

	content, err = ioutil.ReadAll(resp3.Body)
	if err != nil {
		fmt.Println("ERROR_11: ", err.Error())
		return
	}

	temp = string(content)
	if !strings.Contains(temp, "bg.form.addInput(form,\"ids\",\"") {
		fmt.Println("ERROR_12: GET ids Failed")
		return
	}
	temp = temp[strings.Index(temp, "bg.form.addInput(form,\"ids\",\"")+29 : strings.Index(temp, "bg.form.addInput(form,\"ids\",\"")+50]
	ids := temp[:strings.Index(temp, "\");")]

	time.Sleep(1 * time.Second)
	// 第四次请求
	formValues = make(url.Values)
	formValues.Set("ignoreHead", "1")
	formValues.Set("showPrintAndExport", "1")
	formValues.Set("setting.kind", "std")
	formValues.Set("startWeek", "")
	formValues.Set("semester.id", "30")
	formValues.Set("ids", ids)
	req, err = http.NewRequest(http.MethodPost, "http://219.216.96.4/eams/courseTableForStd!courseTable.action", strings.NewReader(formValues.Encode()))
	if err != nil {
		fmt.Println("ERROR_13: ", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0")
	resp4, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR_14: ", err.Error())
		return
	}
	defer resp4.Body.Close()

	content, err = ioutil.ReadAll(resp4.Body)
	if err != nil {
		fmt.Println("ERROR_15: ", err.Error())
		return
	}

	temp = string(content)
	if !strings.Contains(temp, "课表格式说明") {
		fmt.Println("ERROR_16: Get Courses Failed")
		return
	}

	// 利用正则匹配有效信息
	reg1 := regexp.MustCompile(`TaskActivity\(actTeacherId.join\(','\),actTeacherName.join\(','\),"(.*)","(.*)\(.*\)","(.*)","(.*)","(.*)",null,null,assistantName,"",""\);((?:\s*index =\d+\*unitCount\+\d+;\s*.*\s)+)`)
	reg2 := regexp.MustCompile(`\s*index =(\d+)\*unitCount\+(\d+);\s*`)
	coursesStr := reg1.FindAllStringSubmatch(temp, -1)
	for _, courseStr := range coursesStr {
		var course Course
		course.courseID = courseStr[1]
		course.courseName = courseStr[2]
		course.roomID = courseStr[3]
		course.roomName = courseStr[4]
		course.weeks = courseStr[5]
		for _, indexStr := range strings.Split(courseStr[6], "table0.activities[index][table0.activities[index].length]=activity;") {
			if !strings.Contains(indexStr, "unitCount") {
				continue
			}
			var courseTime CourseTime
			courseTime.dayOfTheWeek, _ = strconv.Atoi(reg2.FindStringSubmatch(indexStr)[1])
			courseTime.timeOfTheDay, _ = strconv.Atoi(reg2.FindStringSubmatch(indexStr)[2])
			course.courseTimes = append(course.courseTimes, courseTime)
		}
		myCourses = append(myCourses, course)
		fmt.Println(course)
	}

	fmt.Printf("\n注销中。。。\n")
	time.Sleep(1 * time.Second)
	// 第五次请求
	req, err = http.NewRequest(http.MethodGet, "http://219.216.96.4/eams/logout.action", nil)
	if err != nil {
		fmt.Println("ERROR_17: ", err.Error())
		return
	}

	resp5, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR_18: ", err.Error())
		return
	}
	defer resp5.Body.Close()

	fmt.Printf("\n生成中。。。\n")
	// 生成ics文件用于导入
	var icsData string
	icsData = `BEGIN:VCALENDAR
PRODID:-//nian//getMyCourses 20190324//EN
VERSION:2.0
CALSCALE:GREGORIAN
METHOD:PUBLISH
X-WR-CALNAME:myCourses
X-WR-TIMEZONE:Asia/Shanghai
BEGIN:VTIMEZONE
TZID:Asia/Shanghai
X-LIC-LOCATION:Asia/Shanghai
BEGIN:STANDARD
TZOFFSETFROM:+0800
TZOFFSETTO:+0800
TZNAME:CST
DTSTART:19700101T000000
END:STANDARD
END:VTIMEZONE` + "\n"

	// 本学期第一周开始时间
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Println("ERROR_19: ", err.Error())
		return
	}
	// 2019-03-03，校历第一周周日
	SchoolStartDay := time.Date(2019, time.March, 3, 0, 0, 0, 0, location)

	num := 0
	for _, course := range myCourses {
		var weekDay, st, en int
		weekDay = course.courseTimes[0].dayOfTheWeek
		st = 12
		en = -1
		// 课程上下课时间
		for _, courseTime := range course.courseTimes {
			if st > courseTime.timeOfTheDay {
				st = courseTime.timeOfTheDay
			}
			if en < courseTime.timeOfTheDay {
				en = courseTime.timeOfTheDay
			}
		}

		// debug信息
		num++
		fmt.Println("")
		fmt.Println(num)
		fmt.Println(course.courseName)
		fmt.Println("周" + strconv.Itoa(weekDay) + " 第" + strconv.Itoa(st+1) + "-" + strconv.Itoa(en+1) + "节")

		// 统计要上课的周
		var periods []string
		var startWeek []int
		byday := dayOfWeek[weekDay]
		for i := 0; i < 53; i++ {
			if course.weeks[i] != '1' {
				continue
			}
			if i+1 >= 53 {
				startWeek = append(startWeek, i)
				periods = append(periods, "RRULE:FREQ=WEEKLY;WKST=SU;COUNT=1;INTERVAL=1;BYDAY="+byday)
				// debug信息
				fmt.Println("第" + strconv.Itoa(i) + "周")
				continue
			}
			if course.weeks[i+1] == '1' {
				// 连续周合并
				var j int
				for j = i + 1; j < 53; j++ {
					if course.weeks[j] != '1' {
						break
					}
				}
				startWeek = append(startWeek, i)
				periods = append(periods, "RRULE:FREQ=WEEKLY;WKST=SU;COUNT="+strconv.Itoa(j-i)+";INTERVAL=1;BYDAY="+byday)
				// debug信息
				fmt.Println("第" + strconv.Itoa(i) + "-" + strconv.Itoa(j-1) + "周")
				i = j - 1
			} else {
				// 单双周合并
				var j int
				for j = i + 1; j+1 < 53; j += 2 {
					if course.weeks[j] == '1' || course.weeks[j+1] == '0' {
						break
					}
				}
				startWeek = append(startWeek, i)
				periods = append(periods, "RRULE:FREQ=WEEKLY;WKST=SU;COUNT="+strconv.Itoa((j+1-i)/2)+";INTERVAL=2;BYDAY="+byday)
				// debug信息
				if i%2 == 0 {
					fmt.Printf("双")
				} else {
					fmt.Printf("单")
				}
				fmt.Println(strconv.Itoa(i) + "-" + strconv.Itoa(j-1) + "周")
				i = j - 1
			}
		}

		// 生成EVENT
		for i := 0; i < len(periods); i++ {
			var eventData string
			eventData = `BEGIN:VEVENT` + "\n"
			startDate := SchoolStartDay.AddDate(0, 0, (startWeek[i]-1)*7+weekDay+1)

			eventData = eventData + `DTSTART;TZID=Asia/Shanghai:` + startDate.Format("20060102T") + classStartTime[st] + "\n"
			eventData = eventData + `DTEND;TZID=Asia/Shanghai:` + startDate.Format("20060102T") + classEndTime[en] + "\n"
			eventData = eventData + periods[i] + "\n"
			eventData = eventData + `DTSTAMP:` + time.Now().Format("20060102T150405Z") + "\n"
			eventData = eventData + `UID:` + "\n"
			eventData = eventData + `CREATED:` + time.Now().Format("20060102T150405Z") + "\n"
			eventData = eventData + `DESCRIPTION:` + "\n"
			eventData = eventData + `LAST-MODIFIED:` + time.Now().Format("20060102T150405Z") + "\n"
			eventData = eventData + `LOCATION:` + course.roomName + "\n"
			eventData = eventData + `SEQUENCE:0
STATUS:CONFIRMED` + "\n"
			eventData = eventData + `SUMMARY:` + course.courseName + "\n"

			eventData = eventData + `TRANSP:OPAQUE
END:VEVENT` + "\n"
			icsData = icsData + eventData
		}
	}
	icsData = icsData + `END:VCALENDAR`
	//fmt.Println(icsData)

	// 写入文件
	err = ioutil.WriteFile("./myCourses.ics", []byte(icsData), 0644)
	if err != nil {
		fmt.Println("ERROR_20: ", err.Error())
		return
	}
}
