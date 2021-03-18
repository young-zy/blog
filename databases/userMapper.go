package databases

import (
	"context"

	"blog/models"
)

//// search for a user in database by username
//func GetUser(c context.Context, username string) (user *models.User, err error) {
//	user = &models.User{}
//	err = DefaultDb.WithContext(c).Where(&models.User{Username: username}).First(user).Error
//	return
//}

func (tx *Transaction) GetUser(c context.Context, username string) (user *models.User, err error) {
	user = &models.User{}
	err = tx.tx.WithContext(c).Where(&models.User{Username: username}).First(user).Error
	return
}

// add a new user to database
func (tx *Transaction) AddUser(c context.Context, user *models.User) error {
	return tx.tx.WithContext(c).Create(user).Error
}

//// updates the user, Id won't be updated
//func UpdateUser(user *models.User) error {
//	return DefaultDb.Model(&models.User{}).Where("id = ?", user.Id).Updates(user).Error
//}

// updates the user, Id won't be updated
func (tx *Transaction) UpdateUser(c context.Context, user *models.User) error {
	return tx.tx.WithContext(c).Model(&models.User{}).Where("id = ?", user.Id).Updates(user).Error
}

// delete a user by id
func (tx *Transaction) DeleteUser(c context.Context, ID int) error {
	return tx.tx.WithContext(c).Delete(&models.User{Id: ID}).Error
}
