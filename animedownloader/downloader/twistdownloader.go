package downloader

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"joeyli.dev/animedownloader/decoders"
)

const TwistAccessToken string = "1rj2vRtegS8Y60B3w3qNZm5T2Q0TN2NR"

const TwistBaseUrl string = "https://twist.moe/"
const TwistGetAnimeUrl string = "https://twist.moe/api/anime"

func DownloadTwistEpisodes(animeName string, animeSlug string, startEpisode int, endEpisode int, status chan string) {
	if animeSlug == "" {
		animeSlug = TwistGetAnimeSlug(animeName)
	}
	if animeSlug == "" {
		status <- "ANIMENOTFOUND"
		return
	}
	episodesToDownload := TwistGetEpisodeDownloadUrls(animeSlug, startEpisode, endEpisode, status)
	headers := http.Header{}
	headers.Add("referer", TwistBaseUrl)
	DownloadUrls(episodesToDownload, headers, status)
}

func TwistGetAnimeSlug(animeName string) string {
	req, _ := http.NewRequest(http.MethodGet, TwistGetAnimeUrl, nil)
	req.Header.Add("x-access-token", TwistAccessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	var animes []struct {
		Title string `json:"title"`
		Slug  struct {
			Slug string `json:"slug"`
		} `json:"slug"`
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(bytes, &animes)

	for _, anime := range animes {
		titleAlphaNumerical := AlphaNumericalRegex.ReplaceAllString(anime.Title, "")
		if strings.Contains(strings.ToLower(titleAlphaNumerical), strings.ToLower(animeName)) {
			return anime.Slug.Slug
		}
	}
	return ""
}

func TwistGetEpisodeDownloadUrls(animeSlug string, startEpisode int, endEpisode int, status chan string) map[int]string {
	var episodesToDownload map[int]string = make(map[int]string)

	req, _ := http.NewRequest(http.MethodGet, "https://twist.moe/api/anime/"+animeSlug+"/sources", nil)
	req.Header.Add("x-access-token", TwistAccessToken)
	resp, _ := http.DefaultClient.Do(req)
	bytes, _ := ioutil.ReadAll(resp.Body)

	var episodeSources []struct {
		Number int    `json:"number"`
		Source string `json:"source"`
	}
	json.Unmarshal(bytes, &episodeSources)

	for _, episode := range episodeSources {
		if episode.Number >= startEpisode && episode.Number <= endEpisode {
			episodeDownloadLink := "https://edge-57.cdn.bunny.sh" + decoders.DecodeTwist(episode.Source)
			episodesToDownload[episode.Number] = episodeDownloadLink
			status <- "LINKGOTTEN " + strconv.Itoa(episode.Number)
		}
	}

	return episodesToDownload
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}
