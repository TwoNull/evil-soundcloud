package soundcloud

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/grafov/m3u8"
	jsoniter "github.com/json-iterator/go"
)

const WEB_CLIENTID = "Tl7CY6xVpYugZsGNqmzUhDCRX3urIPNv"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type HydrationData struct {
	Hydratable string      `json:"hydratable"`
	Data       interface{} `json:"data"`
}

type HLSData struct {
	Url string `json:"url"`
}

func getSCPlaylist(url string) ([]HydrationData, error) {
	var data []HydrationData
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bodyText, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	bodyString := string(bodyText)
	startIndex := strings.Index(bodyString, "window.__sc_hydration") + 24
	if startIndex < 24 {
		return nil, errors.New("hydration data not present")
	}
	endIndex := strings.Index(bodyString[startIndex:], "}];") + startIndex + 2
	if endIndex < (startIndex + 2) {
		return nil, errors.New("hydration data end separator not present")
	}
	err = json.UnmarshalFromString(bodyString[startIndex:endIndex], &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func getTrackData(ids string) ([]interface{}, error) {
	var trackData []interface{}
	req, err := http.NewRequest("GET", "https://api-v2.soundcloud.com/tracks?ids="+ids+"&client_id="+WEB_CLIENTID+"&%5Bobject%20Object%5D=&app_version=1718276310&app_locale=en", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.1")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "api-v2.soundcloud.com")
	req.Header.Set("Origin", "https://soundcloud.com")
	req.Header.Set("Referer", "https://soundcloud.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&trackData)
	if err != nil {
		return nil, err
	}
	return trackData, nil
}

func getHLSPlaylist(url string, trackAuthorization string) (*m3u8.MediaPlaylist, error) {
	var hlsData HLSData
	req, err := http.NewRequest("GET", url+"?client_id="+WEB_CLIENTID+"&track_authorization="+trackAuthorization, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", "api-v2.soundcloud.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Origin", "https://soundcloud.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://soundcloud.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&hlsData)
	if err != nil {
		return nil, err
	}

	req2, err := http.NewRequest("GET", hlsData.Url, nil)
	if err != nil {
		return nil, err
	}
	req2.Header.Set("Accept", "*/*")
	req2.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req2.Header.Set("Connection", "keep-alive")
	req2.Header.Set("Host", "cf-hls-media.sndcdn.com")
	req2.Header.Set("Origin", "https://soundcloud.com")
	req2.Header.Set("Referer", "https://soundcloud.com/")
	req2.Header.Set("Sec-Fetch-Dest", "empty")
	req2.Header.Set("Sec-Fetch-Mode", "cors")
	req2.Header.Set("Sec-Fetch-Site", "cross-site")
	req2.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req2.Header.Set("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req2.Header.Set("sec-ch-ua-mobile", "?0")
	req2.Header.Set("sec-ch-ua-platform", `"macOS"`)
	res2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return nil, err
	}
	defer res2.Body.Close()
	p, listType, err := m3u8.DecodeFrom(res2.Body, true)
	if err != nil {
		return nil, err
	}
	switch listType {
	case m3u8.MEDIA:
		return p.(*m3u8.MediaPlaylist), nil
	}
	return nil, errors.New("invalid playlist type")
}

func addSegmentData(f *os.File, url string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "cf-hls-media.sndcdn.com")
	req.Header.Set("Origin", "https://soundcloud.com")
	req.Header.Set("Referer", "https://soundcloud.com/")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	req.Header.Set("sec-ch-ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}
	return nil
}
