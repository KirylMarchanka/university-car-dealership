package car

import (
	"car_dealership/internal/manufacturer"
	"errors"
)

type Car struct {
	Id             int64                     `json:"id"`
	ManufacturerId int64                     `json:"manufacturer_id"`
	Name           string                    `json:"name"`
	Fuel           string                    `json:"fuel"`
	FuelCapacity   float32                   `json:"fuel_capacity"`
	Engine         string                    `json:"engine"`
	EnginePower    float32                   `json:"engine_power"`
	EngineCapacity float32                   `json:"engine_capacity"`
	MaxSpeed       int32                     `json:"max_speed"`
	Acceleration   float32                   `json:"acceleration"`
	Manufacturer   manufacturer.Manufacturer `json:"manufacturer"`
}

func Validate(
	manufacturerId int64,
	name string,
	fuel string,
	fuelCapacity float32,
	engine string,
	enginePower float32,
	engineCapacity float32,
	maxSpeed int32,
	acceleration float32,
) error {
	// Manufacturer ID validation (Check if it exists in the database).
	if !manufacturer.ExistsById(manufacturerId) {
		return errors.New("manufacturer does not exist")
	}

	// Name, Fuel, and Engine validation (between 2 and 255 characters).
	if len(name) < 2 || len(name) > 255 {
		return errors.New("name must be between 2 and 255 characters")
	}

	if len(fuel) < 2 || len(fuel) > 255 {
		return errors.New("fuel must be between 2 and 255 characters")
	}

	if len(engine) < 2 || len(engine) > 255 {
		return errors.New("engine must be between 2 and 255 characters")
	}

	if fuelCapacity <= 0 {
		return errors.New("fuel capacity must be greater than 0")
	}

	if enginePower <= 0 {
		return errors.New("engine power must be greater than 0")
	}

	if engineCapacity <= 0 {
		return errors.New("engine capacity must be greater than 0")
	}

	// MaxSpeed validation (between 0 and 400 kph).
	if maxSpeed < 0 || maxSpeed > 400 {
		return errors.New("max speed must be between 0 and 400 kph")
	}

	// Acceleration validation (between 0 and 60 seconds).
	if acceleration < 0 || acceleration > 60 {
		return errors.New("acceleration must be between 0 and 60 seconds")
	}

	return nil
}

func New(
	manufacturerId int64,
	name string,
	fuel string,
	fuelCapacity float32,
	engine string,
	enginePower float32,
	engineCapacity float32,
	maxSpeed int32,
	acceleration float32,
) *Car {
	id := store(manufacturerId, name, fuel, fuelCapacity, engine, enginePower, engineCapacity, maxSpeed, acceleration)
	if id == 0 {
		return nil
	}

	return &Car{
		Id:             id,
		ManufacturerId: manufacturerId,
		Name:           name,
		Fuel:           fuel,
		FuelCapacity:   fuelCapacity,
		Engine:         engine,
		EnginePower:    enginePower,
		EngineCapacity: engineCapacity,
		MaxSpeed:       maxSpeed,
		Acceleration:   acceleration,
		Manufacturer:   manufacturer.GetById(manufacturerId),
	}
}

func Update(car *Car) error {
	err := car.update()
	if err != nil {
		return err
	}

	return nil
}

func Delete(car *Car) error {
	err := car.delete()
	if err != nil {
		return err
	}

	return nil
}

func SelectCars(manufacturerID int64, name, fuel, orderBy, orderDirection string) ([]Car, error) {
	return selectCars(manufacturerID, name, fuel, orderBy, orderDirection)
}

func Find(id int64) (Car, error) {
	return findById(id)
}
