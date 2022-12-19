package main

import (
	"fmt"
	"github.com/shomali11/slacker"
	"github.com/joho/godotenv"
	"log"
	"os"
	"context"
	"encoding/json"
	witai "github.com/wit-ai/wit-go/v2"
	"github.com/tidwall/gjson"
	"github.com/krognol/go-wolfram"
)


func printCommandEvents(analyticChannel <-chan *slacker.CommandEvent){

	for event := range analyticChannel{
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Parameters)
		fmt.Println(event.Command)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main(){
	godotenv.Load(".env")

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))
	clientWit := witai.NewClient(os.Getenv("WIT_AI_TOKEN"))
	clientWolf := &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}

	go printCommandEvents(bot.CommandEvents())

	bot.Command("query - <message>", &slacker.CommandDefinition{
		Description: "Ask anything to wolfram",
		Examples: []string{"What is the color of water"},
		Handler: func(botCtx slacker.BotContext, r slacker.Request, w slacker.ResponseWriter){
			query := r.Param("message")
			
			msg, _ := clientWit.Parse(&witai.MessageRequest{
				Query: query,
			})

			data, _ := json.MarshalIndent(msg, "", "    ")
			strData := string(data[:])
			value := gjson.Get(strData, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			strValue := value.String()
			res, err := clientWolf.GetSpokentAnswerQuery(strValue, wolfram.Metric, 1000)
			if err != nil {
				fmt.Println("There was an error")
			}
			fmt.Println(strValue)
			w.Reply(res)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)

	if err != nil{
		log.Fatal(err)
	}

}