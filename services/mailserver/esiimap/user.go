package esiimap

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
	"github.com/emersion/go-imap/backend"
	"golang.org/x/oauth2"
)

type User struct {
	username    string
	token       oauth2.TokenSource
	backend     *Backend
	characterID int32
	mailboxes   map[string]*Mailbox
}

func NewUser(username string, token oauth2.TokenSource, backend *Backend, characterID int32) *User {
	user := &User{
		username:    username,
		token:       token,
		backend:     backend,
		characterID: characterID,
		mailboxes:   make(map[string]*Mailbox),
	}
	user.loadMailboxes()
	return user
}

func (u *User) Username() string {
	return u.username
}

func (u *User) loadMailboxes() error {
	// Retreive all the mailboxes from ESI
	auth := context.WithValue(context.Background(), goesi.ContextOAuth2, u.token)
	boxes, _, err := u.backend.esi.ESI.MailApi.GetCharactersCharacterIdMailLabels(auth, u.characterID, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	// Create and load all the mailboxes in the background
	for _, box := range boxes.Labels {
		ucn := strings.ToUpper(box.Name)
		u.mailboxes[ucn] = NewMailbox(ucn, box.LabelId, u, box.UnreadCount)
		go u.mailboxes[ucn].loadMailbox()
	}

	// Retreive mailing lists
	mailingLists, _, err := u.backend.esi.ESI.MailApi.GetCharactersCharacterIdMailLists(auth, u.characterID, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	// Cache the mail lists
	if len(mailingLists) > 0 {
		go u.cacheMailingLists(mailingLists)
	}

	return nil
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

func (u *User) ListMailboxes(subscribed bool) (mailboxes []backend.Mailbox, err error) {
	for _, box := range u.mailboxes {
		mailboxes = append(mailboxes, box)
	}
	return mailboxes, nil
}

func (u *User) GetMailbox(name string) (backend.Mailbox, error) {
	mailbox, ok := u.mailboxes[name]
	if !ok {
		log.Printf("Cant find mailbox %s", name)
		return mailbox, errors.New("No such mailbox")
	}
	return mailbox, nil
}

func (u *User) CreateMailbox(name string) error {
	return errors.New("You cannot create mailboxes")
}

func (u *User) DeleteMailbox(name string) error {
	return errors.New("You cannot delete mailboxes")
}

func (u *User) RenameMailbox(existingName, newName string) error {
	return errors.New("You cannot rename mailboxes")
}

func (u *User) Logout() error {
	return nil
}
