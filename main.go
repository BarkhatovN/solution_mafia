package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func getRoles(players int) []string {
	var roles []string
	if players == 8 {
		roles = []string{"Мафия", "Дон", "Шериф", "Мирный житель", "Мирный житель", "Мирный житель", "Мирный житель", "Мирный житель"}
	} else if players == 9 {
		roles = []string{"Мафия", "Мафия", "Дон", "Шериф", "Мирный житель", "Мирный житель", "Мирный житель", "Мирный житель", "Мирный житель"}
	} else if players == 10 {
		roles = []string{"Мафия", "Мафия", "Дон", "Шериф", "Мирный житель", "Мирный житель", "Мирный житель", "Мирный житель", "Мирный житель", "Мирный житель"}
	}
	return roles
}

func shuffleRoles(players int) []string {
	roles := getRoles(players)
	if roles == nil {
		return nil
	}
	shuffled := make([]string, len(roles))
	copy(shuffled, roles)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})
	return shuffled
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Health check passed"))
}

func main() {
	// Получение переменных окружения
	botToken := os.Getenv("TELEGRAM_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_TOKEN is not set")
	}

	telegramURL := os.Getenv("TELEGRAM_URL")
	if telegramURL == "" {
		log.Fatal("TELEGRAM_URL is not set")
	}
	bot, err := tgbotapi.NewBotAPI(botToken)

	wh, err := tgbotapi.NewWebhook(telegramURL + "/" + botToken)
	_, err = bot.Request(wh)
	if err != nil {
		log.Fatal(err)
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}

	// Запуск прослушивания обновлений в горутине
	go func() {
		updates := bot.ListenForWebhook("/" + bot.Token)

		for update := range updates {
			log.Printf("%+v\n", update)

			if update.Message != nil {
				// Обработка сообщений
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
				if update.Message.Text == "/start" {
					msg.Text = "Привет! Введи количество игроков (от 8 до 10), и я раздам роли."
				} else {
					players := 0
					_, err := fmt.Sscanf(update.Message.Text, "%d", &players)
					if err != nil || players < 8 || players > 10 {
						msg.Text = "Введите корректное число игроков (от 8 до 10)."
					} else {
						assignedRoles := shuffleRoles(players)
						responseText := "Роли распределены:\n"
						for i, role := range assignedRoles {
							responseText += fmt.Sprintf("Игрок %d: %s\n", i+1, role)
						}
						msg.Text = responseText
					}
				}
				bot.Send(msg)
			}
		}
	}()

	// Health Check endpoint для Cloud Run
	http.HandleFunc("/health", healthCheck)

	// Получаем порт из переменной окружения или по умолчанию используем 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s...\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
