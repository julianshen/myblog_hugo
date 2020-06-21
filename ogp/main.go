package main

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gohugoio/hugo/parser/pageparser"
	"github.com/h2non/bimg"
	"github.com/julianshen/text2img"
	flag "github.com/spf13/pflag"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
	"gopkg.in/yaml.v2"
)

var basePath = flag.StringP("base", "b", "", "Base path")
var postName = flag.StringP("file", "f", "", "File name (under $basePath/content/post)")
var l = log.New(os.Stderr, "", 0)

func betterOGImage(images []*bimg.Image, title string) (string, error) {
	fileName := *basePath + "/static/images/posts/" + filepath.Base(*postName) + ".jpg"
	uri := "/images/posts/" + filepath.Base(*postName) + ".jpg"
	if len(images) > 0 {
		options := bimg.Options{
			Width:   1200,
			Height:  630,
			Crop:    true,
			Enlarge: true,
			Gravity: bimg.GravitySmart,
		}
		t, err := images[0].Process(options)
		if err != nil {
			return "", err
		}
		bimg.Write(fileName, t)
	} else {
		d, err := text2img.NewDrawer(text2img.Params{
			FontPath: *basePath + "/fonts/default.ttf",
		})

		if err != nil {
			return "", err
		}

		img, err := d.Draw(title)

		if err != nil {
			return "", err
		}

		file, err := os.Create(fileName)
		defer file.Close()

		if err != nil {
			return "", err
		}
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 100})

		if err != nil {
			return "", err
		}
	}

	return uri, nil
}

func getImages(content string) []*bimg.Image {
	markdown := goldmark.New()
	n := markdown.Parser().Parse(text.NewReader([]byte(content)))
	images := make([]*bimg.Image, 0)
	ast.Walk(n, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if n.Kind() == ast.KindImage && entering {
			var buffer []byte
			var err error
			var imgPath string
			dest := string((n.(*ast.Image)).Destination)

			if strings.HasPrefix(dest, "http://") || strings.HasPrefix(dest, "https://") {
				imgPath = dest
				resp, err := http.Get(dest)

				if err != nil {
					l.Println(err)
				} else {
					buffer, err = ioutil.ReadAll(resp.Body)

					if err != nil {
						return ast.WalkContinue, err
					}
				}
			} else {
				imgPath = *basePath + "/static" + dest
				buffer, err = bimg.Read(imgPath)

				if err != nil {
					return ast.WalkContinue, err
				}
			}

			if err == nil {
				img := bimg.NewImage(buffer)
				size, err := img.Size()

				if err != nil {
					l.Printf("Error get image size %s : %v", imgPath, err)
				} else if size.Width > 150 && size.Height > 150 {
					images = append(images, img)
				} else {
					l.Printf("Image: %s is smaller than 150x150\n", imgPath)
				}
			} else {
				return ast.WalkContinue, err
			}
		}
		return ast.WalkContinue, nil
	})

	return images
}

func main() {
	flag.Parse()

	if postName == nil || *postName == "" {
		panic("Need provide post name")
	}

	fileName := *basePath + "/content/post/" + *postName
	f, err := os.Open(fileName)
	defer f.Close()

	if err != nil {
		panic(err)
	}

	page, err := pageparser.ParseFrontMatterAndContent(f)

	if err != nil {
		panic(err)
	}

	images := getImages(string(page.Content))

	ogImgPath, err := betterOGImage(images, page.FrontMatter["title"].(string))

	if err == nil {
		if _, ok := page.FrontMatter["images"]; !ok {
			page.FrontMatter["images"] = []string{ogImgPath}
		}
	}

	fmt.Println("---")
	d, _ := yaml.Marshal(&page.FrontMatter)
	fmt.Println(string(d))
	fmt.Println("---")
	fmt.Println(string(page.Content))
}
