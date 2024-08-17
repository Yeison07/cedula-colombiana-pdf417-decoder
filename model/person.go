package model

import (
	"errors"
	"fmt"
	"time"
)

type Location struct {
	Department       string `json:"department"`
	DepartmentCode   string `json:"departmentCode"`
	Municipality     string `json:"municipality"`
	MunicipalityCode string `json:"municipalityCode"`
}

type Person struct {
	DocumentNumber string    `json:"documentNumber"`
	FirstName      string    `json:"firstName"`
	MiddleName     string    `json:"middleName"`
	LastName       string    `json:"lastName"`
	SecondLastName string    `json:"secondLastName"`
	Gender         string    `json:"gender"`
	BloodType      string    `json:"bloodType"`
	Location       Location  `json:"location"`
	Birthdate      time.Time `json:"birthdate"`
}

func NewPerson(documentNumber, lastName, secondLastName, middleName, firstName, gender string, birthdate time.Time, location Location) (*Person, error) {

	if documentNumber == "" {
		return nil, errors.New("document number cannot be empyt, please check your data input")
	}

	if firstName == "" || lastName == "" {
		return nil, errors.New("the firstName and lastName cannot be empyt, please check your data input")
	}

	if gender != "m" && gender != "f" {
		return nil, errors.New("gender must be 'm' (masculine) or 'f' (feminine), please check your data input")
	}

	return &Person{
		DocumentNumber: documentNumber,
		LastName:       lastName,
		SecondLastName: secondLastName,
		MiddleName:     middleName,
		FirstName:      firstName,
		Gender:         gender,
		Birthdate:      birthdate,
		Location:       location,
	}, nil

}

func (p Person) GetBirthDay() int {
	return p.Birthdate.Day()
}

func (p Person) GetBirthMonth() int {
	return int(p.Birthdate.Month())
}

func (p Person) GetBirthYeart() int {
	return p.Birthdate.Year()
}

func (p Person) String() string {
	return fmt.Sprintf("DocumentNumber: %s, FirstName: %s, MiddleName: %s, LastName: %s, SecondLastName: %s, Gender: %s, Birthdate: %s, Location: {Department: %s, DepartmentCode: %s, Municipality: %s, MunicipalityCode: %s}",
		p.DocumentNumber, p.FirstName, p.MiddleName, p.LastName, p.SecondLastName, p.Gender, p.Birthdate.Format("2006-01-02"),
		p.Location.Department, p.Location.DepartmentCode, p.Location.Municipality, p.Location.MunicipalityCode)
}
