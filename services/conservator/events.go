package conservator

import (
	"database/sql"
	"log"
	"strings"

	"github.com/antihax/goesi"
)

func (c *Conservator) checkAllUsers() {
	c.services.Range(func(ki, vi interface{}) bool {
		service := vi.(Service)
		members, err := service.Server.GetMembers()
		if err != nil {
			return false
		}
		for _, m := range members {
			if err := c.checkUser(m.ID, m.Name, service.IntegrationID, m.Roles); err != nil {
				log.Println(err)
			}
		}
		return true
	})
}

func (c *Conservator) handleNewMember(memberID, memberName, serverID string) {

}

func (c *Conservator) handleMessage(memberID, memberName, serverID string) {

}

func (c *Conservator) checkUser(memberID, memberName string, integrationID int32, roles []string) error {
	server, err := c.getService(integrationID)
	if err != nil {
		return err
	}
	if inSlice("auth", strings.Split(server.Services, ",")) {
		if server.Options.Auth.Members != "" {
			if characterName, err := c.getMemberStatus(memberID, server.EntityID); err != nil {
				return err
			} else if characterName != "" { // Found them
				server.checkAddRoles(memberID, server.Options.Auth.Members, roles)
			} else {
				server.checkRemoveRoles(memberID, server.Options.Auth.Members, roles)
			}
		}

		if server.Options.Auth.PlusTen != "" {
			if characterName, err := c.getPlusTenStatus(memberID, server.EntityID); err != nil {
				return err
			} else if characterName != "" { // Found them
				server.checkAddRoles(memberID, server.Options.Auth.PlusTen, roles)
			} else {
				server.checkRemoveRoles(memberID, server.Options.Auth.PlusTen, roles)
			}
		}

		if server.Options.Auth.PlusFive != "" {
			if characterName, err := c.getPlusFiveStatus(memberID, server.EntityID); err != nil {
				return err
			} else if characterName != "" { // Found them
				server.checkAddRoles(memberID, server.Options.Auth.PlusFive, roles)
			} else {
				server.checkRemoveRoles(memberID, server.Options.Auth.PlusFive, roles)
			}
		}

		if server.Options.Auth.Militia != "" && server.FactionID > 0 {
			if characterName, err := c.getMilitiaStatus(memberID, server.FactionID); err != nil {
				return err
			} else if characterName != "" { // Found them
				server.checkAddRoles(memberID, server.Options.Auth.Militia, roles)
			} else {
				server.checkRemoveRoles(memberID, server.Options.Auth.Militia, roles)
			}
		}

		if server.Options.Auth.AlliedMilitia != "" && server.FactionID > 0 {
			if characterName, err := c.getMilitiaStatus(memberID, goesi.FactionAllies[server.FactionID]); err != nil {
				return err
			} else if characterName != "" { // Found them
				server.checkAddRoles(memberID, server.Options.Auth.AlliedMilitia, roles)
			} else {
				server.checkRemoveRoles(memberID, server.Options.Auth.AlliedMilitia, roles)
			}
		}
	}
	return nil
}

func (c *Conservator) setMemberStatus(memberID string, characterID int32, integrationID int32) error {
	if _, err := c.db.Exec(`
		INSERT INTO evedata.integrationCharacters (integrationID,characterID,integrationUserID)
		VALUES (?,?,?) ON DUPLICATE KEY UPDATE integrationID=integrationID;
		`, integrationID, characterID, memberID); err != nil {
		return err
	}
	return nil
}

func (c *Conservator) getMemberStatus(memberID string, entity int32) (string, error) {
	ref := ""
	if err := c.db.QueryRowx(`
		SELECT characterName
			FROM evedata.integrationCharacters C
			INNER JOIN evedata.crestTokens T ON T.characterID = C.characterID
			WHERE T.lastCode <= 200 AND T.authCharacter = 1 AND integrationUserID = ? AND (allianceID = ? OR corporationID = ?) LIMIT 1;`, memberID, entity, entity).Scan(&ref); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return ref, nil
}

func (c *Conservator) getPlusFiveStatus(memberID string, entity int32) (string, error) {
	ref := ""
	if err := c.db.QueryRowx(`
		SELECT characterName
			FROM evedata.integrationCharacters C
			INNER JOIN evedata.crestTokens T ON T.characterID = C.characterID
			INNER JOIN evedata.entityContacts E ON E.contactID = T.allianceID OR E.contactID = T.corporationID OR E.contactID = T.tokenCharacterID
			WHERE T.lastCode <= 200 AND T.authCharacter = 1 AND integrationUserID = ? AND entityID = ? AND standing = 10 LIMIT 1;`, memberID, entity).Scan(&ref); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return ref, nil
}

func (c *Conservator) getPlusTenStatus(memberID string, entity int32) (string, error) {
	ref := ""
	if err := c.db.QueryRowx(`
		SELECT characterName
			FROM evedata.integrationCharacters C
			INNER JOIN evedata.crestTokens T ON T.characterID = C.characterID
			INNER JOIN evedata.entityContacts E ON E.contactID = T.allianceID OR E.contactID = T.corporationID OR E.contactID = T.tokenCharacterID
			WHERE T.lastCode <= 200 AND T.authCharacter = 1 AND integrationUserID = ?  AND entityID = ? AND standing = 10 LIMIT 1;`, memberID, entity).Scan(&ref); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return ref, nil
}

func (c *Conservator) getMilitiaStatus(memberID string, militia int32) (string, error) {
	ref := ""
	if err := c.db.QueryRowx(`
		SELECT characterName
			FROM evedata.integrationCharacters C
			INNER JOIN evedata.crestTokens T ON T.characterID = C.characterID
			WHERE T.lastCode <= 200 AND T.authCharacter = 1 AND integrationUserID = ? AND factionID = ? LIMIT 1;`, memberID, militia).Scan(&ref); err != nil && err != sql.ErrNoRows {
		return "", err
	}
	return ref, nil
}
