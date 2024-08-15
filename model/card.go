package model

import (
	"github.com/marcsello/ponyhug2-backend/db"
)

type CardCopyVisibleByPlayer struct {
	CopyID     int32   `json:"copy_id"`
	CopyKey    string  `json:"copy_key"`
	Place      int16   `json:"place"`
	Name       string  `json:"name"`
	Source     *string `json:"source"`
	WearLevel  int16   `json:"wear_level"`
	ImgUrl     string  `json:"img_url"`
	CopiedFrom *string `json:"copied_from"`
}

func CardCopyVisibleByPlayerFromDBPlayerCardsRow(row db.GetPlayerCardsRow) CardCopyVisibleByPlayer {
	return CardCopyVisibleByPlayer{
		CopyID:     row.ID,
		CopyKey:    row.Key,
		Place:      row.Place,
		Name:       row.Name,
		Source:     row.Source,
		WearLevel:  row.WearLevel,
		ImgUrl:     row.ImageUrl,
		CopiedFrom: row.Name_2,
	}
}

func CardCopyVisibleByPlayerFromDBGetCardCopyRow(row db.GetCardCopyRow) CardCopyVisibleByPlayer {
	return CardCopyVisibleByPlayer{
		CopyID:     row.ID,
		CopyKey:    row.Key,
		Place:      row.Place,
		Name:       row.Name,
		Source:     row.Source,
		WearLevel:  row.WearLevel,
		ImgUrl:     row.ImageUrl,
		CopiedFrom: row.Name_2,
	}
}

func CardCopiesVisibleByPlayerFromDBPlayerCardsRows(rows []db.GetPlayerCardsRow) []CardCopyVisibleByPlayer {
	result := make([]CardCopyVisibleByPlayer, len(rows))
	for i, row := range rows {
		result[i] = CardCopyVisibleByPlayerFromDBPlayerCardsRow(row)
	}
	return result
}

type CardBaseForAdmins struct {
	ID        int16   `json:"id"`
	Key       *string `json:"key"`
	Name      string  `json:"name"`
	Source    *string `json:"source"`
	Place     int16   `json:"place"`
	BaseID    *int16  `json:"base_id"`
	WearLevel *int16  `json:"wear_level"`
	ImageUrl  *string `json:"image_url"`
}

func CardBaseForAdminsFromDB(row db.GetCardBasesRow) CardBaseForAdmins {
	return CardBaseForAdmins{
		ID:        row.ID,
		Key:       row.Key,
		Name:      row.Name,
		Source:    row.Source,
		Place:     row.Place,
		BaseID:    row.BaseID,
		WearLevel: row.WearLevel,
		ImageUrl:  row.ImageUrl,
	}
}

func CardBasesForAdminsFromDB(rows []db.GetCardBasesRow) []CardBaseForAdmins {
	result := make([]CardBaseForAdmins, len(rows))
	for i, row := range rows {
		result[i] = CardBaseForAdminsFromDB(row)
	}
	return result
}

type BareCardBase struct {
	ID     int16   `json:"id"`
	Key    *string `json:"key"`
	Name   string  `json:"name"`
	Source *string `json:"source"`
	Place  int16   `json:"place"`
}

func BareCardBaseFromDBCardBase(c db.CardBase) BareCardBase {
	return BareCardBase{
		ID:     c.ID,
		Key:    c.Key,
		Name:   c.Name,
		Source: c.Source,
		Place:  c.Place,
	}
}

type CardWearImg struct {
	BaseID    int16  `json:"base_id"`
	WearLevel int16  `json:"wear_level"`
	ImageUrl  string `json:"image_url"`
}

func CardWearImgFromDB(c db.CardWearImg) CardWearImg {
	return CardWearImg{
		BaseID:    c.BaseID,
		WearLevel: c.WearLevel,
		ImageUrl:  c.ImageUrl,
	}
}
