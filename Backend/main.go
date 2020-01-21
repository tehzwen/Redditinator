package main

import (
	"fmt"
	"os/exec"

	"./reddit"
)

func main() {
	//run the python script and retrieve the output (json), can then cast it to struct
	c := exec.Command("py", "./blobAnalysis.py", "That is rough")

	out, err := c.Output()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println(string(out))

	//01/01/2017 https://www.unixtimestamp.com/index.php
	after := "1512086400"
	before := "1514764800"

	rClient := reddit.Reddit{
		TotalPosts: make(map[string][]reddit.SubredditPost),
	}
	rClient.GetAllSubredditData("", before, after, "ProgrammerHumor")

	fmt.Println("All done!")
	//print out the lengths of all the data retrieved by this reddit client
	for key := range rClient.TotalPosts {
		fmt.Printf("Total posts for %s: %d\n", key, len(rClient.TotalPosts[key]))
	}
}
