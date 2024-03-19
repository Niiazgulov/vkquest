package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DataBase struct {
	DB *sql.DB
}

func NewDB(dbPath string) (QuestDB, error) {
	db, err := sql.Open("pgx", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id SERIAL PRIMARY KEY,
			user_name VARCHAR, 
			balance INTEGER)
		`)
	if err != nil {
		return nil, fmt.Errorf("unable to CREATE TABLE users in DB: %w", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS quests (
		quest_id SERIAL PRIMARY KEY,
		quest_name VARCHAR UNIQUE, 
		cost INTEGER)
	`)
	if err != nil {
		return nil, fmt.Errorf("unable to CREATE TABLE quests in DB: %w", err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS history (
		id SERIAL PRIMARY KEY,
		quest_id INTEGER,
		quest_name VARCHAR,
		user_id INTEGER)
	`)
	if err != nil {
		return nil, fmt.Errorf("unable to CREATE TABLE history in DB: %w", err)
	}

	return &DataBase{DB: db}, nil
}

func (d *DataBase) AddUser(user UserQuestJSON) (string, error) {
	query := `INSERT INTO users (user_name, balance) VALUES ($1, $2) RETURNING user_id`
	row := d.DB.QueryRow(query, user.User_name, user.Balance)
	var imageID string
	if err := row.Scan(&imageID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("[USERS DB] Error while SAVING new USER info: %w", err)
		}
		return "", fmt.Errorf("[USERS DB] Error while SAVING new USER info: %w", err)
	}

	return imageID, nil
}

func (d *DataBase) AddQuest(quest UserQuestJSON) (string, error) {
	query := `INSERT INTO quests (quest_name, cost) VALUES ($1, $2) RETURNING quest_id`
	row := d.DB.QueryRow(query, quest.Quest_name, quest.Cost)
	var questID string
	err := row.Scan(&questID)
	if err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
		return "", ErrNotUniqueQuest
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("[QUESTS DB] Error while SAVING new QUEST info: %w", err)
		}
		return "", fmt.Errorf("[QUESTS DB] Error while SAVING new QUEST info: %w", err)
	}

	return questID, nil
}

func (d *DataBase) CompleteQuest(userquest UserQuestJSON) error {
	query := `SELECT cost, quest_name FROM quests WHERE quest_id = $1`
	rows, err := d.DB.Query(query, userquest.Quest_id)
	if err != nil {
		return fmt.Errorf("unable to return cost, quest_name from DB: %w", err)
	}
	var cost int
	var quest_name string
	for rows.Next() {
		err = rows.Scan(&cost, &quest_name)
		if err != nil {
			return err
		}
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	query2 := `SELECT quest_name FROM history WHERE user_id = $1`
	rows2, err := d.DB.Query(query2, userquest.User_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("Everything is OK")
		}
		return fmt.Errorf("unable to scan History DB: %w", err)
	}
	var quest_name2 string
	for rows2.Next() {
		err = rows2.Scan(&quest_name2)
		if err != nil {
			return err
		}
		if quest_name == quest_name2 {
			return ErrNotUniqueAction
		}
	}
	err = rows2.Err()
	if err != nil {
		return err
	}

	query3 := `UPDATE users SET balance = balance + $1 WHERE user_id = $2`
	_, err = d.DB.Exec(query3, cost, userquest.User_id)
	if err != nil {
		return fmt.Errorf("[USERS DB] Error while updating Balance: %w", err)
	}

	query4 := `INSERT INTO history (user_id, quest_name, quest_id) VALUES ($1, $2, $3)`
	_, err = d.DB.Exec(query4, userquest.User_id, quest_name, userquest.Quest_id)
	if err != nil {
		return fmt.Errorf("[HISTORY DB] Error while SAVING new HISTORY operation info: %w", err)
	}

	return nil
}

func (d *DataBase) GetAllInfo(user_id int) ([]UserQuestJSON, error) {
	query := `SELECT balance, user_name FROM users WHERE user_id = $1`
	rows, err := d.DB.Query(query, user_id)
	if err != nil {
		return nil, fmt.Errorf("unable to return balance from DB: %w", err)
	}
	var balance int
	var user_name string
	for rows.Next() {
		err = rows.Scan(&balance, &user_name)
		if err != nil {
			return nil, err
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	query2 := `SELECT quest_name FROM history WHERE user_id = $1`
	rows2, err := d.DB.Query(query2, user_id)
	if err != nil {
		return nil, fmt.Errorf("unable to return records from DB: %w", err)
	}

	records := []UserQuestJSON{}
	record := UserQuestJSON{User_id: user_id, User_name: user_name, Balance: balance}
	records = append(records, record)

	for rows2.Next() {
		record := UserQuestJSON{}
		err = rows2.Scan(&record.Quest_name)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	err = rows2.Err()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (d DataBase) Close() {
	d.DB.Close()
}
