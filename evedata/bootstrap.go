package evedata

import "github.com/antihax/evedata/appContext"

func bootstrap(c *appContext.AppContext) error {
	if err := bootstrapRefID(c); err != nil {
		return err
	}
	return nil
}

func bootstrapRefID(c *appContext.AppContext) error {
	refTypes, err := c.ESI.EVEAPI.RefTypesXML()
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
