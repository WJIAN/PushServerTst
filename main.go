package main

import (
	"PushServer/slog"
	"PushServer/rediscluster"
)


func main() {
	slog.Init("")
	redisPool = rediscluster.NewRedisPool()

	slog.Infoln("Start loop push")
	go loopPush()


	slog.Infoln("Start http")
	StartHttp(httpAddr)
	

}
