package client

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func Start() {
	cfg, err := NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	client, err := NewClient(cfg.BaseURLAddress)
	if err != nil {
		log.Fatal("Error while creating Client")
	}

	var actionchoice string
	fmt.Println("Введите цифру нужного действия:\n 1. Создать нового пользователя \n 2. Создать новое задание \n 3. Произошло событие - выполнение задания \n 4. История заданий и баланс пользователя")
	fmt.Fscan(os.Stdin, &actionchoice)

	switch actionchoice {
	case "1":
		fmt.Println("Введите имя пользователя:")
		var user string
		fmt.Fscan(os.Stdin, &user)
		userID, err := client.AddUser(user)
		if err != nil {
			log.Fatal("Error while AddUser")
		}
		fmt.Printf("Добавлен пользователь %s с ID: %d \n", user, userID)

	case "2":
		fmt.Println("Введите название задания:")
		var quest string
		fmt.Fscan(os.Stdin, &quest)
		fmt.Println("Введите стоимость награды за задание:")
		var cost string
		fmt.Fscan(os.Stdin, &cost)

		costint := checkStrInt(cost)
		questID, err := client.AddQuest(quest, costint)
		if err != nil {
			log.Fatal("Error while AddQuest")
		}
		fmt.Printf("Добавлено задание %s с ID: %d \n", quest, questID)

	case "3":
		fmt.Println("Введите ID пользователя, завершившего задание:")
		var userID string
		fmt.Fscan(os.Stdin, &userID)
		userIDint := checkStrInt(userID)

		fmt.Println("Введите ID завершенного задание:")
		var questID string
		fmt.Fscan(os.Stdin, &questID)
		questIDint := checkStrInt(questID)

		quest_name, err := client.CompleteQuest(userIDint, questIDint)
		if err != nil {
			log.Fatal("Error while CompleteQuest")
		}
		fmt.Printf("Пользователь с ID %d успешно завершил задание %s с ID: %d \n", userIDint, quest_name, questIDint)

	case "4":
		fmt.Println("Введите ID пользователя для просмотра его истории и баланса:")
		var userID string
		fmt.Fscan(os.Stdin, &userID)

		users, err := client.GetAllInfo(userID)
		if err != nil {
			log.Fatal("Error while GetAllInfo")
		}

		balance := users[0].Balance

		fmt.Printf("Пользователь с ID %s имеет баланс: %d \n", userID, balance)
		fmt.Println("Завершенные задания:")
		for i, v := range users[1:] {
			fmt.Printf("%d. %s \n", i+1, v.Quest_name)
		}

	default:
		fmt.Println("Попробуйте еще раз. Нужно ввести цифру от 1 до 4")
	}
}

func checkStrInt(str string) int {
	res := 0
	var err error
	for res <= 0 {
		res, err = strconv.Atoi(str)
		if err != nil {
			fmt.Println("Введите ЧИСЛО:")
			fmt.Fscan(os.Stdin, &str)
			continue
		}
		if res > 0 {
			break
		}
		fmt.Println("Введите число БОЛЬШЕ 0:")
		fmt.Fscan(os.Stdin, &str)
	}

	return res
}
