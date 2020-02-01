package client

import (
	"log"
)

func logMain(s string) { //主进程中的log输出
	log.Println("Main Process->" + s)
}

func logSub(s string) { //子进程中的log输出
	log.Println("Sub Process->" + s)
}
