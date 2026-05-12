package models

import "time"

type InviteCode struct {
	ID          int       `json:"id"`
	Code        string    `gorm:"uniqueIndex" json:"code"`
	Description string    `json:"description"`
	Enabled     bool      `json:"enabled"`
	UsedCount   int       `json:"usedCount"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (invite *InviteCode) Add() error {
	return DB.Create(invite).Error
}

func (invite *InviteCode) FindByCode() error {
	return DB.Where("code = ?", invite.Code).First(invite).Error
}

func (invite *InviteCode) FindByID() error {
	return DB.Where("id = ?", invite.ID).First(invite).Error
}

func (invite *InviteCode) List() ([]InviteCode, error) {
	var invites []InviteCode
	err := DB.Order("id desc").Find(&invites).Error
	return invites, err
}

func (invite *InviteCode) Update() error {
	return DB.Model(&InviteCode{}).Where("id = ?", invite.ID).Updates(invite).Error
}

func (invite *InviteCode) Delete() error {
	return DB.Delete(invite).Error
}
