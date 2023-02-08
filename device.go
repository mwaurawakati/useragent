package useragent

import (
	"embed"
	"encoding/json"
)

var (
	//go:embed data
	f                  embed.FS
	smartphoneDevIDs = loadSmartPhoneDevIDs()
	//TABLET_DEV_IDS     = loadPackageJsonData("/data/tablet_dev_id.json")
)

func loadSmartPhoneDevIDs() devIDs {
	file, err := f.ReadFile("data/smartphone_dev_id.json")
	if err != nil {
		panic(err)
	}
	var devIDs devIDs
	err = json.Unmarshal(file, &devIDs)
	if err != nil {
		panic(err)
	}
	return devIDs

}
