package main

import (
	cf "curseForgeTelegram/curseForge"
	tg "curseForgeTelegram/telegram"
	"encoding/json"
	"io/ioutil"
	"net/http"
	// "net/url"
	"strings"
	"strconv"
	"regexp"
	"time"
	"fmt"
	"os"
)

var bot tg.Bot
type users__struct struct {
	username string
	PreviousMessageIdBot int
}

func main() {


	bot = tg.Bot{
		Token: "tokenbon",
		Chatid: "none!",
	}

	users := map[string]users__struct{}
	if _,err := os.Stat("users.json"); err != nil {
		ioutil.WriteFile("users.json",[]byte("{}"),0644)
	}
	file,_ := ioutil.ReadFile("users.json")
	json.Unmarshal(file,&users)
	usersSave := func() {
		js,_ := json.Marshal(users)
		ioutil.WriteFile("users.json",js,0644)
	}

	addUser := func(chatid ,username string) {
		_,b := users[chatid]
		if !b {
			users[chatid] = users__struct{}
			users[chatid] = users__struct{
				username: username,
				PreviousMessageIdBot:  0,
			}
			usersSave()
		}
	}

	l := time.Now().UnixMilli()-3
	last_message := 0
	for {
		if time.Now().UnixMilli()-l >= 250 {
			l = time.Now().UnixMilli()
			// https://api.telegram.org/bot6528981379:AAEvlUdCVMCE9tfVs_tks1J5vPLuL98-M5Q/getUpdates?&offset=-1
			r,_ := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?&offset=-1",bot.Token))
			
			pmib := 0
			if _,f := users[bot.Chatid]; f {
				users[bot.Chatid] = users__struct{
					PreviousMessageIdBot: bot.PreviousMessageIdBot,
				}
				usersSave()
				pmib = users[bot.Chatid].PreviousMessageIdBot
			}

			text,_ := ioutil.ReadAll(r.Body)
			j := GetUpdates_message{}
			json.Unmarshal(text,&j)
			if j.Ok && len(j.Result) == 1 && last_message != j.Result[0].Update_id && last_message != 0 {
				last_message = j.Result[0].Update_id

				if j.Result[0].Message.Text != "" {//message text by user
					bot.Chatid = strconv.Itoa(j.Result[0].Message.Chat.Id)
					addUser(bot.Chatid,j.Result[0].Message.Chat.Username)

					fmt.Printf(">>>%s enter text:%s\n",j.Result[0].Message.Chat.Username,j.Result[0].Message.Text)
					if j.Result[0].Message.Text == "/start" {
						bot.SendMessage(`
Hi, here you can download mods from curseforge!
Just write the name of the mod at the bottom

dev vespan`,[][]map[string]string{})
					} else {
						searchMod(j.Result[0].Message.Text)
					}

				} else {

					j := GetUpdates_callback_query{}
					json.Unmarshal(text,&j)
					bot.Chatid = strconv.Itoa(j.Result[0].Callback_query.Message.Chat.Id)
					addUser(bot.Chatid,j.Result[0].Callback_query.Message.Chat.Username)

					fmt.Printf(">>>%s click button:%s\n",j.Result[0].Callback_query.Message.Chat.Username,j.Result[0].Callback_query.Data)
					m := strings.Split(j.Result[0].Callback_query.Data,":")
					if m[0] == "modSelect" {
						bot.DeleteMessage(pmib)

						files,_ := cf.GetFiles(m[1])
						res_inline_keyboard := [][]map[string]string{}
						for f := 0; f < len(files.Data); f++ {
							find := false
							for g := 0; g < len(files.Data[f].GameVersions); g++ {
								
								for l := 0; l < len(res_inline_keyboard); l++ {
									if res_inline_keyboard[l][0]["text"] == files.Data[f].GameVersions[g] {
										find = true
									}
								}
								if !find && len(regexp.MustCompile(`^(\d+)\.(\d+)`).FindStringSubmatch(files.Data[f].GameVersions[g])) >= 3 {
									res_inline_keyboard = append(res_inline_keyboard,[]map[string]string{
										{
											"text": files.Data[f].GameVersions[g],
											"callback_data": "modVersions:"+ fmt.Sprintf("%s:%s:%s",m[1],files.Data[f].GameVersions[g],"0"),
										},
									})
								}
							}
						}
						bot.SendMessage("https://www.curseforge.com/minecraft/mc-mods/"+m[2]+"\n\nmod game versions:",res_inline_keyboard)
					
					} else if m[0] == "modVersions" {
						bot.DeleteMessage(pmib)

						files,_ := cf.GetFiles(m[1])
						var limit int
						var limitBreak int
						var app int
						l,_ := strconv.ParseInt(m[3],10,0)
						limit = int(l)
						limitBreak = limit + 50
						fmt.Println("limit>",limit,"len files",len(files.Data))
						res_inline_keyboard := [][]map[string]string{}
						for f := 0; f < len(files.Data); f++ {
							for g := 0; g < len(files.Data[f].GameVersions); g++ {
								if files.Data[f].GameVersions[g] == m[2] {
									app += 1
									if app >= limit && app <= limitBreak {
										res_inline_keyboard = append(res_inline_keyboard,[]map[string]string{
											{
												"text": files.Data[f].FileName,
												"callback_data": "modGet:"+ fmt.Sprintf("%s:%s:%s",m[1],strconv.Itoa(files.Data[f].Id),files.Data[f].FileName),
											},
										})
									}
								}
							}
							if len(res_inline_keyboard) >= 50 {
								res_inline_keyboard = append(res_inline_keyboard,[]map[string]string{
									{
										"text": "more..",
										"callback_data": "modVersions:"+ fmt.Sprintf("%s:%s:%s",m[1],m[2],strconv.Itoa(int(limitBreak))),
									},
								})
								break
							}
						}
						bot.SendMessage("files:",res_inline_keyboard)

					} else if m[0] == "modGet" {
						bot.DeleteMessage(pmib)

						bot.SendMessage(m[3]+"\nlink download:\n" + fmt.Sprintf("https://www.curseforge.com/api/v1/mods/%s/files/%s/download",m[1],m[2]),[][]map[string]string{})

						cf.Download(m[1],m[2],m[3],"temp/")

						bot.SendDocument("temp/"+m[3])
						os.Remove("temp/"+m[3])

					} else {
						fmt.Println("NIL>")
						fmt.Println(m)
					}
				}

			} else if j.Ok && len(j.Result) == 1 && last_message != j.Result[0].Update_id && last_message == 0 {
				fmt.Println("set last_message != 0")
				last_message = j.Result[0].Update_id
			}
		}
	}

}
func searchMod(s string) {
	ss := s
	ss = strings.ReplaceAll(s," ","+") 
	res,_ := cf.SearchMod(ss)
	if len(res.Data) == 0 {
		bot.SendMessage("couldn't find any mods :(",[][]map[string]string{})
	} else {
		res_inline_keyboard := [][]map[string]string{}
		for d := 0;d < len(res.Data); d++ {
			if res.Data[d].Class.Name == "Mods" {
				res_inline_keyboard = append(res_inline_keyboard,[]map[string]string{
					{
						"text": fmt.Sprintf("%s",res.Data[d].Name),
						"callback_data": "modSelect:"+strconv.Itoa(res.Data[d].Id)+":"+res.Data[d].Slug,
					},
				})
			}
		}
	    bot.SendMessage(s,res_inline_keyboard)
	}
}


type GetUpdates_callback_query struct {
	Ok bool
	Result []struct {
		Update_id int
		Callback_query struct {
			Id string
			From struct {
				Id int
				Is_bot bool
				First_name string
				Username string
				Language_code string
			}
			Message struct {
				Message_id int
				From struct {
					Id int
					Is_bot bool
					First_name string
					Username string
					Language_code string
				}
				Chat struct{
					Id int
					Is_bot bool
					First_name string
					Username string
					Type string
				}
				Date int
				Text string
				Reply_markup struct {
					Inline_keyboard [][]map[string]string
				}
			}
			Chat_instance string
			Data string
		}
	}
}
type GetUpdates_message struct {
	Ok bool
	Result []struct {
		Update_id int
		Message struct {
			Message_id int
			From struct {
				Id int
				Is_bot bool
				First_name string
				Username string
				Language_code string
			}
			Chat struct{
				Id int
				Is_bot bool
				First_name string
				Username string
				Type string
			}
			Date int
			Text string
		}
	}
}
