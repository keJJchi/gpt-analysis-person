package bot

import (
	"fmt"
	"log"
	"strings"

	gpt "kejjchibot/chatgpt"
	"kejjchibot/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocolly/colly"
)

var clientDataMap map[int64]*utils.ClientData // Карта для хранения данных клиентов

func StartBot(token string) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	clientDataMap = make(map[int64]*utils.ClientData) // Инициализируем карту перед использованием

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // игнорируем все не-сообщения
			continue
		}

		HandleMessage(bot, update.Message)
	}
}

func HandleMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message) {
	if message.IsCommand() {
		switch message.Command() {
		case "start":
			clientData := &utils.ClientData{}
			clientDataMap[message.Chat.ID] = clientData

			msg := tgbotapi.NewMessage(message.Chat.ID, "Привет! Я бот, рад тебя видеть! Кто наша цель?")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}
		case "stop":
			msg := tgbotapi.NewMessage(message.Chat.ID, "Программа остановлена.")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println(err)
			}

			// Удаляем данные клиента из карты
			delete(clientDataMap, message.Chat.ID)
		}
	} else {
		clientData := clientDataMap[message.Chat.ID]
		if clientData != nil {
			switch {
			case clientData.Name == "":
				clientData.Name = message.Text
				msg := tgbotapi.NewMessage(message.Chat.ID, "Из какой он страны?")
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			case clientData.Country == "":
				clientData.Country = message.Text
				msg := tgbotapi.NewMessage(message.Chat.ID, "Какие ссылки посмотреть?")
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			case clientData.Webside == "":
				clientData.Webside = message.Text
				msg := tgbotapi.NewMessage(message.Chat.ID, "В какой компании он работает?")
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}
			case clientData.Company == "":
				clientData.Company = message.Text

				// Здесь можешь выполнить действия с полученными данными
				// Например, отправить сообщение с общей информацией о клиенте
				response := "Спасибо за информацию!\n\n"
				response += "Имя: " + clientData.Name + "\n"
				response += "Страна: " + clientData.Country + "\n"
				response += "Компания: " + clientData.Company + "\n"
				response += "Ссылки: " + clientData.Webside + "\n"
				msg := tgbotapi.NewMessage(message.Chat.ID, response)
				_, err := bot.Send(msg)
				if err != nil {
					log.Println(err)
				}

				urls := strings.Split(clientData.Webside, ",")
				// Создаем словарь для хранения текста с каждой ссылки
				// texts := make(map[string]string)

				var responseHistory string
				var isError bool // Флаг для отслеживания ошибки

				for _, url := range urls {
					url = strings.TrimSpace(url) // Удаляем пробелы перед ссылкой
					texts, err := GetCleanTextFromURL(url)
					if err != nil {
						// Отправляем сообщение об ошибке в чат
						errorMsg := fmt.Sprintf("Ошибка при получении текста с веб-сайта %s: %v", url, err)
						msgErr := tgbotapi.NewMessage(message.Chat.ID, errorMsg)
						_, err := bot.Send(msgErr)
						if err != nil {
							log.Println(err)
						}
						// continue
						isError = true // Устанавливаем флаг ошибки
						break          // Прерываем цикл при первой ошибке
					}

					for _, text := range texts {
						// texts[url] = text
						clientData.Content = text
						// clientData.Texts[url] = text
						fmt.Println(clientData.Content)

						// Связь с GPT-моделью для получения ответа
						responseFromGpt, err := gpt.ChatGpt(clientData)
						if err != nil {
							log.Println(err)
							continue
						}

						// Добавляем новый ответ к истории ответов
						responseHistory += responseFromGpt

						// Очищаем clientData.Content перед обработкой следующей части текста
						fmt.Println(responseHistory)
						clientData.Content = ""
					}
				}

				if !isError {
					// После завершения цикла отправляем запрос к GPT, принимая на вход только ответы с предыдущих запросов
					clientData.Content = responseHistory
					finalResponseFromGpt, err := gpt.ChatGptFinal(clientData)
					if err != nil {
						log.Println(err)
					} else {
						// Отправка окончательного ответа пользователю
						msgFinal := tgbotapi.NewMessage(message.Chat.ID, finalResponseFromGpt)
						_, err = bot.Send(msgFinal)
						if err != nil {
							log.Println(err)
						}
					}
				}

			}
		}
	}
}

func GetCleanTextFromURL(url string) ([]string, error) {
	// Создаем новый коллектор
	c := colly.NewCollector()

	// Создаем срез для хранения извлеченного содержимого
	var content []string

	// Настраиваем обработчик для интересующих нас HTML-тегов
	c.OnHTML("title, h1, h2, h3, p", func(e *colly.HTMLElement) {
		// Извлекаем текстовое содержимое из тега
		text := strings.TrimSpace(e.Text)

		// Проверяем, начинается ли строка с префикса "//"
		if strings.HasPrefix(text, "//") {
			return // Прерываем выполнение обработчика и возвращаемся
		}

		// Добавляем текст в срез
		if text != "" {
			content = append(content, text)
		}
	})

	// Настраиваем обработку тега img и его содержимого
	c.OnHTML("img, ul, li", func(e *colly.HTMLElement) {
		// Удаляем тег img и его содержимое из исходного HTML
		e.DOM.Remove()
	})

	// Настраиваем обработку ошибок
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Ошибка при запросе URL:", r.Request.URL, "Ошибка:", err)
	})

	// Посещаем целевой URL
	err := c.Visit(url) // Замените на желаемый URL
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе URL: %s ошибка: %w", url, err)
	}

	// Обработка извлеченного контента
	var processedContent []string
	for _, text := range content {
		// Удаляем лишние пробелы и переносы строк
		text = strings.Join(strings.Fields(text), " ")

		// Добавляем обработанный текст в новый срез
		processedContent = append(processedContent, text)
	}

	// Выводим обработанный контент
	for _, text := range processedContent {
		fmt.Println(text)
	}

	// Объединяем обработанный контент с помощью символа новой строки
	text := strings.Join(content, "\n")

	var result []string

	// Разделяем текст на несколько переменных по 4000 символов
	for len(text) > 0 {
		if len(text) > 4000 {
			result = append(result, text[:4000])
			text = text[4000:]
		} else {
			result = append(result, text)
			break
		}
	}

	return result, nil
}
