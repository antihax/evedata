package esismtp

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/antihax/goesi/esi"

	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/internal/tokenstore"
	"github.com/antihax/goesi"
	smtp "github.com/emersion/go-smtp"
	"golang.org/x/oauth2"
)

func New(tokenAPI *tokenstore.TokenServerAPI, esi *goesi.APIClient, tokenAuth *goesi.SSOAuthenticator, q *redisqueue.RedisQueue) *Backend {
	return &Backend{tokenAPI, esi, tokenAuth, q}
}

type Backend struct {
	tokenAPI   *tokenstore.TokenServerAPI
	esi        *goesi.APIClient
	tokenAuth  *goesi.SSOAuthenticator
	cacheQueue *redisqueue.RedisQueue
}

func (s *Backend) Login(username, password string) (smtp.User, error) {

	parts := strings.Split(username, "@")
	characterID, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	u, err := s.tokenAPI.GetMailUser(int32(characterID), password)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	ts, err := s.tokenAuth.TokenSource(u.Token)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &User{
		username:    username,
		token:       ts,
		backend:     s,
		characterID: u.TokenCharacterID}, nil
}

// Require clients to authenticate using SMTP AUTH before sending emails
func (bkd *Backend) AnonymousLogin() (smtp.User, error) {
	return nil, smtp.ErrAuthRequired
}

type User struct {
	username    string
	token       oauth2.TokenSource
	backend     *Backend
	characterID int32
}

func (u *User) Send(from string, to []string, r io.Reader) error {

	var ids []int32

	if len(to) > 50 {
		return errors.New("Cannot send to more than 50 recepients at a time")
	}

	// Find all the recepients and validate they are id numbers
	for _, email := range to {
		s := strings.Split(email, "@")
		id, err := strconv.ParseInt(s[0], 10, 32)
		if err != nil {
			return err
		}
		ids = append(ids, int32(id))
	}

	// Lookup the IDs
	_, types, err := u.backend.lookupAddresses(ids)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	// Find the subject
	subjectRegx := regexp.MustCompile(`Subject: (.*)`)
	subjectMatch := subjectRegx.FindStringSubmatch(string(b))
	subject := subjectMatch[1]

	// Find the body
	bodyRegx := regexp.MustCompile(`(?s)\n\n(.*)`)
	bodyMatch := bodyRegx.FindStringSubmatch(string(b))
	body := bodyMatch[1]

	if subject == "" || body == "" {
		return errors.New("Did not understand mail format")
	}

	// Build the recepient list
	var recepients []esi.PostCharactersCharacterIdMailRecipient
	for i, id := range ids {
		recepients = append(recepients,
			esi.PostCharactersCharacterIdMailRecipient{
				RecipientId:   id,
				RecipientType: types[i],
			},
		)
	}

	auth := context.WithValue(context.Background(), goesi.ContextOAuth2, u.token)
	_, _, err = u.backend.esi.ESI.MailApi.PostCharactersCharacterIdMail(auth, u.characterID,
		esi.PostCharactersCharacterIdMailMail{
			ApprovedCost: 0,
			Subject:      subject,
			Body:         body,
			Recipients:   recepients,
		}, nil)

	return err
}

func (u *User) Logout() error {
	return nil
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
