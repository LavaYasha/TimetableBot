package botman

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres() *Postgres {
	return &Postgres{}
}

func (p *Postgres) Connect(arg string) (err error) {
	p.db, err = sqlx.Connect("postgres", arg)
	return err
}

func (p *Postgres) GetTimetable() ([]TimeTable, error) {
	_timeTableList := []TimeTable{}
	err := p.db.Select(&_timeTableList, "select S.subject_name as SubjectName, S.teacher as Teacher,T.lesson_num as LessonNum, lt.start_time as Start,lt.end_time as End,T.lecture_hall as LectureHall,T.subject_type as SubjectType,T.week_type as WeekType,T.week_day as Day from time_table T 	inner join subject S on T.subject_id = S.id inner join lesson_times lt on T.lesson_num = lt.lesson_num")
	return _timeTableList, err
}

func (p *Postgres) GetExclusion() ([]Exclusion, error) {
	ExclusionList := []Exclusion{}
	err := p.db.Select(&ExclusionList, "select S.subject_name as SubjectName, S.teacher as Teacher,T.lesson_num as LessonNum, lt.start_time as Start,lt.end_time as End,T.lecture_hall as LectureHall,T.subject_type as SubjectType,T.week_num as WeekNum,T.week_day as Day from exclusion T 	inner join subject S on T.subject_id = S.id inner join lesson_times lt on T.lesson_num = lt.lesson_num")

	return ExclusionList, err
}
