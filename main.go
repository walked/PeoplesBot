package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type scryfallList struct {
	CardList []card `json:"data"`
}
type legalities struct {
	Vintage string `json:"vintage"`
}

type images struct {
	Small  string `json:"small"`
	Normal string `json:"normal"`
	Large  string `json:"large"`
}
type card struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Object     string     `json:"card"`
	Rarity     string     `json:"rarity"`
	Reserved   bool       `json:"reserved"`
	SetType    string     `json:"set_type"`
	Oracle     string     `json:"oracle_text"`
	Legalities legalities `json:"legalities"`
	Image_uris images     `json:"image_uris"`
	Type       string     `json:"type_line"`
	Manacost   string     `json:"mana_cost"`
}

// type bannedCard struct {
// 	Name string `json:"name"`
// }

func main() {

	// NEW DISCORD SESSION AND HANDLERS
	//
	token := os.Getenv("DISCORD_TOKEN")
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error with Discord session, ", err)
		return
	}
	dg.AddHandler(messageCreate)
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	fmt.Println("Bot now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()

}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// if strings.HasPrefix(m.Content, "?legal") {
	// 	fmt.Println("Content:" + m.Content)
	// 	cardname := strings.TrimPrefix(m.Content, "?legal ")
	// 	fmt.Println(cardname)
	// 	cleanInput(cardname)

	// 	// Pull cardname into a struct of scryfallList type
	// 	sl := queryScryfall(cardname)

	// 	if len(sl.CardList) < 1 {
	// 		// false indicates not found
	// 		s.ChannelMessageSendEmbed(m.ChannelID, generateEmbed(cardname, false, sl))
	// 	} else {
	// 		// true indicates found in scryfall database
	// 		s.ChannelMessageSendEmbed(m.ChannelID, generateEmbed(cardname, true, sl))
	// 	}

	// }
	if strings.HasPrefix(m.Content, "?proxy") || strings.HasPrefix(m.Content, "?proxies") {
		s.ChannelMessageSendEmbed(m.ChannelID, proxies())
	}

	//// REWRITE STARTS HERE
	if strings.HasPrefix(m.Content, "!legal") || strings.HasPrefix(m.Content, "?legal") {
		//fmt.Println("Content:" + m.Content)

		cardname := strings.TrimPrefix(m.Content, "?legal ")
		cardname = strings.TrimPrefix(cardname, "!legal ")
		cleanName, urlName := cleanInput(cardname)

		fmt.Println(cleanName + " " + urlName)
		match, matches := scryFallMatch(cleanName, urlName)

		fmt.Printf("match list length: %v", len(match.CardList))
		if len(match.CardList) > 0 {
			//checkBanned(match.CardList[0].Name)
			//checkLegal(&match.CardList)
			//checkLegality(&match.CardList)
			s.ChannelMessageSendEmbed(m.ChannelID, newEmbed(cleanName, urlName, true, match))
			//checklegality(banned, general)
			//sendembed, cardname, banned, generally legal, list
		} else if len(matches.CardList) > 0 {
			var sl scryfallList
			sl.CardList = append(sl.CardList, matches.CardList[0])
			s.ChannelMessageSendEmbed(m.ChannelID, newEmbed(cleanName, urlName, true, &sl))
			//s.ChannelMessageSendEmbed(m.ChannelID, generateEmbed(cardname, true, match))
		} else {
			s.ChannelMessageSendEmbed(m.ChannelID, newEmbed(cleanName, urlName, false, matches))

		}

	}

}
