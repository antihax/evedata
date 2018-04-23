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

	"github.com/antihax/goesi"
	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend/backendutil"
	message "github.com/emersion/go-message"
)

type Mailbox struct {
	name           string
	id             int32
	user           *User
	unreadCount    uint32
	nextuid        uint32
	firstuid       uint32
	count          uint32
	validity       uint32
	messageHeaders map[uint32]*esi.GetCharactersCharacterIdMail200Ok
	messagesUuid   []*esi.GetCharactersCharacterIdMail200Ok
	loaded         sync.WaitGroup
}

func NewMailbox(ucn string, id int32, u *User, unreadCount int32) *Mailbox {
	return &Mailbox{
		name:           ucn,
		id:             id,
		user:           u,
		unreadCount:    uint32(unreadCount),
		messageHeaders: make(map[uint32]*esi.GetCharactersCharacterIdMail200Ok),
	}
}

func (mbox *Mailbox) WaitForLoad() {
	mbox.loaded.Wait()
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

// Load mail headers into the mailbox
func (mbox *Mailbox) loadMailbox() {
	mbox.loaded.Add(1)
	defer mbox.loaded.Done()
	// Get all mail headers
	lastMailID := int64(2147483647)
	mbox.validity = uint32(time.Now().Unix())

	var unseen, count uint32
	messageHeaders := make(map[int64]*esi.GetCharactersCharacterIdMail200Ok)

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
			log.Println(err)
			continue
		}

		// breakout if this was the last page
		if len(mails) == 0 {
			break
		}

		// Store pointers to the mail headers
		for i, m := range mails {
			if _, ok := messageHeaders[m.MailId]; !ok {
				messageHeaders[m.MailId] = &mails[i]
				mbox.messagesUuid = append(mbox.messagesUuid, &mails[i])
				if !m.IsRead {
					unseen++
				}
				count++
			}

			if m.MailId < lastMailID {
				lastMailID = m.MailId
			}
		}
	}

	// Reverse UUID slice so oldest is first
	for i := len(mbox.messagesUuid)/2 - 1; i >= 0; i-- {
		opp := len(mbox.messagesUuid) - 1 - i
		mbox.messagesUuid[i], mbox.messagesUuid[opp] = mbox.messagesUuid[opp], mbox.messagesUuid[i]

	}

	// Put the messages into the map with correct uid
	for i, m := range mbox.messagesUuid {
		mbox.messageHeaders[uint32(i)] = m
	}

	mbox.nextuid = count + 1
	mbox.firstuid = 0
	mbox.unreadCount = unseen
	mbox.count = count
}

func (mbox *Mailbox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	mbox.WaitForLoad()
	status := imap.NewMailboxStatus(mbox.name, items)
	status.Flags = []string{}
	status.PermanentFlags = []string{"\\*"}
	status.UnseenSeqNum = mbox.firstuid
	for _, name := range items {
		switch name {
		case imap.StatusMessages:
			status.Messages = mbox.count
		case imap.StatusUidNext:
			status.UidNext = mbox.nextuid
		case imap.StatusUidValidity:
			status.UidValidity = mbox.validity
		case imap.StatusRecent:
			status.Recent = mbox.unreadCount
		case imap.StatusUnseen:
			status.Unseen = mbox.unreadCount
		}
	}

	return status, nil
}

func (mbox *Mailbox) ListMessages(uid bool, seqSet *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	defer close(ch)
	mbox.WaitForLoad()
	wg := sync.WaitGroup{}
	sem := make(chan bool, 50)
	for i, m := range mbox.messagesUuid {
		if seqSet.Contains(uint32(i)) {
			sem <- true
			wg.Add(1)
			go func(m *esi.GetCharactersCharacterIdMail200Ok, i int) {
				defer func() { wg.Done(); <-sem }()
				im := imap.NewMessage(uint32(i), items)
				err := mbox.fetchMessage(m, im, uint32(i), items)
				if err != nil {
					log.Println(err)
					return
				}
				ch <- im
			}(m, i)
		}
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
			return mbox.fetchWholeMessage(i, seqNum, m.MailId, items)
		}
	}

	return nil
}

func (mbox *Mailbox) fetchWholeMessage(i *imap.Message, uuid uint32, mailID int64, items []imap.FetchItem) error {
	u := mbox.user

	auth := context.WithValue(context.Background(), goesi.ContextOAuth2, u.token)
	m, _, err := u.backend.esi.ESI.MailApi.GetCharactersCharacterIdMailMailId(auth, u.characterID,
		int64(mailID), nil)
	if err != nil {
		log.Println(err)
		return err
	}

	n, e, err := mbox.makeFakeBody(&m, mailID)
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
			i.Uid = uuid
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

func (mbox *Mailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	mbox.WaitForLoad()

	var ids []uint32
	for i, msg := range mbox.messagesUuid {
		ok, err := mbox.MatchMessage(msg, uint32(i), criteria)
		if err != nil || !ok {
			continue
		}
		ids = append(ids, uint32(i))
	}
	return ids, nil
}

func (mbox *Mailbox) MatchMessage(m *esi.GetCharactersCharacterIdMail200Ok, seqNum uint32, c *imap.SearchCriteria) (bool, error) {
	if !MatchSeqNumAndUid(seqNum, seqNum, c) {
		return false, nil
	}

	if !MatchDate(m.Timestamp, c) {
		return false, nil
	}

	return true, nil
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
		log.Println(err)
		return 0, nil, err
	}

	// Build the To list
	to := []string{}
	for _, r := range m.Recipients {
		to = append(to, fmt.Sprintf("%s <%d>", names[idMap[r.RecipientId]], r.RecipientId))
	}

	// Build our fake mail
	s := fmt.Sprintf(`From: %s <%d@evedata.org>
To: %s
Subject: %s
Date: %s
Message-ID: <%d@evedata.org/>
Content-Type: text/plain

Nothing here i'm afraid
`, names[idMap[m.From]], m.From, strings.Join(to, "; "), m.Subject, m.Timestamp.Format(time.RFC822Z), m.MailId)

	e, err := message.Read(bytes.NewReader([]byte(s)))
	if err != nil {
		log.Println(err)
	}
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
		log.Println(err)
		return 0, nil, err
	}

	// Convert to text/plain
	m.Body = strings.Replace(m.Body, "\n", "<br>", -1) // Hack for breaks..
	plain, err := html2text.FromString(m.Body, html2text.Options{PrettyTables: true})
	if err != nil {
		log.Println(err)
		return 0, nil, err
	}

	// Replace killmails killReport:66991326:b80d548e48c419002cccbe74886b8c05e40af596
	rp := regexp.MustCompile("(?m)killReport:([0-9]+):[a-z0-9]+")
	plain = rp.ReplaceAllString(plain, "https://www.zkillboard.com/kill/$1/")

	// Build the To list
	to := []string{}
	for _, r := range m.Recipients {
		to = append(to, fmt.Sprintf("%s <%d@evedata.org>", names[idMap[r.RecipientId]], r.RecipientId))
	}

	// Build our fake mail
	s := fmt.Sprintf(`From: %s <%d@evedata.org>
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
