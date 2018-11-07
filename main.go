package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/cts-rmm/rmm/rmmcore/ESNI"
)

var path = flag.String("path", " ", "Path to xml media file.")

func main() {
	if *path == "" {
		fmt.Println("Provide the path to xml media file.")
		return
	}
	flag.Parse() // Parse first time to get config path.
	medias := parseXML(*path)

	if len(medias) > 0 {

		mediasCheck := map[string]map[string][]string{}
		for i, media := range medias {

			mediasCheck[media.ID] = map[string][]string{
				"starts": []string{},
				"ends":   []string{},
			}

			if len(medias[i].MediaPoints) > 0 {
				for ii, mediapoint := range media.MediaPoints {
					if ii%2 == 1 {
						if strings.Contains(mediapoint.ID, "start") {
							mediasCheck[media.ID]["starts"] = append(mediasCheck[media.ID]["starts"], mediapoint.ID)
						}
						if strings.Contains(mediapoint.ID, "end") {
							mediasCheck[media.ID]["ends"] = append(mediasCheck[media.ID]["ends"], mediapoint.ID)
						}
					}
				}

			}
		}

		for key, val := range mediasCheck {
			if len(val["starts"]) > len(val["ends"]) {
				if len(val["ends"]) > 0 {
					fmt.Println("Out of Order: ")
					for _, end := range val["ends"] {
						fmt.Println("\t- ", end)
					}
				} else {
					fmt.Println("In Order: ", key)
				}
			}
		}
	}

}
func parseXML(configPath string) []*ESNI.Media {

	file, err := os.Open(configPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	medias := []*ESNI.Media{}
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&medias)
	if err != nil {
		panic(err)
	}
	return medias
}
