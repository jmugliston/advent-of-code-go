package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/atheius/aoc/parsing"
)

var partFlag = flag.String("part", "1", "The part of the day to run (1 or 2)")

func main() {
	flag.Parse()

	_, filename, _, _ := runtime.Caller(0)
	dirname := filepath.Dir(filename)

	path := filepath.Join(dirname, "input", "input.txt")

	input, err := os.ReadFile(path)

	if err != nil {
		panic("Could not find the input file")
	}

	if *partFlag == "1" {
		fmt.Println(Part1(string(input)))
	} else {
		fmt.Println(Part2(string(input)))
	}
}

func decode(digits []int) []int {
	var decoded []int
	for i := 0; i < len(digits); i++ {
		if i%2 == 0 {
			for j := 0; j < digits[i]; j++ {
				decoded = append(decoded, i/2)
			}
		} else {
			for j := 0; j < digits[i]; j++ {
				decoded = append(decoded, -1)
			}
		}
	}
	return decoded
}

func calculateChecksum(defragged []int) int {
	checksum := 0
	for i := 0; i < len(defragged); i++ {
		if defragged[i] == -1 {
			continue
		}
		checksum += defragged[i] * i
	}
	return checksum
}

func Part1(input string) int {
	digits := parsing.ReadDigits(input)

	decoded := decode(digits)

	var defragged = make([]int, len(decoded))
	copy(defragged, decoded)

	left := 0
	right := len(decoded) - 1
	for left < right {
		if defragged[left] >= 0 {
			left += 1
			continue
		}
		if defragged[right] < 0 {
			right -= 1
		}
		defragged[left] = defragged[right]
		defragged[right] = -1
	}

	checksum := calculateChecksum(defragged)

	return checksum
}

func Part2(input string) int {
	digits := parsing.ReadDigits(input)

	decoded := decode(digits)

	var defragged = make([]int, len(decoded))
	copy(defragged, decoded)

	checkedFileIds := make(map[int]bool)
	for {
		var fileId int
		var fileIdx int
		fileSize := 0

		// Scan from the right and find the next file id
		for i := len(defragged) - 1; i >= 0; i-- {
			if defragged[i] == -1 {
				continue
			}
			if checkedFileIds[defragged[i]] {
				continue
			}
			fileId = defragged[i]
			for j := i; j >= 0; j-- {
				if defragged[j] != fileId {
					break
				}
				fileIdx = j
				fileSize += 1
			}
			break
		}

		// Scan from the left and look for a large enough gap
		for i := 0; i < len(defragged); i++ {
			if defragged[i] != -1 {
				continue
			}
			gapSize := 0
			for j := i; j < fileIdx; j++ {
				if defragged[j] != -1 {
					break
				}
				gapSize += 1
			}
			if gapSize >= fileSize {
				for j := i; j < i+fileSize; j++ {
					// Add the file id to the gap
					defragged[j] = fileId
					// Remove the file id from the original location
					defragged[fileIdx+(j-i)] = -1
				}
				break
			}
		}

		checkedFileIds[fileId] = true

		if fileId == 0 {
			break
		}
	}

	checksum := calculateChecksum(defragged)

	return checksum
}
