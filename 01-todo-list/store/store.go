package store

type Store interface {
	Add(task string)
	Delete(id string)
	Update(id string, task string, check bool)
	List() [][]string
	Flush()
}
