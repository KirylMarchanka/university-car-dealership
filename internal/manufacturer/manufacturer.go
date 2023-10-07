package manufacturer

import "errors"

type Manufacturer struct {
	Id   int64
	Name string
}

func Validate(name string) error {
	if len(name) <= 2 || len(name) > 255 {
		return errors.New("incorrect manufacturer name length")
	}

	if existsByName(name) {
		return errors.New("non unique manufacturer name")
	}

	return nil
}

func New(name string) *Manufacturer {
	id := store(name)
	if id == 0 {
		return nil
	}

	return &Manufacturer{
		Id:   id,
		Name: name,
	}
}

func GetById(id int64) Manufacturer {
	id, name, err := getById(id)
	if err != nil {
		return Manufacturer{
			Id:   0,
			Name: "",
		}
	}

	return Manufacturer{
		Id:   id,
		Name: name,
	}
}
