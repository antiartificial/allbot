// allbot.go
package main

import (
	"fmt"
	"github.com/thoj/go-ircevent"
	"strings"
	//"time"
)

var roomNames = []string{"#", "#grue"}
var ownerName = "arti"
var authUsers = []string{"arti", "ciri"}

func main() {
	fmt.Println("Allbot starting...")
	var ircobj = irc.IRC("saruka", "saruka")      //Create new ircobj
	err := ircobj.Connect("irc.sylnt.us:6667") //Connect to server

	if err != nil {
		fmt.Println("Connection failed.")
		return
	}
	
	ircobj.AddCallback("001", func (e *irc.Event) {
        //ircobj.Join(roomName)
		fmt.Println("Joining Channels...")
		for index,element := range roomNames {
			fmt.Println(index, element)
			ircobj.Join(element)
			//ircobj.SendRaw("JOIN " + element)
		}
    })
	
	ircobj.AddCallback("JOIN", func (e *irc.Event) {
		//for _,element := range roomNames {
        //	ircobj.Privmsg(element, "hello")
		//}
		ircobj.Privmsg(e.Arguments[0], "hello")
    })

	ircobj.AddCallback("PRIVMSG", func(event *irc.Event) {
		//event.Message() contains the message
		//event.Nick Contains the sender
		//event.Arguments[0] Contains the channel
		fmt.Println(event.Arguments[0], "<" + event.Nick + "> :", event.Message())
		
		if event.Nick == strings.ToLower(ownerName) {
			fmt.Println("Owner")
		}
		
		if strings.Contains(event.Message(),"!say") {
			var sayMessage = strings.Split(event.Message()," ")
			var resultMessage = sayMessage[1:len(sayMessage)]
			ircobj.Privmsg(event.Arguments[0], strings.Join(resultMessage, " "))
		}
	})
	ircobj.Loop()
}

/*func doSay(ircobj, message, channel) {
	if strings.Contains(message,"!say") {
		var sayMessage = strings.Split(message," ")
		var resultMessage = sayMessage[1:len(sayMessage)]
		ircobj.Privmsg(channel, strings.Join(resultMessage, " "))
	}
}*/