package service

import (
	"net/http"
	"sync"
	"workmate/internal/model"
	"workmate/internal/repository"

	"github.com/gofiber/fiber/v2"
)

type URL struct {
	URLStorage repository.URLStorageRepository
}

func NewURLService(repo repository.URLStorageRepository) *URL {
	return &URL{
		URLStorage: repo,
	}
}

func (u *URL) GetUrlByID(ids []int) ([]model.CheckLinksStatusByUrlResponse, error) {
	links, err := u.URLStorage.GetUrlByIDs(ids)
	if err != nil {
		return nil, err
	}

	var linksStatus struct {
		mu    sync.Mutex
		wg    sync.WaitGroup
		links []model.CheckLinksStatusByUrlResponse
	}

	linksStatus.wg.Add(len(links))
	for _, link := range links {
		go func(persistentStorageData model.PersistentStorageData) {
			defer linksStatus.wg.Done()

			var mu sync.Mutex
			var wg sync.WaitGroup
			var linkStatus = model.CheckLinksStatusByUrlResponse{
				Links: make(map[string]string, len(persistentStorageData.LinkedLinks)),
			}
			linkStatus.LinksNum = link.ID
			for _, url := range persistentStorageData.LinkedLinks {
				wg.Add(1)
				go func(url string) {
					mu.Lock()
					defer mu.Unlock()
					defer wg.Done()
					agent := fiber.Get(url)
					statusCode, _, errs := agent.Bytes()
					if len(errs) > 0 {
						linkStatus.Links[url] = "not available"
						return
					}

					if statusCode == http.StatusNoContent || statusCode == http.StatusOK {
						linkStatus.Links[url] = "available"
					} else {
						linkStatus.Links[url] = "not available"
					}
				}(url)
			}

			linksStatus.mu.Lock()
			linksStatus.links = append(linksStatus.links, linkStatus)
			linksStatus.mu.Unlock()

			wg.Wait()
		}(link)
	}

	linksStatus.wg.Wait()

	return linksStatus.links, nil
}

func (u *URL) CheckLinksStatusByUrl(urls []string) (model.CheckLinksStatusByUrlResponse, error) {
	id, links, err := u.URLStorage.GetLinksByUrl(urls)
	if err != nil {
		return model.CheckLinksStatusByUrlResponse{}, err
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	var linkStatus = model.CheckLinksStatusByUrlResponse{
		Links:    make(map[string]string, len(links)),
		LinksNum: id,
	}

	wg.Add(len(links))
	for _, link := range links {
		go func(url string) {
			defer wg.Done()
			mu.Lock()
			defer mu.Unlock()
			agent := fiber.Get(url)
			statusCode, _, errs := agent.Bytes()
			if len(errs) > 0 {
				linkStatus.Links[url] = "not available"
				return
			}

			if statusCode == http.StatusNoContent || statusCode == http.StatusOK {
				linkStatus.Links[url] = "available"
			} else {
				linkStatus.Links[url] = "not available"
			}
		}(link)
	}

	wg.Wait()

	return linkStatus, nil
}
