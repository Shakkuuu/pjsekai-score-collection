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

// Score struct
type Score struct {
	Name    string
	Perfect string
	Great   string
	Good    string
	Bad     string
	Miss    string
}

// 画像加工用のインターフェイス
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

var (
	openimg, name, sperfect, sgreat, sgood, sbad, smiss string   // s〇〇はstringのスコア
	scanls                                              []string // tesseractで読み込んだ値の配列
	iperfect, igreat, igood, ibad, imiss                int      // i〇〇はintのスコア
	atoierr                                             error
)

func main() {
	fmt.Print("画像ファイル名を入力してください: ")
	fmt.Scanln(&openimg)
	fmt.Print("楽曲名を入力してください: ")
	namebufio := bufio.NewScanner(os.Stdin)
	namebufio.Scan()
	name = namebufio.Text()

	// 指定された画像のオープン
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

	cimg := img1.(SubImager).SubImage(image.Rect(950, 570, 1075, 820)) // image.Rect(x0,y0,x1,y1) ここを変更して端末ごとのスクリーンショットに合わせる

	jpeg.Encode(fso, cimg, &jpeg.Options{Quality: 100})

	// 数字取得
	imgname := "cut.jpg"
	out, err := exec.Command("tesseract", imgname, "-").Output() // execでtesseract実行
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

	// score.txtに出力されたスコアを取り出してリスト化
	scanner := bufio.NewScanner(fp)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		scanls = append(scanls, scanner.Text())
	}
	fmt.Println(scanls)

	// 頭の 0 削除
	iperfect, atoierr = strconv.Atoi(scanls[0])
	igreat, atoierr = strconv.Atoi(scanls[1])
	igood, atoierr = strconv.Atoi(scanls[2])
	ibad, atoierr = strconv.Atoi(scanls[3])
	imiss, atoierr = strconv.Atoi(scanls[4])
	if atoierr != nil {
		fmt.Println(err)
	}

	// stringに戻す
	sperfect = strconv.Itoa(iperfect)
	sgreat = strconv.Itoa(igreat)
	sgood = strconv.Itoa(igood)
	sbad = strconv.Itoa(ibad)
	smiss = strconv.Itoa(imiss)

	// 構造体に埋め込み
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

	// scorelist.txtに保存。ファイルがない場合は作成
	addf, err := os.OpenFile("scorelist.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer addf.Close()
	fmt.Fprintln(addf, msg)
}
