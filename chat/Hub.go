package chat

type Hub struct {
	Rooms map[string]*Room
}

func NewHub() *Hub {
	return &Hub{
		Rooms: make(map[string]*Room),
	}
}

func (h *Hub) RegisterRoom(room *Room) {
	h.Rooms[room.name] = room
}

func (h *Hub) RemoveRoom(room *Room) {
	delete(h.Rooms, room.name)
}
