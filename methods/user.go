package methods

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
}

type UserDatabase map[string]User

func convertMapInSlice[K string, V any](baseMap map[K]V) []V {
	var result []V

	for _, value := range baseMap {
		result = append(result, value)
	}

	return result
}

func (db UserDatabase) FindAll() []User {
	return convertMapInSlice(db)
}

func (db UserDatabase) FindById(id string) (User, error) {
	user, ok := db[id]

	if !ok {
		return User{}, errors.New("user not found")
	}

	return user, nil
}

func (db UserDatabase) Insert(firstName, lastName, bio string) (User, error) {
	randomId, err := uuid.NewRandom()

	if err != nil {
		return User{}, errors.New("fail to generate random id")
	}

	newUser := User{
		Id:        randomId.String(),
		FirstName: firstName,
		LastName:  lastName,
		Biography: bio,
	}

	db[newUser.Id] = newUser

	return newUser, nil
}

func (db UserDatabase) Update(id, firstName, lastName, bio string) (User, error) {
	_, ok := db[id]

	if !ok {
		return User{}, errors.New("user not found")
	}

	db[id] = User{
		Id:        id,
		FirstName: firstName,
		LastName:  lastName,
		Biography: bio,
	}

	updatedUser := db[id]

	return updatedUser, nil
}

func (db UserDatabase) Delete(id string) (User, error) {
	user, ok := db[id]

	if !ok {
		return User{}, errors.New("user not found")
	}

	delete(db, id)

	return user, nil
}
