package main

import (
	"fmt"
	"os/exec"

    "io/ioutil"
    "net/http"
    "strings"
	//"encoding/json"

	//"strconv"


	"PushServer/slog"
	"PushServer/util"
)

func getSinceGetui() {
	// 获取之前已经获取到的getuicliends
	data, err := util.GetFile(SinceGetuiFile)
	if err != nil {
		slog.Errorln("getSinceGetui get getui cid err", err)
	}

	cids := strings.Split(string(data), "\n")
	for _, c := range(cids) {
		slog.Infoln("getSinceGetui", c)
		if c != "" {
			if _, ok := sendManGt[c]; !ok {
				sendManGt[c] = &sendState {
					count: 0,
					rvLast: "",
				}
			}
		}
	}

	slog.Infof("getSinceGetui getui sendManGt len:%d", len(sendManGt))


}

func restGetuipush(clientid string, data []byte) string {
	fun := "restGetuipush"
	slog.Infof("%s cid:%s len:%d data:%s", fun, clientid, len(data), data)
	out, err := gettuipush(clientid, data)
	if err != nil {
		return fmt.Sprintf("%s", err)
	} else {
		return fmt.Sprintf("%s", out)
	}


}

func gettuipush(clientid string, data []byte) ([]byte, error){
	out, err := exec.Command("python", "/home/fengguangxiang/shawn_tool/edjtool/getui/push.py", clientid, string(data)).CombinedOutput()
	//out, err := exec.Command("date").CombinedOutput()
	if err != nil {
		//log.Fatal(err)
		//log.Println("ERR", err)
		return []byte(""), err
	} else {
		return out, err
	}

}


// Method: GET
// Uri: /ack/taskid
func ack(w http.ResponseWriter, r *http.Request) {
	fun := "rest.ack"
	if r.Method != "GET" {
		//writeRestErr(w, "method err")
		http.Error(w, "method err", 405)
		return
	}

	slog.Infof("%s %s", fun, r.URL.Path)
	path := strings.Split(r.URL.Path, "/")
	//slog.Info("%q", path)

	if len(path) != 3 {
		//writeRestErr(w, "uri err")
		http.Error(w, "uri invalid", 400)
		return
	}

	// path[0] "", path[1] push
	taskid := path[2]

	slog.Infof("%s taskid:%s", fun, taskid)


}


// Method: POST
// Uri: /push/CLIENT_ID
func sub(w http.ResponseWriter, r *http.Request) {
	fun := "rest.sub"
	//debug_show_request(r)
	if r.Method != "POST" {
		//writeRestErr(w, "method err")
		http.Error(w, "method err", 405)
		return
	}

	slog.Infof("%s %s", fun, r.URL.Path)
	path := strings.Split(r.URL.Path, "/")
	//slog.Info("%q", path)

	if len(path) != 3 {
		//writeRestErr(w, "uri err")
		http.Error(w, "uri invalid", 400)
		return
	}

	// path[0] "", path[1] push
	clientid := path[2]


	if _, ok := sendManGt[clientid]; !ok {
		sendManGt[clientid] = &sendState {
			count: 0,
			rvLast: "",
		}

	}


	slog.Infof("sub getui sendManGt len:%d", len(sendManGt))

	cids := make([]string, 0)
	for c, _ := range(sendManGt) {
		cids = append(cids, c)
	}

	d1 := strings.Join(cids, "\n")
    err := ioutil.WriteFile(SinceGetuiFile, []byte(d1), 0644)

	if err != nil {
		slog.Infof("sub getui write file err:%s", err)

	}



}


// Method: POST
// Uri: /push/CLIENT_ID
// Data: push data
func push(w http.ResponseWriter, r *http.Request) {
	fun := "rest.push"
	//debug_show_request(r)
	if r.Method != "POST" {
		//writeRestErr(w, "method err")
		http.Error(w, "method err", 405)
		return
	}

	slog.Infof("%s %s", fun, r.URL.Path)
	path := strings.Split(r.URL.Path, "/")
	//slog.Info("%q", path)

	if len(path) != 3 {
		//writeRestErr(w, "uri err")
		http.Error(w, "uri invalid", 400)
		return
	}

	// path[0] "", path[1] push
	clientid := path[2]


	data, err := ioutil.ReadAll(r.Body);
	if err != nil {
		//writeRestErr(w, "data err")
		er := fmt.Sprintf("body read err:%s", err)
		http.Error(w, er, 501)
		return
	}

	if len(data) == 0 {
		//writeRestErr(w, "data empty")
		http.Error(w, "data empty", 400)
		return
	}


	out, err := gettuipush(clientid, data)
	if err != nil {
		fmt.Fprintf(w, "%s", err)
	} else {
		fmt.Fprintf(w, "%s", out)
	}

	//slog.Debugf("%s msgid:%d link:%s", fun, msgid, link)
	//js, _ := json.Marshal(&RestReturn{Msgid: msgid, Link: link})
	//fmt.Fprintf(w, "%s", js)


}


func StartHttp(httpport string) {
	http.HandleFunc("/getui/", push)
	http.HandleFunc("/getuisub/", sub)
	http.HandleFunc("/getuiack/", ack)

	err := http.ListenAndServe(httpport, nil) //设置监听的端口
	if err != nil {
		slog.Panicf("StartHttp ListenAndServe: %s", err)
	}

}


