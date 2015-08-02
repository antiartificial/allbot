// allbot2.go
package main

import (
	"fmt"
	"strings"
	irc "github.com/fluffle/goirc/client"
	//"encoding/json"
	//"io/ioutil"
	"os"
	"io"
	"log"
	"bufio"
)

var debug = true;
var botName = "saruka"
var ircServer = "irc.sylnt.us:6667"
var channels = []string{"#grue"}
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
		
	c.HandleFunc("ctcp",
		func(conn *irc.Conn, line *irc.Line) {
			fmt.Printf("[%s] %s : %s\n", line.Cmd, line.Nick, line.Args )	
			
			if debug {
				fmt.Printf("Line args: [%s]\n", line.Args[0])
			}
			
			if line.Args[0] == "VERSION" {
				fmt.Printf("CTCP Reply sent to: [%s]\n", line.Nick)
				conn.CtcpReply(line.Nick, "VERSION", "allbot")
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
		fmt.Printf("Event Join fired: [%s]\n", line)
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
	
	// Determine what type of privmsg it is (channel||to bot)
	if channel {
		if debug {
			fmt.Printf("Line args: [%s]\n", line.Args[1])
		}
		
		var lastLine = strings.Split(line.Args[1], " ")

		if lastLine[0] == strings.ToLower("!say") {
			triggerSay(conn, lastLine, line.Args[0]) // why can't we pass line.Args?
		}
		
		if lastLine[0] == strings.ToLower("!ll") {
			var temp = findNickLastLine("testfile.txt", lastLine[1])
			fmt.Println("line:"+temp)
			conn.Privmsg(line.Args[0], temp)
		}
		
		// Silence Bot
		if lastLine[0] == strings.ToLower("!muzzle") {
			if muzzle == false {
				muzzle = true
				channelSay(conn, line.Args[0], "Muzzle enabled")
			} else {
				muzzle = false
				channelSay(conn, line.Args[0], "Muzzle disabled")
			}
			fmt.Printf("Muzzle set to: [%s]\n", muzzle)
		}
		
		if lastLine[0] == strings.ToLower("!op") {
			if line.Nick != botName && line.Nick == strings.Contains(users, line.Nick) {
				//TODO do op
			}
			//TODO do unauthorized message?
		}
		
		/*select {
		case strings.ToLower("!say"):
			triggerSay(conn, lastLine, line.Args[0]) // why can't we pass line.Args?
			fallthrough
		default:
			fmt.Println("No match.")
		}*/
	} else {
		// PrivMSG is to bot.
		if debug {
			fmt.Printf("Line args: [%s]\n", line.Args[1])
		}
		//TODO bot privmsg handling
	}
}

func channelSay(conn *irc.Conn, channel string, text string) {
	if debug {
		fmt.Printf("Channel Say fired:\n")
	}
	
	conn.Privmsg(channel, text)
}

func triggerSay(conn *irc.Conn, lastLine []string, channel string) {
	if debug {
		fmt.Printf("Trigger Say fired: [%s]\n", lastLine)
	}
	
	//TODO ponder if sanity checking is better per function or @higher level
	conn.Privmsg(channel, strings.Join(lastLine[1:len(lastLine)], " "))
}

// Returns last line containing nick string
func findNickLastLine(path string, nick string) (line string) {
	if debug {
		fmt.Printf("FindNickLastLine fired: [%s] [%s]\n", path, nick)
	}
	
    f, err := os.Open(path)
    if err != nil {
        log.Fatal(err)
    }
    bf := bufio.NewReader(f)
	
	var lastLineWithNick string
	
    for {
        switch line, err := bf.ReadString('\n'); err {
        case nil:
            // valid line, echo it.  note that line contains trailing \n.
			if strings.Contains(line, nick) { //TODO regex?
				lastLineWithNick = line
            	fmt.Print(line)
			}
        case io.EOF:
            if line > "" {
                // last line of file missing \n, but still valid
                fmt.Println(line)
            }
        default:
            log.Fatal(err)
        }
    }
	
	fmt.Printf("Last line with [%s] is [%s]\n", nick, lastLineWithNick)
	return lastLineWithNick
	
}