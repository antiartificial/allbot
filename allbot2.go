// allbot2.go
package main

import (
	"fmt"
	"strings"
	irc "github.com/fluffle/goirc/client"
	//"encoding/json"
	//"io/ioutil"
)

var debug = true;
var botName = "saruka"
var ircServer = "irc.sylnt.us:6667"
var channels = []string{"#", "#grue"}
var users = []string{"arti", "ciri"} //TODO Struct with permissions/types
var autoJoin = true
var muzzle = false
var doOnJoin = true

//TODO arrays/maps how do they work
/*type IRCSetting struct {
	BotName = "saruka"
	IRCServer = "irc.sylnt.us:6667"
	Channels = []string{"#", "#grue"}
	users = []string{"arti", "ciri"}
}*/

func main() {
	fmt.Println("Allbotv2 loading...")
	
	// TODO persistent settings
	/*r, err := ioutil.ReadFile("settings.json")
	if err != nil {
		panic(err)
	}
	setting := IRCSetting{}*/
	
	// Do our connection settings
    cfg := irc.NewConfig(botName)
    cfg.SSL = false
    cfg.Server = ircServer
    cfg.NewNick = func(n string) string { return n + "^" }
    c := irc.Client(cfg)
	
	// Begin handlers
	c.HandleFunc("connected",
        func(conn *irc.Conn, line *irc.Line) {
			fmt.Printf("Connected to %s...\n", cfg.Server)
			
			if autoJoin {
				if len(channels) > 0 {
					for index,channel := range channels {
						fmt.Printf("Joining [%s] [%d/%d]\n", channel, index, len(channels)-1)
						conn.Join(channel)
					}
				}
			} else {
				fmt.Println("What's next?")
			}
		})
		
	c.HandleFunc("join",
		func(conn *irc.Conn, line *irc.Line) {
			if doOnJoin {
				//TODO send greeting
				eventJoin(conn, line)
			}
		})
		
	c.HandleFunc("privmsg",
		func(conn *irc.Conn, line *irc.Line) {
			fmt.Printf("[%s] %s : %s\n", line.Cmd, line.Nick, line.Args )
			
			// determine if channel message or addressing bot
			if line.Args[0] != botName {
				//fmt.Println("^---- Channel message -----^")	
				if muzzle == false {
					eventPrivmsg(conn, line, true)
				}	
			} else {
				//fmt.Println("^---- Priavte message -----^")
				eventPrivmsg(conn, line, false)
			}
			
		})
		
    // And a signal on disconnect
    quit := make(chan bool)
    c.HandleFunc("disconnected",
        func(conn *irc.Conn, line *irc.Line) {
			// TODO reconnect upon disconnect signal
			quit <- true 
		})
	
	// Tell client to connect.
    if err := c.Connect(); err != nil {
        fmt.Printf("Connection error: %s\n", err)
    }
	
    // Wait for disconnect
    <-quit
}

func eventJoin(conn *irc.Conn, line *irc.Line) {
	if debug {
		fmt.Printf("Event Join fired: [%s]", line)
	}
	
	if botName != line.Nick {
		conn.Notice(line.Src, line.Nick+", welcome");
	}
}

// Called when incoming privmsg
// channel boolean evaluates true for channel privmsgs, false for bot privmsgs
func eventPrivmsg(conn *irc.Conn, line *irc.Line, channel bool) {
	if debug {
		fmt.Printf("Event PrivMSG fired: [%s] [%b]\n", line, channel)
	
	}
	if channel {
		if debug {
			fmt.Printf("Line args: [%s]\n", line.Args[1])
		}
		
		var lastLine = strings.Split(line.Args[1], " ")
		l := lastLine[0];
		select {
		case l == strings.ToLower("!say"):
			triggerSay(conn, lastLine, line.Args[0]) // why can't we pass line.Args?
			fallthrough
		default:
			fmt.Println("No match.")
		}
	} else {
		
	}
}

func triggerSay(conn *irc.Conn, lastLine []string, channel string) {
	if debug {
		fmt.Printf("Trigger Say fired: [%s]\n", lastLine)
	}
	
	//TODO ponder if sanity checking is better per function or @higher level
	conn.Privmsg(channel, strings.Join(lastLine[1:len(lastLine)], " "))
}