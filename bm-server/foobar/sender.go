package foobar

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

func InitClientIncoming() {
	// ClientCh = make(chan string)
	// go func(c chan string) {
	// 	uuid := <- c
	// 	processIncomingClientMessage(uuid)
	// }(ClientCh)

	logrus.Info("Stared incoming client channel")
}

func ProcessIncomingClientMessage(uuid string) {
	fmt.Printf("Started Processing message: %s\n", uuid)
	time.Sleep(5 * time.Second)
	go process(uuid)
	fmt.Printf("Finished Processing message: %s\n", uuid)
}

func process(uuid string) {

}
