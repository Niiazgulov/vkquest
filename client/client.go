package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Client struct {
	endpoint string
}

func NewClient(baseurl string) (ClientInterf, error) {
	return &Client{baseurl}, nil
}

func (c *Client) AddUser(user string) (int, error) {
	users := UserQuestJSON{User_name: user}
	data, err := json.Marshal(users)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	url := c.endpoint + "adduser"
	userJSON, err := requestServer(data, http.MethodPost, url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return userJSON.User_id, nil
}

func (c *Client) AddQuest(quest string, cost int) (int, error) {
	quests := UserQuestJSON{Quest_name: quest, Cost: cost}
	data, err := json.Marshal(quests)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	url := c.endpoint + "addquest"
	questJSON, err := requestServer(data, http.MethodPost, url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return questJSON.Quest_id, nil
}

func (c *Client) CompleteQuest(user_id int, quest_id int) (string, error) {
	ids := UserQuestJSON{User_id: user_id, Quest_id: quest_id}
	data, err := json.Marshal(ids)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	url := c.endpoint + "action"
	questJSON, err := requestServer(data, http.MethodPost, url)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return questJSON.Quest_name, nil
}
func (c *Client) GetAllInfo(user_id string) ([]UserQuestJSON, error) {
	data := url.Values{}
	data.Set("id", user_id)

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodGet, c.endpoint+user_id, bytes.NewBufferString(data.Encode()))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var users []UserQuestJSON
	err = json.Unmarshal(body, &users)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return users, nil
}

func requestServer(data []byte, method string, url string) (UserQuestJSON, error) {
	client := &http.Client{}
	request, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var uqJSON UserQuestJSON
	err = json.Unmarshal(body, &uqJSON)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return uqJSON, nil
}
