package main

import (
	"fmt"
	"go-nosql-db/command"
	"go-nosql-db/util"
)

func main() {

	result := command.BuildCommand().Run()
	prettyPrint := util.JsonPrettyPrint(result)
	if len(prettyPrint) == 0 {
		fmt.Println("Record Not found")
	} else {
		fmt.Println(prettyPrint)
	}

}
