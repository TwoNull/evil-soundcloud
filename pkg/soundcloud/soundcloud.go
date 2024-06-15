package soundcloud

import (
	"fmt"
	"log"
	"os"

	"github.com/grafov/m3u8"
)

func GetTracks(url string, outputDir string) {
	hyData, err := getSCPlaylist(url)
	if err != nil {
		log.Fatal(err)
	}
	tracks := hyData[len(hyData)-1].Data.(map[string]interface{})["tracks"].([]interface{})
	// Patch missing data
	trackIds := ""
	trackIndices := []int{}
	for i, t := range tracks {
		if t.(map[string]interface{})["media"] == nil {
			trackIds += fmt.Sprint(int(t.(map[string]interface{})["id"].(float64))) + ","
			trackIndices = append(trackIndices, i)
		}
	}
	trackData, err := getTrackData(trackIds)
	if err != nil {
		log.Fatal(err)
	}
	for i, v := range trackIndices {
		tracks[v] = trackData[i]
	}
	// Get HLS playlists for tracks
	for _, t := range tracks {
		trackAuthorization := t.(map[string]interface{})["track_authorization"].(string)
		hlsUrl := t.(map[string]interface{})["media"].(map[string]interface{})["transcodings"].([]interface{})[0].(map[string]interface{})["url"].(string)
		playlist, err := getHLSPlaylist(hlsUrl, trackAuthorization)
		if err != nil {
			log.Fatal(err)
		}
		err = saveTrack(t, playlist, outputDir)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func saveTrack(track interface{}, playlist *m3u8.MediaPlaylist, outputDir string) error {
	f, err := os.Create(outputDir + "/" + track.(map[string]interface{})["title"].(string) + ".mp3")
	if err != nil {
		return err
	}
	defer f.Close()
	for _, v := range playlist.GetAllSegments() {
		err = addSegmentData(f, v.URI)
		if err != nil {
			return err
		}
	}
	return nil
}
