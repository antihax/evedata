package eveConsumer

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/antihax/evedata/models"
	"github.com/garyburd/redigo/redis"
)

func init() {
	addConsumer("wallets", walletsConsumer, "EVEDATA_walletQueue")
	addTrigger("wallets", walletsTrigger)
}

// Perform contact sync for wardecs
func walletsTrigger(c *EVEConsumer) (bool, error) {

	// Gather characters for update. Group for optimized updating.
	rows, err := c.ctx.Db.Query(
		`SELECT characterID, tokenCharacterID FROM evedata.crestTokens WHERE 
		walletCacheUntil < UTC_TIMESTAMP() AND lastStatus NOT LIKE "%400 Bad Request%" AND 
		scopes LIKE "%characterWalletRead%";`)
	if err != nil {
		log.Printf("Wallets: Failed query: %v", err)
		return false, err
	}
	defer rows.Close()

	r := c.ctx.Cache.Get()
	defer r.Close()
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

	return true, err
}

func walletsConsumer(c *EVEConsumer, redisPtr *redis.Conn) (bool, error) {
	r := *redisPtr
	ret, err := r.Do("SPOP", "EVEDATA_walletQueue")
	if err != nil {
		return false, err
	} else if ret == nil {
		return false, nil
	}

	v, err := redis.String(ret, err)
	if err != nil {
		return false, fmt.Errorf("error collecting redis string %s string: %s", err, v)
	}

	dest := strings.Split(v, ":")

	if len(dest) != 2 {
		return false, errors.New("Invalid wallet string")
	}

	char, err := strconv.ParseInt(dest[0], 10, 64)
	if err != nil {
		return false, fmt.Errorf("%s string: %s", err, v)
	}
	tokenChar, err := strconv.ParseInt(dest[1], 10, 64)
	if err != nil {
		return false, fmt.Errorf("%s string: %s", err, v)
	}

	token, err := c.ctx.TokenStore.GetTokenSource(char, tokenChar)
	if err != nil {
		return false, fmt.Errorf("%s string: %s", err, v)
	}

	var fromID int64
	for {
		wallets, err := c.ctx.ESI.EVEAPI.CharacterWalletJournalXML(token, tokenChar, fromID)
		if err != nil {
			tokenError(char, tokenChar, nil, err)
			return false, fmt.Errorf("%s %d %d", err, tokenChar, fromID)
		}

		tokenSuccess(char, tokenChar, 200, "OK")

		// there are no entries in this journal page.
		if len(wallets.Entries) == 0 {
			break
		}

		tx, err := models.Begin()
		if err != nil {
			return false, err
		}

		var values []string

		for _, wallet := range wallets.Entries {
			if wallet.RefID < fromID || fromID == 0 {
				fromID = wallet.RefID
			}

			values = append(values, fmt.Sprintf("(%d,%d,%d,%d,%d,%d,%q,%f,%f,%q,%d,%f,%q)",
				tokenChar, wallet.RefID, wallet.RefTypeID, wallet.OwnerID1, wallet.OwnerID2,
				wallet.ArgID1, wallet.ArgName1, wallet.Amount, wallet.Balance,
				wallet.Reason, wallet.TaxReceiverID.Int64, wallet.TaxAmount.Float64, wallet.Date.UTC().Format(models.SQLTimeFormat)))
			if err != nil {
				log.Printf("Wallets: %v %d\n", err, wallet.ArgID1)
				break
			}
		}

		stmt := fmt.Sprintf(`INSERT INTO evedata.walletJournal
								(characterID, refID, refTypeID, ownerID1, ownerID2,
								argID1, argName1, amount, balance,
								reason, taxReceiverID, taxAmount, date)
								VALUES %s ON DUPLICATE KEY UPDATE characterID=characterID;`, strings.Join(values, ",\n"))

		_, err = tx.Exec(stmt)
		if err != nil {
			tx.Rollback()
			return false, err
		}

		_, err = tx.Exec(`UPDATE evedata.crestTokens SET walletCacheUntil = ?
							WHERE characterID = ? AND tokenCharacterID = ?`,
			wallets.CachedUntil.UTC(), char, tokenChar)
		if err != nil {
			log.Printf("Wallets: %v\n", err)
			break
		}

		err = models.RetryTransaction(tx)
		if err != nil {
			log.Printf("%s", err)
			return false, err
		}
	}

	fromID = 0
	for {
		transactions, err := c.ctx.ESI.EVEAPI.CharacterWalletTransactionXML(token, (int64)(tokenChar), fromID)
		if err != nil || transactions == nil {
			tokenError(char, tokenChar, nil, err)
			return false, err
		}

		tokenSuccess(char, tokenChar, 200, "OK")

		// there are no entries in this journal page.
		if len(transactions.Entries) == 0 {
			break
		}

		tx, err := models.Begin()
		if err != nil {
			return false, err
		}

		var values []string

		for _, transaction := range transactions.Entries {
			if transaction.TransactionID < fromID || fromID == 0 {
				fromID = transaction.TransactionID
			}

			values = append(values, fmt.Sprintf("(%d,%d,%d,%d,%f,%d,%d,%q,%q,%d,%q,%d)",
				tokenChar, transaction.TransactionID, transaction.Quantity, transaction.TypeID, transaction.Price,
				transaction.ClientID, transaction.StationID, transaction.TransactionType,
				transaction.TransactionFor, transaction.JournalTransactionID, transaction.TransactionDateTime.UTC().Format(models.SQLTimeFormat), transaction.ClientTypeID))
		}

		stmt := fmt.Sprintf(`INSERT INTO evedata.walletTransactions
								(characterID, transactionID, quantity, typeID, price,
								clientID,  stationID, transactionType,
								transactionFor, journalTransactionID, transactionDateTime, clientTypeID)
								VALUES %s ON DUPLICATE KEY UPDATE characterID=characterID;`, strings.Join(values, ",\n"))

		_, err = tx.Exec(stmt)
		if err != nil {
			tx.Rollback()
			return false, err
		}

		_, err = tx.Exec(`UPDATE evedata.crestTokens SET walletCacheUntil = ?
							WHERE characterID = ? AND tokenCharacterID = ?`,
			transactions.CachedUntil.UTC(), char, tokenChar)
		if err != nil {
			log.Printf("Wallets: %v\n", err)
			break
		}

		err = models.RetryTransaction(tx)
		if err != nil {
			log.Printf("%s", err)
			return false, err
		}
	}

	return true, err
}
