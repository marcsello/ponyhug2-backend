package model

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/marcsello/ponyhug2-backend/db"
)

type PlayerSelf struct {
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin"`
}

func PlayerSelfFromDB(p db.Player) PlayerSelf {
	return PlayerSelf{
		Name:    p.Name,
		IsAdmin: p.IsAdmin,
	}
}

type PlayerData struct {
	ID         int32            `json:"id"`
	Name       string           `json:"name"`
	Registered pgtype.Timestamp `json:"registered"` // this can be marshaled it seems
	IsAdmin    bool             `json:"is_admin"`
}

func PlayerDataFromDB(p db.Player) PlayerData {
	return PlayerData{
		ID:         p.ID,
		Name:       p.Name,
		Registered: p.Registered,
		IsAdmin:    p.IsAdmin,
	}
}

func PlayersDataFromDB(players []db.Player) []PlayerData {
	pd := make([]PlayerData, len(players))
	for i, p := range players {
		pd[i] = PlayerDataFromDB(p)
	}
	return pd
}
