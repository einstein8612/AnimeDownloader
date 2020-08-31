package downloader

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

const FourAnimeBaseUrl string = "https://4anime.to/"
const FourAnimeSearchUrl string = "https://4anime.to/wp-admin/admin-ajax.php"

const FourAnimeCDNDownloadAtOnce int = 2

func DownloadFourAnimeEpisodes(animeName string, animeSku string, startEpisode int, endEpisode int, status chan string) {
	animeBaseUrl := FourAnimeGetAnimeBaseUrl(animeName, animeSku)
	if animeBaseUrl == "" {
		status <- "Anime not found"
		return
	}
	episodeUrls := FourAnimeGetEpisodeUrls(animeBaseUrl)
	episodesToDownload := FourAnimeGetEpisodeDownloadUrls(startEpisode, endEpisode, episodeUrls, status)
	isGoogleStorage := strings.HasPrefix(GetFirstEntryInMap(episodesToDownload), "https://storage.googleapis.com/")
	if isGoogleStorage {
		DownloadUrls(episodesToDownload, http.Header{}, status)
		return
	}
}

func FourAnimeGetAnimeBaseUrl(animeName string, animeSku string) string {
	var animeBaseUrl string
	if animeSku != "" {
		animeBaseUrl = FourAnimeBaseUrl + "anime/" + animeSku
	}

	payload := url.Values{}
	payload.Add("action", "ajaxsearchlite_search")
	payload.Add("aslp", animeName)
	payload.Add("asid", "1")
	payload.Add("options", "qtranslate_lang=0&set_intitle=None&customset%5B%5D=anime")

	resp, err := http.Post(FourAnimeSearchUrl, "application/x-www-form-urlencoded", strings.NewReader(payload.Encode()))
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	animeNames := doc.Find(".name")
	if len(animeNames.Nodes) == 0 {
		return animeBaseUrl
	}
	animeBaseUrl = doc.Find(".name").First().AttrOr("href", "")
	return animeBaseUrl
}

func FourAnimeGetEpisodeUrls(animeBaseUrl string) map[int]string {
	var episodes map[int]string = make(map[int]string)

	resp, _ := http.Get(animeBaseUrl)
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	doc.Find("ul.episodes").Children().Each(func(index int, selection *goquery.Selection) {
		linkObject := selection.Children().First()
		episodeNum, _ := strconv.Atoi(linkObject.Text())
		episodeLink := linkObject.AttrOr("href", "")
		episodes[episodeNum] = episodeLink
	})

	return episodes
}

func FourAnimeGetEpisodeDownloadUrls(startEpisode int, endEpisode int, episodeUrls map[int]string, status chan string) map[int]string {
	var episodesToDownload map[int]string = make(map[int]string)

	var wgGetUrls sync.WaitGroup
	for episodeNum, episodeLink := range episodeUrls {
		if episodeNum >= startEpisode && episodeNum <= endEpisode {
			wgGetUrls.Add(1)
			go func(episodeLink string, episodeNum int) {
				req, _ := http.NewRequest(http.MethodGet, episodeLink, nil)
				req.Header.Add("user-agent", USER_AGENT)
				resp, _ := http.DefaultClient.Do(req)
				doc, _ := goquery.NewDocumentFromReader(resp.Body)
				episodeDownloadLink := doc.Find("source").First().AttrOr("src", "")
				episodesToDownload[episodeNum] = episodeDownloadLink
				status <- "LINKGOTTEN " + strconv.Itoa(episodeNum)
				wgGetUrls.Done()
			}(episodeLink, episodeNum)
		}
	}
	wgGetUrls.Wait()

	return episodesToDownload
}
