package helper

import "fmt"

func IntToString(i int) string {
	return fmt.Sprintf("%d", i)
}

func StringToInt(s string) (int, error) {
	var i int
	_, err := fmt.Sscanf(s, "%d", &i)
	if err != nil {
		return 0, fmt.Errorf("failed to convert string to int: %w", err)
	}
	return i, nil
}
