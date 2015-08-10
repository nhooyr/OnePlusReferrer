package main

import (
	"bufio"
	"bytes"
	"code.google.com/p/go-imap/go1/imap"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"net/url"
	"regexp"
	"sync"
	"time"
)

var userList []string

func insertDots(str string, at int) {
	if at == 0 {
		userList = append(userList, str)
		return
	}
	newStr := str[:at] + "." + str[at:]
	for i := 0; at > i; i = i + 2 {
		insertDots(newStr, i)
	}
}

func sendRequests(user string) {
	defer wg.Done()
	reqTimeout := 15
	var successCounter int
	for i, user := range userList[1000:] {
		log.Println("current index in userList", i)
		log.Println("current reqTimeout", reqTimeout)
		for {
			log.Println("sending request git remote add origin https://github.com/aubble/goReferrer.git
			git push -u origin masterfor", user)
			c, err := tls.Dial("tcp", "invites.oneplus.net:443", nil)
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second * time.Duration(reqTimeout))
				continue
			}
			_, err = c.Write(confirmURL.ReplaceAll(reqTemplate, []byte("https://invites.oneplus.net/index.php?r=share/signup&success_jsonpCallback=success_jsonpCallback&email="+user+"%40gmail.com&koid=IR0T3R&_=1439186651073"))) // might have to replace cache buster
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second * time.Duration(reqTimeout))
				continue
			}
			resp, err := http.ReadResponse(bufio.NewReader(c), nil)
			if err != nil {
				log.Println(err)
				time.Sleep(time.Second * time.Duration(reqTimeout))
				continue
			}
			if resp.StatusCode == 200 {
				log.Println("successfully sent request for", user)
				if successCounter == 10 {
					successCounter = 0
					reqTimeout--
				} else {
					successCounter++
				}
				break
			} else {
				log.Println("god help the hamsters for", user)
				reqTimeout++
				successCounter = 0
			}
			time.Sleep(time.Second * time.Duration(reqTimeout))
		}
		time.Sleep(time.Second * time.Duration(reqTimeout))
	}
}

func spoofLinks(user string) {
	defer wg.Done()
	for {
		connectNet(user)
	}
}

func connectNet(user string) {
	netTimeout := time.Second * 30
	cl, err := imap.DialTLS("imap.gmail.com", nil)
	if err != nil {
		log.Println(err)
		time.Sleep(netTimeout)
		return
	}
	defer cl.Logout(-1)
	log.Println("server says hello:", cl.Data[0].Info)
	cl.Data = nil
	log.Println("logging in")
	if cl.State() == imap.Login {
		_, err = cl.Login(user+"@gmail.com", "PASSHERE")
		if err != nil {
			log.Println(err)
			time.Sleep(netTimeout)
			return
		}
	}
	log.Println("logged in")
	log.Println("selecting INBOX")
	_, err = cl.Select("INBOX", false)
	if err != nil {
		log.Println(err)
		time.Sleep(netTimeout)
		return
	}
	log.Println("selected INBOX")
	set, _ := imap.NewSeqSet("")
searchLoop:
	for {
		var successCounter int
		conTimeout := 5
		log.Println("searching INBOX")
		cmd, err := imap.Wait(cl.Search("SUBJECT \"Confirm your email\""))
		if err != nil {
			log.Println(err)
			break searchLoop
		}
		searchResults := cmd.Data[0].SearchResults()
		if searchResults == nil {
			log.Println("no email[s] found")
			time.Sleep(time.Second * 5)
			break searchLoop
		} else {
			log.Println("found", len(searchResults), "emails to spoof, results:", searchResults)
			set.AddNum(searchResults...)
		}
		var rsp *imap.Response
		link := regexp.MustCompile(`https:=\r\n\/\/invites.oneplus.net\/confirm\/[0-9A-Z]*.`)
		var UIDCounter int
		log.Println("beginning spoof confirmation loop")
		cmd, err = cl.Fetch(set, "RFC822")
		if err != nil {
			log.Println(err)
			break searchLoop
		}
		setDone, _ := imap.NewSeqSet("")
	getMailLoop:
		for cmd.InProgress() {
			cl.Recv(-1)
			for _, rsp = range cmd.Data {
				message := imap.AsBytes(rsp.MessageInfo().Attrs["RFC822"])
				if msg, err := mail.ReadMessage(bytes.NewReader(message)); msg != nil {
					log.Println("current conTimeout", conTimeout)
					body, err := ioutil.ReadAll(msg.Body)
					if err != nil {
						log.Println(err)
						break searchLoop
					}
					linkRawM := link.FindSubmatch(body)
					var linkRaw []byte
					if len(linkRawM) > 0 {
						linkRaw = linkRawM[0]
					} else {
						log.Println("found nolink mail")
						continue getMailLoop
					}
					log.Println(string(linkRaw))
					url, err := url.Parse(string(linkRaw[:6]) + string(linkRaw[9:]))
					if err != nil {
						log.Println(err)
						break searchLoop
					}
					log.Println(url)
					req := confirmURL.ReplaceAll(reqTemplate, []byte(url.Path)[:len(url.Path)-1])
					for {
						log.Println("attempting to spoof confirmation for", searchResults[UIDCounter])
						c, err := tls.Dial("tcp", url.Host+":443", nil)
						if err != nil {
							log.Println(err)
							break searchLoop
						}
						_, err = c.Write(req)
						if err != nil {
							log.Println(err)
							break searchLoop
						}
						resp, err := http.ReadResponse(bufio.NewReader(c), nil)
						if err != nil {
							log.Println(err)
							break searchLoop
						}
						if resp.StatusCode == 302 {
							redirectReg := regexp.MustCompile(`https:\/\/oneplus\.net\/invites\?kid=[0-9A-Z]*`)
							redirectURL, err := resp.Location()
							if err != nil {
								log.Println(err)
								break searchLoop
							}
							log.Println(resp.Location())
							if redirectReg.MatchString(redirectURL.String()) {
								log.Println("successfully spoofed confirmation for", searchResults[UIDCounter])
								setDone.AddNum(searchResults[UIDCounter])
								UIDCounter++
								if successCounter == 10 {
									conTimeout--
									successCounter = 0
								}
								break
							} else {
								log.Println("server trying to migitate spoofing")
								conTimeout++
								successCounter = 0
							}
						} else {
							log.Println("gg rip hamsters")
							conTimeout++
							successCounter = 0
						}
						time.Sleep(time.Second * time.Duration(conTimeout))
					}
				} else {
					log.Println(err)
				}
			}
			time.Sleep(time.Second * time.Duration(conTimeout))
			cmd.Data = nil
		}
		log.Println("deleting mail", setDone)
		_, err = imap.Wait(cl.Store(setDone, "+FLAGS", "\\DELETED"))
		if err != nil {
			log.Println(err)
			break searchLoop
		}
		_, err = imap.Wait(cl.Expunge(nil))
		if err != nil {
			log.Println(err)
			break searchLoop
		}
	}
}

var wg sync.WaitGroup
var reqTemplate []byte
var confirmURL = regexp.MustCompile("confirmURL")

func main() {
	user := "referralreferralreferralreferr"
	for i := 0; len(user) > i; i = i + 2 {
		insertDots(user, i)
	}
	log.Println(userList)
	var err error
	reqTemplate, err = ioutil.ReadFile("request")
	if err != nil {
		log.Println(err)
		return
	}
	wg.Add(2)
	go sendRequests(user)
	go spoofLinks(user)
	wg.Wait()
}

//todo closing and make seperate functions, resume index
