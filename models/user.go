package models

type User struct {
	ID         uint   `json:"id" gorm:"primary_key"`
	Email      string `json:"email" gorm:"unique"`
	Password   string `json:"password"`
	Role       string `json:"role"`
	FounderID  uint   `json:"founderID"`
	InvestorID uint   `json:"investorID"`
}

type PublicUser struct {
	ID         uint   `json:"id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	FounderID  uint   `json:"founderID"`
	InvestorID uint   `json:"investorID"`
}

func (u User) GetPublicUser() PublicUser {
	return PublicUser{
		ID:         u.ID,
		Email:      u.Email,
		Role:       u.Role,
		FounderID:  u.FounderID,
		InvestorID: u.InvestorID,
	}
}
