package artifice

import "log"

func (s *Artifice) startup() error {
	err := s.loadKills()
	if err != nil {
		return err
	}
	err = s.loadFinishedWars()
	if err != nil {
		return err
	}
	return nil
}

func (s *Artifice) loadKills() error {
	var known []int64
	if err := s.db.Select(&known, `SELECT id FROM evedata.killmails WHERE hash != "";`); err != nil {
		return err
	}

	log.Printf("Known Kills: %d\n", len(known))
	err := s.inQueue.SetWorkCompletedInBulk("evedata_known_kills", known)
	if err != nil {
		return err
	}

	return nil
}

func (s *Artifice) loadFinishedWars() error {
	var known []int64
	if err := s.db.Select(&known, `SELECT id FROM evedata.wars WHERE timeFinished < UTC_TIMESTAMP();`); err != nil {
		return err
	}

	log.Printf("Finished Wars: %d\n", len(known))
	err := s.inQueue.SetWorkCompletedInBulk("evedata_war_finished", known)
	if err != nil {
		return err
	}

	return nil
}
