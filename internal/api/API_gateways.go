package api

type AgifyGateway interface {
	GetAge(name string) (int, error)
	//GetAge(name string, ch chan int)
}

type GenderizeGateway interface {
	GetGender(name string) (string, error)
}

type NationalizeGateway interface {
	GetNationality(surname string) (string, error)
}
