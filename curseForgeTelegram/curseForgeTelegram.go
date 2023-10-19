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

func main() {
	bot = tg.Bot{
		Token: "0",
		Chatid: "1",
	}

	checkUpdates()

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

func checkUpdates() {
	l := time.Now().Unix()-3
	last_message := 0
	for {
		if time.Now().Unix()-l >= 2 {
			l = time.Now().Unix()
			r,_ := http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates?chat_id=%s&offset=-1",bot.Token,bot.Chatid))
			text,_ := ioutil.ReadAll(r.Body)
			j := GetUpdates_message{}
			json.Unmarshal(text,&j)
			if j.Ok && len(j.Result) == 1 && last_message != j.Result[0].Update_id && last_message != 0 {
				last_message = j.Result[0].Update_id

				if j.Result[0].Message.Text != "" {//message text by user

					fmt.Println("enter text by user:",j.Result[0].Message.Text)
					searchMod(j.Result[0].Message.Text)

				} else {
					j := GetUpdates_callback_query{}
					json.Unmarshal(text,&j)
					fmt.Println(">>>>set button!",j.Result[0].Callback_query.Data)
					m := strings.Split(j.Result[0].Callback_query.Data,":")
					if m[0] == "modSelect" {
						bot.DeleteMessage(bot.PreviousMessageIdBot)

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
											"callback_data": "modVersions:"+ fmt.Sprintf("%s:%s",m[1],files.Data[f].GameVersions[g]),
										},
									})
								}
							}
						}
						bot.SendMessage("https://www.curseforge.com/minecraft/mc-mods/"+m[2]+"\n\nmod game versions:",res_inline_keyboard)
					} else if m[0] == "modVersions" {
						bot.DeleteMessage(bot.PreviousMessageIdBot)

						files,_ := cf.GetFiles(m[1])
						res_inline_keyboard := [][]map[string]string{}
						for f := 0; f < len(files.Data); f++ {
							for g := 0; g < len(files.Data[f].GameVersions); g++ {
								if files.Data[f].GameVersions[g] == m[2] {
									res_inline_keyboard = append(res_inline_keyboard,[]map[string]string{
										{
											"text": files.Data[f].FileName,
											"callback_data": "modGet:"+ fmt.Sprintf("%s:%s:%s",m[1],strconv.Itoa(files.Data[f].Id),files.Data[f].FileName),
										},
									})
								}
							}
						}
						bot.SendMessage("files:",res_inline_keyboard)

					} else if m[0] == "modGet" {
						bot.DeleteMessage(bot.PreviousMessageIdBot)

						res_inline_keyboard := [][]map[string]string{}
						bot.SendMessage(m[3]+"\nlink download:\n" + fmt.Sprintf("https://www.curseforge.com/api/v1/mods/%s/files/%s/download",m[1],m[2]),res_inline_keyboard)

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
