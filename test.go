messageScene := scene.Message(func(bot *api.Bot, message update.Message) {
		// users, _ := bot.GetUsers(message.UserID)
		user := bot.GetUser(types.User{
			UserID: message.UserID,
			Fields: []string{"screen_name"},
		})
		m := regexp.MustCompile(`^/(\S+)`).FindStringSubmatch(message.Text)
		log.Println("command",m)
		if len(m) == 2 {
			switch m[1] {
			case "leaders":
				msg := ""
				for _,v := range cfg.Leaders {
					msg = msg + fmt.Sprintf(`
📝 %s - %s
📝 %s
📝
📝 %s

`, v.VK, v.Nickname, v.Server, v.Mafia)
				}
				if len(cfg.Leaders) == 0 {
					msg = "Пока нет лидеров"
				}
				bot.SendMessage(message.ChatID, msg)
			default: 
				if (findToLeaders(user.ScreenName) || findToAdmins(user.ScreenName)) {
					switch m[1] {
					case "makeleader":
						if findToAdmins(user.ScreenName) {
							m = regexp.MustCompile(`^/\S+\s+(\S+)\s+(\S+)\s+(\S+)\s+(.+)`).FindStringSubmatch(message.Text)
							if len(m) == 5 {
								query := getArizonaQuery()
								queryFind := false
								for _,v := range query.Query {
									if strings.ToLower(strings.ReplaceAll(v.Name," ","")) == m[3] || strconv.Itoa(v.Number) == m[3] {
										m[3] = fmt.Sprintf("%s[%d]",v.Name,v.Number)
										queryFind = true
										break
									}
								}
								if queryFind {
									VKID := strings.Split(m[1],"|") 
									m[1] = strings.ReplaceAll(VKID[1],"]","")
									cfg.Leaders = append(cfg.Leaders, struct {
										Nickname string
										VK string
										Server string
										Mafia string
									}{
										m[2],
										strings.ReplaceAll(m[1],"@",""),
										m[3],
										m[4],

									})
									cfgString,_ := json.Marshal(cfg)
									ioutil.WriteFile("./config.json",cfgString,0644)
									bot.SendMessage(message.ChatID, "@"+strings.ReplaceAll(m[1],"@","")+`bla-bla-bla`)
								} else {
									bot.SendMessage(message.ChatID, "Введите правильно сервер, название или номер сервера\n!Пишите название сервера с БЕЗ пробелов!")
								}
							} else {
								bot.SendMessage(message.ChatID, "/makeleader [@idvk] [Nick_Name] [server(NameOrNumber:Brainburg[05])] [Мафия]")
							}
						}
					case "removeleader":
						if findToAdmins(user.ScreenName) {
							m = regexp.MustCompile(`^/\S+\s+(\S+)\s+(\S+)`).FindStringSubmatch(message.Text)
							if len(m) == 3 {
								query := getArizonaQuery()
								queryFind := false
								for _,v := range query.Query {
									if strings.ToLower(strings.ReplaceAll(v.Name," ","")) == m[2] || strconv.Itoa(v.Number) == m[2] {
										for kk, vv := range cfg.Leaders {
											if vv.Server == fmt.Sprintf("%s[%d]",v.Name,v.Number) && strings.ToLower(vv.Nickname) == strings.ToLower(m[1]) {
												bot.SendMessage(message.ChatID, "Лидер был удален!")
												cfg.Leaders = append(cfg.Leaders[:kk], cfg.Leaders[kk+1:]...)
												cfgString,_ := json.Marshal(cfg)
												ioutil.WriteFile("./config.json",cfgString,0644)
												queryFind = true
												break
											}
										}
									}
								}
								if !queryFind {
									bot.SendMessage(message.ChatID, "Ошибка, не смог найти лидера!\n/removeleader [Nick_Name] [server(NameOrNumber:Brainburg[05])]")
								}
							} else {
								bot.SendMessage(message.ChatID, "/removeleader [Nick_Name] [server(NameOrNumber:Brainburg[05])]")
							}
						}
					}
				}
			}
		}

	})


////////////
func findToLeaders(whoFind string) bool {
	for _, v := range cfg.Leaders {
		if strings.ToLower(whoFind) == strings.ToLower(strings.ReplaceAll(v.VK,"@","")) || strings.ToLower(whoFind) == strings.ToLower(v.Nickname) {
			return true
		}
	}
	return false
}
func findToAdmins(whoFind string) bool {
	for _, v := range cfg.Admins {
		if strings.ToLower(whoFind) == strings.ToLower(strings.ReplaceAll(v,"@","") ) {
			return true
		}
	}
	return false
}