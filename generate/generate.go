package generate

import (
	"fmt"
	"github.com/google/uuid"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 课程具体时间，周几第几节
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

// 作息时间表，浑南上课时间
var ClassStartTimeHunnan = []string{
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

// 作息时间表，浑南下课时间
var classEndTimeHunnan = []string{
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

// 作息时间表，南湖上课时间
var ClassStartTimeNanhu = []string{
	"080000",
	"090000",
	"101000",
	"111000",
	"140000",
	"150000",
	"161000",
	"171000",
	"183000",
	"193000",
	"203000",
	"213000",
}

// 作息时间表，南湖下课时间
var classEndTimeNanhu = []string{
	"085000",
	"095000",
	"110000",
	"120000",
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

// 从html源码生成ics文件内容
func GenerateIcs(html string) (string, error) {
	fmt.Println("\n生成ics文件中。。。")
	// 利用正则匹配有效信息
	var myCourses []Course
	reg1 := regexp.MustCompile(`TaskActivity\(actTeacherId.join\(','\),actTeacherName.join\(','\),"(.*)","(.*)\(.*\)","(.*)","(.*)","(.*)",null,null,assistantName,"",""\);((?:\s*index =\d+\*unitCount\+\d+;\s*.*\s)+)`)
	reg2 := regexp.MustCompile(`\s*index =(\d+)\*unitCount\+(\d+);\s*`)
	coursesStr := reg1.FindAllStringSubmatch(html, -1)
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
	}

	// 生成ics文件头
	var icsData string
	icsData = `BEGIN:VCALENDAR
PRODID:-//nian//getMyCourses 20190522//EN
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
	// 2019-03-03，校历第一周周日
	location := time.FixedZone("UTC+8", 8*60*60)
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

		// 生成ics文件中的EVENT
		for i := 0; i < len(periods); i++ {
			var eventData string
			eventData = `BEGIN:VEVENT` + "\n"
			startDate := SchoolStartDay.AddDate(0, 0, (startWeek[i]-1)*7+weekDay+1)

			if strings.Contains(course.roomName, "浑南") {
				eventData = eventData + `DTSTART;TZID=Asia/Shanghai:` + startDate.Format("20060102T") + ClassStartTimeHunnan[st] + "\n"
				eventData = eventData + `DTEND;TZID=Asia/Shanghai:` + startDate.Format("20060102T") + classEndTimeHunnan[en] + "\n"
			} else {
				eventData = eventData + `DTSTART;TZID=Asia/Shanghai:` + startDate.Format("20060102T") + ClassStartTimeNanhu[st] + "\n"
				eventData = eventData + `DTEND;TZID=Asia/Shanghai:` + startDate.Format("20060102T") + classEndTimeNanhu[en] + "\n"
			}
			eventData = eventData + periods[i] + "\n"
			eventData = eventData + `DTSTAMP:` + time.Now().Format("20060102T150405Z") + "\n"
			eventData = eventData + `UID:` + uuid.New().String() + "\n"
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

	fmt.Println("\n生成ics文件完成。")
	return icsData, nil
}
