package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func getGeneralLegality(cl *[]card, n string) string {
	legalRarity := false
	vintageLegality := false
	for _, v := range *cl {
		if strings.EqualFold(v.Name, n) {
			if v.Legalities.Vintage == "legal" || v.Legalities.Vintage == "restricted" {
				vintageLegality = true
			} else {
				vintageLegality = false
				break
			}
			if v.Rarity == "common" || v.Rarity == "uncommon" {
				legalRarity = true
			}
			if v.Reserved == true {
				legalRarity = false
			}

		}
	}
	fmt.Printf("Checking General Legality for: %v \n", n)
	if legalRarity == true && vintageLegality == true {
		fmt.Println("Legal")
		return "legal"
	} else {
		fmt.Println("Not legal")
	}
	return "not legal"
	// fmt.Printf("Found: %v with Rarity: %v \n", v.Name, v.Rarity)
	// fmt.Println(v.Legalities.Vintage)
}

func getBannedList(n string) string {
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

	legal := true

	for _, v := range banList {
		if strings.EqualFold(v.Name, n) {
			legal = false

		}
	}

	if legal == false {
		return "banned"
	}
	return "not banned"
	// fmt.Printf("Found: %v with Rarity: %v \n", v.Name, v.Rarity)
	// fmt.Println(v.Legalities.Vintage)
}

func queryScryfall(n string) *scryfallList {
	var sl scryfallList
	urlCardname := strings.ReplaceAll(n, " ", "+")
	response, err := http.Get("https://api.scryfall.com/cards/search?unique=prints&q=" + urlCardname + "&pretty=true")
	responseData, err := ioutil.ReadAll(response.Body)
	//fmt.Println(string(responseData))

	err = json.Unmarshal(responseData, &sl)
	if err != nil {
		fmt.Print(err)
	}

	var culled scryfallList
	for i, v := range sl.CardList {
		if strings.EqualFold(n, v.Name) {
			culled.CardList = append(culled.CardList, sl.CardList[i])
		}
	}
	return &culled
}

func generateEmbed(c string, found bool, sl *scryfallList) *discordgo.MessageEmbed {
	var embed discordgo.MessageEmbed

	if !found {
		embed.Title = "Not Found"
		embed.Color = 15158332
		embed.Description = c + " not found in Scryfall Database"
		return &embed
	} else {
		general := getGeneralLegality(&sl.CardList, sl.CardList[0].Name)
		banned := getBannedList(c)
		legalCard := false
		if general == "legal" && banned == "not banned" {
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
		urlCardname := strings.ReplaceAll(c, " ", "+")
		footer := discordgo.MessageEmbedFooter{
			Text: "\n Banlist Status: " + banned +
				"\nLegality otherwise: " + general +
				"\nReserved List: " + strconv.FormatBool(sl.CardList[0].Reserved),
		}

		embed = discordgo.MessageEmbed{
			Title: sl.CardList[0].Name,
			URL:   "https://scryfall.com/search?as=grid&order=released&q=%21\"" + urlCardname + "\"&unique=prints",
			Color: 15158332,
			Description: "**Proletariat Legal:** " + strconv.FormatBool(legalCard) +
				"\n\n`" + sl.CardList[0].Oracle + "`\n",
			Thumbnail: &thumbnail,
			Footer:    &footer,
		}
	}

	return &embed
}

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

// func generateEmbed(c string, l scryfallList) *discordgo.MessageEmbed {
// 	var embed discordgo.MessageEmbed
// }
