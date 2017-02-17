package evedata

import "github.com/antihax/evedata/appContext"

func bootstrap(c *appContext.AppContext) error {
	if err := boostrapRefID(c); err != nil {
		return err
	}
	return nil
}

func boostrapRefID(c *appContext.AppContext) error {
	refTypes, err := c.EVE.RefTypesXML()
	if err != nil {
		return err
	}

	for _, r := range refTypes.RefTypes {
		_, err := c.Db.Exec(`INSERT INTO evedata.walletJournalRefType (refTypeID, refTypeName) 
            VALUES(?,?) ON DUPLICATE KEY UPDATE refID = refID`,
            r.RefTypeID, r.RefTypeName)
		if err != nil {
			return err
		}
	}

	return nil
}
