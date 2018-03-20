package artifice

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/antihax/goesi"
	"github.com/antihax/goesi/esi"
)

func init() {
	registerTrigger("mailNotifications", mailNotifications, time.NewTicker(time.Second*100))
}

func mailNotifications(s *Artifice) error {

	type mailRecipient struct {
		CharacterID      int32  `db:"characterID"`
		TokenCharacterID int32  `db:"tokenCharacterID"`
		LastStatus       string `db:"lastStatus"`
		CharacterName    string `db:"characterName"`
	}
	recipient := mailRecipient{}

	if err := s.db.QueryRowx(
		`	SELECT characterID, tokenCharacterID, characterName, lastStatus FROM evedata.crestTokens
			WHERE lastCode = 999 AND mailedError = 0 AND lastStatus LIKE "%invalid_token%"
			LIMIT 1;`).StructScan(&recipient); err != nil {
		// Ignore this error.
		if strings.Contains(err.Error(), "no rows in result set") {
			return nil
		}
		return err
	}

	_, err := s.db.Exec(`UPDATE evedata.crestTokens SET mailedError = 1 WHERE characterID = ? AND tokenCharacterID = ?`,
		recipient.CharacterID, recipient.TokenCharacterID)
	if err != nil {
		return err
	}

	mail := esi.PostCharactersCharacterIdMailMail{
		Recipients: []esi.PostCharactersCharacterIdMailRecipient{
			{
				RecipientId:   recipient.CharacterID,
				RecipientType: "character",
			},
		},
		Subject: "EVEData.org: Character Update Failure!",
		Body: `Hi, we ran into an issue with one of your characters and need you to reauthenticate the character with us.

Please login with the character receiving this evemail and update any characters with a bad status. After logging in, you can simply click <i>Add Character</i> and overwrite the character with a bad status.

Login here: <a href="https://www.evedata.org/account">https://www.evedata.org/account</a>

This failure could happen for a number of reasons including:
	- The account password was changed,
	- The authentication is more than a year old, or
	- The character was transfered.

Leaving the character in a failed state will affect services provided by EVEData.org.

Thanks,

EveDataRules`,
	}

	s.mail <- mail
	return nil
}

func (s *Artifice) mailCorporationChangeWithShares(characterID int32) {
	mail := esi.PostCharactersCharacterIdMailMail{
		Recipients: []esi.PostCharactersCharacterIdMailRecipient{
			{
				RecipientId:   characterID,
				RecipientType: "character",
			},
		},
		Subject: "EVEData.org: Corporation Change Detected!",
		Body: `Hi, we noticed you just changed corporations and wanted to inform that this character is sharing data to other entities.

Please log into our site with your main character and verify that you wish to continue sharing data with these entities.
View shares here: <a href="https://www.evedata.org/shares">https://www.evedata.org/shares</a>.

For security purposes, we do not divulge details in evemail.

Thanks,

EveDataRules`,
	}
	s.mail <- mail
}

func (s *Artifice) mailRunner() {
	throttle := time.Tick(time.Second * 12)
	auth := context.WithValue(context.Background(), goesi.ContextOAuth2, *s.token)
	for {
		m := <-s.mail

		mailID, _, err := s.esi.ESI.MailApi.PostCharactersCharacterIdMail(auth, s.tokenCharID, m, nil)
		log.Printf("Mailed %v about %s. Mail %d %s\n", m.Recipients, m.Subject, mailID, err)
		<-throttle
	}
}
