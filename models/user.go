package models

type User struct {
	UUID         string  `json:"uuid" gorm:"primary_key"`
	ID           *uint   `json:"id"`
	Name         string  `json:"name"`
	Email        string  `json:"email" gorm:"unique"`
	Password     string  `json:"password"`
	Role         string  `json:"role"`
	FounderUUID  *string `json:"founder_uuid"`
	FounderID    *uint   `json:"founder_id"`
	InvestorUUID *string `json:"investor_uuid"`
	InvestorID   *uint   `json:"investor_id"`
}

type PublicUser struct {
	UUID         string  `json:"uuid"`
	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Role         string  `json:"role"`
	FounderUUID  *string `json:"founder_uuid"`
	InvestorUUID *string `json:"investor_uuid"`
}

func (u User) GetPublicUser() PublicUser {
	return PublicUser{
		UUID:         u.UUID,
		Email:        u.Email,
		Role:         u.Role,
		FounderUUID:  u.FounderUUID,
		InvestorUUID: u.InvestorUUID,
	}
}
