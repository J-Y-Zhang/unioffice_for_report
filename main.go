package main

import (
	"bufio"
	"fmt"
	"github.com/unidoc/unioffice/common"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/unidoc/unioffice/common/license"
)

var (
	doc *document.Document
)

func init() {
	err := license.SetMeteredKey(key)
	if err != nil {
		log.Fatalln(err)
	}

	doc = document.New()
}

func main() {
	defer func() {
		doc.SaveToFile("./生成结果/" + 姓名 + 学号 + ".docx")
		doc.Close()
	}()

	// 1.加目录并换页
	生成自动更新的目录并换页()

	// 2.加页码
	加页码()

	// 3.姓名学号等基本信息
	生成居中加粗的一段文字(fmt.Sprintf("姓名: %v 学号: %v 班级: %v", 姓名, 学号, 班级), 15)
	fmt.Println("生成姓名学号等成功")

	// 4.实验名称
	生成居中加粗的一段文字(实验名称, -1, "Title")
	fmt.Println("生成实验名称成功")

	// 5.目标及要求
	生成标题(1, "1", "实验目标及要求")
	从文件中生成开头不自动空格的一段文字("./文字内容/目标及要求.txt")
	fmt.Println("生成实验目标及要求生成")

	// 6.主要内容
	生成标题(1, "2", "实验主要内容")
	从文件中生成开头不自动空格的一段文字("./文字内容/主要内容.txt")
	fmt.Println("生成实验主要内容成功")

	// 7.1.实验代码
	生成标题(1, "3", "实验代码及运行结果")
	生成标题(4, "3.1", "实验代码")
	fmt.Println("生成实验代码成功")

	// 7.2.运行结果
	从文件中生成代码("./文字内容/实验代码.txt")
	生成标题(4, "3.2", "运行结果")
	rd, err := ioutil.ReadDir("./图片内容")
	if err != nil {
		log.Fatalln("打开图片内容目录错误")
	}

	i := 1
	for _, imgFile := range rd {
		if !imgFile.IsDir() {
			imgPath := "./图片内容/" + imgFile.Name()
			生成图片(imgPath)
			生成居中加粗的一段文字("图"+strconv.Itoa(i)+": "+imgFile.Name(), 9)
		}
		i++
	}
	fmt.Println("生成实验结果图已成功, 请查收报告")

}

func 生成自动更新的目录并换页() {
	doc.Settings.SetUpdateFieldsOnOpen(true)
	doc.AddParagraph().AddRun().AddField(document.FieldTOC)
	doc.AddParagraph().Properties().AddSection(wml.ST_SectionMarkNextPage)

	fmt.Println("生成目录已成功")
}

func 加页码() {
	ftr := doc.AddFooter()
	para := ftr.AddParagraph()
	para.SetAlignment(wml.ST_JcCenter)

	run := para.AddRun()
	run.Properties().SetFontFamily("Times New Roman")
	run.Properties().SetSize(9)
	run.AddField(document.FieldCurrentPage)
	doc.BodySection().SetFooter(ftr, wml.ST_HdrFtrDefault)

	fmt.Println("生成页码已成功")
}

func 生成居中加粗的一段文字(文字 string, size measurement.Distance, style ...string) {
	para := doc.AddParagraph()
	para.SetAlignment(wml.ST_JcCenter)
	for _, val := range style {
		para.SetStyle(val)
	}

	run := para.AddRun()
	if size > 0 {
		run.Properties().SetSize(size)
	}
	run.Properties().SetFontFamily("宋体")
	run.Properties().SetBold(true)
	run.AddText(文字)
	run.AddBreak()
}

func 生成标题(标题等级 int, 标题前缀, 标题内容 string) {
	para := doc.AddParagraph()
	para.SetStyle("Heading" + strconv.Itoa(标题等级))
	run := para.AddRun()
	run.Properties().SetBold(true)
	run.Properties().SetFontFamily("宋体")
	run.AddText(标题前缀 + ". " + 标题内容)
}

func 生成开头自动空格的一段文字(段落全部内容 string) {
	para := doc.AddParagraph()
	para.Properties().SetFirstLineIndent(0.5 * measurement.Inch)
	run := para.AddRun()
	run.Properties().SetFontFamily("宋体")
	run.AddText(段落全部内容)
	run.AddBreak()
}

func 从文件中生成开头不自动空格的一段文字(文件路径 string) {
	file, err := os.Open(文件路径)
	if err != nil {
		fmt.Printf("打开文件: %v 失败, 错误信息: %v\n", 文件路径, err)
		return
	}
	defer file.Close()

	para := doc.AddParagraph()
	run := para.AddRun()
	run.Properties().SetFontFamily("宋体")
	br := bufio.NewReader(file)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		str := string(line)
		run.AddText(str)
		run.AddBreak()
	}
}

func 从文件中生成代码(文件路径 string) {
	file, err := os.Open(文件路径)
	if err != nil {
		fmt.Printf("打开文件: %v 失败, 错误信息: %v\n", 文件路径, err)
		return
	}
	defer file.Close()
	para := doc.AddParagraph()
	run := para.AddRun()
	run.Properties().SetFontFamily("Courier New")
	run.Properties().SetHighlight(wml.ST_HighlightColorLightGray)

	br := bufio.NewReader(file)
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		str := string(line)
		run.AddText(str)
		run.AddBreak()
	}
}

func 生成图片(imgPath string) {
	img, err := common.ImageFromFile(imgPath)
	if err != nil {
		log.Fatalf("unable to create image: %s", err)
	}
	imgRef, err := doc.AddImage(img)
	if err != nil {
		log.Fatalf("unable to add image to document: %s", err)
	}

	para := doc.AddParagraph()
	para.SetAlignment(wml.ST_JcCenter)
	run := para.AddRun()
	inl, err := run.AddDrawingInline(imgRef)
	if err != nil {
		log.Fatalf("unable to add inline image: %s", err)
	}
	inl.SetSize(4*measurement.Inch, 4*measurement.Inch)
	run.AddBreak()
}
