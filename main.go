package main

import (
	"bufio"
	"fmt"
	"os"

	"os/exec"
)

type Score struct {
	Perfect string
	Great   string
	Good    string
	Bad     string
	Miss    string
}

func main() {
	var img string
	fmt.Scan(&img)
	imgname := img
	out, err := exec.Command("tesseract", imgname, "-", "-l", "eng").Output()
	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile("test.txt", out, 0644)

	fp, err := os.Open("test.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	var bbb []string
	for scanner.Scan() {
		fmt.Println(scanner.Text())
		bbb = append(bbb, scanner.Text())
	}
	fmt.Println(bbb)
	sc := Score{
		Perfect: bbb[0],
		Great:   bbb[1],
		Good:    bbb[2],
		Bad:     bbb[3],
		Miss:    bbb[4],
	}

	fmt.Printf("Perfect: %s, Great: %s, Good: %s, Bad: %s, Miss: %s", sc.Perfect, sc.Great, sc.Good, sc.Bad, sc.Miss)
}
