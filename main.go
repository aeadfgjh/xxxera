package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/spf13/afero"
)

type Image struct {
	Name string
	Path string
}

type PartialImage struct {
	Image
	ax, ay, bx, by uint
}

type CSVParseWarning struct {
	data string
	path string
}

func (w CSVParseWarning) Error() string {
	return fmt.Sprintf("invalid CSV record in %s: %s", w.path, w.data)
}

// loadImageCSV loads and returns image definitions from IO buffer.
// CSV record must be in either form of:
//  1. image definition
//     `name,image path`
//  2. partial image definition
//     `name,image path,ax,ay,bx,by`
//     * ax, ay, bx, by must be unsigned integers
//  3. comment (ignored)
//     `;This is comment`
//
// lines not statisfying any of these rule is considered as invalid and
// warnings are generated. Third return value is list of warnings. The list
// return value is critical error such as opening file, invalid syntax. etc.
func loadImageCSV(r io.Reader) ([]Image, []PartialImage, []error, error) {
	csvReader := csv.NewReader(r)
	csvReader.Comment = ';'
	csvReader.FieldsPerRecord = -1
	csvReader.ReuseRecord = true

	var images []Image
	var partialImages []PartialImage
	var warnings []error

	for {
		record, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return nil, nil, nil, err
		}

		switch len(record) {
		case 2:
			images = append(images, Image{Name: record[0], Path: record[1]})
		case 6:
			coords, err := parseUintN(record[2], record[3], record[4], record[5])
			if err != nil {
				warnings = append(warnings, fmt.Errorf("failed to parse image coordinate: %w", err))
				continue
			}
			partialImages = append(partialImages, PartialImage{
				Image: Image{Name: record[0], Path: record[1]},
				ax:    coords[0], ay: coords[1], bx: coords[2], by: coords[3],
			})
		default:
			warnings = append(warnings, fmt.Errorf("invalid number of image CSV field: %v", record))
		}
	}

	return images, partialImages, warnings, nil
}

func parseUintN(args ...string) ([]uint, error) {
	results := make([]uint, 0, len(args))
	for _, arg := range args {
		n, err := strconv.ParseUint(arg, 10, 64)
		if err != nil {
			return nil, err
		}
		results = append(results, uint(n))
	}
	return results, nil
}

func main() {
	fs := afero.NewOsFs()

	csvPath := "./test/list.csv"
	file, err := fs.OpenFile(csvPath, os.O_RDONLY, 0440)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	images, partialImages, warnings, err := loadImageCSV(file)
	if err != nil {
		log.Fatal(err)
	}

	if len(warnings) > 0 {
		log.Printf("warnings for image CSV file %s\n", csvPath)
		for _, warning := range warnings {
			log.Println(warning)
		}
	}

	fmt.Println("Images:")
	for _, image := range images {
		fmt.Println(image)
	}

	fmt.Println("Partial images:")
	for _, image := range partialImages {
		fmt.Println(image)
	}
}
