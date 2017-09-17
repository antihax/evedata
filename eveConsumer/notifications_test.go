package eveConsumer

import "testing"

func TestNotificationTrigger(t *testing.T) {
	_, err := notificationsTrigger(eC)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestNotificationConsumer(t *testing.T) {
	r := ctx.Cache.Get()
	defer r.Close()
	for {
		work, err := notificationsConsumer(eC, &r)
		if err != nil {
			t.Error(err)
			return
		}
		if work == false {
			break
		}
	}
}
