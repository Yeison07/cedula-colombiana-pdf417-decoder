package utils

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/rs/zerolog/log"
	"github.com/yeison07/cedula-colombiana-pdf417-decoder/model"
)

const KEYWORD_CEDULA = "pubdsk"

func FormatterScannerInput(input string) *model.Person {
	cleanedInput := strings.ReplaceAll(input, "\x00", " ")
	inputContainsKeyWord := findIndexEndOf(cleanedInput, KEYWORD_CEDULA)
	indexNameStart := 0

	if inputContainsKeyWord == -1 {
		indexNameStart = findNextLetterOcurrence(cleanedInput, 0)
	} else {
		indexNameStart = findNextLetterOcurrence(cleanedInput, inputContainsKeyWord)
	}

	documentNumber, err := getDocumentNumber(cleanedInput, indexNameStart)
	if err != nil {
		log.Error().Err(err).Str("input", cleanedInput).Msg("Error al obtener el numero de documento")
		return nil
	}

	indexNameEnd := findIndexStartOf(cleanedInput, "0m")
	gender := "m"
	if indexNameEnd == -1 {
		indexNameEnd = findIndexStartOf(cleanedInput, "0f")
		gender = "f"
	}

	names := cleanedInput[indexNameStart:indexNameEnd]
	lastName, secondlastname, middleName, firstName := splitName(names)

	birthdayString := cleanedInput[indexNameEnd+2 : indexNameEnd+10]
	birthday, err := time.Parse("20060102", birthdayString)
	if err != nil {
		log.Error().Err(err).Str("input", cleanedInput).Msg("Error al parsear la fecha: " + birthdayString)
		return nil
	}

	departmentCode := cleanedInput[indexNameEnd+10 : indexNameEnd+13]
	departmentCode = strings.TrimLeft(departmentCode, "0")
	municipalityCode := cleanedInput[indexNameEnd+13 : indexNameEnd+16]
	location := model.Location{
		DepartmentCode:   departmentCode,
		MunicipalityCode: municipalityCode,
	}
	if ciudades, departmentExists := model.Data[departmentCode]; departmentExists {
		if values, municipalityExists := ciudades[municipalityCode]; municipalityExists {
			location.Department = values[0]
			location.Municipality = values[1]
		} else {
			log.Error().Msg("La ciudad con el código especificado no fue encontrada.")
		}
	} else {
		log.Error().Msg("El departamento con el código especificado no fue encontrado.")
	}

	person, err := model.NewPerson(documentNumber, lastName, secondlastname, middleName, firstName, gender, birthday, location)
	if err != nil {
		log.Error().Err(err).Str("input", cleanedInput).Msg("Error al crear la instancia de Person")
		return nil
	}
	return person
}

func findIndexStartOf(str string, substring string) int {
	return strings.Index(strings.ToLower(str), substring)
}

func findIndexEndOf(str string, substring string) int {
	index := strings.Index(strings.ToLower(str), strings.ToLower(substring))
	if index == -1 {
		return -1
	}
	return index + len(substring)
}

func findNextLetterOcurrence(str string, start int) int {
	for i := start; i < len(str); i++ {
		if unicode.IsLetter(rune(str[i])) && unicode.IsLetter(rune(str[i+1])) {
			return i
		}

	}
	return -1
}

/*
Separa los nombres y verifica a cada tipo de nombre corresponde
*/
func splitName(fullname string) (lastname, secondlastname, middlename, firstname string) {
	parts := strings.Fields(fullname)

	switch len(parts) {
	case 4:
		lastname = parts[0]
		secondlastname = parts[1]
		middlename = parts[3]
		firstname = parts[2]
	case 3:
		lastname = parts[0]
		secondlastname = parts[1]
		firstname = parts[2]
	case 2:
		lastname = parts[0]
		firstname = parts[1]
	default:
		log.Error().Msgf("La cadena no contiene una cantidad de nombres correcta. Cantidad: %d", len(parts))
	}
	return
}

/*
Recibe la cadena original y a partir del indice donde comienza los nombres extrae el numero de documento.
Limpia los ceros (0) a la izquierda que pudiesen contener los documentos menores a 10 digitos y retorna el numero de documento formateado.
*/
func getDocumentNumber(input string, index int) (string, error) {
	if index < 10 {
		return "", fmt.Errorf("el índice proporcionado no es suficientemente grande para obtener los primeros diez dígitos hacia atrás")
	}

	start := index - 1
	digitsCount := 0
	var digits strings.Builder
	for start >= 0 && digitsCount < 10 {
		char := input[start]
		if unicode.IsDigit(rune(char)) {
			digits.WriteByte(char)
			digitsCount++
		} else if digitsCount > 0 {
			break
		}
		start--
	}
	digitsString := digits.String()
	reversedDigits := reverseString(digitsString)
	trimmed := strings.TrimLeft(reversedDigits, "0")

	return trimmed, nil
}

func reverseString(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
