package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Vaansh/gore/internal/model"
	"github.com/Vaansh/gore/internal/platform"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/url"
	"os"
)

type ChannelShortsListResponse struct {
	Kind  string `json:"kind"`
	Etag  string `json:"etag"`
	Items []struct {
		Kind   string `json:"kind"`
		Etag   string `json:"etag"`
		ID     string `json:"id"`
		Shorts []struct {
			VideoID    string `json:"videoId"`
			Title      string `json:"title"`
			Thumbnails []struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnails"`
			ViewCount int `json:"viewCount"`
		} `json:"shorts"`
		NextPageToken string `json:"nextPageToken"`
	} `json:"items"`
}

type YoutubeClient struct {
	apiKey string
}

func NewYoutubeClient() *YoutubeClient {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading environment variables file")
	}

	return &YoutubeClient{
		apiKey: os.Getenv("YOUTUBE_API_KEY"),
	}
}

// Videos
func (c *YoutubeClient) FetchVideo(videoURL string) (map[string]interface{}, error) {
	resp, err := http.Get(videoURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (c *YoutubeClient) FetchLatestVideoByChannel(channelID string) (string, error) {
	paginatedVideos := c.PaginatedVideosAPI(channelID)
	resp, err := http.Get(paginatedVideos.String())
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	item := data["items"].([]interface{})[0]
	itemMap := item.(map[string]interface{})
	videoID := itemMap["id"].(map[string]interface{})["videoId"].(string)

	return videoID, nil
}

// Shorts
func (c *YoutubeClient) FetchLatestShortByChannel(channelId string) (model.Post, error) {
	paginatedShorts := PaginatedShortsAPI(channelId)

	resp, err := http.Get(paginatedShorts.String())
	if err != nil {
		fmt.Println(err)
		return model.Post{}, err
	}

	defer resp.Body.Close()

	var response ChannelShortsListResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
		return model.Post{}, err
	}

	if len(response.Items) < 1 || len(response.Items[0].Shorts) < 1 {
		return model.Post{}, errors.New("invalid response")
	}

	author, err := c.FetchChannelName(channelId)
	return *model.NewPost(response.Items[0].Shorts[0].VideoID, response.Items[0].Shorts[0].Title,
		author, platform.YOUTUBE), nil
}

func (c *YoutubeClient) FetchPaginatedShortsByChannel(channelId string) ([]model.Post, string, error) {
	paginatedShorts := PaginatedShortsAPI(channelId)

	resp, err := http.Get(paginatedShorts.String())
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	defer resp.Body.Close()

	var response ChannelShortsListResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		fmt.Println(err)
		return nil, "", err
	}

	var posts []model.Post
	var nextPageToken string

	for _, item := range response.Items {
		for _, short := range item.Shorts {
			author, err := c.FetchChannelName(channelId)
			if err != nil {
			}
			posts = append(posts, *model.NewPost(short.VideoID, short.Title, author, platform.YOUTUBE))
		}
		nextPageToken = item.NextPageToken
	}

	return posts, nextPageToken, nil
}

// Channels
func (c *YoutubeClient) FetchChannelName(channelId string) (string, error) {
	paginatedVideos := c.ChannelInfoAPI(channelId)
	resp, err := http.Get(paginatedVideos.String())
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	item := data["items"].([]interface{})[0]
	itemMap := item.(map[string]interface{})
	channelName := itemMap["snippet"].(map[string]interface{})["title"].(string)
	return channelName, nil
}

func (c *YoutubeClient) PaginatedVideosAPI(channelId string) *url.URL {
	return &url.URL{
		Scheme:     "https",
		Host:       "www.googleapis.com",
		Path:       "/youtube/v3/search",
		ForceQuery: true,
		RawQuery:   fmt.Sprintf("key=%s&channelId=%s&part=snippet,id&order=date&maxResults=25&type=video", c.apiKey, channelId),
	}
}

func (c *YoutubeClient) ChannelInfoAPI(channelId string) *url.URL {
	return &url.URL{
		Scheme:     "https",
		Host:       "www.googleapis.com",
		Path:       "/youtube/v3/channels",
		ForceQuery: true,
		RawQuery:   fmt.Sprintf("key=%s&id=%s&part=snippet,id", c.apiKey, channelId),
	}
}

func PaginatedShortsAPI(videoId string) *url.URL {
	return &url.URL{
		Scheme:     "https",
		Host:       "yt.lemnoslife.com",
		Path:       "channels",
		ForceQuery: true,
		RawQuery:   fmt.Sprintf("part=shorts&id=%s&order=date", videoId),
	}
}