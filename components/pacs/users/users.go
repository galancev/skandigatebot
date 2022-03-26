package users

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"skandigatebot/base"
	pc "skandigatebot/components/pacs/config"
	u "skandigatebot/models/user"
	"skandigatebot/models/user/role"
	"strconv"
)

func UpdateUsers() {
	log.Print("Run update users")

	pacsUsers, err := getPACSUsers()

	if err != nil {
		fmt.Println(err)

		return
	}

	users, err := getUsers()

	if err != nil {
		fmt.Println(err)

		return
	}

	syncUsers(pacsUsers, users)

	log.Print("Finish update users")
}

func getPACSUsers() (PACSUserResponse, error) {
	conf := pc.New()

	client := &http.Client{}
	URL := conf.Host + "/data.cgx?cmd={\"Command\":\"GetAllGateUsers\"}"

	req, err := http.NewRequest("GET", URL, nil)
	req.SetBasicAuth(conf.User, conf.Password)
	resp, err := client.Do(req)
	if err != nil && resp.StatusCode != http.StatusOK {
		fmt.Println(err)

		return PACSUserResponse{}, err
	}

	pacUsersResponse := &PACSUserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&pacUsersResponse)
	if err != nil {
		fmt.Println(err)

		return PACSUserResponse{}, err
	}

	return *pacUsersResponse, err
}

func getUsers() ([]u.User, error) {
	users, err := u.GetUsers()

	return users, err
}

func syncUsers(pacsUsers PACSUserResponse, users []u.User) {
	var isFound bool

	usersToDelete := make(map[uint]u.User)
	usersToUpdate := make(map[uint]u.User)
	usersToInsert := make(map[uint]u.User)

	for _, user := range users {
		usersToDelete[user.Phone] = user
	}

	for _, pacsUser := range pacsUsers.Users {
		isFound = false

		for _, user := range users {
			// 0 - phone
			// 1 - name
			phoneInt, err := strconv.Atoi(pacsUser[0])

			if err != nil {
				continue
			}

			phone := uint(phoneInt)

			if phone == user.Phone {
				delete(usersToDelete, phone)

				if pacsUser[1] != user.FirstName {
					user.FirstName = pacsUser[1]
					usersToUpdate[phone] = user
				}

				isFound = true
			}
		}

		if !isFound {
			phoneInt, err := strconv.Atoi(pacsUser[0])

			if err != nil {
				continue
			}

			phone := uint(phoneInt)

			newUser := u.User{
				Phone:     phone,
				FirstName: pacsUser[1],
				RoleId:    role.User,
			}

			usersToInsert[phone] = newUser
		}
	}

	err := base.GetDB().Transaction(func(tx *gorm.DB) error {
		for _, user := range usersToDelete {
			base.GetDB().Delete(&user)
		}
		for _, user := range usersToUpdate {
			base.GetDB().Updates(&user)
		}
		for _, user := range usersToInsert {
			base.GetDB().Create(&user)
		}

		return nil
	})

	if err != nil {
		return
	}

}
