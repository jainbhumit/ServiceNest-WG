// package service
//
// import (
//
//	"bufio"
//	"fmt"
//	"github.com/fatih/color"
//	"github.com/google/uuid"
//	"golang.org/x/crypto/ssh/terminal"
//	"os"
//	"serviceNest/model"
//	"serviceNest/repository"
//	"serviceNest/util"
//	"strings"
//
// )
//
//	func SignUp(userRepo *repository.UserRepository) (*model.User, error) {
//		reader := bufio.NewReader(os.Stdin)
//
//		fmt.Print("Enter Name: ")
//		name, _ := reader.ReadString('\n')
//		name = strings.TrimSpace(name)
//
//		var email string
//
//		for {
//			fmt.Print("Enter Email: ")
//			email, _ = reader.ReadString('\n')
//			email = strings.TrimSpace(email)
//
//			// Validate email
//			if err := util.ValidateEmail(email); err != nil {
//				color.Red("%s", err)
//				continue
//			}
//			// Check if email already exists
//			existingUser, _ := userRepo.GetUserByEmail(email)
//			if existingUser != nil {
//				return nil, fmt.Errorf("email already registered. Please use a different email address.")
//			}
//			break
//		}
//		var hiddenPassword string
//
//		for {
//			fmt.Print("Enter Password: ")
//			password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
//			hiddenPassword := string(password)
//			hiddenPassword = strings.TrimSpace(hiddenPassword)
//			//hiddenPassword, _ = reader.ReadString('\n')
//			//hiddenPassword = strings.TrimSpace(hiddenPassword)
//			if err := util.ValidatePassword(hiddenPassword); err != nil {
//				color.Red("%s", err)
//				continue
//			}
//			break
//		}
//		////newPassword, err := bcrypt.GenerateFromPassword([]byte(hiddenPassword), 10)
//		//if err != nil {
//		//	color.Red("%s", err)
//		//}
//		var role string
//		for {
//			fmt.Print(`Enter Role
//
// Press 1-Householder
// Press 2-ServiceProvider
// `)
//
//			fmt.Println()
//			var choice int
//			fmt.Scanf("%d", &choice)
//			switch choice {
//			case 1:
//				role = "Householder"
//			case 2:
//				role = "ServiceProvider"
//			default:
//				color.Red("Invalid choice")
//				continue
//			}
//			break
//		}
//
//		fmt.Print("Enter Address: ")
//		address, _ := reader.ReadString('\n')
//		address = strings.TrimSpace(address)
//
//		var contact string
//		for {
//			fmt.Print("Enter Contact: ")
//			contact, _ = reader.ReadString('\n')
//			contact = strings.TrimSpace(contact)
//			// Validate phone number
//			if err := util.ValidatePhoneNumber(contact); err != nil {
//				color.Red("%s", err)
//				continue
//			}
//			break
//		}
//
//		//fmt.Print("Enter Latitude: ")
//		//var lat float64
//		//fmt.Scanf("%f", &lat)
//		//
//		////fmt.Print("Enter Longitude: ")
//		//var lon float64
//		//fmt.Scanf("%f", &lon)
//
//		// Generate UUID for the user ID
//		user := model.User{
//			ID:       uuid.New().String(),
//			Name:     name,
//			Email:    email,
//			Password: hiddenPassword,
//			Role:     role,
//			Address:  address,
//			Contact:  contact,
//			//Latitude:  lat,
//			//Longitude: lon,
//		}
//
//		if err := userRepo.SaveUser(user); err != nil {
//			return nil, err
//		}
//
//		fmt.Println("User registered successfully!")
//		return &user, nil
//	}
//
//	func Login(userRepo *repository.UserRepository) (*model.User, error) {
//		reader := bufio.NewReader(os.Stdin)
//
//		fmt.Print("Enter Email: ")
//		email, _ := reader.ReadString('\n')
//		email = strings.TrimSpace(email)
//
//		fmt.Print("Enter Password: ")
//		password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
//		hiddenPassword := string(password)
//		hiddenPassword = strings.TrimSpace(hiddenPassword)
//		//password, _ := reader.ReadString('\n')
//		//password = strings.TrimSpace(password)
//		user, err := userRepo.GetUserByEmail(email)
//
//		if err != nil {
//			return nil, err
//		}
//		if hiddenPassword != user.Password {
//			return nil, fmt.Errorf("invalid credentials")
//		}
//		//ok := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
//		//if ok != nil {
//		//	return nil, fmt.Errorf("invalid credentials")
//		//}
//
//		fmt.Println("Login successful!")
//		return user, nil
//	}
package service

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"serviceNest/interfaces"
	"serviceNest/model"
	"serviceNest/util"
	"strings"
)

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

	role, err := getRole()
	if err != nil {
		return nil, err
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
		ID:       uuid.New().String(),
		Name:     name,
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
		Address:  address,
		Contact:  contact,
	}

	if err := userRepo.SaveUser(user); err != nil {
		return nil, err
	}

	fmt.Println("User registered successfully!")
	return &user, nil
}

func getInput(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
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

		if err := util.ValidateEmail(email); err != nil {
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
