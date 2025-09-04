package models

type User struct {
	UUID string `json:"uuid" gorm:"type:uuid;primary_key"`
	ID   *uint  `json:"id" gorm:"unique;index"` // legacy

	Name     string  `json:"name"`
	Email    string  `json:"email" gorm:"unique;not null;default:null"`
	Password *string `json:"password"`
	Role     string  `json:"role"`

	FounderUUID *string  `json:"founder_uuid"`
	FounderID   *uint    `json:"founder_id"` // legacy
	Founder     *Founder `json:"founder" gorm:"foreignKey:FounderUUID;references:UUID"`

	InvestorUUID *string   `json:"investor_uuid"`
	InvestorID   *uint     `json:"investor_id"` // legacy
	Investor     *Investor `json:"investor" gorm:"foreignKey:InvestorUUID;references:UUID"`

	UserPicture *UserPicture `json:"user_picture" gorm:"foreignKey:UserUUID;references:UUID"`
}

type UserPicture struct {
	UserUUID string `json:"user_uuid" gorm:"type:uuid;primary_key"`
	Picture  []byte `json:"picture" gorm:"type:bytea"`
}

type PublicUser struct {
	UUID string `json:"uuid"`

	Name         string  `json:"name"`
	Email        string  `json:"email"`
	Role         string  `json:"role"`
	FounderUUID  *string `json:"founder_uuid"`
	InvestorUUID *string `json:"investor_uuid"`
}

func (u User) GetPublicUser() PublicUser {
	return PublicUser{
		UUID:         u.UUID,
		Name:         u.Name,
		Email:        u.Email,
		Role:         u.Role,
		FounderUUID:  u.FounderUUID,
		InvestorUUID: u.InvestorUUID,
	}
}

type UserCreationRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Role  string `json:"role" binding:"required"`
}
