package login

import (
	"crypto/sha1"
	"encoding/hex"
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

// 树维教务系统登录：http://219.216.96.4/eams/loginExt.action
func LoginViaSupwisdom(username string, password string) (*cookiejar.Jar, error) {
	fmt.Println("\n树维教务系统登录中。。。")

	// Cookie 自动维护
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	// http 请求客户端
	var client http.Client
	client.Jar = cookieJar

	// 第一次请求
	req, err := http.NewRequest(http.MethodGet, "http://219.216.96.4/eams/loginExt.action", nil)
	if err != nil {
		return nil, err
	}

	// 发送
	resp1, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()

	// 读取
	content, err := ioutil.ReadAll(resp1.Body)
	if err != nil {
		return nil, err
	}

	// 检查
	temp := string(content)
	if !strings.Contains(temp, "CryptoJS.SHA1(") {
		return nil, errors.New("登录页面打开失败，请检查 http://219.216.96.4/eams/loginExt.action")
	}

	// 对密码进行SHA1哈希
	temp = temp[strings.Index(temp, "CryptoJS.SHA1(")+15 : strings.Index(temp, "CryptoJS.SHA1(")+52]
	password = temp + password
	bytes := sha1.Sum([]byte(password))
	password = hex.EncodeToString(bytes[:])

	// 第二次请求
	time.Sleep(1 * time.Second)
	formValues := make(url.Values)
	formValues.Set("username", username)
	formValues.Set("password", password)
	formValues.Set("session_locale", "zh_CN")

	req, err = http.NewRequest(http.MethodPost, "http://219.216.96.4/eams/loginExt.action", strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0")

	// 发送
	resp2, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	// 读取
	content, err = ioutil.ReadAll(resp2.Body)
	if err != nil {
		return nil, err
	}

	// 检查
	temp = string(content)
	if !strings.Contains(temp, "personal-name") {
		return nil, errors.New("登录失败，请检查用户名和密码")
	}

	temp = temp[strings.Index(temp, "class=\"personal-name\">")+23 : strings.Index(temp, "class=\"personal-name\">")+60]
	fmt.Println(temp[:strings.Index(temp, ")")+1])

	fmt.Println("树维教务系统登录完成。")
	return cookieJar, nil
}

// 统一身份认证：https://pass.neu.edu.cn
func LoginViaTpass(username string, password string) (*cookiejar.Jar, error) {
	fmt.Println("\n统一身份认证登录中。。。")

	// Cookie 自动维护
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	// http 请求客户端
	var client http.Client
	client.Jar = cookieJar

	// 第一次请求
	req, err := http.NewRequest(http.MethodGet, "http://219.216.96.4/eams/localLogin!tip.action", nil)
	if err != nil {
		return nil, err
	}

	// 发送
	resp1, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp1.Body.Close()

	// 第二次请求
	req, err = http.NewRequest(http.MethodGet, "https://pass.neu.edu.cn/tpass/login?service=http%3A%2F%2F219.216.96.4%2Feams%2FhomeExt.action", nil)
	if err != nil {
		return nil, err
	}

	// 发送
	resp2, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp2.Body.Close()

	// 读取
	content, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		return nil, err
	}

	// 检查
	temp := string(content)
	if !strings.Contains(temp, "<form id=\"loginForm\" action=\"/tpass/login") {
		return nil, errors.New("登录页面打开失败，请检查 https://pass.neu.edu.cn")
	}

	// 提取表单信息
	reg1 := regexp.MustCompile(`id="loginForm" action="(/tpass/login;jsessionid=[^"]*)`)
	form_action := "https://pass.neu.edu.cn" + reg1.FindStringSubmatch(temp)[1]
	reg2 := regexp.MustCompile(`id="lt" name="lt" value="([^"]*)`)
	form_lt := reg2.FindStringSubmatch(temp)[1]
	form_rsa := username + password + form_lt
	form_ul := len(username)
	form_pl := len(password)
	reg3 := regexp.MustCompile(`name="execution" value="([^"]*)`)
	form_execution := reg3.FindStringSubmatch(temp)[1]
	reg4 := regexp.MustCompile(`name="_eventId" value="([^"]*)`)
	form__eventId := reg4.FindStringSubmatch(temp)[1]

	// 第三次请求
	time.Sleep(1 * time.Second)
	formValues := make(url.Values)
	formValues.Set("rsa", form_rsa)
	formValues.Set("ul", strconv.Itoa(form_ul))
	formValues.Set("pl", strconv.Itoa(form_pl))
	formValues.Set("lt", form_lt)
	formValues.Set("execution", form_execution)
	formValues.Set("_eventId", form__eventId)

	req, err = http.NewRequest(http.MethodPost, form_action, strings.NewReader(formValues.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:66.0) Gecko/20100101 Firefox/66.0")

	// 发送
	resp3, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp3.Body.Close()

	// 读取
	content, err = ioutil.ReadAll(resp3.Body)
	if err != nil {
		return nil, err
	}

	// 检查
	temp = string(content)
	if !strings.Contains(temp, "personal-name") {
		return nil, errors.New("登录失败，请检查用户名和密码")
	}

	temp = temp[strings.Index(temp, "class=\"personal-name\">")+23 : strings.Index(temp, "class=\"personal-name\">")+60]
	fmt.Println(temp[:strings.Index(temp, ")")+1])

	fmt.Println("统一身份登录认证登录完成。")
	return cookieJar, nil
}
