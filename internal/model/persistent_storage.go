package model

type PersistentStorageData struct {
	ID          int      `json:"id"`
	LinkedLinks []string `json:"links"`
}
