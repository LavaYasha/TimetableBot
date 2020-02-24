package botman

//https://api.telegram.org/bot<token>/METHOD_NAME
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Botman struct {
	token  string
	db     *Postgres
	apiUrl string
}

func NewBotman(botToken string, db *Postgres) *Botman {
	return &Botman{
		token:  botToken,
		db:     db,
		apiUrl: "https://api.telegram.org/bot" + botToken,
	}
}

func (b *Botman) Run() error {
	offset := 0

	const timeFormat = "2006.01.02-15.04.05"
	beginningDate, _ := time.Parse(timeFormat, "2020.02.10-00.00.00")

	for {
		updates, err := b.getUpdates(offset)
		if err != nil {
			fmt.Printf("can't get updates: %v", err)
		}
		for _, update := range updates {
			err = b.respond(update, beginningDate)
			if err != nil {
				return err
			}
			offset = update.UpdateId + 1
		}
		if len(updates) != 0 {
			fmt.Println(updates)
		}
	}
}

func (b *Botman) getUpdates(offset int) ([]Update, error) {
	resp, err := http.Get(b.apiUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
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

func (b *Botman) respond(update Update, beginningDate time.Time) error {
	const shortTimeFormat = "15:04"
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	messageText := update.Message.Text

	timeTableList, err := b.db.GetTimetable()
	if err != nil {
		return err
	}
	exclusion, err := b.db.GetExclusion()
	if err != nil {
		return err
	}
	sort(timeTableList)
	sort2(exclusion)

	for i := 0; i < len(timeTableList); i++ {
		startTime, err := time.Parse(time.RFC3339, timeTableList[i].Start)
		if err != nil{
			return err
		}
		endTime, err := time.Parse(time.RFC3339, timeTableList[i].End)
		if err != nil{
			return err
		}
		timeTableList[i].Start = startTime.Format(shortTimeFormat)
		timeTableList[i].End = endTime.Format(shortTimeFormat)
	}

	for i := 0; i < len(exclusion); i++ {
		startTime, err := time.Parse(time.RFC3339, exclusion[i].Start)
		if err != nil{
			return err
		}
		endTime, err := time.Parse(time.RFC3339, exclusion[i].End)
		if err != nil{
			return err
		}
		exclusion[i].Start = startTime.Format(shortTimeFormat)
		exclusion[i].End = endTime.Format(shortTimeFormat)
	}

	//log.Println(_timeTableList)
	//log.Println(_exclusion)
	now := time.Now()
	newMessage := ""
	var	massageArr []string
	subtime := now.Sub(beginningDate)
	passedTime := int(subtime.Hours() / 24.0)
	weekNumber := (passedTime / 7) + 1

	switch messageText {
	case "/start":
		{
			newMessage = "Простенький бот для расписания\n>набери \"/\" для начала работы"
			massageArr = append(massageArr, newMessage)
		}
		break
	case "Сегодня", "с", "С", "/1":
		{
			requaredDay := now
			newMessage += findLesson(timeTableList, exclusion, requaredDay, beginningDate)
			massageArr = append(massageArr, newMessage)
		}
		break
	case "Завтра", "з", "З", "/2":
		{
			requaredDay := now.AddDate(0, 0, 1)
			newMessage += findLesson(timeTableList, exclusion, requaredDay, beginningDate)
			massageArr = append(massageArr, newMessage)
		}
		break
	case "послезавтра", "после", "по", "п", "/3":
		{
			requaredDay := now.AddDate(0, 0, 2)
			newMessage += findLesson(timeTableList, exclusion, requaredDay, beginningDate)
			massageArr = append(massageArr, newMessage)
		}
		break
	case "Неделя", "/4":
		{
			requaredWeek := weekNumber
			massageArr = searchLessonToWeek(timeTableList, exclusion, requaredWeek, beginningDate)
		}
		break
	case "/5":
		{
			requaredWeek := weekNumber + 1
			massageArr = searchLessonToWeek(timeTableList, exclusion, requaredWeek, beginningDate)
		}
		break
	default:
		{
			newMessage = "Сила программирования!!!"
		}
		break
	}
	newMessage = "\n\n" + "Дата запроса: " + now.Format("2 Jan 2006 15:04:05")
	massageArr = append(massageArr, newMessage)
	for i := 0; i < len(massageArr); i++ {
		botMessage.Text = massageArr[i]
		//newBotMessage := ""
		//log.Println(now.Format("Mon Jan 02"), int(passedTime), biginingDate, now.Weekday(), weekType)

		buf, err := json.Marshal(botMessage)

		if err != nil {
			return err
		}
		_, err = http.Post(b.apiUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))

		if err != nil {
			return err
		}
	}
	return nil
}

func findLesson(timetableList []TimeTable, exclusion []Exclusion, requestingTime time.Time, beginingDate time.Time) (message string) {

	weekType := ""
	subtime := requestingTime.Sub(beginingDate)
	passedTime := int(subtime.Hours() / 24.0)
	weekNumber := (passedTime / 7) + 1
	message += "------------------------------------------------------------------------\n"
	message += "++++ " + requestingTime.Weekday().String() + " (" + requestingTime.Format("2 Jan 2006") + ") " + "++++" + "\n"
	message += "------------------------------------------------------------------------\n"
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
	message += "------------------------------------------------------------------------\n"
	thereIsLesson := false
	for j := 0; j < len(exclusion); j++ {
		if exclusion[j].Day == requestingTime.Weekday().String() && exclusion[j].WeekNum == weekNumber {
			thereIsLesson = true
			message += exclusion[j].SubjectName + "\n" + exclusion[j].Teacher + " (" + exclusion[j].Start + " - " + exclusion[j].End + ") "
			message += "\nАудитория: " + strconv.Itoa(exclusion[j].LectureHall) + "\n"
			message += exclusion[j].SubjectType + "\n"
			message += "------------------------------------------------------------------------\n"
		}
	}
	for i := 0; i < len(timetableList); i++ {
		if timetableList[i].Day == requestingTime.Weekday().String() && timetableList[i].WeekType == weekType {
			thereIsLesson = true
			message += timetableList[i].SubjectName + "\n" + timetableList[i].Teacher + " (" + timetableList[i].Start + " - " + timetableList[i].End + ") "
			message += "\nАудитория: " + strconv.Itoa(timetableList[i].LectureHall) + "\n"
			message += timetableList[i].SubjectType + "\n"
			message += "------------------------------------------------------------------------\n"
		}
	}

	if !thereIsLesson {
		message += "Пар нет, вот это да!\n"
		//message += "==================================\n"
		message += "------------------------------------------------------------------------\n"
	}
	return message
}

func searchLessonToWeek(timetableList []TimeTable, exclusion []Exclusion, requestingWeek int, beginingDate time.Time) (arrmessage []string) {
	message := ""
	requestingTime := beginingDate.AddDate(0, 0, 7*(requestingWeek-1))
	for i := 0; i < 7; i++ {
		message = findLesson(timetableList, exclusion, requestingTime, beginingDate)
		if i != 6 {
			message += "\n\n"
		}
		arrmessage = append(arrmessage, message)
		requestingTime = requestingTime.AddDate(0, 0, 1)
	}

	return arrmessage
}

func sort(list []TimeTable) {
	for i := 0; i < len(list); i++ {
		if i < len(list)-1 {
			if list[i].LessonNum > list[i+1].LessonNum {
				list[i].LessonNum, list[i+1].LessonNum = list[i+1].LessonNum, list[i].LessonNum
			}
		}
	}
}

func sort2(list []Exclusion) {
	for i := 0; i < len(list); i++ {
		if i < len(list)-1 {
			if list[i].LessonNum > list[i+1].LessonNum {
				list[i].LessonNum, list[i+1].LessonNum = list[i+1].LessonNum, list[i].LessonNum
			}
		}
	}
}
