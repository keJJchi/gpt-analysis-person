package chatgpt

import (
	"context"
	"fmt"
	"kejjchibot/utils"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
)

func ChatGpt(clientData *utils.ClientData) (string, error) {
	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		return "", fmt.Errorf("missing API_KEY")
	}

	client := gpt3.NewClient(apiKey)

	ctx := context.Background()

	response, err := getResponse(ctx, client, clientData)
	if err != nil {
		return "", err
	}

	return response, nil
}

func getResponse(ctx context.Context, client gpt3.Client, clientData *utils.ClientData) (response string, err error) {
	tonal := "Оцени тональность текста, по отношению к клиенту и его работе, который представлен ниже по этим критериям: Позитив, Нейтрал, Негатив, Нерелевантно. В Ответ дать одним словов. Текст:\n"
	question := tonal + "Клиент: " + clientData.Name + "\n" + " Страна: " + clientData.Country + "\n" + " Работает в компании: " + clientData.Company + "\n" + clientData.Content

	sb := strings.Builder{}

	fmt.Println(clientData.Content)
	// fmt.Println(question)

	err = client.CompletionStreamWithEngine(
		ctx,
		gpt3.TextDavinci003Engine,
		gpt3.CompletionRequest{
			Prompt: []string{
				question,
			},
			MaxTokens:   gpt3.IntPtr(20), //количество токенов в ответе
			Temperature: gpt3.Float32Ptr(0),
		},
		func(resp *gpt3.CompletionResponse) {
			text := resp.Choices[0].Text

			sb.WriteString(text)
		},
	)
	if err != nil {
		return "", err
	}
	response = sb.String()
	response = strings.TrimLeft(response, "\n")
	fmt.Println(sb)
	fmt.Println(response)

	return response, nil
}

func ChatGptFinal(clientData *utils.ClientData) (string, error) {
	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		return "", fmt.Errorf("missing API_KEY")
	}
	fmt.Println(apiKey)

	client := gpt3.NewClient(apiKey)

	ctx := context.Background()

	response, err := getResponseFinal(ctx, client, clientData)
	if err != nil {
		return "", err
	}

	return response, nil
}

func getResponseFinal(ctx context.Context, client gpt3.Client, clientData *utils.ClientData) (response string, err error) {
	// tonal := "Дай среднее значение по тональностям из предложенного. В Ответе дать одно слово, к какой тональности относится текст. Текст:\n"
	tonal := "Оцени как часто в предложенном тексте встречается, на русском или английском, слова: Позитивно. Негативно. Нейтрально. Нерелевантно. Выдай в ответе одно слово, которые встречается чаще. Если нет этих слов напиши 'Ошибка' Текс:\n"
	question := tonal + clientData.Content

	sb := strings.Builder{}

	// fmt.Println(clientData)
	fmt.Println(question)

	err = client.CompletionStreamWithEngine(
		ctx,
		gpt3.TextDavinci003Engine,
		gpt3.CompletionRequest{
			Prompt: []string{
				question,
			},
			MaxTokens:   gpt3.IntPtr(200), //количество токенов в ответе
			Temperature: gpt3.Float32Ptr(0),
		},
		func(resp *gpt3.CompletionResponse) {
			text := resp.Choices[0].Text

			sb.WriteString(text)
		},
	)
	if err != nil {
		return "", err
	}
	response = sb.String()
	response = strings.TrimLeft(response, "\n")
	fmt.Println(sb)
	fmt.Println(response)

	return response, nil
}
