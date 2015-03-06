// Copyright (c) 2015 Bertrand Janin <b@janin.com>
// Use of this source code is governed by the ISC license in the LICENSE file.

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/SlyMarbo/rss"
)

var (
	errDisconnected = errors.New("server disconnected")
)

type ByDate []*rss.Item

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].Date.Before(a[j].Date) }

func broadcast(notifications []string) error {
	log.Printf("connecting...")
	conn, err := connect()
	if err != nil {
		return err
	}

	sendLine(conn, fmt.Sprintf("NICK %s", cfg.IRCNickname))
	sendLine(conn, fmt.Sprintf("USER %s localhost "+
		"127.0.0.1 :%s\r\n", cfg.IRCNickname,
		cfg.IRCNickname))

	bufReader := bufio.NewReader(conn)

	for {
		data, err := bufReader.ReadString('\n')
		if err == io.EOF {
			return errDisconnected
		}
		if err != nil {
			log.Fatal(err)
		}
		code := parseServerMessageCode(data)

		if code > 1 {
			return errors.New(fmt.Sprintf("server message code: %d", code))
		}

		break
	}

	time.Sleep(500 * time.Millisecond)

	for _, c := range cfg.Channels {
		for _, n := range notifications {
			sendLine(conn, fmt.Sprintf("PRIVMSG %s :%s", c, n))
		}
	}

	return nil
}

func main() {
	var notifications []string

	parseCommandLine()
	err := parseConfigFile()
	if err != nil {
		log.Fatal("config error: ", err.Error())
	}

	for _, topic := range cfg.Topics {
		var latest time.Time
		hasUpdates := false

		f, err := os.Open(topic + ".latest")
		if err != nil {
			if !os.IsNotExist(err) {
				log.Fatal(err)
			}
			latest = time.Now()
			hasUpdates = true
		} else {
			defer f.Close()
			data := make([]byte, 64)
			_, err = f.Read(data)
			if err != nil {
				log.Fatal(err)
			}
			s := strings.Trim(string(data), "\n\r\t\x00 ")
			latest, err = time.Parse(time.RFC3339, s)
			if err != nil {
				log.Fatal(err)
			}
		}

		url := fmt.Sprintf("http://status.aws.amazon.com/rss/%s.rss", topic)
		feed, err := rss.Fetch(url)
		if err != nil {
			if strings.Contains(err.Error(), "no feeds found") {
				continue
			}
			log.Fatal(topic, err)
		}

		sort.Sort(ByDate(feed.Items))

		for _, item := range feed.Items {
			if item.Date.After(latest) {
				var text string
				if item.Content != "" {
					text = item.Content
				} else {
					text = item.Title
				}
				n := fmt.Sprintf("%s: %s", topic, text)
				notifications = append(notifications, n)
				latest = item.Date
				hasUpdates = true
			}
		}

		if hasUpdates {
			f, err := os.Create(topic + ".latest")
			if err != nil {
				log.Fatal(err)
			}
			data := latest.Format(time.RFC3339)
			f.Write([]byte(data))
		}
	}

	if len(notifications) > 0 {
		log.Printf("NOTIFICATIONS: %d", len(notifications))
		err = broadcast(notifications)
		if err != nil {
			log.Fatal(err)
		}
	}

	os.Exit(0)
}
