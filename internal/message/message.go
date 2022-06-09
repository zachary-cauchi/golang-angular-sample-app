package message

import (
	"errors"
	"sync"

	"github.com/rs/xid"
)

var (
	list []Message

	mtx sync.RWMutex

	once sync.Once
)

type Message struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

func init() {
	// Initialises the list ONLY ONCE, even if somehow called again.
	once.Do(initialiseList)
}

func initialiseList() {
	list = []Message{}
}

func Get() []Message {
	return list
}

// Add a message to the list with the given text.
func Add(text string) string {
	t := newMessage(text)

	mtx.Lock()
	list = append(list, t)
	mtx.Unlock()

	return t.ID
}

// Delete a message with the given id.
func Delete(id string) error {
	location, err := findMessageLocation(id)

	if err != nil {
		return err
	}

	removeMessageByLocationOrdered(location)

	return nil
}

func newMessage(txt string) Message {
	return Message{
		ID:   xid.New().String(),
		Text: txt,
	}
}

// Finds a message by the given id.
func findMessageLocation(id string) (int, error) {
	mtx.RLock()

	defer mtx.RUnlock()

	for i, t := range list {
		if isMatchingID(t.ID, id) {
			return i, nil
		}
	}

	return 0, errors.New("Could not find message based on id")
}

// Removes a Message object at the specified index.
// This maintains the list order at the cost of speed.
func removeMessageByLocationOrdered(i int) {
	mtx.Lock()
	list = append(list[:i], list[i+i:]...)
	mtx.Unlock()
}

// Removes a Message object at the specified index.
// This does not maintain the list order but is faster as a result.
func removeMessageByLocationUnordered(i int) {
	mtx.Lock()
	// Copy the last element to the element to be deleted.
	list[i] = list[len(list)-1]
	list[len(list)-1] = Message{}
	// Resize the list, deleting the last element.
	list = list[:len(list)-1]
	mtx.Unlock()
}

func isMatchingID(a string, b string) bool {
	return a == b
}
