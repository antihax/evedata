package esiimap

import (
	"context"
	"log"
	"strings"

	"github.com/antihax/goesi/esi"
)

// Precache tries to get ahead of id to name lookups
func (s *Backend) precacheLookup() {
	for {
		var (
			ids   []int32
			count int
		)
		knownId := make(map[int32]bool)
		for id := range s.cacheLookup {
			count++
			if !knownId[id] {
				ids = append(ids, id)
				knownId[id] = true
				if count > 200 {
					break
				}
			}
		}
		go s.lookupAddresses(ids)
	}
}

// Precache tries to get ahead of id to name lookups
func (s *Backend) precacheMailingLists() {
	for {
		mailingList := <-s.cacheMailingList
		n, _ := s.cacheQueue.GetCache("addressName", mailingList)
		t, _ := s.cacheQueue.GetCache("addressType", mailingList)
		if n == "" || t == "" {
			s.cacheQueue.SetCache("addressName", mailingList, "## Unknown Mailing List ##")
			s.cacheQueue.SetCache("addressType", mailingList, "mailing_list")
		}
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func (s *Backend) lookupAddresses(ids []int32) ([]string, []string, error) {
	names, err := s.cacheQueue.GetCacheInBulk("addressName", ids)
	if err != nil {
		return nil, nil, err
	}
	types, err := s.cacheQueue.GetCacheInBulk("addressType", ids)
	if err != nil {
		return nil, nil, err
	}

	missing := []int32{}
	missingIdx := []int{}

	for i := range ids {
		if names[i] == "" || types[i] == "" {
			missing = append(missing, ids[i])
			missingIdx = append(missingIdx, i)
		}
	}

	if len(missing) > 0 {
		lookup, _, err := s.esi.ESI.UniverseApi.PostUniverseNames(context.Background(), missing, nil)
		if err != nil {
			log.Printf("%#v %v\n", err, missing)
			if strings.Contains(err.Error(), "404") {
				for i, missingID := range missing {
					lookup, _, err := s.esi.ESI.UniverseApi.PostUniverseNames(context.Background(), []int32{missingID}, nil)
					if err != nil {
						if strings.Contains(err.Error(), "404") {
							names[missingIdx[i]] = "## Unknown Mailing List ##"
							types[missingIdx[i]] = "mailing_list"
						} else {
							return nil, nil, err
						}
					} else {
						for _, e := range lookup {
							names[missingIdx[i]] = e.Name
							types[missingIdx[i]] = e.Category
						}
					}
				}
			} else {
				return nil, nil, err
			}
		} else {
			for i, e := range lookup {
				names[missingIdx[i]] = e.Name
				types[missingIdx[i]] = e.Category
			}
		}

		err = s.cacheQueue.SetCacheInBulk("addressName", ids, names)
		if err != nil {
			return nil, nil, err
		}
		err = s.cacheQueue.SetCacheInBulk("addressType", ids, types)
		if err != nil {
			return nil, nil, err
		}
	}

	return names, types, nil
}

func (u *User) cacheMailingLists(mailingLists []esi.GetCharactersCharacterIdMailLists200Ok) {
	names := []string{}
	types := []string{}
	ids := []int32{}

	for _, e := range mailingLists {
		names = append(names, e.Name)
		types = append(types, "mailing_list")
		ids = append(ids, e.MailingListId)
	}

	if err := u.backend.cacheQueue.SetCacheInBulk("addressName", ids, names); err != nil {
		log.Println(err)
		return
	}
	if err := u.backend.cacheQueue.SetCacheInBulk("addressType", ids, types); err != nil {
		log.Println(err)
		return
	}
	return
}
