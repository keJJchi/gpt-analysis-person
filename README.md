# Analysis about a person or company using GPT chat

This is a Telegram bot that utilizes the GPT-3 language model to provide responses to user messages. The bot can interact with users in a conversation, ask questions, and provide answers based on the provided text.

# Getting Started

## Prerequisites
- Go 1.20 or later
- GPT-3 API Key
- Telegram Bot Token

## Installation 

1. To install gpt-analysis-person, you need to have Go (version 1.20 or later) installed on your system. Then, run the following command:
```sh
go get github.com/keJJchi/gpt-analysis-person
```
This will download the source code and install the gpt-analysis-person executable in your `$GOPATH/bin` directory.

2. Set up environment variables:
```sh
export TG_KEY=your_telegram_bot_token
export API_KEY=your_gpt3_api_key
```
3. Build and run the application:
```sh
   go build
./gpt-analysis-person
```
## Telegram Bot Commands
The Telegram bot supports the following commands:

- `/start` - Start a new conversation with the bot.
- `/stop` - Stop the conversation and clear client data.

## How it Works
The application consists of two main functionalities: `ChatGpt()` and `ChatGptFinal()`, both of which utilize the GPT-3 API to generate responses based on the provided text.

 - `ChatGpt()`: This function receives client data and performs tone evaluation on the provided text. It sends a request to the GPT-3 API with the question and collects the response. The tone evaluation is based on the occurrences of specific words in the text (e.g., positive, negative, neutral, irrelevant).

 - `ChatGptFinal()`: This function is an extension of ChatGpt() and gathers responses from multiple parts of the text. It combines the responses and sends them to the GPT-3 API for final tone evaluation. It looks for the frequency of specific words (e.g., positive, negative, neutral, irrelevant) in the combined responses.

 - `GetCleanTextFromURL()`: This function extracts clean text from a given URL by visiting the website and parsing the HTML content. It removes unnecessary tags and returns the extracted text as a list of strings.

## Telegram Bot Usage
1. Start a conversation with the bot by sending the `/start` command.
2. The bot will prompt you for client information (name, country, company, and links to check).
3. After providing the required information, the bot will start processing the provided URLs and generate responses.
4. Finally, the bot will provide a summary of the results based on the tone evaluation of the combined responses.

## Contributing
Contributions to this project are welcome. If you find any issues or have suggestions for improvement, feel free to open an issue or submit a pull request.
