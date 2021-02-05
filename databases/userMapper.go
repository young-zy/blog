package databases

import (
	"blog/models"
)

// search for a user in database by username
func GetUser(username string) (user *models.User, err error) {
	user = &models.User{}
	err = db.Where(&models.User{Username: username}).First(user).Error
	return
}

// add a new user to database
func AddUser(user *models.User) error {
	return db.Create(user).Error
}

// updates the user, Id won't be updated
func UpdateUser(user *models.User) error {
	return db.Model(&models.User{}).Updates(user).Error
}

// delete a user by id
func DeleteUser(ID int) error {
	return db.Delete(&models.User{Id: ID}).Error
}
