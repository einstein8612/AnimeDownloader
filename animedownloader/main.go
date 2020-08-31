package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"joeyli.dev/animedownloader/downloader"
)

func main() {
	status := make(chan string)
	go func() {
		for {
			out := <-status
			fmt.Println(out)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Please type in a query for what anime you wish to download > ")
	scanner.Scan()
	animeName := scanner.Text()
	fmt.Print("Please type in the slug of the anime you're trying to watch (Can leave empty) > ")
	scanner.Scan()
	animeSlug := scanner.Text()
	fmt.Print("Please type in what service you'd like to watch on (4Anime: 1, Twist: 2, AnimeFreak: 3) > ")
	scanner.Scan()
	animeServiceInt, _ := strconv.Atoi(scanner.Text())
	var animeService downloader.AnimeSource
	if animeServiceInt == 1 {
		animeService = downloader.FourAnime
	} else if animeServiceInt == 2 {
		animeService = downloader.Twist
	} else if animeServiceInt == 3 {
		animeService = downloader.AnimeFreak
	}
	fmt.Print("Please type in which episodes you'd like to download or type ALL/leave empty to download all. (Example: 1-5) > ")
	scanner.Scan()
	episodes := strings.Split(scanner.Text(), "-")

	var startEpisode int
	var endEpisode int

	if episodes[0] == "" || episodes[0] == "ALL" {
		startEpisode = 0
		endEpisode = math.MaxInt32
	} else {
		if len(episodes) != 2 {
			fmt.Println("That format wasn't able to be read.")
			return
		}
		var err error
		startEpisode, err = strconv.Atoi(episodes[0])
		endEpisode, err = strconv.Atoi(episodes[1])
		if err != nil {
			fmt.Println("That format wasn't able to be read.")
			return
		}
	}
	downloader.DownloadEpisodes(animeName, animeSlug, startEpisode, endEpisode, animeService, status)

	time.Sleep(time.Hour)
}
