package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	utils "ascii/utlis"
)

type Args struct {
	banner string
	output string
	input  string
}

var (
	errOutputFileNotDefined = errors.New("output file is not defined")
	errOuputAlreadyDefiend  = errors.New("output argument is already defined")
	errUsage                = errors.New("usage : go run . [OPTION] [STRING] [BANNER]")
	errInputNotDefiend      = errors.New("input string not defiend")
)

func main() {
	args := os.Args[1:]

	if !utils.IsValidArgs(args) {
		fmt.Println("Usage : go run . [OPTION] [STRING] [BANNER]")
		return
	}
	Args, err := HandleArgs(args)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("input:", Args.input, "; output:", Args.output, ";Banner:", Args.banner)
	input := Args.input
	// default banner if banner doesn't exits
	data, err := utils.HandleBanner(Args.banner)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := utils.IsValidInput(input, data); err != nil {
		fmt.Println(err)
		return
	}

	output := PrintAscii(input, data)

	if Args.output == "" {
		fmt.Print(output)
	} else {
		file, err := os.Create(Args.output)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		// Write to the file
		_, err = file.WriteString(output)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}

func PrintAscii(input string, data []string) string {
	output := ""
	lines := strings.Count(input, "\\n")
	words := strings.Split(input, "\\n")
	for j := 0; j < len(words); j++ {
		// Handle new line
		if words[j] == "" {
			if lines > 0 {
				output += "\n"
				lines--
			}
			continue
		}

		for i := 0; i < 8; i++ {
			for _, char := range words[j] {
				pos := utils.StartIndex(char)
				output += data[pos+i]
			}
			output += "\n"
		}
	}
	return output
}

func HandleArgs(args []string) (Args, error) {
	res := Args{}
	isInputDefined := false
	res.banner = "standard"
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--output=") {
			if res.output != "" {
				return res, errOuputAlreadyDefiend
			}
			if res.input != "" {
				return res, errUsage
			}
			if len(arg) <= 9 {
				return res, errOutputFileNotDefined
			}
			filename := arg[9:]
			fmt.Println(filename)
			res.output = filename
		} else if res.input == "" && !isInputDefined {
			res.input = arg
			isInputDefined = true
		} else {
			res.banner = arg
		}
	}
	if res.input == "" {
		return res, errInputNotDefiend
	}
	return res, nil
}


