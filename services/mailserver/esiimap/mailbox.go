package esiimap

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/antihax/goesi/esi"
	"github.com/antihax/goesi/optional"
	"github.com/jaytaylor/html2text"

	imap "github.com/antihax/go-imap"
	"github.com/antihax/go-imap/backend/backendutil"
	"github.com/antihax/goesi"
	message "github.com/emersion/go-message"
)

type Mailbox struct {
	name        string
	id          int32
	user        *User
	unreadCount int32
	lastMailID  int64
}

func (mbox *Mailbox) Name() string {
	return mbox.name
}

func (mbox *Mailbox) Info() (*imap.MailboxInfo, error) {
	info := &imap.MailboxInfo{
		Delimiter: "/",
		Name:      mbox.name,
	}
	return info, nil
}

func (mbox *Mailbox) getMailUIDSince(since time.Time) ([]uint32, error) {
	var allMails []uint32
	lastMailID := int64(2147483647)

	auth := context.WithValue(context.Background(), goesi.ContextOAuth2, mbox.user.token)
	for {
		mails, _, err := mbox.user.backend.esi.ESI.MailApi.GetCharactersCharacterIdMail(
			auth,
			mbox.user.characterID,
			&esi.GetCharactersCharacterIdMailOpts{
				Labels:     optional.NewInterface([]int32{mbox.id}),
				LastMailId: optional.NewInt64(lastMailID),
			},
		)
		if err != nil {
			return nil, err
		}

		// breakout if this was the last page
		if len(mails) == 0 {
			break
		}

		// Find the last ID and consolidate mails
		for _, m := range mails {
			if m.MailId < lastMailID {
				lastMailID = m.MailId
			}
			if m.Timestamp.After(since) {
				allMails = append(allMails, uint32(m.MailId))
			} else {
				return allMails, nil
			}
		}
	}

	return allMails, nil
}

func (mbox *Mailbox) getMailHeaders(start int64, end int64) ([]esi.GetCharactersCharacterIdMail200Ok, error) {
	var allMails []esi.GetCharactersCharacterIdMail200Ok
	lastMailID := int64(2147483647)
	// If we have an end, work our way down
	if end > 0 {
		lastMailID = end
	}

	auth := context.WithValue(context.Background(), goesi.ContextOAuth2, mbox.user.token)
	for {
		mails, _, err := mbox.user.backend.esi.ESI.MailApi.GetCharactersCharacterIdMail(
			auth,
			mbox.user.characterID,
			&esi.GetCharactersCharacterIdMailOpts{
				Labels:     optional.NewInterface([]int32{mbox.id}),
				LastMailId: optional.NewInt64(lastMailID),
			},
		)
		if err != nil {
			return nil, err
		}

		// breakout if this was the last page
		if len(mails) == 0 {
			break
		}

		// Find the last ID and consolidate mails
		for _, m := range mails {
			if m.MailId < lastMailID {
				lastMailID = m.MailId
			}
			if end == 0 || m.MailId <= end {
				if start == 0 || m.MailId >= start {
					allMails = append(allMails, m)
				}
			}
		}
	}

	return allMails, nil
}

func (mbox *Mailbox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {

	status := imap.NewMailboxStatus(mbox.name, items)
	status.Flags = []string{}
	status.PermanentFlags = []string{"\\*"}
	status.UnseenSeqNum = 0

	mails, err := mbox.getMailHeaders(0, 0)
	if err != nil {
		return nil, err
	}

	var unseen, messages, nextuid uint32
	for _, m := range mails {
		if !m.IsRead {
			unseen++
		}
		messages++
		if uint32(m.MailId) > nextuid {
			nextuid = uint32(m.MailId)
		}
	}

	for _, name := range items {
		switch name {
		case imap.StatusMessages:
			status.Messages = nextuid // Hack for Win10 Mail Client
		case imap.StatusUidNext:
			status.UidNext = nextuid
		case imap.StatusUidValidity:
			status.UidValidity = 1
		case imap.StatusRecent:
			status.Recent = unseen
		case imap.StatusUnseen:
			status.Unseen = unseen
		}
	}

	return status, nil
}

func (mbox *Mailbox) ListMessages(uid bool, seqSet *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	defer close(ch)

	wg := sync.WaitGroup{}
	sem := make(chan bool, 50)

	for _, seq := range seqSet.Set {
		sem <- true
		wg.Add(1)
		go func(seq imap.Seq) {
			defer func() { wg.Done(); <-sem }()
			if seq.Start == seq.Stop {
				i := imap.NewMessage(seq.Start, items)
				err := mbox.fetchWholeMessage(i, seq.Start, items)
				if err != nil {
					log.Println(err)
					return
				}
				ch <- i
			} else {
				mails, err := mbox.getMailHeaders(int64(seq.Start), int64(seq.Stop))
				if err != nil {
					log.Println(err)
					return
				}
				for _, m := range mails {
					sem <- true
					wg.Add(1)
					go func(m *esi.GetCharactersCharacterIdMail200Ok) {
						defer func() { wg.Done(); <-sem }()
						i := imap.NewMessage(uint32(m.MailId), items)
						err := mbox.fetchMessage(m, i, uint32(m.MailId), items)
						if err != nil {
							log.Println(err)
							return
						}
						ch <- i
					}(&m)
				}
			}
		}(seq)
	}
	wg.Wait()
	return nil
}

func (mbox *Mailbox) fetchMessage(m *esi.GetCharactersCharacterIdMail200Ok, i *imap.Message, seqNum uint32, items []imap.FetchItem) error {
	for _, item := range items {
		n, e, err := mbox.makeFakeHeader(m)
		if err != nil {
			log.Println(err)
			return err
		}

		switch item {
		case imap.FetchEnvelope:
			i.Envelope, _ = backendutil.FetchEnvelope(e.Header)
		case imap.FetchFlags:
			i.Flags = []string{}
			if m.IsRead {
				i.Flags = append(i.Flags, "Seen")
			}
		case imap.FetchInternalDate:
			i.InternalDate = m.Timestamp
		case imap.FetchRFC822Size:
			i.Size = uint32(n) // We're lying
		case imap.FetchUid:
			i.Uid = uint32(seqNum)
		case imap.FetchRFC822Header:
			section, err := imap.ParseBodySectionName(item)
			if err != nil {
				log.Println(err)
				break
			}
			l, _ := backendutil.FetchBodySection(e, section)
			i.Body[section] = l
		default:
			return mbox.fetchWholeMessage(i, seqNum, items)
		}
	}

	return nil
}

func (mbox *Mailbox) fetchWholeMessage(i *imap.Message, seqNum uint32, items []imap.FetchItem) error {
	u := mbox.user

	auth := context.WithValue(context.Background(), goesi.ContextOAuth2, u.token)
	m, _, err := u.backend.esi.ESI.MailApi.GetCharactersCharacterIdMailMailId(auth, u.characterID,
		int64(seqNum), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	n, e, err := mbox.makeFakeBody(&m, int64(seqNum))
	if err != nil {
		log.Println(err)
		return err
	}

	for _, item := range items {
		switch item {
		case imap.FetchEnvelope:
			i.Envelope, _ = backendutil.FetchEnvelope(e.Header)
		case imap.FetchBody, imap.FetchBodyStructure:
			i.BodyStructure, _ = backendutil.FetchBodyStructure(e, item == imap.FetchBodyStructure)
		case imap.FetchFlags:
			i.Flags = []string{}
			if m.Read {
				i.Flags = append(i.Flags, "Seen")
			}
		case imap.FetchInternalDate:
			i.InternalDate = m.Timestamp
		case imap.FetchRFC822Size:
			i.Size = uint32(n)
		case imap.FetchUid:
			i.Uid = uint32(seqNum)
		default:
			section, err := imap.ParseBodySectionName(item)
			if err != nil {
				log.Println(err)
				break
			}
			l, _ := backendutil.FetchBodySection(e, section)
			i.Body[section] = l
		}
	}

	return nil
}

func (mbox *Mailbox) makeFakeHeader(m *esi.GetCharactersCharacterIdMail200Ok) (int, *message.Entity, error) {
	// Make a list of all IDs and a map to the resulting array
	idMap := make(map[int32]int)
	ids := []int32{m.From}
	seen := make(map[int32]bool)
	idMap[m.From] = 0
	seen[m.From] = true
	i := 0
	for _, r := range m.Recipients {
		if !seen[r.RecipientId] && r.RecipientType != "mailing_list" {
			i++
			ids = append(ids, r.RecipientId)
			idMap[r.RecipientId] = i
			seen[r.RecipientId] = true
		}
	}

	// Lookup IDs to names
	names, _, err := mbox.user.backend.lookupAddresses(ids)
	if err != nil {
		return 0, nil, err
	}

	// Build the To list
	to := []string{}
	for _, r := range m.Recipients {
		to = append(to, fmt.Sprintf("%s <%d>", names[idMap[r.RecipientId]], r.RecipientId))
	}

	// Build our fake mail
	s := fmt.Sprintf(`From: %s <%d>
To: %s
Subject: %s
Date: %s
Message-ID: <%d@evedata.org/>
Content-Type: text/plain


Nothing here i'm afraid
`, names[idMap[m.From]], m.From, strings.Join(to, "; "), m.Subject, m.Timestamp.Format(time.RFC822Z), m.MailId)

	e, err := message.Read(bytes.NewReader([]byte(s)))
	return len(s), e, err
}

func (mbox *Mailbox) makeFakeBody(m *esi.GetCharactersCharacterIdMailMailIdOk, id int64) (int, *message.Entity, error) {
	// Make a list of all IDs and a map to the resulting array
	idMap := make(map[int32]int)
	ids := []int32{m.From}
	seen := make(map[int32]bool)
	idMap[m.From] = 0
	seen[m.From] = true
	i := 0
	for _, r := range m.Recipients {
		if !seen[r.RecipientId] {
			i++
			ids = append(ids, r.RecipientId)
			idMap[r.RecipientId] = i
			seen[r.RecipientId] = true
		}
	}

	// Lookup IDs to names
	names, _, err := mbox.user.backend.lookupAddresses(ids)
	if err != nil {
		return 0, nil, err
	}

	// Convert to text/plain
	m.Body = strings.Replace(m.Body, "\n", "<br>", -1) // Hack for breaks..
	plain, err := html2text.FromString(m.Body, html2text.Options{PrettyTables: true})
	if err != nil {
		return 0, nil, err
	}

	// Replace killmails killReport:66991326:b80d548e48c419002cccbe74886b8c05e40af596
	rp := regexp.MustCompile("(?m)killReport:([0-9]+):[a-z0-9]+")
	plain = rp.ReplaceAllString(plain, "https://www.zkillboard.com/kill/$1/")

	// Build the To list
	to := []string{}
	for _, r := range m.Recipients {
		to = append(to, fmt.Sprintf("%s <%d>", names[idMap[r.RecipientId]], r.RecipientId))
	}

	// Build our fake mail
	s := fmt.Sprintf(`From: %s <%d>
To: %s
Subject: %s
Date: %s
Message-ID: <%d@evedata.org/>
Content-Type: text/plain

%s
`, names[idMap[m.From]], m.From, strings.Join(to, "; "), m.Subject, m.Timestamp.Format(time.RFC822Z), id, plain)

	e, err := message.Read(bytes.NewReader([]byte(s)))
	return len(s), e, err
}

func (mbox *Mailbox) SetSubscribed(subscribed bool) error {
	return nil
}

func (mbox *Mailbox) Check() error {
	return nil
}
func (mbox *Mailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	if !criteria.Since.IsZero() {
		uid, err := mbox.getMailUIDSince(criteria.Since)
		fmt.Printf("%v\n", uid)
		return uid, err
	}

	return nil, errors.New("not supported")
}

func (mbox *Mailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	return errors.New("not supported")
}

func (mbox *Mailbox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, op imap.FlagsOp, flags []string) error {
	return nil
}

func (mbox *Mailbox) CopyMessages(uid bool, seqset *imap.SeqSet, destName string) error {
	return nil
}

func (mbox *Mailbox) Expunge() error {
	return nil
}
