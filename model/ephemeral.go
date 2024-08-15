package model

// ephemeral types generally does not directly represent an entry in the db
// they are helper structs used for some specific api calls

type PlayerRegistrationSuccess struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}
