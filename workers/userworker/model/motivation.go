package model

type Motivation struct {
	ID         int
	Motivation string
}

func NewMotivationModel(
	id int,
	motivation string,
) Motivation {
	return Motivation{
		ID:         id,
		Motivation: motivation,
	}
}
