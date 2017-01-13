package eveConsumer

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
)

// Perform contact sync for wardecs
func (c *EVEConsumer) walletShouldUpdate() {
	r := c.ctx.Cache.Get()
	defer r.Close()

	// Gather characters for update. Group for optimized updating.
	rows, err := c.ctx.Db.Query(
		`SELECT characterID, tokenCharacterID FROM evedata.crestTokens WHERE 
		walletCacheUntil < UTC_TIMESTAMP() AND lastStatus NOT LIKE "%Invalid refresh token%" AND 
		scopes LIKE "%characterWalletRead%";`)
	if err != nil {
		log.Printf("Wallets: Failed query: %v", err)
		return
	}

	// Loop updatable characters
	for rows.Next() {
		var (
			char      int64 // Source char
			tokenChar int64 // Token Char
		)

		err = rows.Scan(&char, &tokenChar)
		if err != nil {
			log.Printf("Wallets: Failed scan: %v", err)
			continue
		}
		_, err = r.Do("SADD", "EVEDATA_walletQueue", fmt.Sprintf("%d:%d", char, tokenChar))
		if err != nil {
			log.Printf("Wallets: Failed scan: %v", err)
			continue
		}
	}
	rows.Close()
}

func (c *EVEConsumer) walletsCheckQueue(r redis.Conn) error {
	ret, err := r.Do("SPOP", "EVEDATA_walletQueue")
	if err != nil {
		return err
	} else if ret == nil {
		return nil
	}

	v, err := redis.String(ret, err)
	if err != nil {
		return err
	}

	dest := strings.Split(v, ":")

	if len(dest) != 2 {
		return errors.New("Invalid wallet string")
	}

	char, err := strconv.ParseInt(dest[0], 10, 64)
	if err != nil {
		return err
	}
	tokenChar, err := strconv.ParseInt(dest[1], 10, 64)
	if err != nil {
		return err
	}

	token, err := c.getToken(char, tokenChar)
	if err != nil {
		return err
	}

	var fromID int64 = 0
	for {
		wallets, err := c.ctx.EVE.CharacterWalletJournalXML(token, (int64)(tokenChar), fromID)
		if err != nil {
			syncError(char, tokenChar, nil, err)
		} else {
			syncSuccess(char, tokenChar, 200, "OK")
		}
		// there are no entries in this journal page.
		if len(wallets.Entries) == 0 {
			break
		}

		for {
			tx, err := c.ctx.Db.Beginx()
			if err != nil {
				return err
			}
			for _, wallet := range wallets.Entries {
				if wallet.RefID < fromID || fromID == 0 {
					fromID = wallet.RefID
				}

				_, err := tx.Exec(`INSERT IGNORE INTO evedata.walletJournal
								(refID, refTypeID, ownerID1, ownerID2,
								argID1, argName1, amount, balance,
								reason, taxReceiverID, taxAmount, date)
								VALUES (?,?,?,?,?,?,?,?,?,?,?,?);`,
					wallet.RefID, wallet.RefTypeID, wallet.OwnerID1, wallet.OwnerID2,
					wallet.ArgID1, wallet.ArgName1, wallet.Amount, wallet.Balance,
					wallet.Reason, wallet.TaxReceiverID, wallet.TaxAmount, wallet.Date.UTC())
				if err != nil {
					log.Printf("Wallets: %v\n", err)
					break
				}
			}

			_, err = tx.Exec(`UPDATE evedata.crestTokens SET walletCacheUntil = ?
							WHERE characterID = ? AND tokenCharacterID = ?`,
				wallets.CachedUntil.UTC(), char, tokenChar)
			if err != nil {
				log.Printf("Wallets: %v\n", err)
				break
			}

			err = tx.Commit()
			if err != nil {
				log.Printf("Wallets: %v\n", err)
			} else {
				break

			}
		}
	}

	return err
}
