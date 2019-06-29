package game

import (
	"fmt"
	"sync"
	"time"
)

type QuestLog struct {
	items []QuestLogItem

	sync.RWMutex
}

func (q *QuestLog) Len() int {
	q.RLock()
	defer q.RUnlock()
	return len(q.items)
}

func (q *QuestLog) AddItem(i QuestLogItem) {
	q.Lock()
	defer q.Unlock()
	q.items = append(q.items, i)
}

func (q *QuestLog) GetItem(i int) (QuestLogItem, error) {
	q.RLock()
	defer q.RUnlock()

	var item QuestLogItem
	if i < 0 || i >= len(q.items) {
		return item, fmt.Errorf("bad index (%d) must be > 0 and < %d",
			i, len(q.items))
	}

	return q.items[i], nil
}

func (q *QuestLog) RecordPlayerMovement(x, y, oldX, oldY int) {
	i := QuestLogItem{
		time: time.Now(),
		icon: '*',
		msg: fmt.Sprintf(
			"You move from (%3d,%3d) to (%3d,%3d)",
			oldX, oldY, x, y),
	}
	q.AddItem(i)
}

type QuestLogItem struct {
	time time.Time
	icon rune
	msg  string
}

func (i QuestLogItem) String() string {
	dateStr := i.time.Format("15:04:05")
	return fmt.Sprintf("%s %-6s - %s", string(i.icon), dateStr, i.msg)
}
