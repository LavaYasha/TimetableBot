package modules

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func ConnectDataBase(arg string) (sqlx.DB, error){
	db, DBerr := sqlx.Connect("postgres", arg)
	if DBerr != nil {
		log.Fatalln(DBerr)
	}
	return *db, nil
}

func GetTimetable(db sqlx.DB) ([]TimeTable){
	_timeTableList := []TimeTable{}
	db.Select(&_timeTableList, "select S.subject_name as SubjectName, S.teacher as Teacher, T.start_time as Start, T.end_time as End, T.lecture_hall as LectureHall, T.subject_type as SubjectType, T.week_type as WeekType, T.week_day as Day from time_table T inner join subject S on T.subject_id = S.id ")
	return _timeTableList
}

func GetExclusion(db sqlx.DB) ([]Exclusion){
	ExclusionList := []Exclusion{}
	//db.Select(&ExclusionList, "select S.subject_name as SubjectName, S.teacher as Teacher, T.start_time as Start, T.end_time as End, T.lecture_hall as LectureHall, T.subject_type as SubjectType, T.week_type as WeekType, T.week_day as Day from time_table T inner join subject S on T.subject_id = S.id ")
	return ExclusionList
}