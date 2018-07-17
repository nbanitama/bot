package model

type ChatbotLog struct {
	DateTimestamp  string `json:"date_timestamp"`
	FromUid        string `json:"from_uid"`
	IntentName     string `json:"intent_name"`
	ChatfromUser   string `json:"chat_from_user"`
	Score          string `json:"score"`
	ChatfromBot    string `json:"chat_from_bot"`
	UserSays       string `json:"user_says"`
	ActualIntent   string `json:"actual_intent"`
	Status         int    `json:"status"`
	IsAdditionToDF bool   `json:"addition_to_df"`
	PIC            string `json:"pic"`
}

type AjaxForm struct {
	Draw            int          `json:"draw"`
	RecordsTotal    int          `json:"recordsTotal"`
	RecordsFiltered int          `json:"recordsFiltered"`
	Data            []ChatbotLog `json:"data"`
	Error           string       `json:"error"`
}
