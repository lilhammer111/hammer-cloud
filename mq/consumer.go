package mq

import "log"

var done chan bool

func StartConsume(qName, cName string, callback func(msg []byte) bool) {
	// 1. get a non-buffer channel
	messages, err := channel.Consume(qName, cName, true, false, false, false, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// 2. Cyclic acquisition of message channels
	done = make(chan bool)
	go func() {
		for msg := range messages {
			// 3.invoke to callback function
			suc := callback(msg.Body)
			if !suc {
				// todo send the task to another queue for retrying
			}

		}
	}()

	<-done
	err = channel.Close()
	if err != nil {
		log.Println(err)
		return
	}
}
