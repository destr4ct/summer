package keyboard

import tga "github.com/go-telegram-bot-api/telegram-bot-api/v5"

// Клавиатура (Ресурсы)
// 1.
// 		- Просто отправляет список ресурсов, по которым делается summary
// 2. Добавить
//  	- Добавляет ссылку/ссылки
// 3. Удалить
//		- Удаляет ссылку под номером #n
// 4. Добавить ключевые слова
// 5. Назад
//		- Возвращает в прошлое меню

var ResourcesKeyboard = tga.NewReplyKeyboard(
	tga.NewKeyboardButtonRow(
		tga.NewKeyboardButton("Ресурсы"),
		tga.NewKeyboardButton("Ключевые слова"),
	),
	tga.NewKeyboardButtonRow(
		tga.NewKeyboardButton("Добавить"),
		tga.NewKeyboardButton("Удалить"),
	),
	tga.NewKeyboardButtonRow(
		tga.NewKeyboardButton("Назад"),
	),
)
