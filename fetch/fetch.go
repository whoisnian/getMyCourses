package fetch

import (
	"errors"
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

// 获取课程表所在页面源代码
func FetchCourses(cookieJar *cookiejar.Jar) (string, error) {
	fmt.Println("\n获取课表详情中。。。")

	// http 请求客户端
	var client http.Client
	client.Jar = cookieJar

	// 第一次请求
	time.Sleep(1 * time.Second)
	req, err := http.NewRequest(http.MethodGet, "http://219.216.96.4/eams/courseTableForStd.action", nil)
	if err != nil {
		return "", err
	}

	// 发送
	resp1, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp1.Body.Close()

	// 读取
	content, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return "", err
	}

	// 检查
	temp := string(content)
	if !strings.Contains(temp, "bg.form.addInput(form,\"ids\",\"") {
		return "", errors.New("获取ids失败")
	}
	temp = temp[strings.Index(temp, "bg.form.addInput(form,\"ids\",\"")+29 : strings.Index(temp, "bg.form.addInput(form,\"ids\",\"")+50]
	ids := temp[:strings.Index(temp, "\");")]
	semesterId := resp1.Cookies()[0].Value

	// 第二次请求
	time.Sleep(1 * time.Second)
	formValues := make(url.Values)
	formValues.Set("ignoreHead", "1")
	formValues.Set("showPrintAndExport", "1")
	formValues.Set("setting.kind", "std")
	formValues.Set("startWeek", "")
	formValues.Set("semester.id", semesterId)
	formValues.Set("ids", ids)

	req, err = http.NewRequest(http.MethodPost, "http://219.216.96.4/eams/courseTableForStd!courseTable.action", strings.NewReader(formValues.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0")

	// 发送
	resp2, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp2.Body.Close()

	// 读取
	content, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		return "", err
	}

	// 检查
	temp = string(content)
	if !strings.Contains(temp, "课表格式说明") {
		return "", errors.New("获取课表失败")
	}

	fmt.Println("获取课表详情完成。")
	return temp, nil
}

// 获取当前教学周
func FetchLearnWeek(cookieJar *cookiejar.Jar) (int, error) {
	fmt.Println("\n获取当前教学周中。。。")

	// http 请求客户端
	var client http.Client
	client.Jar = cookieJar

	req, err := http.NewRequest(http.MethodGet, "http://219.216.96.4/eams/homeExt.action", nil)
	if err != nil {
		return 0, err
	}

	// 发送
	resp1, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp1.Body.Close()

	// 读取
	content, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return 0, err
	}

	// 检查
	temp := string(content)
	if !strings.Contains(temp, "教学周") {
		return 0, errors.New("获取教学周失败")
	}
	temp = temp[strings.Index(temp, "id=\"teach-week\">") : strings.Index(temp, "教学周")+10]

	reg := regexp.MustCompile(`学期\s*<font size="\d+px">(\d+)<\/font>\s*教学周`)
	res := reg.FindStringSubmatch(temp)
	if len(res) < 2 {
		return 0, errors.New(temp + " 中未匹配到教学周")
	}

	fmt.Println("获取当前教学周完成。")
	return strconv.Atoi(res[1])
}
