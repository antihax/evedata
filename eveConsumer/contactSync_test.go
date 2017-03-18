package eveConsumer

import (
	"log"
	"testing"
	"time"

	"github.com/antihax/evedata/models"
)

func TestContactSyncConsumer(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()

	err := models.UpdateCorporation(12000, "wardec corp", "WARSRUS", 1200, 0, 500001, "", 22, time.Now())
	if err != nil {
		t.Error(err)
		return
	}
	err = models.UpdateCorporation(12001, "deced corp", "HELP ME", 1200, 0, 0, "", 22, time.Now())
	if err != nil {
		t.Error(err)
		return
	}
	err = models.UpdateCorporation(12002, "pending corp", "oh crap", 1200, 0, 0, "", 22, time.Now())
	if err != nil {
		t.Error(err)
		return
	}
	err = models.UpdateCorporation(12003, "fac war enemy", "muahahaha", 1200, 0, 500004, "", 22, time.Now())
	if err != nil {
		t.Error(err)
		return
	}

	_, err = models.RetryExec(`INSERT INTO evedata.wars
				(id, timeStarted,timeDeclared,openForAllies,cacheUntil,aggressorID,defenderID,mutual)
				VALUES(?,?,?,?,?,?,?,?)
				ON DUPLICATE KEY UPDATE 
					openForAllies=VALUES(openForAllies), 
					mutual=VALUES(mutual), 
					cacheUntil=VALUES(cacheUntil);`,
		10002,
		time.Now().UTC().Add(-time.Hour*24).Format(models.SQLTimeFormat),
		time.Now().UTC().Add(-time.Hour*48).Format(models.SQLTimeFormat),
		false, time.Now(), 12000,
		12001, false)
	if err != nil {
		log.Fatal(err)
		return
	}

	_, err = models.RetryExec(`INSERT INTO evedata.wars
				(id, timeDeclared,openForAllies,cacheUntil,aggressorID,defenderID,mutual)
				VALUES(?,?,?,?,?,?,?)
				ON DUPLICATE KEY UPDATE 
					openForAllies=VALUES(openForAllies), 
					mutual=VALUES(mutual), 
					cacheUntil=VALUES(cacheUntil);`,
		10003,
		time.Now().UTC().Format(models.SQLTimeFormat),
		false, time.Now(), 12000,
		12002, false)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = models.UpdateCharacter(1200, "war dude", 1, 1, 12000, 0, 1, "male", -10, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	err = models.UpdateCharacter(1201, "alt dude", 1, 1, 12004, 0, 1, "male", -10, time.Now())
	if err != nil {
		log.Fatal(err)
		return
	}

	// Add a fake contact sync to the characters created above.
	err = models.AddContactSync(1200, 1200, 1201)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = contactSyncTrigger(eC)
	if err != nil {
		t.Error(err)
		return
	}

	for {
		work, err := contactSyncConsumer(eC, &r)
		if err != nil {
			t.Error(err)
			return
		}
		if work == false {
			break
		}
	}
}
