package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jonlaing/htmlmeta"
)

// func getGeneralLegality(cl *[]card, n string) string {
// 	legalRarity := false
// 	vintageLegality := false
// 	for _, v := range *cl {
// 		if strings.EqualFold(v.Name, n) {
// 			if v.Legalities.Vintage == "legal" || v.Legalities.Vintage == "restricted" {
// 				vintageLegality = true
// 			} else {
// 				vintageLegality = false
// 				break
// 			}
// 			if v.Rarity == "common" || v.Rarity == "uncommon" {
// 				legalRarity = true
// 			}
// 			if v.Reserved == true {
// 				legalRarity = false
// 			}

// 		}
// 	}
// 	fmt.Printf("Checking General Legality for: %v \n", n)
// 	if legalRarity == true && vintageLegality == true {
// 		fmt.Println("Legal")
// 		return "legal"
// 	} else {
// 		fmt.Println("Not legal")
// 	}
// 	return "not legal"
// 	// fmt.Printf("Found: %v with Rarity: %v \n", v.Name, v.Rarity)
// 	// fmt.Println(v.Legalities.Vintage)
// }

// func getBannedList(n string) string {
// 	bannedCards, err := os.Open("banlist.json")
// 	if err != nil {
// 		fmt.Printf("Fatal Error Reading banlist.json: %v", err)
// 	}
// 	defer bannedCards.Close()
// 	byteBanList, err := ioutil.ReadAll(bannedCards)

// 	var banList []card
// 	err = json.Unmarshal(byteBanList, &banList)

// 	if err != nil {
// 		fmt.Printf("Fatal Error unmarshalling banlist.json: %v", err)
// 	}

// 	legal := true

// 	for _, v := range banList {
// 		if strings.EqualFold(v.Name, n) {
// 			legal = false

// 		}
// 	}

// 	if legal == false {
// 		return "banned"
// 	}
// 	return "not banned"
// 	// fmt.Printf("Found: %v with Rarity: %v \n", v.Name, v.Rarity)
// 	// fmt.Println(v.Legalities.Vintage)
// }

// func queryScryfall(n string) *scryfallList {
// 	var sl scryfallList
// 	urlCardname := strings.ReplaceAll(n, " ", "+")
// 	response, err := http.Get("https://api.scryfall.com/cards/search?unique=prints&q=" + urlCardname + "&pretty=true")
// 	responseData, err := ioutil.ReadAll(response.Body)
// 	//fmt.Println(string(responseData))

// 	err = json.Unmarshal(responseData, &sl)
// 	if err != nil {
// 		fmt.Print(err)
// 	}

// 	var culled scryfallList
// 	for i, v := range sl.CardList {
// 		tmp := strings.ReplaceAll(v.Name, "'", "")

// 		//fmt.Println(tmp)
// 		if strings.EqualFold(strings.ReplaceAll(n, "'", ""), tmp) {
// 			culled.CardList = append(culled.CardList, sl.CardList[i])
// 		}
// 	}
// 	fmt.Println(len(culled.CardList))
// 	fmt.Println(culled.CardList[0].Name)
// 	fmt.Println(len(culled.CardList))
// 	return &culled
// }

// func generateEmbed(c string, found bool, sl *scryfallList) *discordgo.MessageEmbed {
// 	var embed discordgo.MessageEmbed

// 	if !found {
// 		embed.Title = "Not Found"
// 		embed.Color = 15158332
// 		embed.Description = c + " not found in Scryfall Database"
// 		return &embed
// 	} else {
// 		general := getGeneralLegality(&sl.CardList, sl.CardList[0].Name)
// 		banned := getBannedList(c)
// 		legalCard := false
// 		if general == "legal" && banned == "not banned" {
// 			legalCard = true
// 		} else {
// 			legalCard = false
// 		}
// 		var thumbnail discordgo.MessageEmbedThumbnail

// 		if len(sl.CardList) > 0 {
// 			thumbnail.URL = sl.CardList[0].Image_uris.Normal
// 		} else {
// 			thumbnail.URL = ""
// 		}
// 		urlCardname := strings.ReplaceAll(c, " ", "+")
// 		footer := discordgo.MessageEmbedFooter{
// 			Text: "\n Banlist Status: " + banned +
// 				"\nLegality otherwise: " + general +
// 				"\nReserved List: " + strconv.FormatBool(sl.CardList[0].Reserved),
// 		}

// 		embed = discordgo.MessageEmbed{
// 			Title: sl.CardList[0].Name,
// 			URL:   "https://scryfall.com/search?as=grid&order=released&q=%21\"" + urlCardname + "\"&unique=prints",
// 			Color: 15158332,
// 			Description: "**Proletariat Legal:** " + strconv.FormatBool(legalCard) +
// 				"\n\n`" + sl.CardList[0].Oracle + "`\n",
// 			Thumbnail: &thumbnail,
// 			Footer:    &footer,
// 		}
// 	}

// 	return &embed
// }

func proxies() *discordgo.MessageEmbed {
	var embed discordgo.MessageEmbed

	var thumbnail discordgo.MessageEmbedThumbnail
	thumbnail.URL = "https://thepeoplesformat.com/discord_logo.png"

	footer := discordgo.MessageEmbedFooter{
		Text: "As always, if you have any questions or concerns, please contact a member of the Admin team and we'll be glad to assist",
	}

	embed = discordgo.MessageEmbed{
		Title:       "Proletariat Statement on Proxies",
		Color:       15158332,
		Description: "Proletariat is not an offical Wizards of the Coast format. We at Proletariat MTG are entirely welcoming of legible proxy cards with recognizable color and art printed. *That said*, for any events please check with your Tournmanet Organizer on the acceptability of proxies for competitive play.",
		Thumbnail:   &thumbnail,
		Footer:      &footer,
	}

	return &embed
}

///// REBUILD STARTS HERE

func decklistEmbed(legalList bool, id string) *discordgo.MessageEmbed {
	var embed discordgo.MessageEmbed

	var thumbnail discordgo.MessageEmbedThumbnail
	thumbnail.URL = "https://thepeoplesformat.com/discord_logo.png"

	footer := discordgo.MessageEmbedFooter{
		Text: "Please note at this time this function only checks individual card legality.",
	}
	color := 15158332
	if legalList {
		color = 3066993
	}

	title := id
	response, err := http.Get("https://www.mtggoldfish.com/deck/" + id)
	if err != nil {
		fmt.Printf("%s", err)
	} else {
		defer response.Body.Close()
		meta := htmlmeta.Extract(response.Body)
		title = meta.Title
	}

	embed = discordgo.MessageEmbed{
		Title:       title,
		URL:         "https://www.mtggoldfish.com/deck/" + id,
		Color:       color,
		Description: "Legal List: " + strconv.FormatBool(legalList),
		Thumbnail:   &thumbnail,
		Footer:      &footer,
	}

	return &embed
}

// Clean input from any given query to remove special characters and add replace spaces with plus symbols
func cleanInput(n string) (name, urlName string) {
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	if err != nil {
		log.Fatal(err)
	}

	cleanedInput := reg.ReplaceAllString(n, "")
	urlString := strings.ReplaceAll(cleanedInput, " ", "+")

	//fmt.Printf("Input String: %s\n", cleanedInput)

	return cleanedInput, urlString
}
func goldfish(id string) bool {
	main := make(map[string]int)
	side := make(map[string]int)

	response, err := http.Get("https://www.mtggoldfish.com/deck/download/" + id)
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err)
	} else {
		//fmt.Print(string(responseData))
	}
	list := string(responseData)

	//fmt.Print(list)
	scanner := bufio.NewScanner(strings.NewReader(list))
	sideboard := false
	for scanner.Scan() {

		//fmt.Println(scanner.Text())
		if scanner.Text() != "" && sideboard == false {
			qty, card := processLine(scanner.Text())
			main[card] = qty
			fmt.Print("MAINBOARD ")
			fmt.Printf("Qty: %d :: Card: %s \n", qty, card)

		} else if scanner.Text() != "" && sideboard == true {
			qty, card := processLine(scanner.Text())
			side[card] = qty
			fmt.Print("SIDEBOARD ")
			fmt.Printf("Qty: %d :: Card: %s \n", qty, card)

		} else if scanner.Text() == "" {
			sideboard = true
		}
		//processLine(scanner.Text())
	}
	return legalList(main, side)

}
func deckbox(id string) bool {
	main := make(map[string]int)
	side := make(map[string]int)

	response, err := http.Get("https://deckbox.org/sets/" + id + "/export")
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err)
	} else {
		//fmt.Print(string(responseData))
	}
	list := string(responseData)

	//fmt.Print(list)
	scanner := bufio.NewScanner(strings.NewReader(list))
	sideboard := false
	for scanner.Scan() {

		//fmt.Println(scanner.Text())
		if scanner.Text() != "" && scanner.Text() != "Sideboard:" && sideboard == false {
			qty, card := processLine(scanner.Text())
			main[card] = qty
			fmt.Print("MAINBOARD ")
			fmt.Printf("Qty: %d :: Card: %s \n", qty, card)

		} else if scanner.Text() != "" && scanner.Text() != "Sideboard:" && sideboard == true {
			qty, card := processLine(scanner.Text())
			side[card] = qty
			fmt.Print("SIDEBOARD ")
			fmt.Printf("Qty: %d :: Card: %s \n", qty, card)

		} else if scanner.Text() == "" {
			sideboard = true
		}
		//processLine(scanner.Text())
	}
	return legalList(main, side)

}

func legalList(m, s map[string]int) bool {
	legal := true

	mainQty := 0
	sideQty := 0

	for c, q := range m {
		cleanName, urlName := cleanInput(c)
		mainQty += q
		match, _ := scryFallMatch(cleanName, urlName)
		if len(match.CardList) > 0 {
			ban, gen := checkLegality(&match.CardList)
			if q > 4 {
				if !strings.Contains(match.CardList[0].Type, "Basic") {
					legal = false
					break
				}
			}
			if ban == true || gen == false {
				legal = false
				break
			}

		}
	}
	for c, q := range s {
		sideQty += q
		cleanName, urlName := cleanInput(c)
		match, _ := scryFallMatch(cleanName, urlName)
		if len(match.CardList) > 0 {
			ban, gen := checkLegality(&match.CardList)
			if ban == true || gen == false {
				legal = false
			}

		}
	}

	if mainQty < 60 {
		legal = false
	}
	if sideQty > 15 {
		legal = false
	}
	if legal {
		//mergemap(m, s)
		for k, v := range mergemap(m, s) {
			if v > 4 {
				cleanName, urlName := cleanInput(k)
				match, _ := scryFallMatch(cleanName, urlName)
				if len(match.CardList) > 0 {
					if !strings.Contains(match.CardList[0].Type, "Basic") {
						legal = false
					}
				}
			}
		}
	}
	return legal
}

// func mergeMap(main map[string]int, side map[string]int) map[string]int {
// 	for k,v := range b{
// 		a[k] += v
// 	}
// 	reutrn main
// }
func mergemap(a map[string]int, b map[string]int) map[string]int {
	for k, v := range b {
		a[k] += v
	}
	return a
}

func processLine(ln string) (qty int, name string) {
	count := 0
	cardname := ""
	line := strings.SplitN(ln, " ", 2)
	qty, err := strconv.Atoi(line[0])
	if err != nil {
		fmt.Print("ERROR " + err.Error())

	} else {
		cardname = line[1]
		count = qty
	}
	return count, cardname
}
func scryFallMatch(n string, urlString string) (match, matches *scryfallList) {
	var sl scryfallList
	response, err := http.Get("https://api.scryfall.com/cards/search?unique=prints&q=name:/^" + urlString + "/&pretty=true")
	responseData, err := ioutil.ReadAll(response.Body)

	err = json.Unmarshal(responseData, &sl)
	if err != nil {
		fmt.Print(err)
	}
	if len(sl.CardList) == 0 {
		response, err := http.Get("https://api.scryfall.com/cards/search?unique=prints&q=" + urlString + "&pretty=true")
		responseData, err := ioutil.ReadAll(response.Body)

		err = json.Unmarshal(responseData, &sl)
		if err != nil {
			fmt.Print(err)
		}
	}
	var culled scryfallList
	for i, v := range sl.CardList {
		c, _ := cleanInput(v.Name)
		if strings.EqualFold(n, c) {
			culled.CardList = append(culled.CardList, sl.CardList[i])
		}

	}

	fmt.Println(len(culled.CardList))
	//fmt.Println(culled.CardList[0].Name)
	fmt.Println(len(sl.CardList))
	return &culled, &sl
}

func checkLegality(cl *[]card) (banned, generallLegal bool) {

	return checkBanned((*cl)[0].Name), checkLegal(cl)
}

func checkBanned(n string) bool {
	bannedCards, err := os.Open("banlist.json")
	if err != nil {
		fmt.Printf("Fatal Error Reading banlist.json: %v", err)
	}
	defer bannedCards.Close()
	byteBanList, err := ioutil.ReadAll(bannedCards)

	var banList []card
	err = json.Unmarshal(byteBanList, &banList)

	if err != nil {
		fmt.Printf("Fatal Error unmarshalling banlist.json: %v", err)
	}

	banned := false

	for _, v := range banList {
		if strings.EqualFold(v.Name, n) {
			banned = true

		}
	}

	fmt.Printf("[BANLIST]:: Checked: %v with Status: %t \n", n, banned)
	return banned
	// fmt.Printf("Found: %v with Rarity: %v \n", v.Name, v.Rarity)
	// fmt.Println(v.Legalities.Vintage)
}

func checkLegal(cl *[]card) bool {
	legalRarity := false
	vintageLegality := false
	reserved := false
	for _, v := range *cl {

		if v.Legalities.Vintage == "legal" || v.Legalities.Vintage == "restricted" {
			vintageLegality = true
		} else {
			vintageLegality = false
			break
		}
		if v.Rarity == "common" || v.Rarity == "uncommon" {
			legalRarity = true
		}
		if v.Reserved {
			reserved = true
		}

	}
	fmt.Printf("[LEGAL]:: Found: %v with: Rarity Legality: %t \n Vintage Legality: %t \n Reserved List: %t \n", (*cl)[0].Name, legalRarity, vintageLegality, reserved)
	if legalRarity && vintageLegality && !reserved {
		return true
	}
	return false

}

func newEmbed(c string, urlName string, found bool, sl *scryfallList) *discordgo.MessageEmbed {
	var embed discordgo.MessageEmbed

	if !found {
		embed.Title = "Not Found"
		embed.Color = 15158332
		embed.Description = c + " not found in Scryfall Database"
		return &embed
	} else {
		banned, general := checkLegality(&sl.CardList)

		legalCard := false
		if general == true && banned == false {
			legalCard = true
		} else {
			legalCard = false
		}
		var thumbnail discordgo.MessageEmbedThumbnail

		if len(sl.CardList) > 0 {
			thumbnail.URL = sl.CardList[0].Image_uris.Normal
		} else {
			thumbnail.URL = ""
		}

		footer := discordgo.MessageEmbedFooter{
			Text: "\n Banlist Status: " + strconv.FormatBool(banned) +
				"\nLegality otherwise: " + strconv.FormatBool(general) +
				"\nReserved List: " + strconv.FormatBool(sl.CardList[0].Reserved),
		}

		color := 15158332
		if legalCard {
			color = 3066993
		}
		embed = discordgo.MessageEmbed{
			Title: sl.CardList[0].Name,
			URL:   "https://scryfall.com/search?as=grid&order=released&q=%21\"" + urlName + "\"&unique=prints",
			Color: color,
			Description: sl.CardList[0].Manacost +
				"\n*" + sl.CardList[0].Type + "*" +
				"\n**Proletariat Legal:** " + strconv.FormatBool(legalCard) +
				"\n\n`" + sl.CardList[0].Oracle + "`\n",
			Thumbnail: &thumbnail,
			Footer:    &footer,
		}
	}

	return &embed
}
