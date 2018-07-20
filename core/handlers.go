package core

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Page struct {
	data string
}

type ChatbotLog struct {
	DateTimestamp      string `json:"date_timestamp"`
	FromUid            string `json:"from_uid"`
	IntentName         string `json:"intent_name"`
	ChatfromUser       string `json:"chat_from_user"`
	Score              string `json:"score"`
	ChatfromBot        string `json:"chat_from_bot"`
	UserSays           string `json:"user_says"`
	ActualIntent       string `json:"actual_intent"`
	Status             int    `json:"status"`
	IsAdditionToDF     bool   `json:"addition_to_df"`
	PIC                string `json:"pic"`
	HashId             string `json:"hash_id"`
	SuggestedNewIntent string `json:"suggested_new_intent"`
}

type AjaxForm struct {
	Draw            int          `json:"draw"`
	RecordsTotal    int          `json:"recordsTotal"`
	RecordsFiltered int          `json:"recordsFiltered"`
	Data            []ChatbotLog `json:"data"`
	Error           string       `json:"error"`
}

type Response struct {
	ResponseCode string `json:"response_status"`
}

type IntentData struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

type IntentSelect struct {
	Results    []IntentData `json:"results"`
	TotalCount int          `json:"total_count"`
}

func (c *TaskModule) HandlerPopulateData(w http.ResponseWriter, r *http.Request) {
	log.Println("handling populating data")
}

func (c *TaskModule) HandlerPopulateIntent(w http.ResponseWriter, r *http.Request) {
	log.Println("handling populating intent")
}

func (c *TaskModule) HandlerGetDialogFlow(w http.ResponseWriter, r *http.Request) {
	log.Println("handling getting dialog flow")
}

func (c *TaskModule) HandlerShowForm(w http.ResponseWriter, r *http.Request) {
	log.Println("handling showing form")

	root, err := os.Getwd()
	if err != nil {
		log.Println(err)
		log.Println("Error when getting get working directory...")
		return
	}

	data := Page{"nothing"}
	filepath := root + "/resources/html/form.html"
	t, _ := template.ParseFiles(filepath)
	t.Execute(w, data)
}

func (c *TaskModule) HandlerFormDatatables(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling querying data from datatable..")

	queryValues := r.URL.Query()

	drawStr := queryValues.Get("draw")
	draw, _ := strconv.Atoi(drawStr)
	startStr := queryValues.Get("start")
	lengthStr := queryValues.Get("length")
	search := strings.Replace(strings.ToLower(queryValues.Get("search[value]")), "_", "\\_", -1)
	columnStr := queryValues.Get("order[0][column]")
	column, _ := strconv.Atoi(columnStr)
	orderStr := queryValues.Get("order[0][dir]")

	var isFilter bool
	if search == "" {
		isFilter = false
	} else {
		isFilter = true
	}

	totalRecord, err := getCount("", false)
	if err != nil {

		log.Fatal(err)
		return
	}

	totalFilteredRecord, err := getCount(search, isFilter)
	if err != nil {
		log.Fatal(err)
		return
	}
	length, _ := strconv.Atoi(lengthStr)
	start, _ := strconv.Atoi(startStr)
	var dataList []ChatbotLog
	if totalFilteredRecord-start < length {
		diff := (totalFilteredRecord - start)
		diffStr := strconv.Itoa(diff)
		dataList, err = getList(search, startStr, diffStr, column, orderStr)
	} else {
		dataList, err = getList(search, startStr, lengthStr, column, orderStr)
	}

	if err != nil {
		log.Println(err)
		return
	}

	resp := AjaxForm{draw, totalRecord, totalFilteredRecord, dataList, ""}
	data, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Write(data)
}

func (c *TaskModule) HandlerGetFormAjax(w http.ResponseWriter, r *http.Request) {
	log.Println("handling getting form ajax")
}

func (c *TaskModule) HandlerPostFormAjax(w http.ResponseWriter, r *http.Request) {
	log.Println("handling posting form ajax")
	var payload ChatbotLog

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println(err)
	}

	_, err = postData(&payload)

	_, err = addSuggestIntentData(&payload)

	var respCode string
	if err != nil {
		log.Println(err)
		respCode = "404"
	} else {
		respCode = "200"
	}

	resp := Response{respCode}
	response, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Write(response)
}

func (c *TaskModule) HandlerGetIntentAjax(w http.ResponseWriter, r *http.Request) {
	log.Println("handling getting intent ajax")

	queryValues := r.URL.Query()

	search := queryValues.Get("term")

	data, length, err := getIntentList(search)

	resp := IntentSelect{data, length}
	response, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Write(response)
}

func (c *TaskModule) HandlerGetSuggestIntentAjax(w http.ResponseWriter, r *http.Request) {
	log.Println("handling getting suggest-intent ajax")

	queryValues := r.URL.Query()

	search := queryValues.Get("term")

	data, length, err := getSuggestIntentList(search)

	resp := IntentSelect{data, length}
	response, err := json.Marshal(resp)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Write(response)
}

func postData(data *ChatbotLog) (sql.Result, error) {
	query := "update topbot_ops_chat_log " +
		"set an_user_says = $1, an_actual_intent_name=$2, an_new_intent_name = $3 " +
		"where hash_id = $4"

	db := postgresConnection.GetConnection()

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return stmt.Exec(data.UserSays, data.ActualIntent, data.SuggestedNewIntent, data.HashId)
}

func addSuggestIntentData(data *ChatbotLog) (sql.Result, error) {
	query := "INSERT INTO topbot_suggested_intent_list(created_time, updated_time, suggested_intent_name) VALUES (now(), now(), $1)"

	db := postgresConnection.GetConnection()

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println(data.SuggestedNewIntent)
	return stmt.Exec(data.SuggestedNewIntent)
}

func getCount(search string, isFiltered bool) (int, error) {
	query := "SELECT count(1) FROM topbot_ops_chat_log "
	if isFiltered {
		query += "where lower(intent_name) like '%" + search + "%' or lower(resolved_query) LIKE '%" + search + "%' or lower(coalesce(parsed_message, '')) LIKE '%" + search + "%' or lower(coalesce(an_pic, '')) LIKE '%" + search + "%' or lower(from_uid) LIKE '%" + search + "%'"
	}

	var result int
	err = postgresConnection.ExecuteQueryInt(query, &result)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return result, nil
}

func getIntentCount(search string) (int, error) {
	query := "SELECT count(1) FROM topbot_intent_list where lower(intent_name) like '%" + search + "%' "

	var result int
	err = postgresConnection.ExecuteQueryInt(query, &result)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return result, nil
}

func getSuggestIntentCount(search string) (int, error) {
	query := "SELECT count(1) FROM topbot_suggested_intent_list where lower(suggested_intent_name) like '%" + search + "%' "

	var result int
	err = postgresConnection.ExecuteQueryInt(query, &result)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return result, nil
}

func getSuggestIntentList(search string) ([]IntentData, int, error) {
	length, err := getSuggestIntentCount(search)
	if err != nil {
		log.Println(err)
		return make([]IntentData, 0), 0, err
	}
	intents := make([]IntentData, length)
	query := "SELECT suggested_intent_name, suggested_intent_name FROM topbot_suggested_intent_list where lower(suggested_intent_name) like '%" + search + "%'"

	rows, err := postgresConnection.ExecuteQuery(query)

	if err != nil {
		log.Println("Error occurred when querying db")
		return nil, 0, err
	} else {
		defer rows.Close()
	}

	i := 0
	for rows.Next() {
		c := IntentData{}
		err := rows.Scan(&c.Id, &c.Text)

		if err != nil {
			log.Println(err)
			return nil, 0, err
		}

		intents[i] = c
		i++
	}

	return intents, length, nil
}

func getIntentList(search string) ([]IntentData, int, error) {
	length, err := getIntentCount(search)
	if err != nil {
		log.Println(err)
		return make([]IntentData, 0), 0, err
	}
	intents := make([]IntentData, length)
	query := "SELECT intent_name, intent_name FROM topbot_intent_list where lower(intent_name) like '%" + search + "%'"

	rows, err := postgresConnection.ExecuteQuery(query)

	if err != nil {
		log.Println("Error occurred when querying db")
		return nil, 0, err
	} else {
		defer rows.Close()
	}

	i := 0
	for rows.Next() {
		c := IntentData{}
		err := rows.Scan(&c.Id, &c.Text)

		if err != nil {
			log.Println(err)
			return nil, 0, err
		}

		intents[i] = c
		i++
	}

	return intents, length, nil
}

func getList(search string, startStr string, lengthStr string, column int, order string) ([]ChatbotLog, error) {
	length, _ := strconv.Atoi(lengthStr)

	dataList := make([]ChatbotLog, length)

	query := "SELECT hash_id, from_uid, intent_name, score, resolved_query, coalesce(parsed_message, ''), coalesce(an_user_says, '')," +
		"coalesce(an_actual_intent_name, ''), coalesce(an_status, 1), coalesce(an_add_to_df, true), coalesce(an_pic, ''), coalesce(an_new_intent_name,''), coalesce(bot_timestamp,'') FROM topbot_ops_chat_log " +
		"WHERE lower(intent_name) LIKE '%" + search + "%' or lower(resolved_query) LIKE '%" + search + "%' or lower(coalesce(parsed_message, '')) LIKE '%" + search + "%' or lower(coalesce(an_pic, '')) LIKE '%" + search + "%' or lower(from_uid) LIKE '%" + search + "%'" +
		"ORDER BY "

	if column == 0 {
		query += " bot_timestamp "
	} else if column == 1 {
		query += " from_uid "
	} else if column == 2 {
		query += " intent_name "
	} else if column == 3 {
		query += " resolved_query "
	} else if column == 4 {
		query += " score "
	} else if column == 5 {
		query += " parsed_message "
	} else if column == 6 {
		query += " an_user_says "
	} else if column == 7 {
		query += " an_actual_intent_name "
	} else if column == 8 {
		query += " an_new_intent_name "
	} else if column == 9 {
		query += " an_status "
	} else if column == 10 {
		query += " an_add_to_df "
	} else if column == 11 {
		query += " an_pic "
	}

	query += " " + order + " LIMIT " + lengthStr + " OFFSET " + startStr

	log.Println("query : " + query)
	rows, err := postgresConnection.ExecuteQuery(query)

	if err != nil {
		log.Println("Error occurred when querying db")
		return nil, err
	} else {
		defer rows.Close()
	}

	i := 0
	for rows.Next() {
		c := ChatbotLog{}
		err := rows.Scan(&c.HashId, &c.FromUid, &c.IntentName, &c.Score, &c.ChatfromUser, &c.ChatfromBot,
			&c.UserSays, &c.ActualIntent, &c.Status, &c.IsAdditionToDF, &c.PIC, &c.SuggestedNewIntent, &c.DateTimestamp)

		if err != nil {
			return nil, err
		}
		//log.Printf("%+v", c)

		dataList[i] = c
		i++
	}

	return dataList, nil
}

func getDataList(search string) []ChatbotLog {
	data := make([]ChatbotLog, 3)

	if search == "" {
		data[0] = ChatbotLog{
			DateTimestamp:  "12 January 2018",
			FromUid:        "123456",
			IntentName:     "Noba",
			ChatfromUser:   "Noba 12",
			Score:          "100",
			ChatfromBot:    "coba kontak novan",
			UserSays:       "Ini masalah",
			ActualIntent:   "Actual yaa",
			Status:         3,
			IsAdditionToDF: true,
			PIC:            "Novando",
			HashId:         "123456",
		}

		data[1] = ChatbotLog{
			DateTimestamp:  "12 February 2018",
			FromUid:        "123456",
			IntentName:     "Noooba",
			ChatfromUser:   "Noba 2",
			Score:          "100",
			ChatfromBot:    "coba kontak novan",
			UserSays:       "Ini masalah",
			ActualIntent:   "Actual yaa",
			Status:         3,
			IsAdditionToDF: true,
			PIC:            "Novao",
			HashId:         "56789",
		}
		data[2] = ChatbotLog{
			DateTimestamp:  "12 March 2018",
			FromUid:        "123456",
			IntentName:     "Noba Noba",
			ChatfromUser:   "Noba 345",
			Score:          "100",
			ChatfromBot:    "coba kontak novan",
			UserSays:       "Ini masalah",
			ActualIntent:   "Actual yaa",
			Status:         3,
			IsAdditionToDF: true,
			PIC:            "Banitama",
			HashId:         "98765",
		}
	} else {
		data[0] = ChatbotLog{
			DateTimestamp:  "12 May 2018",
			FromUid:        "12345236",
			IntentName:     "banasa",
			ChatfromUser:   "Noba 12dsdsd",
			Score:          "100",
			ChatfromBot:    "saas kontak novan",
			UserSays:       "Inisdsd masalah",
			ActualIntent:   "Actual yaa",
			Status:         3,
			IsAdditionToDF: true,
			PIC:            "Novando",
			HashId:         "111111111111",
		}

		data[1] = ChatbotLog{
			DateTimestamp:  "12 December 2018",
			FromUid:        "123456",
			IntentName:     "resqs",
			ChatfromUser:   "Noba 2",
			Score:          "100",
			ChatfromBot:    "coba kontak novan",
			UserSays:       "Ini masalah",
			ActualIntent:   "Actual yaa",
			Status:         3,
			IsAdditionToDF: true,
			PIC:            "Novao",
			HashId:         "88888",
		}
		data[2] = ChatbotLog{
			DateTimestamp:  "12 July 2018",
			FromUid:        "1234590",
			IntentName:     "alesz",
			ChatfromUser:   "Noba 345",
			Score:          "100",
			ChatfromBot:    "coba kontak novan",
			UserSays:       "Ini masalah",
			ActualIntent:   "Actual yaa",
			Status:         3,
			IsAdditionToDF: true,
			PIC:            "Banitama",
			HashId:         "44444444",
		}
	}
	return data
}
