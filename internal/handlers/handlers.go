package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Niiazgulov/vkquest/internal/storage"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func CreateUserHandler(repo storage.QuestDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tempStorage storage.UserQuestJSON
		if err := json.NewDecoder(r.Body).Decode(&tempStorage); err != nil {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		user := storage.UserQuestJSON{User_name: tempStorage.User_name, Balance: 0}
		userID, err := repo.AddUser(user)
		if err != nil {
			http.Error(w, "CreateUserHandler: can't add new user", http.StatusBadRequest)
			return
		}
		userIDint, err := strconv.Atoi(userID)
		if err != nil {
			http.Error(w, "CreateUserHandler: can't convert userID", http.StatusBadRequest)
			return
		}
		result := storage.UserQuestJSON{User_id: userIDint}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&result)
	}
}

func CreateQuestHandler(repo storage.QuestDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tempStorage storage.UserQuestJSON
		if err := json.NewDecoder(r.Body).Decode(&tempStorage); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		quest := storage.UserQuestJSON{Quest_name: tempStorage.Quest_name, Cost: tempStorage.Cost}
		questID, err := repo.AddQuest(quest)
		if err != nil {
			if errors.Is(err, storage.ErrNotUniqueQuest) {
				http.Error(w, "Can't add new quest - such quest already exists!", http.StatusBadRequest)
				return
			}
			http.Error(w, "CreateQuestHandler: can't add new quest", http.StatusBadRequest)
			return
		}
		questIDint, err := strconv.Atoi(questID)
		if err != nil {
			http.Error(w, "CreateQuestHandler: can't convert questID", http.StatusBadRequest)
			return
		}
		result := storage.UserQuestJSON{Quest_id: questIDint}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(&result)
	}
}

func NewActionHandler(repo storage.QuestDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var tempStorage storage.UserQuestJSON
		if err := json.NewDecoder(r.Body).Decode(&tempStorage); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		quest_name, err := repo.CompleteQuest(tempStorage)
		if err != nil {
			if errors.Is(err, storage.ErrNotUniqueAction) {
				http.Error(w, "Impossible - such user already made this operation!", http.StatusBadRequest)
				return
			}
			http.Error(w, "NewActionHandler: can't add new action", http.StatusBadRequest)
			return
		}
		result := storage.UserQuestJSON{Quest_name: quest_name}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&result)
	}
}

func GetUserInfoHandler(repo storage.QuestDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user_id := chi.URLParam(r, "id")
		user_id2, err := strconv.Atoi(user_id)
		if err != nil {
			http.Error(w, "GetUserInfoHandler: can't convert to int", http.StatusBadRequest)
			return
		}
		log.Println(user_id)
		allQuests, err := repo.GetAllInfo(user_id2)
		if err != nil {
			http.Error(w, "GetUserInfoHandler: can't get all info", http.StatusBadRequest)
			return
		}
		response, err := json.Marshal(allQuests)
		if err != nil {
			http.Error(w, "Status internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
