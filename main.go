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
	Perfect string
	Great   string
	Good    string
	Bad     string
	Miss    string
}

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

var (
	openimg, name, sperfect, sgreat, sgood, sbad, smiss string
	bbb                                                 []string
	iperfect, igreat, igood, ibad, imiss                int
	atoierr                                             error
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

	// 頭の 0 削除
	iperfect, atoierr = strconv.Atoi(bbb[0])
	igreat, atoierr = strconv.Atoi(bbb[1])
	igood, atoierr = strconv.Atoi(bbb[2])
	ibad, atoierr = strconv.Atoi(bbb[3])
	imiss, atoierr = strconv.Atoi(bbb[4])
	if atoierr != nil {
		fmt.Println(err)
	}

	// stringに戻す
	sperfect = strconv.Itoa(iperfect)
	sgreat = strconv.Itoa(igreat)
	sgood = strconv.Itoa(igood)
	sbad = strconv.Itoa(ibad)
	smiss = strconv.Itoa(imiss)

	sc := Score{
		Name:    name,
		Perfect: sperfect,
		Great:   sgreat,
		Good:    sgood,
		Bad:     sbad,
		Miss:    smiss,
	}

	msg := "楽曲名: " + sc.Name + ", Perfect: " + sc.Perfect + ", Great: " + sc.Great + ", Good: " + sc.Good + ", Bad: " + sc.Bad + ", Miss: " + sc.Miss
	// fmt.Printf("楽曲名: %s, Perfect: %d, Great: %d, Good: %d, Bad: %d, Miss: %d", sc.Name, sc.Perfect, sc.Great, sc.Good, sc.Bad, sc.Miss)
	fmt.Println(msg)

	// addf, err := os.OpenFile("scorelist.txt", os.O_APPEND|os.O_WRONLY, 0644)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Fprintln(addf, "aa")
}
