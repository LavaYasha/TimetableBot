package modules

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type TimeTable struct {
	SubjectName string
	Teacher     string
	Start       string
	End         string
	LectureHall int
	SubjectType string
	WeekType    string
	Day         string
}

type Exclusion struct {
SubjectName string
Teacher     string
Start       string
End         string
LectureHall int
SubjectType string
WeekNum 	int
Day         string
}
