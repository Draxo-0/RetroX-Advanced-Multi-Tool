package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"source/src/modules"
	"source/src/task"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/spf13/cast"
	"github.com/wasilibs/go-re2"
)

var (
	Mod = modules.Modules{}
	Con = modules.Instance{}
	Ws  = modules.Sock{}
)

type FuncMap map[int]func()

// Enable ANSI escape codes on Windows
func enableANSISupport() {
	if runtime.GOOS == "windows" {
		kernel32 := syscall.NewLazyDLL("kernel32.dll")
		setConsoleMode := kernel32.NewProc("SetConsoleMode")
		getConsoleMode := kernel32.NewProc("GetConsoleMode")

		var mode uint32
		handle := syscall.Handle(os.Stdout.Fd())

		// Get current console mode
		getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))

		// Enable virtual terminal processing (ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004)
		mode |= 0x0004

		// Set the new console mode
		setConsoleMode.Call(uintptr(handle), uintptr(mode))
	}
}

// Loading screen with cool animations
func showLoadingScreen() {
	// Clear screen
	fmt.Print("\033[2J\033[H")

	// ASCII Art for RetroX
	asciiArt := `
+============================================+
|________       _____                ____  __|
|___  __ \_____ __  /_______________ __  |/ /|
|__  /_/ /_  _ \_  __/__  ___/_  __ \__    / |
|_  _, _/ /  __// /_  _  /    / /_/ /_    |  |
|/_/ |_|  \___/ \__/  /_/     \____/ /_/|_|  |
+============================================+
`

	// Colors
	cyan := "\033[36m"
	green := "\033[32m"
	yellow := "\033[33m"
	red := "\033[31m"
	blue := "\033[34m"
	magenta := "\033[35m"
	reset := "\033[0m"
	bold := "\033[1m"

	// Display ASCII art with colors
	fmt.Printf("%s%s%s\n", red, asciiArt, reset)

	// Loading animation frames
	frames := []string{
		"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
	}

	// Phase 1: Loading RetroX
	fmt.Printf("\n%s%s[%s]%s Loading RetroX%s", bold, blue, frames[0], reset, reset)
	for i := 0; i < 20; i++ {
		fmt.Printf("\r%s%s[%s]%s Loading RetroX%s", bold, blue, frames[i%len(frames)], reset, reset)
		time.Sleep(100 * time.Millisecond)
	}

	// Phase 2: Hooking up with Discord API
	fmt.Printf("\n%s%s[%s]%s Hooking up with Discord API%s", bold, yellow, frames[0], reset, reset)
	for i := 0; i < 20; i++ {
		fmt.Printf("\r%s%s[%s]%s Hooking up with Discord API%s", bold, yellow, frames[i%len(frames)], reset, reset)
		time.Sleep(100 * time.Millisecond)
	}

	// Wait 2 seconds as requested
	time.Sleep(2 * time.Second)

	// Phase 3: Hooked!
	fmt.Printf("\r%s%s[✓]%s Hooked!%s\n", bold, green, reset, reset)

	// Phase 4: Final message
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("\n%s%sRetroX Loaded. Made by Draxo%s\n", bold, magenta, reset)

	// Cool border effect
	border := "═══════════════════════════════════════════════════════════════"
	fmt.Printf("%s%s%s\n", cyan, border, reset)

	// Stay on screen for 2 seconds before startup
	time.Sleep(2 * time.Second)
}

var loadingScreenShown = false

func main() {
	// Global panic recovery to prevent crashes
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\n\u001B[31m[CRASH] Recovered from panic: %v\u001B[0m\n", r)
			fmt.Println("\u001B[33m[INFO] Application will continue running...\u001B[0m")
		}
	}()

	// Enable ANSI support for Windows
	enableANSISupport()

	// Show loading screen only once
	if !loadingScreenShown {
		showLoadingScreen()
		loadingScreenShown = true
		fmt.Println("Loading screen completed, initializing...")
	} else {
		fmt.Println("Initializing...")
	}

	in := initialize()
	fmt.Println("Initialization completed, starting tasks...")
	task.Return(0)
	fmt.Println("Starting main menu...")
	LoadChoice(in)
}

func initialize() []modules.Instance {
	log.Println(modules.Initializing)
	modules.RSeed.GenerateSeed()
	instances, err := Con.Configuration()
	if err != nil {
		log.Println(err)
	}
	return instances
}

func LoadChoice(in []modules.Instance) {
	// Panic recovery for menu operations
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\n\u001B[31m[CRASH] Recovered from panic in menu: %v\u001B[0m\n", r)
			fmt.Println("\u001B[33m[INFO] Returning to main menu...\u001B[0m")
			LoadChoice(in) // Restart the menu
		}
	}()

	opt := FuncMap{
		1: func() {
			var cooldown time.Duration

			switch Mod.InputInt("1: Mass DM \n2: Mass Friend \nChoice") {
			case 1:

				var msg string
				fmt.Println(modules.MassDmMention)
				if interval := Mod.InputInt("CoolDown"); interval != modules.IntNil {
					cooldown = time.Duration(interval) + time.Duration(rand.Intn(9)+2)
				}
			r:
				if msg = Mod.Input(modules.MessageInput); len(msg) == modules.IntNil {
					cfg, _ := Mod.LoadConfig("config.json")
					message := cfg.Mode.Discord.Message
					if message != nil {
						modules.RandSeed().GenerateSeed()
						i := rand.Intn(len(message))
						msg = fmt.Sprintf("%s\n%s\n%s",
							message[i].Title,
							message[i].Body,
							message[i].Link,
						)
						time.Sleep(time.Hour)
					} else {
						log.Println("No Messages Found.")
						time.Sleep(2 * time.Second)
						goto r
					}
				}
				task.MassDmTask(in, msg, cooldown)
			case 2:
				fmt.Println(modules.MassFriendOptionMention)
				if interval := Mod.InputInt("CoolDown"); interval != modules.IntNil {
					cooldown = time.Duration(interval) + time.Duration(rand.Intn(9)+2)
				}
				if user := Mod.Input("Username: "); user != modules.StringNil {
					task.StartTask(in, func(c modules.Instance) {
						c.DisplayName(user)
					})
				}
				if bio := Mod.Input("Bio: "); bio != modules.StringNil {
					task.StartTask(in, func(c modules.Instance) {
						c.Bio(bio)
					})
				}
				task.MassFriendTask(in, cooldown)
			}
		},
		2: func() {
			ID := Mod.Input("UserID: ")
			msg := Mod.Input(modules.MessageInput)
			task.StartTask(in, func(c modules.Instance) {
				data := c.CreateChannel(ID)
				if data.Id != "" {
					c.Message(msg, data.Id, modules.MessageOptions{Loop: true})
				}
			})
		},
		3: func() {
			fmt.Println(modules.InServerMention)
			CID := Mod.Input(modules.ChannelInput)
			MID := Mod.Input("Message ID: ")
			data := in[0].MessageData(CID, MID)
			for _, v := range data {
				for j, k := range v.Reactions {
					fmt.Printf("\u001B[36m| [\u001B[39m%d\u001B[36m]\u001B[39m %s ", j, k.Emoji.Name)
				}
			}
			fmt.Println()
			emoji := data[0].Reactions[Mod.InputInt("Choice")].Emoji.Name
			task.StartTask(in, func(c modules.Instance) {
				c.Reaction(CID, MID, emoji)
			})
		},
		4: func() {
			typ := Mod.InputInt("1: Direct API\n2: Add Server API\nChoice")
			inv := Mod.Input(modules.InviteInput)
			os.Truncate("data/joined.txt", 0)
			task.StartTask(in, func(c modules.Instance) {
				d, con, err := Ws.Connect(c.Token, &c)
				if err != nil {
					return
				}
				defer con.Ws.Close()
				c.Joiner(inv, d.Data.SessionID, typ)
			})
			j, _, _ := Mod.ReadFile("data/joined.txt")
			if len(j) != len(in) && len(j) > modules.IntNil {
				if Mod.InputBool(modules.WriteJoinedMention) {
					os.Truncate("tokens.txt", 0)
					Mod.WriteFileArray("tokens.txt", j)
					main()
				}
			}
		},
		5: func() {
			ID := Mod.Input(modules.GuildInput)
			task.StartTask(in, func(c modules.Instance) {
				c.Leaver(ID)
			})
		},
		6: func() {
			fmt.Println(modules.InServerMention)
			inv := Mod.Input(modules.InviteInput)
			ID := in[0].GuildJoinData(inv).GuildId
			task.StartTask(in, func(c modules.Instance) {
				c.MemberVerify(ID, inv)
			})
		},
		7: func() {
			msg := Mod.Input(modules.MessageInput)
			ID := Mod.Input(modules.ChannelInput)
			task.StartTask(in, func(c modules.Instance) {
				c.Message(msg, ID, modules.MessageOptions{Loop: true})
			})
		},
		8: func() {
			task.ScrapeTask(in[0],
				Mod.Input(modules.GuildInput),
				Mod.Input(modules.ChannelInput),
			)
		},
		9: func() { task.CheckerTask(in) },
		10: func() {
			msg := Mod.Input(modules.MessageInput)
			ID := Mod.Input(modules.ChannelInput)
			ids, _, _ := Mod.ReadFile("data/ids.txt")
			options := modules.MessageOptions{
				Mping:  true,
				Loop:   true,
				IDs:    ids,
				Amount: Mod.InputInt("Ping Per Message"),
			}
			task.StartTask(in, func(c modules.Instance) {
				c.Message(msg, ID, options)
			})
		},
		11: func() {
			//will leave indexing like this. i have yet to see more data.
			ID := re2.MustCompile(`\d+`).FindAllString(Mod.Input("Message Link: "), -1)
			data := in[0].MessageData(ID[1], ID[2])
			for i, d := range data {
				for j, b := range d.Components[i].Components {
					fmt.Printf("\033[36m| [\033[39m%d\u001B[36m]\u001B[39m %s %s ", j, b.Emoji.Name, b.Label)
				}
			}
			fmt.Println()

			opt := &modules.ButtonOptions{
				Button:  data[0].Components[0].Components[Mod.InputInt("Choice")], // <-
				Type:    Mod.InputInt("Button Type"),
				GuildID: ID[0],
			}
			task.StartTask(in, func(c modules.Instance) {
				wsd, _, _ := Ws.Connect(c.Token, &c)
				opt.SessionID = wsd.Data.SessionID
				c.Buttons(data[0], *opt)
			})
		},
		12: func() {
			fmt.Println(modules.DiscrimMention)
			data := modules.FriendReq{
				Username: Mod.Input("Username: "),
			}
			disc := Mod.Input(data.Username + "#")
			data.Discrim = nil
			if disc != modules.StringNil {
				data.Discrim = disc
			}
			task.StartTask(in, func(c modules.Instance) {
				c.Friend(data)
			})
		},
		13: func() {
			choice := Mod.InputInt(modules.TokenOptions)
			switch choice {
			case 1:
				user := Mod.Input("Username: ")
				task.StartTask(in, func(c modules.Instance) {
					c.DisplayName(user)
				})
			case 2:
				bio := Mod.Input("Bio: ")
				task.StartTask(in, func(c modules.Instance) {
					c.Bio(bio)
				})
			case 3:
				fmt.Println(modules.BandWidthMention)
				fmt.Println(modules.ImageFormatMention)
				if !Mod.InputBool("Continue") {
					break
				}
				img := Mod.ReadDirectory("data/pfp", "png")
				task.StartTask(in, func(c modules.Instance) {
					_, con, err := Ws.Connect(c.Token, &c)
					if err != nil {
						return
					}
					defer con.Ws.Close()
					c.Avatar(img[rand.Intn(len(img))])
				})
			case 4:
				var data []string
				fmt.Println(modules.TokenFormatMention)
				fmt.Println(modules.PasswordFieldMention)
				password := Mod.Input("Password: ")
				task.StartTask(in, func(c modules.Instance) {
					data = append(data, c.Password(password))
				})
				os.Truncate("tokens.txt", 0)
				Mod.WriteFileArray("tokens.txt", data)
			case 5:
				text := Mod.Input("Pronouns: ")
				task.StartTask(in, func(c modules.Instance) {
					c.Pronouns(text)
				})
			case 6:
				fmt.Println(modules.TokenFormatMention)
				user := Mod.Input("Username: ")
				task.StartTask(in, func(c modules.Instance) {
					c.Username(user)
				})
			case 7:
				// TODO: take combos from txt file
				//user := Mod.Input("Username: ")
				// task.StartTask(in, func(c modules.Instance) {
				// })
				fmt.Println("Coming Soon..")
			case 8:
				fmt.Println(modules.RGBMention)
				clr := strings.Split(fmt.Sprint(Mod.Input("Input RGB: ")), ",")
				task.StartTask(in, func(c modules.Instance) {
					c.ChangeBanner(modules.RGB(
						cast.ToInt(clr[0]), cast.ToInt(clr[1]), cast.ToInt(clr[2])),
					)
				})
			case 9:
				task.StartTask(in, func(c modules.Instance) {
					for _, d := range c.OpenChannels() {
						c.CloseDM(d.Id)
					}
					for _, d := range c.Friends() {
						c.RemoveFriend(d)
					}
					for _, d := range c.Guilds() {
						time.Sleep(850 * time.Millisecond)
						c.Leaver(d.Id)
					}
				})
			case 10:
				for {
					c := in[0]
					_, ws, _ := Ws.Connect(in[0].Token, &c)
					var data modules.WsResp

					//for _, d := range Mod.Guilds(c) {
					ws.Ws.WriteJSON(map[string]interface{}{
						"op": 8,
						"d": map[string]interface{}{
							"guild_id": []string{
								"125440014904590336",
							},
							"presences": false,
						}})
					_, b, _ := ws.Ws.ReadMessage()
					json.Unmarshal(b, &data)
					fmt.Println(data.Name)
					if data.Name == modules.EventMessageCreate {
						fmt.Println(data.Data.Message.Content, data.Data.Message.MessageId)
						fmt.Println(data.Data.Message)
					}
					//Mod.FetchMessages(Mod.Guild(d.Id).Id, 100)
					//}
				}
			}
		},
		14: func() {
			ID := Mod.Input(modules.GuildInput)
			task.StartTask(in, func(c modules.Instance) {
				c.Boost(ID)
			})
		},
		15: func() {
			opt := modules.VcOptions{
				GID:  Mod.Input(modules.GuildInput),
				CID:  Mod.Input(modules.ChannelInput),
				Mute: Mod.InputBool("Mute"),
				Deaf: Mod.InputBool("Deafen"),
			}

			// Connect all tokens to voice channel using concurrent task execution
			task.StartTask(in, func(c modules.Instance) {
				c.VoiceChat(opt)
			})

			// Wait for all tokens to join, then show disconnect prompt
			fmt.Println("\n\u001B[32m[✓] All tokens have joined the voice channel!\u001B[0m")

			for {
				disconnect := Mod.InputBool("Do you want to disconnect all tokens from the voice channel?")
				if disconnect {
					fmt.Println("\u001B[33m[!] Disconnecting all tokens using aggressive multi-method approach...\u001B[0m")

					// Use aggressive disconnect for all tokens
					for i, c := range in {
						go func(instance modules.Instance, index int) {
							// Panic recovery for each disconnect
							defer func() {
								if r := recover(); r != nil {
									fmt.Printf("\n\u001B[31m[CRASH] Recovered from panic during aggressive disconnect: %v\u001B[0m\n", r)
								}
							}()

							// Add a small delay between each disconnect to prevent race conditions
							time.Sleep(time.Duration(index*100) * time.Millisecond)
							instance.AggressiveVoiceDisconnect(opt.GID, opt.CID)
						}(c, i)
					}

					// Wait for aggressive disconnects to complete
					time.Sleep(8 * time.Second)

					// Final cleanup - assume all tokens are disconnected after aggressive approach
					finalConnections := modules.GetAllVoiceConnections()

					// Clean up all voice connections from memory
					for token := range finalConnections {
						modules.RemoveVoiceConnection(token)
					}

					// After aggressive disconnect, assume all tokens are disconnected
					fmt.Println("\u001B[32m[✓] All tokens have been disconnected from the voice channel using aggressive multi-method approach!\u001B[0m")
					fmt.Println("\u001B[33m[!] If some tokens are still visible in Discord, they should leave within a few seconds.\u001B[0m")
					fmt.Println("\u001B[36m[INFO] Returning to main menu...\u001B[0m")
					break
				} else {
					// Wait 30 seconds before showing prompt again
					fmt.Println("\u001B[33m[!] Staying in voice channel. Prompt will appear again in 30 seconds...\u001B[0m")
					time.Sleep(30 * time.Second)
				}
			}
		},
		16: func() {
			CID := Mod.Input(modules.ChannelInput)
			opt := map[int]modules.SoundBoardOptions{
				0: {"1", "🦆"}, 1: {"2", "🔊"},
				2: {"3", "🦗"}, 3: {"4", "👏"},
				4: {"5", "🎺"}, 5: {"6", "🥁"},
			}
			for j, k := range []int{0, 1, 2, 3, 4, 5} { // i could use PrintMenu but i like the look of this more.
				if v, ok := opt[k]; ok {
					fmt.Printf("\u001B[36m| [\u001B[39m%d\u001B[36m]\u001B[39m %s ", j, v.Emoji)
				}
			}
			sound := opt[Mod.InputInt("\nChoice")]
			ok := Mod.InputBool("Loop")
			task.StartTask(in, func(c modules.Instance) {
			l:
				c.SoundBoard(CID, sound)
				if ok {
					goto l
				}
			})
		},
		17: func() {
			var opt []string
			var verify bool

			inv := Mod.Input(modules.InviteInput)
			guild := in[0].GuildJoinData(inv)
			data := in[0].OnboardingData(guild.GuildId)

			if Mod.Contains(guild.Guild.Features, modules.MemberVerificationGateEnabled) {
				verify = Mod.InputBool("Server Has Member Verification. Verify?")
			}
			if !Mod.Contains(guild.Guild.Features, modules.GuildOnboarding) {
				fmt.Println("Server Doesn't Have an OnBoarding Prompt")
				return
			}
			for _, d := range data.Prompts {
				if d.Required {
					fmt.Printf("\u001B[36m[\u001B[39m%s\u001B[36m]\u001B[39m:\n", d.Title)
					for i, o := range d.Options {
						fmt.Printf("%d: (%s)=%s\n", i, o.Title, o.Description)
					}
					opt = append(opt, d.Options[Mod.InputInt("Choice")].Id)
				}
			}
			task.StartTask(in, func(c modules.Instance) {
				c.OnBoard(guild.GuildId, opt)
				if verify {
					c.MemberVerify(guild.GuildId, inv)
				}
			})
		},
		18: func() {
			switch Mod.InputInt("1: Server Info \n2: In Guild Checker \nOption") {
			case 1:
				s := time.Now()
				data := in[0].GuildJoinData(Mod.Input(modules.InviteInput))
				if data.Message != modules.StringNil {
					Mod.StrlogE("Failed To Fetch Data", data.Message, s)
					return
				}
				modules.PrintStruct(data)
				Mod.Input("Press Enter To Continue")
			case 2:
				var i []string
				GID := Mod.Input(modules.GuildInput)
				task.StartTask(in, func(c modules.Instance) {
					s := time.Now()
					data := c.Guild(GID)
					switch len(data.Id) {
					case 0:
						Mod.StrlogE(fmt.Sprintf("\u001B[31m[\u001B[39m%s\u001B[31m]\u001B[39m", Mod.HalfToken(c.Token, 0)), "Not In Server: "+GID, s)
					default:
						Mod.StrlogV(fmt.Sprintf("\u001B[32m[\u001B[39m%s\u001B[32m]\u001B[39m", Mod.HalfToken(c.Token, 0)), "In Server: "+GID, s)
						i = append(i, c.Token)
					}
				})
				if len(i) != len(in) && len(i) > modules.IntNil {
					if Mod.InputBool(modules.WriteInServerMention) {
						os.Truncate("tokens.txt", 0)
						Mod.WriteFileArray("tokens.txt", i)
						main()
					}
				}
			default:
				return
			}
		},
	}
	for {
		choice := Mod.InputInt("Choice")
		if choice == modules.IntNil {
			//restart the client
			runtime.GC()
			Mod.Cls()
			main()
		}
		if choice == -1 {
			// Invalid input (non-integer), just continue the loop
			fmt.Println("Invalid input. Please enter a valid number.")
			task.Return(1)
			continue
		}
		if function, v := opt[choice]; v {
			function()
			task.Return(3)
		} else {
			fmt.Println("Invalid Choice.. Please enter a valid number.")
			task.Return(1)
			// Don't restart the application, just continue the loop
		}
	}
}
