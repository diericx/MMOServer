package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func getConsoleInput() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		//do soemthing with it
		var keywords []string
		keywords = findAllKeywords(input)
		//if the user types give ("give player scraps 100")
		if keywords[0] == "give" {
			//make sure they input the correct amount of parameters
			if len(keywords) == 4 {
				playerID := findPlayerIndexByID(keywords[1])
				if playerID >= 0 {
					if keywords[2] == "scraps" {
						//try to convert value to string
						i, err := strconv.ParseInt(keywords[3], 10, 32)
						if err != nil {
							fmt.Println("\"" + keywords[3] + "\" is not a valid integer!")
						} else {
							players[playerID].scraps += int32(i)
							fmt.Println("Succesfully gave player \"" + keywords[1] + "\" " + keywords[3] + " scraps!")
						}
					} else {
						fmt.Println("\"" + keywords[2] + "\" is not a known command!")
					}
				} else {
					fmt.Println("Player with id \"" + keywords[1] + "\" not found!")
				}
			} else {
				fmt.Println("Not enough parameters supplied for \"give\" command!")
			}
		} else {
			fmt.Println("\"" + keywords[0] + "\" is not a known command!")
		}
	}
}

func findAllKeywords(val string) []string {
	input := val
	foundParams := make([]string, 0)

	for {
		i := strings.Index(input, " ")
		if i == -1 {
			foundParams = append(foundParams, input)
			break
		} else {
			beforeSpace := input[:i]
			afterSpace := input[i+1:]
			foundParams = append(foundParams, beforeSpace)
			input = afterSpace
		}

	}

	return foundParams
}
