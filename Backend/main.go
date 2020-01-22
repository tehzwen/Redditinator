package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync"

	"./reddit"
)

func main() {
	//detect what OS we are running
	const GOOS string = runtime.GOOS
	var pythonString string

	if GOOS == "linux" {
		pythonString = "python3"
	} else if GOOS == "windows" {
		pythonString = "py"
	}

	//run the python script and retrieve the output (json), can then cast it to struct
	c := exec.Command(pythonString, "./blobAnalysis.py", "That is rough")

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

	wg := sync.WaitGroup{}
	subreddits := []string{"alberta", "ProgrammerHumor", "witcher"}

	//loop through and get data for each of the subreddits
	for i := 0; i < len(subreddits); i++ {
		go rClient.GetAllSubredditData("", before, after, subreddits[i], &wg)
		wg.Add(1)
	}
	wg.Wait()

	fmt.Println("All done!")
	//print out the lengths of all the data retrieved by this reddit client
	for key := range rClient.TotalPosts {
		fmt.Printf("Total posts for %s: %d\n", key, len(rClient.TotalPosts[key]))
	}
}
