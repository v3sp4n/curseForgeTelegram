package telegram

import (
    "mime/multipart"
    "encoding/json"
    "io/ioutil"
    "net/http"
    "net/url"
    "strconv"
    // "strings"
    "bytes"
    "fmt"
    "io"
    "os"
)

type Bot struct {
    Token string
    Chatid string
    PreviousMessageIdBot int
}

func (v Bot) SendDocument(path string) {
    var buffer bytes.Buffer
    file,_ := os.Open(path)
    writer := multipart.NewWriter(&buffer)
    part, _ := writer.CreateFormFile("document", path)
    io.Copy(part, file)
    writer.Close()
    req, _ := http.NewRequest("POST", fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument?chat_id=%s", v.Token, v.Chatid), &buffer)
    req.Header.Set("Content-Type", writer.FormDataContentType())
    client := &http.Client{}
    resp, _ := client.Do(req)
    defer resp.Body.Close()
}

func (v *Bot) SendMessage(text string, keyboard [][]map[string]string) {
    k := map[string]interface{}{
        "inline_keyboard": keyboard,
    }
    kj,_ := json.Marshal(k)
    r,_ := http.PostForm("https://api.telegram.org/bot"+v.Token+"/sendMessage", url.Values{
        "chat_id": {v.Chatid},
        "text": {text},
        "reply_markup": {string(kj)},
    })
    b,_ := ioutil.ReadAll(r.Body)
    type s struct {
        Ok bool
        Result struct {
            Message_id int
        }
    }
    j := s{}
    json.Unmarshal(b,&j)
    v.PreviousMessageIdBot = j.Result.Message_id
}

func (v Bot) DeleteMessage(message_id int) {
    http.PostForm("https://api.telegram.org/bot"+v.Token+"/deleteMessage",url.Values{
        "chat_id": {v.Chatid},
        "message_id": {strconv.Itoa(message_id)},
    })
}