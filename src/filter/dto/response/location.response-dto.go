package response

type LocationResponseDto struct {
	Id        string `copier:"UUID" json:"id" xml:"id" example:"uuid"`
	Continent string `copier:"ContinentName" json:"continent" xml:"continent" example:"Asia"`
	Country   string `copier:"CountryName" json:"country" xml:"country" example:"Qatar"`
	ISO       string `copier:"CountryIsoCode" json:"iso" xml:"iso" example:"QA"`
}
