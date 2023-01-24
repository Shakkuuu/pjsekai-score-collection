package main

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
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

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func main() {
	var openimg string
	fmt.Scan(&openimg)

	f, err := os.Open("./img/" + openimg)
	if err != nil {
		fmt.Println("open:", err)
		return
	}
	defer f.Close()

	img1, _, err := image.Decode(f)
	if err != nil {
		fmt.Println("decode:", err)
		return
	}

	fso, err := os.Create("cut.jpg")
	if err != nil {
		fmt.Println("create:", err)
		return
	}
	defer fso.Close()

	cimg := img1.(SubImager).SubImage(image.Rect(950, 570, 1075, 820))

	jpeg.Encode(fso, cimg, &jpeg.Options{Quality: 100})

	// 数字取得
	// var img2 string
	// fmt.Scan(&img2)
	// imgname := img2
	imgname := "cut.jpg"
	out, err := exec.Command("tesseract", imgname, "-", "-l", "eng").Output()
	if err != nil {
		fmt.Println(err)
	}
	os.WriteFile("score.txt", out, 0644)

	fp, err := os.Open("score.txt")
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
