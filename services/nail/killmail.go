package nail

import (
	"fmt"
	"log"
	"strings"

	"github.com/antihax/evedata/internal/datapackages"
	"github.com/antihax/evedata/services/vanguard/models"

	"github.com/antihax/evedata/internal/gobcoder"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("killmail", spawnKillmailConsumer)
}

func spawnKillmailConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.killmailHandler)))
}

func (s *Nail) killmailHandler(message *nsq.Message) error {
	killmail := datapackages.Killmail{}
	err := gobcoder.GobDecoder(message.Body, &killmail)
	if err != nil {
		log.Println(err)
		return err
	}

	mail := killmail.Kill
	err = s.doSQL(`
		INSERT INTO evedata.killmails
		(id,solarSystemID,killTime,victimCharacterID,victimCorporationID,victimAllianceID,
		attackerCount,factionID,damageTaken,x,y,z,shipType,warID,hash) 
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE id=id;
		`, mail.KillmailId, mail.SolarSystemId, mail.KillmailTime.Format(models.SQLTimeFormat), mail.Victim.CharacterId, mail.Victim.CorporationId, mail.Victim.AllianceId,
		len(mail.Attackers), mail.Victim.FactionId, mail.Victim.DamageTaken, mail.Victim.Position.X, mail.Victim.Position.Y, mail.Victim.Position.Z, mail.Victim.ShipTypeId,
		mail.WarId, killmail.Hash)
	if err != nil {
		log.Println(err)
		return err
	}

	var attackers []interface{}
	for _, a := range mail.Attackers {
		attackers = append(attackers, mail.KillmailId, a.CharacterId, a.CorporationId, a.AllianceId, a.ShipTypeId, a.SecurityStatus)
	}
	if len(attackers) > 0 {
		err = s.doSQL(fmt.Sprintf(`INSERT INTO evedata.killmailAttackers
			(id,characterID,corporationID,allianceID,shipType,securityStatus)
			VALUES %s ON DUPLICATE KEY UPDATE id=id;
			`, joinParameters(6, len(mail.Attackers))), attackers...)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	// Mark the killmail complete to prevent duplicates
	err = s.outQueue.SetWorkCompleted("evedata_known_kills", int64(mail.KillmailId))
	if err != nil {
		log.Println(err)
	}

	return nil
}

func joinParameters(nParam, nEntries int) string {
	s := strings.Repeat("?,", nParam)
	s = "(" + s[:len(s)-1] + "),"

	s = strings.Repeat(s, nEntries)
	s = s[:len(s)-1]

	return s
}
