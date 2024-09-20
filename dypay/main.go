package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/preorder", func(c *gin.Context) {

		timestamp := strconv.Itoa(int(time.Now().Unix()))

		url := "https://developer.toutiao.com/api/apps/ecpay/v1/create_order"
		paramsMap := map[string]interface{}{
			"app_id":       "ttb905cfb8263a12dc01",
			"out_order_no": "noncestr" + timestamp,
			"total_amount": 1,
			"subject":      "抖音商品",
			"body":         "抖音商品",
			"valid_time":   300,
			"notify_url":   "https://app.mujuapi.cn/api/pay/notify",
		}
		sing := getSign(paramsMap)

		postData := map[string]interface{}{
			"app_id":       "ttb905cfb8263a12dc01",
			"out_order_no": "noncestr" + timestamp,
			"total_amount": 1,
			"subject":      "抖音商品",
			"body":         "抖音商品",
			"valid_time":   300,
			"notify_url":   "https://app.mujuapi.cn/api/pay/notify",
			"sign":         sing,
		}

		// 将POST数据编码为JSON
		jsonBytes, err := json.Marshal(postData)

		if err != nil {
			panic(err)
		}
		// 创建POST请求
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
		if err != nil {
			panic(err)
		}
		req.Header.Set("Content-Type", "application/json") // 指定请求体内容类型为JSON

		// 执行请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		// 读取响应体
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}
		var result map[string]interface{}
		json.Unmarshal(respBody, &result)
		c.JSON(200, result)

	})

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

const salt = "JEwfiklDXIfPw4nl0lT15cpmJlreMPIMQOJbTOpT"

func getSign(paramsMap map[string]interface{}) string {
	var paramsArr []string
	for k, v := range paramsMap {
		if k == "other_settle_params" {
			continue
		}
		value := strings.TrimSpace(fmt.Sprintf("%v", v))
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") && len(value) > 1 {
			value = value[1 : len(value)-1]
		}
		value = strings.TrimSpace(value)
		if value == "" || value == "null" {
			continue
		}
		switch k {
		// app_id, thirdparty_id, sign 字段用于标识身份，不参与签名
		case "app_id", "thirdparty_id", "sign":
		default:
			paramsArr = append(paramsArr, value)
		}
	}
	paramsArr = append(paramsArr, salt)
	sort.Strings(paramsArr)
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(paramsArr, "&"))))
}
