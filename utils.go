package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/sacOO7/gowebsocket"
)

func GetBots(number int) (Response, error) {
	bots := Response{}

	resp, err := http.Get("https://quizizz.com/_api/ratelimit/v2/adminRecommend/2?page=2&pageSize=" + strconv.Itoa(number) + "&_=1632765352057")
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return Response{}, errors.New("response failed")
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	if err := json.Unmarshal(out, &bots); err != nil {
		return Response{}, err
	}

	return bots, nil
}

func Spam(bots Response, config Config) {
	log.Println("Started sending bots")
	fmt.Println("Press ctrl + c to cancel")

	var botNames string
	for i, bot := range bots.Data.Quizzes {
		if bot.CreatedBy.Id == "" {
			continue
		}

		switch config.Mode {
		case 0:
			botNames = fmt.Sprintf("%s %s", bot.CreatedBy.FirstName, bot.CreatedBy.LastName)
		case 1:
			botNames = fmt.Sprintf("%s %d", config.CustomName, i)
		}

		socket := gowebsocket.New("wss://socket.quizizz.com/_gsocket/sockUpg/?experiment=authRevamp&EIO=4&transport=websocket")
		socket.OnConnected = func(socket gowebsocket.Socket) {
			socket.SendText("40")
		}
		socket.Connect()

		socket.SendText(`420["v5/join",{"roomHash":"` + config.RoomHash + `","player":{"id":"` + botNames + `","origin":"web","isGoogleAuth":false,"avatarId":2,"startSource":"gameCode|typed","name":"` + strconv.Itoa(i) + `","mongoId":"` + bot.CreatedBy.Id + `","userAgent":"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/94.0.4606.61 Safari/537.36","uid":"f6f597d7-1bbe-42e7-ac3e-600487250ed1","expName":"ratelimiter_exp","expSlot":"6"},"powerupInternalVersion":"13","__cid__":"v5/join.|1.1632599434062"}] `)
		time.Sleep(time.Duration(config.Delay) * time.Millisecond)

		socket.Close()
	}
	log.Println("all bots sent")
	os.Exit(0)
}

func GetRoomHash(pin string) (string, error) {
	var data = []byte(fmt.Sprintf(`{"roomCode":"%s"}`, pin))

	req, err := http.NewRequest("POST", "https://game.quizizz.com/play-api/v5/checkRoom", bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("invalid status code")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	room := RoomResp{}

	if err := json.Unmarshal(body, &room); err != nil {
		return "", err
	}
	if room.Room.Hash == "" {
		return "", errors.New("invalid pin")
	}

	return room.Room.Hash, nil
}
