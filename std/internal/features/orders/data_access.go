package orders

type DataAccess struct {
}

func NewDataAccess() *DataAccess {
	return &DataAccess{}
}

func (db *DataAccess) InsertOrder(order Order) error {
	return nil
}
