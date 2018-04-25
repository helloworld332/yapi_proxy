package main

import (
	"encoding/json"
	"flag"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	log "yapi_proxy/common/alog"
)

func loadConfig(path string) (map[string]interface{}, error) {

	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error("read file error: ", err.Error())
		return nil, err
	}

	var config map[string]interface{}
	if err = json.Unmarshal(file, &config); err != nil {
		log.Error("Unmarshal config error: ", err.Error())
		return nil, err
	}

	return config, nil
}

var gConfig map[string]interface{}

func main() {
	confFilePath := flag.String("conf_path", "../conf/config.json", "-conf_path=xxx the config file path.")
	flag.Parse()
	config, err := loadConfig(*confFilePath)
	if err != nil {
		log.Error(err.Error())
		return
	}
	gConfig = config

	listenPort := config["listen_port"].(string)

	log.Debug("Start http proxy serving on port:", listenPort)

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		//proxy.OnRequest(goproxy.UrlIs("/self/getssl")).DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

		if req.URL.Path != "/self/getssl" {
			return
		}
		log.Debug("get ssl")
		b, err := ioutil.ReadFile("../static/ca.cer")
		if err != nil {
			log.Fatal(err.Error())
		}

		//resp := goproxy.NewResponse(req, "application/octet-stream", http.StatusAccepted, string(b))

		userAgent := strings.ToLower(req.UserAgent())
		fileName := "ca.cer"
		if strings.Index(userAgent, "msie") != -1 { //fileName用urlencode编码后正常
			w.Header().Set("Content-Disposition", "attachment; filename="+url.QueryEscape(fileName))
		} else if strings.Index(userAgent, "firefox") != -1 { //正常
			w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		} else if strings.Index(userAgent, "chrome") != -1 { //不能显示文件名,文件名只显示"file",而且没有扩展名
			w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		} else {
			w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(b)

		//return nil, resp
	})
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {

		log.Debugf("old http host: %s, scheme: %s, path: %s, url.Host: %s", req.Host, req.URL.Scheme, req.URL.Path, req.URL.Host)
		newReq, ok := getNewRequest(req)
		if ok {
			req = newReq
			log.Debugf("new http host: %s, scheme: %s, path: %s, url.Host: %s", req.Host, req.URL.Scheme, req.URL.Path, req.URL.Host)
		}

		return req, nil
	})
	proxy.Verbose = true

	//http.Handle("/", http.StripPrefix("/static/", http.FileServer(http.Dir("../static"))))

	log.Fatal(http.ListenAndServe(":"+listenPort, proxy))

	//proxyHandler := handler.New(config)

	//log.Fatal(http.ListenAndServe(":"+listenPort, proxyHandler))
}

func getNewRequest(req *http.Request) (*http.Request, bool) {
	if gConfig == nil {
		log.Error("gConfig is nil")
		return req, false
	}

	urlList, ok := gConfig["list"].(map[string]interface{})
	if !ok {
		log.Error("gConfig>list is not map interface")
		return req, false
	}

	apiMap, ok := urlList[req.URL.Host].(map[string]interface{})
	if !ok {
		host := stripPort(req.URL.Host)
		apiMap, ok = urlList[host].(map[string]interface{})
		if !ok {
			return req, false
		}
	}

	newPath, ok := apiMap[req.URL.Path].(string)
	if !ok {
		return req, false
	}

	newUrl, err := url.Parse(newPath)
	if err != nil {
		log.Error(err.Error())
		return req, false
	}

	req.URL.Scheme = newUrl.Scheme
	req.URL.Host = newUrl.Host
	req.URL.Path = newUrl.Path
	req.Host = newUrl.Host

	return req, true
}

func stripPort(s string) string {
	ix := strings.IndexRune(s, ':')
	if ix == -1 {
		return s
	}
	return s[:ix]
}
