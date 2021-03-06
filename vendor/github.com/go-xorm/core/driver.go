package core

type Driver interface {
	Parse(string, string) (*Uri, error)
}

var (
	drivers = map[string]Driver{}
)

func RegisterDriver(driverName string, driver Driver) {
	if driver == nil {
		panic("core: Register driver is nil")
	}
	if _, dup := drivers[driverName]; dup {
		panic("core: Register called twice for driver " + driverName)
	}
	drivers[driverName] = driver
}

func QueryDriver(driverName string) Driver {
	return drivers[driverName]
}

func RegisteredDriverSize() int {
	return len(drivers)
}
