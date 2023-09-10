package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DataBase struct {
}
type User struct {
	TgId int64
	Score int
	QuestNum int64
}
type Question struct {
	Id int64
	Desc string
	VariantA string
	VariantB string
	VariantC string
	VariantD string
	Correct int
	Answer string
}

var dbInfo string

func (d *DataBase) InitInfo(host, port,user,password,dbname,sslmode string) {
	dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
}

func (d *DataBase) GetQuestionById(id int64) (Question, error) {
	db, err := sql.Open("postgres",dbInfo)
	if err != nil {
		return Question{}, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT * FROM questions WHERE Id = $1`, id)
	if err != nil {
		return Question{}, err
	}
	var question Question
	for rows.Next(){
		err := rows.Scan(&question.Id, &question.Desc,&question.VariantA,&question.VariantB,&question.VariantC,&question.VariantD,&question.Correct,&question.Answer)
		if err != nil {
			return Question{}, err
		}
	}
	return question, nil
}

func (d *DataBase) GetUserById(id int64) (User, error) {
	db, err := sql.Open("postgres",dbInfo)
	if err != nil {
		return User{}, err
	}
	defer db.Close()
	rows, err := db.Query(`SELECT * FROM users WHERE TgId = $1`, id)
	if err != nil {
		return User{}, err
	}
	var user User
	for rows.Next(){
		err := rows.Scan(&user.TgId, &user.Score,&user.QuestNum)
		if err != nil {
			return User{}, err
		}
	}
	return user, nil
}
func (d *DataBase) CreateUser(u User)  error{
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
        return err
    }
    defer db.Close()

	data := `INSERT INTO users (TgId, Score, QuestNum) VALUES ($1, $2, $3)`
	if _, err = db.Exec(data,u.TgId, u.Score,u.QuestNum); err != nil {
		return err
	}
	return nil
}

func (d *DataBase) ChangeUserScore(id int64, newScore int) error{
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()
	data := `UPDATE users SET Score = $1 WHERE TgID = $2`
	if _, err = db.Exec(data,newScore, id); err != nil {
		return err
	}
	return nil
}
func (d *DataBase) ChangeUserQuestNum(id int64, newQuest int64) error{
	db, err := sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}
	defer db.Close()
	data := `UPDATE users SET QuestNum = $1 WHERE TgID = $2`
	if _, err = db.Exec(data,newQuest, id); err != nil {
		return err
	}
	return nil
}