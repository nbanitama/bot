package main

import (
	"log"
	"os"

	"github.com/tokopedia/chatbot-scheduler/task"
)

func main() {
	if len(os.Args) == 2 {
		log.Println("Starting chatbot-scheduler")
		argument := os.Args[1]
		if argument == "pic" {
			task.PopulateData()
		} else if argument == "new-intent" {
			task.PopulateIntentName()
		} else if argument == "dialog-flow-intent" {
			task.PopulateDialogFlowIntent()
		}
	} else {
		log.Println("chatbot-scheduler doesn't start!!")
		log.Println("please check argument...")
	}
}
