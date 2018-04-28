package mailserver

/*
var mailserver *MailServer

func TestMain(m *testing.M) {
	redis := redigohelper.ConnectRedisTestPool()
	defer redis.Close()
	mailserver, err := NewMailServer(redis, "notreal", "reallynotreal")
	if err != nil {
		log.Fatalln(err)
	}

	go func() {
		if err := mailserver.Run(); err != nil {
			log.Fatalln(err)
		}
	}()
	retCode := m.Run()
	time.Sleep(time.Second * 5)
	mailserver.Close()

	time.Sleep(time.Second)

	os.Exit(retCode)
}
*/
