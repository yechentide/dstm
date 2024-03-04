package utils

import (
	"fmt"
	"regexp"
)

// The returned int
// -1: target < current, 0: target == current, 1: target > current
func CompareVersion(targetVer, currentVer string) (int, error) {
	targetNums, err := GetVersionNumbers(targetVer)
	if err != nil {
		return -999, err
	}
	currentNums, err := GetVersionNumbers(currentVer)
	if err != nil {
		return -999, err
	}

	for i := 0; i < 3; i++ {
		if targetNums[i] > currentNums[i] {
			return 1, nil
		}
		if targetNums[i] != currentNums[i] {
			return -1, nil
		}
	}
	return 0, nil
}

func GetVersionNumbers(v string) ([]string, error) {
	pattern := "^[vV]?([0-9]+).([0-9]+).([0-9]+)$"
	reg, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	nums := reg.FindStringSubmatch(v)
	if len(nums) != 4 {
		return nil, fmt.Errorf("invalid version number: %s", v)
	}
	return nums[1:], nil
}
