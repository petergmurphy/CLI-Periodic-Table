package elements

import (
	"encoding/csv"
	"log"
	"os"
	"periodic-table/ui/element"
	"periodic-table/ui/grid"
)

var (
	emptyEntryRange = map[int]int{
		1:   17,
		20:  30,
		38:  48,
		126: 148,
		162: 166,
	}
)

func ReadElements() []grid.Cell {
	// open file
	f, err := os.Open("data/elements.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// convert records to array of structs
	return createElements(data)
}

func createElements(data [][]string) []grid.Cell {
	var elements []grid.Cell
	var count int
	var lastGroups []grid.Cell
	for i, line := range data {
		if skipLines := emptyEntryRange[count]; skipLines != 0 {
			for j := count; j < skipLines; j++ {
				elements = append(elements, element.CreateElement(element.Data{}, true))
				count++
			}
		}

		if i > 0 {
			data := createElementData(line)

			elements = append(elements, element.CreateElement(data, false))
			count++
		}
	}

	elements = append(elements, lastGroups...)

	return elements
}

func createElementData(line []string) element.Data {
	return element.Data{
		AtomicNumber:      line[0],
		Element:           line[1],
		Symbol:            line[2],
		AtomicMass:        line[3],
		NumberOfNeutrons:  line[4],
		NumberOfProtons:   line[5],
		NumberOfElectrons: line[6],
		Period:            line[7],
		Group:             line[8],
		Phase:             line[9],
		Radioactive:       line[10],
		Natural:           line[11],
		Metal:             line[12],
		Nonmetal:          line[13],
		Metalloid:         line[14],
		Type:              line[15],
		AtomicRadius:      line[16],
		Electronegativity: line[17],
		FirstIonization:   line[18],
		Density:           line[19],
		MeltingPoint:      line[20],
		BoilingPoint:      line[21],
		NumberOfIsotopes:  line[22],
		Discoverer:        line[23],
		Year:              line[24],
		SpecificHeat:      line[25],
		NumberOfShells:    line[26],
		NumberOfValence:   line[27],
	}
}
