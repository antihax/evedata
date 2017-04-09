package appContext

import "testing"

func TestAppContext(t *testing.T) {
	app := NewTestAppContext()
	r := app.Cache.Get()
	_, err := r.Do("INFO")
	if err != nil {
		t.Error(err)
		return
	}
}
