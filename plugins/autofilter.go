package plugins

import (
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/Jisin0/Go-Filter-Bot/utils"
	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers"
	"github.com/PaulSonOfLars/gotgbot/v2/ext/handlers/filters/message"
)

var cbEditPattern = regexp.MustCompile(`edit\((.+)\)`)

func autoFilter(b *gotgbot.Bot, ctx *ext.Context) error {
	message := ctx.Message
	settings, err := getSettings(message.Chat.Id)
	if err != nil {
		return err
	}

	// Ignore commands and certain special characters
	if message.Text[0] == '/' || message.Text[0] == ',' || message.Text[0] == '!' || message.Text[0] == '.' || isEmoji(message.Text[0]) {
		return nil
	}

	if len(message.Text) > 2 && len(message.Text) < 100 {
		search := message.Text
		files, offset, totalResults, err := getSearchResults(search, 0, true)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			if settings.SpellCheck {
				return advantageSpellCheck(ctx)
			}
			return nil
		}

		pre := "file"
		if settings.FileSecure {
			pre = "filep"
		}
		req := 0
		if message.From != nil {
			req = message.From.Id
		}

		btn := generateButtons(files, settings.Button, pre, req)

		btn = append([][]gotgbot.InlineKeyboardButton{{{Text: "🔗 ʜᴏᴡ ᴛᴏ ᴅᴏᴡɴʟᴏᴀᴅ 🔗", CallbackData: "howdl"}}}, btn...)
		if offset != "" {
			key := fmt.Sprintf("%d-%d", message.Chat.Id, message.MessageId)
			utils.Temp.GP_BUTTONS[key] = search
			btn = append(btn, []gotgbot.InlineKeyboardButton{
				{Text: fmt.Sprintf("❄️ ᴩᴀɢᴇꜱ 1/%d", int(math.Ceil(float64(totalResults)/6))), CallbackData: "pages"},
				{Text: "➡️ ɴᴇxᴛ", CallbackData: fmt.Sprintf("next_%d_%s_%d", req, key, offset)},
			})
		} else {
			btn = append(btn, []gotgbot.InlineKeyboardButton{{Text: "❄️ ᴩᴀɢᴇꜱ 1/1", CallbackData: "pages"}})
		}

		imdb, err := getPoster(search, files[0].FileName)
		if err != nil {
			return err
		}

		var cap string
		if imdb != nil {
			cap = formatTemplate(settings.Template, imdb, message, search)
		} else {
			cap = fmt.Sprintf("Hᴇʀᴇ Is Wʜᴀᴛ I Fᴏᴜɴᴅ Fᴏʀ Yᴏᴜʀ Qᴜᴇʀʏ %s", search)
		}

		if imdb != nil && imdb.Poster != "" {
			err = sendPoster(message, imdb.Poster, cap, btn)
			if err != nil {
				return err
			}
		} else {
			err = sendText(message, cap, btn)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func generateButtons(files []File, useButton bool, pre string, req int) [][]gotgbot.InlineKeyboardButton {
	var btn [][]gotgbot.InlineKeyboardButton
	for _, file := range files {
		if useButton {
			btn = append(btn, []gotgbot.InlineKeyboardButton{{Text: fmt.Sprintf("[%s] %s", getSize(file.FileSize), file.FileName), CallbackData: fmt.Sprintf("%s#%d#%s", pre, req, file.FileId)}})
		} else {
			btn = append(btn, []gotgbot.InlineKeyboardButton{
				{Text: file.FileName, CallbackData: fmt.Sprintf("%s#%d#%s", pre, req, file.FileId)},
				{Text: getSize(file.FileSize), CallbackData: fmt.Sprintf("%s#%d#%s", pre, req, file.FileId)},
			})
		}
	}
	return btn
}

func sendPoster(message *gotgbot.Message, poster, caption string, buttons [][]gotgbot.InlineKeyboardButton) error {
	opts := &gotgbot.SendPhotoOpts{
		Caption:   caption,
		ParseMode: "HTML",
		ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		},
	}
	_, err := message.ReplyPhoto(poster, opts)
	return err
}

func sendText(message *gotgbot.Message, text string, buttons [][]gotgbot.InlineKeyboardButton) error {
	opts := &gotgbot.SendMessageOpts{
		ParseMode: "HTML",
		ReplyMarkup: &gotgbot.InlineKeyboardMarkup{
			InlineKeyboard: buttons,
		},
	}
	_, err := message.Reply(text, opts)
	return err
}

func formatTemplate(template string, imdb *IMDB, message *gotgbot.Message, search string) string {
	return fmt.Sprintf(template,
		message.Chat.Title,
		message.From.Mention(),
		search,
		imdb.Title,
		imdb.Votes,
		imdb.Aka,
		imdb.Seasons,
		imdb.BoxOffice,
		imdb.LocalizedTitle,
		imdb.Kind,
		imdb.IMDBId,
		imdb.Cast,
		imdb.Runtime,
		imdb.Countries,
		imdb.Certificates,
		imdb.Languages,
		imdb.Director,
		imdb.Writer,
		imdb.Producer,
		imdb.Composer,
		imdb.Cinematographer,
		imdb.MusicTeam,
		imdb.Distributors,
		imdb.ReleaseDate,
		imdb.Year,
		imdb.Genres,
		imdb.Poster,
		imdb.Plot,
		imdb.Rating,
		imdb.URL,
	)
}

func isEmoji(char byte) bool {
	// Implement the function to check if the character is an emoji
	return false
}
