# Analysis about a person or company using GPT chat

## Installation 
To install gpt-analysis-person, you need to have Go (version 1.20 or later) installed on your system. Then, run the following command:
```sh
go get github.com/keJJchi/gpt-analysis-person
```
This will download the source code and install the hcl-templater executable in your `$GOPATH/bin` directory.


## Usage
For the application to work, you will need to enter your API for the bot's telegrams (you can get it from for example @BotFather)
```sh
main.go
bot_token = "TG_KEY"
```
and the API for the GPT chat
```sh
chatgpt.go
apiKey = "API_KEY"
```
## Features

This bot helps to evaluate the sentiment of a text from a site about a person or a company and analyze this text using the GPT chat and give an answer: neutral, positive, negative or irrelevant

At the entrance to the telegram bot, you will need to enter the name of the person, in which company he works, attach links to those sites that he must analyze and write his country.

And in the answer you get what tone each link corresponds to.
