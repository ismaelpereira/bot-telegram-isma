package controllers

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-redis/redis/v7"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/ismaelpereira/telegram-bot-isma/api/clients"
	"github.com/ismaelpereira/telegram-bot-isma/bot/msgs"
	"github.com/ismaelpereira/telegram-bot-isma/config"
	"github.com/ismaelpereira/telegram-bot-isma/types"
	"github.com/txgruppi/parseargs-go"
)

var kbComplete tgbotapi.InlineKeyboardMarkup

func ChecklistHandleUpdate(
	cfg *config.Config,
	redis *redis.Client,
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update,
) error {
	if update.CallbackQuery != nil {
		chatID := strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10)
		if strings.HasPrefix(update.CallbackQuery.Data, "checklist:"+chatID) {
			return checklistCallback(bot, update)
		}
		if strings.HasPrefix(update.CallbackQuery.Data, "checklist:") {
			return checkItem(bot, update)
		}
	}
	query := update.Message.CommandArguments()
	if query == "" {
		text := "Para usar o comando /checklist são necesários vários comandos: \n" +
			"/checklist new <nome da checklist> -- cria uma nova checklist\n" +
			"/checklist add <nome da checklist> <valor1,valor2,valorn...> -- adiciona itens a checklist\n" +
			"/checklist list -- lista todas as checklists ja adicionadas com um menu de navegacao\n" +
			"/checklist delete <nome da checklist> -- deleta a checklist desejada\n"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		_, err := bot.Send(msg)
		return err
	}
	args, err := parseargs.Parse(query)
	if err != nil {
		return err
	}
	if args[0] == "new" {
		chatID := strconv.FormatInt(update.Message.Chat.ID, 10)
		title := strings.Join(args[1:], " ")
		err := clients.NewReminder(chatID, title)
		if err != nil {
			return err
		}
		messageText := "Checklist criado com sucesso!\nCom o titulo: " + title
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
		_, err = bot.Send(msg)
		return err
	}

	if args[0] == "add" {
		chatID := strconv.FormatInt(update.Message.Chat.ID, 10)
		keywords := strings.Join(args[1:], " ")
		itens := strings.Fields(keywords)
		if len(itens) < 2 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID,
				"Você digitou o comando errado. Não foi possível completar a solicitação")
			_, err := bot.Send(msg)
			return err
		}
		checklistTitle := itens[0]
		if strings.Contains(itens[1], " ") {
			itens[1] = strings.ReplaceAll(itens[1], " ", "")
		}
		spew.Dump(itens)
		values := strings.Split(strings.TrimSpace(itens[1]), ",")
		spew.Dump(values)

		var checklist types.Checklist
		objects := make([]types.ChecklistItem, len(values), cap(values))
		checklist.Title = checklistTitle
		for i, itens := range values {
			objects[i].Name = itens
		}
		checklist.Itens = objects
		listJSON, err := json.Marshal(checklist)
		if err != nil {
			return err
		}
		if err = clients.AddReminder(chatID, checklistTitle, listJSON); err != nil {
			return err
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Itens adicionados com sucesso!")
		_, err = bot.Send(msg)
		return err
	}
	if args[0] == "delete" {
		chatID := strconv.FormatInt(update.Message.Chat.ID, 10)
		title := strings.Join(args[1:], " ")
		err := clients.DeleteReminder(chatID, title)
		if err != nil {
			return err
		}
		messageText := "Checklist com o titulo " + title + " deletada com sucesso!"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageText)
		_, err = bot.Send(msg)
		return err
	}

	if args[0] == "list" {
		chatID := strconv.FormatInt(update.Message.Chat.ID, 10)
		list, err := clients.ListReminder(chatID)
		if err != nil {
			return err
		}
		var kb [][]tgbotapi.InlineKeyboardButton
		for i, list := range list {
			kb = append(kb, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i+1)+". "+strings.TrimPrefix(list, "checklist:"+chatID+":"), list),
			))
		}
		if len(kb) < 1 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Você não tem listas ainda, adicione uma para aparecer aqui\n")
			_, err = bot.Send(msg)
			return err
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "**CHECKLISTS**\n")
		kbComplete.InlineKeyboard = kb
		msg.ReplyMarkup = kbComplete
		_, err = bot.Send(msg)
		return err
	}
	return nil
}

func checklistCallback(
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update) error {
	values, err := clients.GetReminder(update.CallbackQuery.Data)
	if err != nil {
		return err
	}
	if len(values) == 0 && cap(values) == 0 {
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Você não adiciou itens nessa lista, adicione um para ativar o botão\n")
		_, err = bot.Send(msg)
		return err
	}
	var data types.Checklist
	err = json.Unmarshal(values, &data)
	if err != nil {
		return err
	}
	message := make([]string, 0, len(data.Itens))
	kb := make([][]tgbotapi.InlineKeyboardButton, 0, len(data.Itens))
	for i, item := range data.Itens {
		var symbol string
		if !item.IsChecked {
			symbol = msgs.IconX
		}
		if item.IsChecked {
			symbol = msgs.IconOk
		}
		message = append(message, strconv.Itoa(i+1)+". "+item.Name+symbol+"\n")
		kb = append(kb, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(strconv.Itoa(i+1)+". "+item.Name, "checklist:"+data.Title+":"+strconv.Itoa(i)),
		))
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "**"+data.Title+"**"+"\n"+strings.Join(message, ""))
	kbComplete.InlineKeyboard = kb
	msg.ReplyMarkup = kbComplete
	_, err = bot.Send(msg)
	return err
}

func checkItem(
	bot *tgbotapi.BotAPI,
	update *tgbotapi.Update) error {
	chatID := strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10)
	update.CallbackQuery.Data = strings.TrimPrefix(update.CallbackQuery.Data, "checklist:")
	lastBin := strings.Index(update.CallbackQuery.Data, ":")
	listTitle := update.CallbackQuery.Data[:lastBin]
	arrayPos := update.CallbackQuery.Data[lastBin:]
	arrayPos = strings.TrimPrefix(arrayPos, ":")
	pos, err := strconv.Atoi(arrayPos)
	if err != nil {
		return err
	}
	data, err := clients.GetReminder("checklist:" + chatID + ":" + listTitle)
	if err != nil {
		return err
	}
	var list types.Checklist
	err = json.Unmarshal(data, &list)
	if err != nil {
		return err
	}
	message := make([]string, 0, len(list.Itens))
	for i, item := range list.Itens {
		var symbol string
		if !item.IsChecked {
			symbol = msgs.IconX
		}
		if item.IsChecked {
			symbol = msgs.IconOk
		}
		message = append(message, strconv.Itoa(i+1)+". "+item.Name+symbol+"\n")
	}
	itemChecked := strings.Replace(message[pos], msgs.IconX, msgs.IconOk, 1)
	message[pos] = strings.Replace(message[pos], message[pos], itemChecked, 1)
	msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, strings.Join(message, ""))
	msg.ReplyMarkup = &kbComplete
	_, err = bot.Send(msg)
	if err != nil {
		return err
	}

	list.Itens[pos].IsChecked = true
	newValues, err := json.Marshal(list)
	if err != nil {
		return err
	}
	if err = clients.AddReminder(chatID, listTitle, newValues); err != nil {
		return err
	}
	return nil
}
