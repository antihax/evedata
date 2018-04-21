package esiimap

import (
	"context"
	"errors"
	"strings"

	"github.com/antihax/goesi"
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

func (u *User) Username() string {
	return u.username
}

func (u *User) ListMailboxes(subscribed bool) (mailboxes []backend.Mailbox, err error) {
	auth := context.WithValue(context.Background(), goesi.ContextOAuth2, u.token)
	boxes, _, err := u.backend.esi.ESI.MailApi.GetCharactersCharacterIdMailLabels(auth, u.characterID, nil)
	if err != nil {
		return nil, err
	}

	for _, box := range boxes.Labels {
		ucn := strings.ToUpper(box.Name)
		if _, ok := u.mailboxes[ucn]; !ok {
			u.mailboxes[ucn] = &Mailbox{
				name:        ucn,
				id:          box.LabelId,
				user:        u,
				unreadCount: box.UnreadCount,
			}

		}
		mailboxes = append(mailboxes, u.mailboxes[ucn])
	}

	return mailboxes, err
}

func (u *User) GetMailbox(name string) (backend.Mailbox, error) {
	mailbox, ok := u.mailboxes[name]
	if !ok {
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
