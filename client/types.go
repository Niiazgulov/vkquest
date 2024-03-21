package client

type ClientInterf interface {
	AddUser(string) (int, error)
	AddQuest(string, int) (int, error)
	CompleteQuest(int, int) (string, error)
	GetAllInfo(string) ([]UserQuestJSON, error)
}

type UserQuestJSON struct {
	User_name  string `json:"user_name,omitempty"`
	User_id    int    `json:"user_id,omitempty"`
	Balance    int    `json:"balance,omitempty"`
	Quest_name string `json:"quest_name,omitempty"`
	Quest_id   int    `json:"quest_id,omitempty"`
	Cost       int    `json:"cost,omitempty"`
}
