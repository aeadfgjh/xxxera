package main

import (
	"fmt"
	"strings"
	"testing"
)

func imageToCSVRecord(i Image) string {
	return fmt.Sprintf("%s,%s\n", i.Name, i.Path)
}

func partialImageToCSVRecord(i PartialImage) string {
	return fmt.Sprintf("%s,%s,%d,%d,%d,%d\n", i.Name, i.Path, i.ax, i.ay, i.bx, i.by)
}

func TestLoadImageCSVGood(t *testing.T) {
	expectedImage := Image{Name: "test image", Path: "test image.webp"}
	expectedPartialImage := PartialImage{
		Image: Image{Name: "test image partial", Path: "test image partial.webp"},
		ax:    1, ay: 2, bx: 3, by: 4,
	}

	input := strings.Join(
		[]string{imageToCSVRecord(expectedImage), partialImageToCSVRecord(expectedPartialImage)},
		"\n",
	)
	images, partialImages, warnings, err := loadImageCSV(strings.NewReader(input))
	if err != nil {
		t.Fatalf("loadImageCSV returned unexpected error: %s", err.Error())
	} else if len(warnings) != 0 {
		t.Fatalf("loadImageCSV returned unexpected warnings: %s", warnings)
	}

	if len(images) != 1 || len(partialImages) != 1 {
		t.Fatalf("loadImageCSV returned unexpected number of records")
	} else if expectedImage != images[0] || expectedPartialImage != partialImages[0] {
		t.Fatalf("loadImageCSV returned unexpected record values")
	}
}
