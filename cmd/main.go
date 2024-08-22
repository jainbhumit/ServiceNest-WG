package main

import (
	"fmt"
	"github.com/fatih/color"
)

func main() {
	for {
		fmt.Println("-----------------------Welcome-----------------------")
		color.Blue("For SignUp press 1\n")
		color.Blue("For Login press 2\n")
		color.Blue("For Exit press 3\n")
		var choice int
		color.Cyan("Enter your choice: ")
		fmt.Scanln(&choice)
		switch choice {
		case 1:

			SignUpUser()
		case 2:

			LoginUser()
		case 3:
			return
		default:
			color.Red("Invalid choice")

		}
	}

}
