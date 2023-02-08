package useragent

import (
	"embed"
	"encoding/json"
)

var (
	//go:embed data
	f                  embed.FS
	SMARTPHONE_DEV_IDS = loadSmartPhoneDevIDs()
	//TABLET_DEV_IDS     = loadPackageJsonData("/data/tablet_dev_id.json")
)

func loadSmartPhoneDevIDs() DevIDs {
	file, err := f.ReadFile("data/smartphone_dev_id.json")
	if err != nil {
		panic(err)
	}
	var devIDs DevIDs
	err = json.Unmarshal(file, &devIDs)
	if err != nil {
		panic(err)
	}
	return devIDs

}
