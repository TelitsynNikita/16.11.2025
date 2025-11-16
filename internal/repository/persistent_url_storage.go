package repository

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"workmate/internal/model"
)

type PersistentURLStorage struct {
	localStorageMutex      sync.Mutex
	localStorage           map[int]string
	persistentStorageMutex sync.Mutex
}

func NewPersistentURLStorage() *PersistentURLStorage {
	return &PersistentURLStorage{
		localStorage: make(map[int]string),
	}
}

func (p *PersistentURLStorage) GetLinksByUrl(urls []string) (int, []string, error) {
	p.localStorageMutex.Lock()
	defer p.localStorageMutex.Unlock()
	if len(p.localStorage) == 0 {
		err := p.ReadFileToLocalStorage()
		if err != nil {
			return 0, nil, err
		}
	}

	sort.Strings(urls)
	urlsJoined := strings.Join(urls, ",")
	encodeUrlStrings := base64.StdEncoding.EncodeToString([]byte(urlsJoined))

	var theMostIndexInMap = 0
	for key, value := range p.localStorage {
		if encodeUrlStrings == value {
			return key, urls, nil
		}

		if theMostIndexInMap < key {
			theMostIndexInMap = key
		}
	}

	p.localStorage[theMostIndexInMap+1] = encodeUrlStrings

	return theMostIndexInMap + 1, urls, nil
}

func (p *PersistentURLStorage) GetUrlByIDs(ids []int) ([]model.PersistentStorageData, error) {
	p.localStorageMutex.Lock()
	defer p.localStorageMutex.Unlock()
	if len(p.localStorage) == 0 {
		err := p.ReadFileToLocalStorage()
		if err != nil {
			return nil, err
		}
	}

	var links []model.PersistentStorageData
	for _, id := range ids {
		encodeString, ok := p.localStorage[id]
		if !ok {
			return nil, fmt.Errorf("there is no data by id: %d", id)
		}

		data, err := base64.StdEncoding.DecodeString(encodeString)
		if err != nil {
			return nil, err
		}

		links = append(links, model.PersistentStorageData{
			ID:          id,
			LinkedLinks: strings.Split(string(data), ","),
		})
	}

	if len(links) == 0 {
		return nil, fmt.Errorf("there is no data by ids: %v", ids)
	}

	return links, nil
}

func (p *PersistentURLStorage) ReadFileToLocalStorage() error {
	p.persistentStorageMutex.Lock()
	defer p.persistentStorageMutex.Unlock()
	file, err := p.readPersistentStorage()
	if err != nil {
		return err
	}

	var data map[int]string
	err = json.Unmarshal(file, &data)
	if err != nil {
		return err
	}

	p.localStorage = data

	return nil
}

func (p *PersistentURLStorage) WriteDataToFileAndLocalStorage() error {
	p.persistentStorageMutex.Lock()
	defer p.persistentStorageMutex.Unlock()
	if len(p.localStorage) == 0 {
		return nil
	}

	err := p.removeAllDataFromPersistentStorage()
	if err != nil {
		return err
	}

	data, err := json.Marshal(p.localStorage)
	if err != nil {
		return err
	}

	err = p.writeDataToPersistentStorage(data)
	if err != nil {
		return err
	}

	return nil
}

func (p *PersistentURLStorage) InitPersistentStorage() error {
	file, err := p.readPersistentStorage()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = p.createPersistentStorage()
			if err != nil {
				return err
			}

			err = p.writeDataToPersistentStorage([]byte("{}"))
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	if len(file) == 0 {
		err = p.writeDataToPersistentStorage([]byte("{}"))
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PersistentURLStorage) createPersistentStorage() error {
	file, err := os.Create("persistent_storage.txt")
	if err != nil || file == nil {
		return err
	}
	defer file.Close()

	return nil
}

func (p *PersistentURLStorage) writeDataToPersistentStorage(data []byte) error {
	err := os.WriteFile("persistent_storage.txt", data, 0666)
	if err != nil {
		return err
	}

	return nil
}

func (p *PersistentURLStorage) removeAllDataFromPersistentStorage() error {
	err := exec.Command("/bin/bash", "-c", "echo > ./persistent_storage.txt").Run()
	if err != nil {
		return err
	}

	return nil
}

func (p *PersistentURLStorage) readPersistentStorage() ([]byte, error) {
	file, err := os.ReadFile("persistent_storage.txt")
	if err != nil {
		return nil, err
	}

	return file, nil
}
