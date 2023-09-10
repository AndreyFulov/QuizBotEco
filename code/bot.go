package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot_token string
	db *DataBase
}

func NewBot(token string, db *DataBase) *TelegramBot{
	return &TelegramBot{
		bot_token: token,
		db: db,
	}
}

func (tg *TelegramBot) Bot() {
	bot, err := tgbotapi.NewBotAPI(tg.bot_token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates{
		if update.Message != nil {
			if (update.Message.Text == "/start") {
				var questnum int64
				questnum = 1
				user, err :=tg.db.GetUserById(update.Message.From.ID)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,"Ошибка! Попробуйте позже!")
					bot.Send(msg)
				}
				if user.TgId != 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,fmt.Sprintf("Вы уже начали, ваш счет %s ", strconv.Itoa(user.Score)))
					bot.Send(msg)
					questnum = user.QuestNum
				}else {
					u := User{
						TgId: update.Message.From.ID,
						Score : 0,
						QuestNum: 1,
					}
					tg.db.CreateUser(u)
					time.Sleep(1 * time.Second)
				}
				question, err := tg.db.GetQuestionById(questnum)
				msg2 := tgbotapi.NewMessage(update.Message.Chat.ID,fmt.Sprintf("💬Вопрос №%s:\n\n%s\n✅Варианты ответов:\n1.%s\n2.%s\n3.%s\n4.%s", strconv.Itoa(int(user.QuestNum)),question.Desc, question.VariantA,question.VariantB,question.VariantC,question.VariantD))
				bot.Send(msg2)
			}
			if(update.Message.Text == "1" || update.Message.Text == "2" || update.Message.Text == "3" || update.Message.Text == "4") {
				a, err := strconv.Atoi(update.Message.Text)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID,"Ошибка! Попробуйте позже!")
					bot.Send(msg)
				}else {
					user, err :=tg.db.GetUserById(update.Message.From.ID)
					if err != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID,"Ошибка! Попробуйте позже!")
						bot.Send(msg)
					}else {
						if user.TgId != 0 {
							question, err := tg.db.GetQuestionById(user.QuestNum)
							if err != nil {
								msg := tgbotapi.NewMessage(update.Message.Chat.ID,"Ошибка! Попробуйте позже!")
								bot.Send(msg)
							}else {
								if question.Correct == a {
									
									tg.db.ChangeUserScore(user.TgId, user.Score + 1)
									msg := tgbotapi.NewMessage(update.Message.Chat.ID,fmt.Sprintf("✅Верно!\n\n📄Объяснение:\n%s\n\nВаш счет:%s",question.Answer, strconv.Itoa(user.Score + 1)))
									bot.Send(msg)
									tg.db.ChangeUserQuestNum(user.TgId, user.QuestNum + 1)
									question, err := tg.db.GetQuestionById(user.QuestNum + 1)
									if err != nil {

									}
									msg2 := tgbotapi.NewMessage(update.Message.Chat.ID,fmt.Sprintf("💬Вопрос №%s:\n\n%s\n✅Варианты ответов:\n1.%s\n2.%s\n3.%s\n4.%s", strconv.Itoa(int(user.QuestNum + 1)),question.Desc, question.VariantA,question.VariantB,question.VariantC,question.VariantD))
									bot.Send(msg2)
								}else {
									msg := tgbotapi.NewMessage(update.Message.Chat.ID,fmt.Sprintf("❌Неверно! Правильный вариант:%s\n\n📄Объяснение:\n%s",strconv.Itoa(question.Correct),question.Answer))
									bot.Send(msg)
									tg.db.ChangeUserQuestNum(user.TgId, user.QuestNum + 1)
									question, err := tg.db.GetQuestionById(user.QuestNum + 1)
									if err != nil {

									}
									msg2 := tgbotapi.NewMessage(update.Message.Chat.ID,fmt.Sprintf("💬Вопрос №%s:\n\n%s\n✅Варианты ответов:\n1.%s\n2.%s\n3.%s\n4.%s", strconv.Itoa(int(user.QuestNum + 1)),question.Desc, question.VariantA,question.VariantB,question.VariantC,question.VariantD))
									bot.Send(msg2)
								}
							}
						}else {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID,"Вы не зарегестрированы, напишите /start")
							bot.Send(msg)
						}
				}
				}

			}
		}
	}
}