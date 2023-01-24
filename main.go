package main

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"os/exec"
	"strconv"
)

type Score struct {
	Name    string
	Perfect int
	Great   int
	Good    int
	Bad     int
	Miss    int
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

var (
	openimg, name                   string
	bbb                             []string
	perfect, great, good, bad, miss int
	atoierr                         error
)

func main() {
	fmt.Print("画像ファイル名を入力してください: ")
	fmt.Scanln(&openimg)
	fmt.Print("楽曲名を入力してください: ")
	fmt.Scanln(&name)

	f, err := os.Open("./img/" + openimg)
	if err != nil {
		fmt.Println("open err:", err)
		return
	}
	defer f.Close()

	img1, _, err := image.Decode(f)
	if err != nil {
		fmt.Println("decode err:", err)
		return
	}

	fso, err := os.Create("cut.jpg")
	if err != nil {
		fmt.Println("create err:", err)
		return
	}
	defer fso.Close()

	cimg := img1.(SubImager).SubImage(image.Rect(950, 570, 1075, 820))

	jpeg.Encode(fso, cimg, &jpeg.Options{Quality: 100})

	// 数字取得
	imgname := "cut.jpg"
	out, err := exec.Command("tesseract", imgname, "-").Output()
	if err != nil {
		fmt.Println("command err", err)
		return
	}
	os.WriteFile("score.txt", out, 0644)

	fp, err := os.Open("score.txt")
	if err != nil {
		fmt.Println("open err", err)
		return
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		bbb = append(bbb, scanner.Text())
	}
	fmt.Println(bbb)

	perfect, atoierr = strconv.Atoi(bbb[0])
	great, atoierr = strconv.Atoi(bbb[1])
	good, atoierr = strconv.Atoi(bbb[2])
	bad, atoierr = strconv.Atoi(bbb[3])
	miss, atoierr = strconv.Atoi(bbb[4])

	if atoierr != nil {
		fmt.Println(err)
	}

	sc := Score{
		Name:    name,
		Perfect: perfect,
		Great:   great,
		Good:    good,
		Bad:     bad,
		Miss:    miss,
	}

	fmt.Printf("楽曲名: %s, Perfect: %d, Great: %d, Good: %d, Bad: %d, Miss: %d", sc.Name, sc.Perfect, sc.Great, sc.Good, sc.Bad, sc.Miss)
}
