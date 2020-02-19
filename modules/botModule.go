package modules
//https://api.telegram.org/bot<token>/METHOD_NAME
import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func StartBot(botToken string, db sqlx.DB) {
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0

	const timeFormat = "2006.01.02-15.04.05"
	beginingDate, _ := time.Parse(timeFormat, "2020.02.10-00.00.00")

	for {

		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("err: ", err.Error())
		}
		for _, update := range updates {
			err = respond(botUrl, update, db, beginingDate)
			offset = update.UpdateId + 1
		}
		fmt.Println(updates)
	}
}

func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func respond(botUrl string, update Update, db sqlx.DB, beginingDate time.Time) error {
	const shortTimeFormat = "15:04"
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	messageText := update.Message.Text

	_timeTableList := GetTimetable(db)
	_exclusion := GetExclusion(db)

	for i := 0; i < len(_timeTableList); i++ {
		startTime, _ := time.Parse(time.RFC3339, _timeTableList[i].Start)
		endTime, _ := time.Parse(time.RFC3339, _timeTableList[i].End)
		_timeTableList[i].Start = startTime.Format(shortTimeFormat)
		_timeTableList[i].End = endTime.Format(shortTimeFormat)
	}

	//log.Println(_timeTableList)
	now := time.Now()
	newMessage := ""

	switch messageText {
	case "/start":
		{
			newMessage = "Простенький бот для расписания\n>набери \"/\" для начала работы"
		}
		break
	case "Сегодня", "с", "С", "/1":
		{
			requaredDay := now
			newMessage += searchLesson(_timeTableList, _exclusion,requaredDay, beginingDate)
		}
		break
	case "Завтра", "з", "З", "/2":
		{
			requaredDay := now.AddDate(0, 0, 1)
			newMessage += searchLesson(_timeTableList, _exclusion,requaredDay, beginingDate)
		}
		break
	case "послезавтра", "после", "по", "п", "/3":
		{
			requaredDay := now.AddDate(0, 0, 2)
			newMessage += searchLesson(_timeTableList, _exclusion,requaredDay, beginingDate)
		}
		break
	default:
		{
			newMessage = "Сила программирования!!!"
		}
		break
	}
	botMessage.Text = newMessage
	//newBotMessage := ""
	//log.Println(now.Format("Mon Jan 02"), int(passedTime), biginingDate, now.Weekday(), weekType)

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

func searchLesson(timetableList []TimeTable, exclusion []Exclusion, requestingTime time.Time, beginingDate time.Time) (message string) {

	weekType := ""
	subtime := requestingTime.Sub(beginingDate)
	passedTime := int(subtime.Hours() / 24.0)
	weekNumber := (passedTime / 7) + 1
	message += "==================================\n"
	if weekNumber%2 == 0 {
		weekType = "even"
	} else {
		weekType = "odd"
	}
	message += strconv.Itoa(weekNumber) + " Неделя"
	if weekType == "even" {
		message += " (чётная)\n"
	} else {
		message += " (нечётная)\n"
	}
	message += "----------------------------------\n"
	thereIsLesson := false
	for i := 0; i < len(timetableList); i++ {
		if timetableList[i].Day == requestingTime.Weekday().String() && timetableList[i].WeekType == weekType {
			thereIsLesson = true
			message += timetableList[i].SubjectName + "\n" + timetableList[i].Teacher + " (" + timetableList[i].Start + " - " + timetableList[i].End + ") "
			message += "\nАудитория: " + strconv.Itoa(timetableList[i].LectureHall) + "\n"
			message += "--------------------------------------------------------------------\n"
		}
	}
	if !thereIsLesson {
		message += "Пар нет, вот это да!\n"
		message += "==================================\n"
	}
	return message
}
