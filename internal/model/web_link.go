package model

type CheckLinksStatusByUrlRequest struct {
	Links []string `json:"links" validate:"required,gt=0,lte=200,dive,gt=0,lte=200"`
}

type CheckLinksStatusByUrlResponse struct {
	Links    map[string]string `json:"links"`
	LinksNum uint              `json:"links_num"`
}
