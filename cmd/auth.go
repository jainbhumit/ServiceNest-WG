package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
	"io"
	"os"
	"serviceNest/interfaces"
	"serviceNest/model"
	"serviceNest/util"
	"strings"
)

var GetUUID = util.GenerateUUID
var ValidEmail = util.ValidateEmail
var inputReader *bufio.Reader = bufio.NewReader(os.Stdin)

func SetInputReader(r io.Reader) {
	inputReader = bufio.NewReader(r)
}
func SignUp(userRepo interfaces.UserRepository) (*model.User, error) {
	_ = bufio.NewReader(os.Stdin)

	name, err := getInput("Enter Name: ")
	if err != nil {
		return nil, err
	}

	email, err := getValidEmail(userRepo)
	if err != nil {
		return nil, err
	}

	password, err := getPassword("Enter Password: ")
	if err != nil {
		return nil, err
	}
	var role string
	for {
		fmt.Print(`Enter Role 
Press 1-Householder
Press 2-ServiceProvider
`)
		var choice string
		fmt.Scanln(&choice)
		switch choice {
		case "1":
			role = "Householder"
		case "2":
			role = "ServiceProvider"
		default:
			color.Red("Invalid choice")
			continue
		}
		break

	}

	address, err := getInput("Enter Address: ")
	if err != nil {
		return nil, err
	}

	contact, err := getValidContact()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.User{
		ID:       GetUUID(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		Address:  address,
		Contact:  contact,
	}

	if err := userRepo.SaveUser(&user); err != nil {
		return nil, err
	}

	fmt.Println("User registered successfully!")
	return &user, nil
}

func getInput(prompt string) (string, error) {
	fmt.Print(prompt)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

func getValidEmail(userRepo interfaces.UserRepository) (string, error) {
	for {
		email, err := getInput("Enter Email: ")
		if err != nil {
			return "", err
		}

		if err := ValidEmail(email); err != nil {
			color.Red("%s", err)
			continue
		}

		existingUser, _ := userRepo.GetUserByEmail(email)
		if existingUser != nil {
			return "", fmt.Errorf("email already registered. Please use a different email address.")
		}
		return email, nil
	}
}

func getPassword(prompt string) (string, error) {
	for {
		fmt.Print(prompt)
		password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		passwordStr := strings.TrimSpace(string(password))
		if err := util.ValidatePassword(passwordStr); err != nil {
			color.Red("%s", err)
			continue
		}
		return passwordStr, nil
	}
}

func getRole() (string, error) {
	for {
		fmt.Print(`Enter Role 
Press 1-Householder
Press 2-ServiceProvider
`)
		var choice int
		_, err := fmt.Scanf("%d", &choice)
		if err != nil {
			return "", err
		}
		switch choice {
		case 1:
			return "Householder", nil
		case 2:
			return "ServiceProvider", nil
		default:
			color.Red("Invalid choice")
			continue
		}

	}
}

func getValidContact() (string, error) {
	for {
		contact, err := getInput("Enter Contact: ")
		if err != nil {
			return "", err
		}

		if err := util.ValidatePhoneNumber(contact); err != nil {
			color.Red("%s", err)
			continue
		}
		return contact, nil
	}
}
func Login(userRepo interfaces.UserRepository) (*model.User, error) {
	email, err := getInput("Enter Email: ")
	if err != nil {
		return nil, err
	}

	password, err := getPassword("Enter Password: ")
	if err != nil {
		return nil, err
	}

	user, err := userRepo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	fmt.Println("Login successful!")
	return user, nil
}
