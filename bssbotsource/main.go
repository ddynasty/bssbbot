package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Syfaro/telegram-bot-api"
	"github.com/gocql/gocql"
)

type Item struct {
	ID          string
	Name        string
	Description string
	Price       int
	Photos      [][]byte
}

func main() {
	session, err := CassandraConn()
	if err != nil {
		log.Panic(err)
	}
	defer session.Close()
	bot, err := tgbotapi.NewBotAPI(teltoken)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	var updchan tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	updchan.Timeout = 60

	// Creating new update channel
	nupch, err := bot.GetUpdatesChan(updchan)
	if err != nil {
		log.Panic(err)
	}

	// Handler for uploading items
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST is supported!", http.StatusForbidden)
			return
		}
		var images [][]byte
		if err = r.ParseMultipartForm(0); nil != err {
			http.Error(w, "Incorrect Multipart Data", http.StatusInternalServerError)
			return
		}
		name := r.FormValue("Name")
		description := r.FormValue("Description")
		price := r.FormValue("Price")
		for _, fheaders := range r.MultipartForm.File {
			for _, hdr := range fheaders {
				image, err := hdr.Open()
				if err != nil {
					log.Panic(err)
				}
				defer image.Close()
				buf := bytes.NewBuffer(nil)
				_, err = io.Copy(buf, image)
				if err != nil {
					log.Panic(err)
				}
				images = append(images, buf.Bytes())
			}
		}
		gocqlUuid := gocql.TimeUUID()
		if err := session.Query("INSERT INTO table1 (id, name, description, price, photos) VALUES (?, ?, ?, ?, ?)", gocqlUuid, name, description, price, images).Exec(); err != nil {
			log.Panic(err)
		}
	})

	// Listening to localhost:8080
	go http.ListenAndServe(":8080", nil)

	// Reading from channel
	for update := range nupch {
		if update.Message != nil {
			if update.Message.Text == "showlist" {
				var items []Item
				m := map[string]interface{}{}
				query := session.Query("SELECT * FROM table1").Iter()
				for query.MapScan(m) {
					items = append(items, Item{
						ID:          m["id"].(gocql.UUID).String(),
						Name:        m["name"].(string),
						Description: m["description"].(string),
						Price:       m["price"].(int),
						Photos:      m["photos"].([][]byte),
					})
					m = map[string]interface{}{}
				}
				for _, item := range items {
					ChatID := update.Message.Chat.ID
					price := strconv.Itoa(item.Price)
					var buttons []tgbotapi.InlineKeyboardButton
					button := tgbotapi.NewInlineKeyboardButtonData("Purchase", item.ID)
					buttons = append(buttons, button)
					markup := tgbotapi.NewInlineKeyboardMarkup(buttons)
					reply := "Name: " + item.Name + "\n" + "Description: " + item.Description + "\n" + "Price: " + price + "\n" + "Purchase Link: "
					text := tgbotapi.NewMessage(ChatID, reply)
					text.ReplyMarkup = markup
					bot.Send(text)
					for i, image := range item.Photos {
						im := tgbotapi.FileBytes{Name: string(i), Bytes: image}
						photo := tgbotapi.NewPhotoUpload(ChatID, im)
						bot.Send(photo)
					}
				}
			}
		}
		if update.CallbackQuery != nil {
			if err := session.Query("DELETE FROM table1 WHERE id=?", update.CallbackQuery.Data).Exec(); err != nil {
				log.Panic(err)
			}
		}
	}
}
