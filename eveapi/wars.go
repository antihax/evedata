package eveapi

import "fmt"

type Wars struct {
	*AnonymousClient
	crestPagedFrame

	Items []struct {
		HRef string
		ID   int
	}
}

type War struct {
	*AnonymousClient

	crestSimpleFrame

	TimeFinished  EVETime
	OpenForAllies bool
	TimeStarted   EVETime
	AllyCount     int
	TimeDeclared  EVETime

	Allies []struct {
		HRef string
		ID   int
		Icon struct {
			HRef string
		}
		Name string
	}
	Aggressor struct {
		ShipsKilled int

		Name string
		HRef string

		Icon struct {
			HRef string
		}
		ID        int
		IskKilled float64
	}
	Mutual bool

	Killmails string

	Defender struct {
		ShipsKilled int

		Name string
		HRef string

		Icon struct {
			HRef string
		}
		ID        int
		IskKilled float64
	}
	ID int
}

func (c *AnonymousClient) Wars(page int) (*Wars, error) {
	w := &Wars{AnonymousClient: c}
	url := c.base.CREST + fmt.Sprintf("wars/?page=%d", page)
	res, err := c.httpClient.Get(url)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	err = decode(res, w)
	if err != nil {
		return nil, err
	}

	w.getFrameInfo(url, res)
	return w, nil
}

func (c *Wars) NextPage() (*Wars, error) {
	w := &Wars{AnonymousClient: c.AnonymousClient}
	if c.Next.HRef == "" {
		return nil, nil
	}
	res, err := c.httpClient.Get(c.Next.HRef)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	err = decode(res, w)
	if err != nil {
		return nil, err
	}

	w.getFrameInfo(c.Next.HRef, res)
	return w, nil
}

func (c *AnonymousClient) War(href string) (*War, error) {
	w := &War{AnonymousClient: c}
	url := href
	res, err := c.httpClient.Get(url)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	err = decode(res, w)
	if err != nil {
		return nil, err
	}

	w.getFrameInfo(res)
	return w, nil
}
