package storage

type UserQuestJSON struct {
	User_name  string `json:"user_name,omitempty"`
	User_id    int    `json:"user_id,omitempty"`
	Balance    int    `json:"balance,omitempty"`
	Quest_name string `json:"quest_name,omitempty"`
	Quest_id   int    `json:"quest_id,omitempty"`
	Cost       int    `json:"cost,omitempty"`
}

type QuestDB interface {
	AddUser(user UserQuestJSON) (string, error)
	AddQuest(quest UserQuestJSON) (string, error)
	CompleteQuest(userquest UserQuestJSON) error
	GetAllInfo(user_id int) ([]UserQuestJSON, error)
	Close()
}
