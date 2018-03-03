package nail

import (
	"fmt"
	"log"
	"strings"

	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/goesi/esi"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("killmail", spawnKillmailConsumer)
}

func spawnKillmailConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(s.wait(nsq.HandlerFunc(s.killmailHandler)))
}

func (s *Nail) killmailHandler(message *nsq.Message) error {
	mail := esi.GetKillmailsKillmailIdKillmailHashOk{}
	err := gobcoder.GobDecoder(message.Body, &mail)
	if err != nil {
		log.Println(err)
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
		INSERT INTO evedata.killmails
		(id,solarSystemID,killTime,victimCharacterID,victimCorporationID,victimAllianceID,
		attackerCount,factionID,damageTaken,x,y,z,shipType,warID) 
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE id=id;
		`, mail.KillmailId, mail.SolarSystemId, mail.KillmailTime, mail.Victim.CharacterId, mail.Victim.CorporationId, mail.Victim.AllianceId,
		len(mail.Attackers), mail.Victim.FactionId, mail.Victim.DamageTaken, mail.Victim.Position.X, mail.Victim.Position.Y, mail.Victim.Position.Z, mail.Victim.ShipTypeId,
		mail.WarId)
	if err != nil {
		log.Println(err)
		return err
	}

	var attackers []interface{}
	for _, a := range mail.Attackers {
		attackers = append(attackers, mail.KillmailId, a.CharacterId, a.CorporationId, a.AllianceId, a.ShipTypeId, a.FinalBlow, a.DamageDone, a.WeaponTypeId, a.SecurityStatus)
	}
	if len(attackers) > 0 {
		_, err = tx.Exec(fmt.Sprintf(`INSERT INTO evedata.killmailAttackers
			(id,characterID,corporationID,allianceID,shipType,finalBlow,damageDone,weaponType,securityStatus)
			VALUES %s ON DUPLICATE KEY UPDATE id=id;
			`, joinParameters(9, len(mail.Attackers))), attackers...)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return err
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
