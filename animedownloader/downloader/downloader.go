package downloader

import (
	"errors"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type AnimeSource int

// Priority. Lower = higher priority
var FourAnime AnimeSource = 1
var Twist AnimeSource = 2
var AnimeFreak AnimeSource = 3

const USER_AGENT = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"

var AlphaNumericalRegex *regexp.Regexp

func init() {
	AlphaNumericalRegex, _ = regexp.Compile("[^a-zA-Z0-9 ]+")
}

func DownloadEpisodes(animeName string, animeSku string, startEpisode int, endEpisode int, animeSource AnimeSource, status chan string) {
	switch animeSource {
	case FourAnime:
		DownloadFourAnimeEpisodes(animeName, animeSku, startEpisode, endEpisode, status)
		break
	case Twist:
		DownloadTwistEpisodes(animeName, animeSku, startEpisode, endEpisode, status)
		break
	case AnimeFreak:
		DownloadAnimeFreakEpisodes(animeName, animeSku, startEpisode, endEpisode, status)
	}
}

func DownloadFile(filepath string, url string, headers http.Header) error {
	// Get the data
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header = headers
	req.Header.Add("user-agent", USER_AGENT)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode == 404 {
		return errors.New("The download could not be found")
	}
	return err
}

func DownloadUrls(episodesToDownload map[int]string, headers http.Header, status chan string) {
	var wgDownloadLinks sync.WaitGroup
	for episodeNum, episodeLink := range episodesToDownload {
		wgDownloadLinks.Add(1)
		go func(episodeLink string, episodeNum int) {
			linkParts := strings.Split(episodeLink, "/")
			fileName := strings.Split(linkParts[len(linkParts)-1], "?")[0]
			err := DownloadFile("./anime/"+fileName, episodeLink, headers)
			if err != nil {
				status <- "LINKDOWNLOADFAILED " + strconv.Itoa(episodeNum)
			} else {
				status <- "LINKDOWNLOADED " + strconv.Itoa(episodeNum)
			}
			wgDownloadLinks.Done()
		}(episodeLink, episodeNum)
	}
	wgDownloadLinks.Wait()
}

func GetFirstEntryInMap(a map[int]string) string {
	for _, b := range a {
		return b
	}
	return ""
}
