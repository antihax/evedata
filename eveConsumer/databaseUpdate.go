package eveConsumer

import "log"

func (c *EVEConsumer) updateDatabase() {
	_, err := c.ctx.Db.Exec(`call updateMarket;`)
	if err != nil {
		log.Printf("EVEConsumer: Failed updateMarket: %v", err)
		return
	}
}
