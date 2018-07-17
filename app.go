package main

import (
	"log"
	"net/http"

	config "github.com/tokopedia/chatbot-scheduler/config/chatbot"
	"github.com/tokopedia/chatbot-scheduler/core"
	logging "gopkg.in/tokopedia/logging.v1"
)

func main() {
	logging.LogInit()
	config.NewMainConfig()
	log.Printf("%+v", config.Main)

	taskModule, err := core.NewTaskModule(&config.Main)
	if err == nil {
		http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("resources/js"))))

		http.HandleFunc("/populate/data", taskModule.HandlerPopulateData)
		http.HandleFunc("/populate/intent", taskModule.HandlerPopulateIntent)
		http.HandleFunc("/populate/dialog-flow-intent", taskModule.HandlerGetDialogFlow)
		http.HandleFunc("/form", taskModule.HandlerShowForm)
		http.HandleFunc("/form/ajax_get", taskModule.HandlerGetFormAjax)
		http.HandleFunc("/form/ajax_post", taskModule.HandlerPostFormAjax)
		http.HandleFunc("/form/data", taskModule.HandlerFormDatatables)

		log.Println("Starting the application...")
		http.ListenAndServe(":8787", nil)
	} else {
		log.Println("Failed to run the app...")
	}
}
