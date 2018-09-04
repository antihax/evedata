package nail

import (
	"github.com/antihax/goesi/esi"

	sq "github.com/Masterminds/squirrel"
	"github.com/antihax/evedata/internal/datapackages"

	"github.com/antihax/evedata/internal/gobcoder"
	"github.com/antihax/evedata/internal/sqlhelper"

	nsq "github.com/nsqio/go-nsq"
)

func init() {
	AddHandler("mutatedItem", func(s *Nail, consumer *nsq.Consumer) {
		consumer.AddConcurrentHandlers(s.wait(nsq.HandlerFunc(s.mutatedItem)), 50)
	})
}

func (s *Nail) mutatedItem(message *nsq.Message) error {
	b := datapackages.ResolveItems{}
	err := gobcoder.GobDecoder(message.Body, &b)
	if err != nil {
		return err
	}

	s.doItemAttributes(b.ItemID, b.Item.DogmaAttributes)
	if err != nil {
		return err
	}

	s.doItemEffects(b.ItemID, b.Item.DogmaEffects)
	if err != nil {
		return err
	}

	return s.doSQL(`INSERT INTO evedata.mutations 
					(typeID, itemID, createdBy, mutatorTypeID, sourceTypeID) 
						VALUES (?,?,?,?,?) ON DUPLICATE KEY UPDATE itemID=itemID;`,
		b.TypeID, b.ItemID, b.Item.CreatedBy, b.Item.MutatorTypeId, b.Item.SourceTypeId,
	)
}

func (s *Nail) doItemAttributes(itemID int64, a []esi.GetDogmaDynamicItemsTypeIdItemIdDogmaAttribute) error {
	sql := sq.Insert("evedata.mutationAttributes").Columns("itemID", "attributeID", "value")
	for _, c := range a {
		sql = sql.Values(itemID, c.AttributeId, c.Value)
	}
	sqlq, args, err := sql.ToSql()
	if err != nil {
		return err
	}
	err = s.doSQL(sqlq+` ON DUPLICATE KEY UPDATE itemID=itemID`, args...)
	if err != nil {
		return err
	}
	return nil
}

func (s *Nail) doItemEffects(itemID int64, a []esi.GetDogmaDynamicItemsTypeIdItemIdDogmaEffect) error {
	sql := sq.Insert("evedata.mutationEffects").Columns("itemID", "effectID", "isDefault")
	for _, c := range a {
		sql = sql.Values(itemID, c.EffectId, sqlhelper.IToB(c.IsDefault))
	}
	sqlq, args, err := sql.ToSql()
	if err != nil {
		return err
	}
	err = s.doSQL(sqlq+` ON DUPLICATE KEY UPDATE itemID=itemID`, args...)
	if err != nil {
		return err
	}
	return nil
}
