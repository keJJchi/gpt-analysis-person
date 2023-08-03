package main

import (
        "kejjchibot/bot"
        "os"
	"fmt"
)

func main() {
        bot_token, ok := os.LookupEnv("TG_KEY")
        if !ok {
		fmt.Println("missing TG_KEY")
		return
        }
	bot.StartBot(bot_token)
}
