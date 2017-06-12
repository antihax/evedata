package nail

import (
	"fmt"
	"log"
	"strings"

	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/goesi/v1"
	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("killmail", spawnKillmailConsumer)
}

func spawnKillmailConsumer(s *Nail, consumer *nsq.Consumer) {
	consumer.AddHandler(nsq.HandlerFunc(s.killmailHandler))
}

func (s *Nail) killmailHandler(message *nsq.Message) error {
	mail := goesiv1.GetKillmailsKillmailIdKillmailHashOk{}
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

	_, err = tx.Exec(`
		INSERT INTO evedata.killmails
		(id,solarSystemID,killTime,victimCharacterID,victimCorporationID,victimAllianceID,
		attackerCount,damageTaken,x,y,z,shipType,warID) 
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE id=id;
		`, mail.KillmailId, mail.SolarSystemId, mail.KillmailTime, mail.Victim.CharacterId, mail.Victim.CorporationId, mail.Victim.AllianceId,
		len(mail.Attackers), mail.Victim.DamageTaken, mail.Victim.Position.X, mail.Victim.Position.Y, mail.Victim.Position.Z, mail.Victim.ShipTypeId,
		mail.WarId)
	if err != nil {
		log.Println(err)
		return err
	}

	var parameters []interface{}
	for _, a := range mail.Attackers {
		parameters = append(parameters, mail.KillmailId, a.CharacterId, a.CorporationId, a.AllianceId, a.ShipTypeId, a.FinalBlow, a.DamageDone, a.WeaponTypeId, a.SecurityStatus)
	}
	_, err = tx.Exec(fmt.Sprintf(`INSERT INTO evedata.killmailAttackers
			(id,characterID,corporationID,allianceID,shipType,finalBlow,damageDone,weaponType,securityStatus)
			VALUES %s ON DUPLICATE KEY UPDATE id=id;
			`, joinParameters(9, len(mail.Attackers))), parameters...)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, i := range mail.Victim.Items {
		parameters = append(parameters[:0], mail.KillmailId, i.ItemTypeId, i.Flag, i.QuantityDestroyed, i.QuantityDropped, i.Singleton)
	}
	_, err = tx.Exec(fmt.Sprintf(`INSERT INTO evedata.killmailItems
			(id,itemType,flag,quantityDestroyed,quantityDropped,singleton)
			VALUES %s ON DUPLICATE KEY UPDATE id=id;;
			`, joinParameters(6, len(mail.Victim.Items))), parameters...)
	if err != nil {
		log.Println(err)
		return err
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
