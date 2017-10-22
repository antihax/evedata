package redisqueue

import (
	"testing"

	"github.com/antihax/evedata/internal/redigohelper"
	"github.com/stretchr/testify/assert"
)

func TestHQ(t *testing.T) {
	pool := redigohelper.ConnectRedisTestPool()
	hq := NewRedisQueue(pool, "test-redisqueue")
	err := hq.QueueWork(
		[]Work{
			Work{"alliance",
				2,
			},
			Work{"alliance",
				3,
			},
			Work{"alliance",
				4,
			},
			Work{"alliance",
				5,
			},
			Work{"alliance",
				6,
			},
		},
	)
	assert.Nil(t, err)
	size, err := hq.Size()
	assert.Nil(t, err)
	assert.Equal(t, 5, size)

	var work []*Work

	for i := 0; i < 5; i++ {
		w, err := hq.GetWork()
		assert.Nil(t, err)
		work = append(work, w)
	}

	check := map[int]bool{2: true, 3: true, 4: true, 5: true, 6: true}

	for i := range work {
		delete(check, work[i].Parameter.(int))
	}
	assert.Empty(t, check)
}

func TestFailure(t *testing.T) {
	pool := redigohelper.ConnectRedisTestPool()
	hq := NewRedisQueue(pool, "test-redisqueue")
	err := hq.SetWorkFailure("testKey", 1)
	assert.Nil(t, err)

	b := hq.CheckWorkFailure("testKey", 1)
	assert.Equal(t, b, true)
}

func TestCompletion(t *testing.T) {
	pool := redigohelper.ConnectRedisTestPool()
	hq := NewRedisQueue(pool, "test-redisqueue")
	err := hq.SetWorkCompleted("testKeyComp", 1)
	assert.Nil(t, err)

	b := hq.CheckWorkCompleted("testKeyComp", 1)
	assert.Equal(t, b, true)
}
