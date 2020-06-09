package email_service

const (
	welcomeTheme = "Добро пожаловать!"
	resetTheme   = "Инструкция по сбросу пароля"
)

const welcomeText = `
	Приветствую!

	Мы надеемся, что вам понравится наше приложение!

	С наилучшими пожеланиями,
	Команда NoName
`

const welcomeHTML = `
	Приветствую!<br/>
	<br/>
	Добро пожаловать на <a href="https://www.example.com">Пример</a>! Мы надеемся, что вам понравится наше приложение!<br/> 
	<br/>
	С наилучшими пожеланиями,<br/>
	Команда NoName
`

const resetTextTmpl = `
	Приветствую!

	Похоже, что Вы запросили сброс пароля. Если это были вы, перейдите по ссылке ниже, чтобы обновить свой пароль:

	URL: %s
	HTTP Verb: PUT
	Body (JSON Payload): { "password": "Ваш новый пароль" }
	
	Если вас спросят о токине, пожалуйста используйте следующее значение:

	%s

	Если вы не запрашивали сброс пароля, можете игнорировать данное сообщение, и ваша учетная запись не будет изменена.

	С наилучшими пожеланиями,
	Команда NoName
`

const resetHTMLTmpl = `
	Приветствую!<br/>
	<br/>
	Похоже, что Вы запросили сброс пароля. 
	Если это были вы, перейдите по ссылке ниже, чтобы обновить свой пароль:<br/>
	<br/>
	URL: %s<br/>
	HTTP Verb: PUT<br/>
	Body (JSON Payload): { "password": "Ваш новый пароль" }<br/>
	<br/>
	Если вас спросят о токине, пожалуйста используйте следующее значение:<br/>
	<br/>
	%s<br/>
	<br/>
	Если вы не запрашивали сброс пароля, можете игнорировать данное сообщение, 
	и ваша учетная запись не будет изменена.<br/>
	<br/>
	С наилучшими пожеланиями,<br/>
	Команда NoName<br/>
`
