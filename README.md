# wb-l0
скрипт для отправки сообщений в канал
```
package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	message := `{
  "order_uid": "1234566",
  "track_number": "WBILMTESTTRACK",
  "entry": "WBIL"
}`
	channel := "channel"

	err = nc.Publish(channel, []byte(message))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Отправлено сообщение на канал '%s': %s\n", channel, message)
}
```
