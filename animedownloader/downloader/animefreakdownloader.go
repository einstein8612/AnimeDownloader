package downloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const AnimeFreakBaseUrl string = "https://www.animefreak.tv/"
const AnimeFreakSearchUrl string = "https://www.animefreak.tv/search/topSearch?q="

func DownloadAnimeFreakEpisodes(animeName string, animeSku string, startEpisode int, endEpisode int, status chan string) {
	animeBaseUrl := AnimeFreakGetAnimeBaseUrl(animeName, animeSku)
	if animeBaseUrl == "" {
		status <- "Anime not found"
		return
	}
	episodeUrls := AnimeFreakGetEpisodeUrls(animeBaseUrl)
	episodesToDownload := AnimeFreakGetEpisodeDownloadUrls(startEpisode, endEpisode, episodeUrls, status)
	DownloadUrls(episodesToDownload, http.Header{}, status)
	fmt.Println("Done")
}

func AnimeFreakGetAnimeBaseUrl(animeName string, animeSku string) string {
	var animeBaseUrl string
	if animeSku != "" {
		animeBaseUrl = AnimeFreakBaseUrl + "watch/" + animeSku
	}

	resp, err := http.Get(AnimeFreakSearchUrl + url.PathEscape(animeName))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	var episodes struct {
		Data []struct {
			SeoName string `json:"seo_name"`
		} `json:"data"`
	}
	json.Unmarshal(bytes, &episodes)

	if len(episodes.Data) == 0 {
		return animeBaseUrl
	}
	animeBaseUrl = AnimeFreakBaseUrl + "watch/" + episodes.Data[0].SeoName
	return animeBaseUrl
}

func AnimeFreakGetEpisodeUrls(animeBaseUrl string) map[int]string {
	var episodes map[int]string = make(map[int]string)

	resp, _ := http.Get(animeBaseUrl)
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	doc.Find("ul.check-list").Eq(1).Children().Each(func(index int, selection *goquery.Selection) {
		episodeLink := selection.Children().Eq(1).AttrOr("href", "")
		linkParts := strings.Split(episodeLink, "-")
		episodeNumString := linkParts[len(linkParts)-1]
		episodeNum, _ := strconv.Atoi(episodeNumString)
		episodes[episodeNum] = episodeLink
	})

	return episodes
}

func AnimeFreakGetEpisodeDownloadUrls(startEpisode int, endEpisode int, episodeUrls map[int]string, status chan string) map[int]string {
	var episodesToDownload map[int]string = make(map[int]string)

	var wgGetUrls sync.WaitGroup
	for episodeNum, episodeLink := range episodeUrls {
		if episodeNum >= startEpisode && episodeNum <= endEpisode {
			wgGetUrls.Add(1)
			go func(episodeLink string, episodeNum int) {
				req, _ := http.NewRequest(http.MethodGet, episodeLink, nil)
				req.Header.Add("user-agent", USER_AGENT)
				resp, _ := http.DefaultClient.Do(req)
				bytes, _ := ioutil.ReadAll(resp.Body)
				body := string(bytes)
				x := body[strings.Index(body, `var file = "`)+12:]
				episodeDownloadLink := x[:strings.Index(x, `";`)]
				episodesToDownload[episodeNum] = episodeDownloadLink
				status <- "LINKGOTTEN " + strconv.Itoa(episodeNum)
				wgGetUrls.Done()
			}(episodeLink, episodeNum)
		}
	}
	wgGetUrls.Wait()

	return episodesToDownload
}
