package kdatabase



type KDBSet struct {
	rows	[][]interface{}
}

func NewKDBSet() *KDBSet {

	return &KDBSet{
		rows:	make([][]interface{}, 0, 128),
	}
}

type KDBResult struct {
	sets	[]*KDBSet
}

func NewKDBResult() *KDBResult {

	return &KDBResult{
		sets:	make([]*KDBSet, 0, 8),
	}
}