package input

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashload/boss/internal/pkg/utils"
)

func Select(msg string, options []string, required bool, def int) (opt int, err error) {

	fmt.Printf("%s\n", msg)
	for index, option := range options {
		fmt.Printf("%d - %s\n", index+1, option)
	}
	for {
		fmt.Printf("Enter a number (Default is %d): ", def+1)
		reader := bufio.NewReader(os.Stdin)
		in, err := reader.ReadString('\n')
		if err != nil {
			return -1, err
		}
		replacer := strings.NewReplacer("\n", "", "\r", "")

		cleanStr := replacer.Replace(in)

		if cleanStr == "" {
			if required && def == -1 {
				fmt.Printf("Input must not be empty. Answer by a number.\n\n")
				continue
			} else {
				return def, nil
			}
		}

		opt, err = strconv.Atoi(replacer.Replace(in))
		if err != nil {
			fmt.Printf("%q is not a valid input. Answer by a number.\n\n", cleanStr)
			continue
		}

		// Check answer is in range of list
		if opt < 1 || len(options) < opt {
			fmt.Printf("%q is not a valid choice. Choose a number from 1 to %d.\n\n",
				cleanStr, len(options))
			continue
		}
		opt--
		break
	}

	return //lint :/
}

func SelectBool(msg string, required bool, def bool) (opt bool, err error) {
	defInt := 0
	if !def {
		defInt = 1
	}

	selected, err := Select(
		msg,
		[]string{"Yes", "No"},
		required,
		defInt,
	)

	if selected == 0 {
		return true, err
	} else {
		return false, err
	}
}

func GetTextOrDef(msg string, def string) string {

	fmt.Printf("%s (Default is '%s'): ", msg, def)
	reader := bufio.NewReader(os.Stdin)
	in, err := reader.ReadString('\n')
	if err != nil {
		utils.CheckError(err)
	}
	replacer := strings.NewReplacer("\n", "", "\r", "")

	cleanStr := replacer.Replace(in)

	if cleanStr == "" {
		return def
	}

	return ""
}
