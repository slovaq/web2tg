package data

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"

	"github.com/disintegration/imaging"
	"github.com/mallvielfrass/fmc"
)

func PrepareImg(file multipart.File, header *multipart.FileHeader) string {
	FileName := "web/files/" + header.Filename
	f, err := os.Create(FileName)
	if err != nil {
		fmc.Printfln("#rbtError: #ybt%s", err.Error())
	}
	// It's idiomatic to defer a `Close` immediately
	// after opening a file.
	defer f.Close()
	fmt.Println("file: ", file)
	buf := bytes.NewBuffer(nil)

	if _, err := io.Copy(buf, file); err != nil {
		fmc.Printfln("#rbtError: #ybt%s", err.Error())
	}

	n2, err := f.Write(buf.Bytes())
	if err != nil {
		fmc.Printfln("#rbtError: #ybt%s", err.Error())
	}
	fmt.Printf("wrote %d bytes\n", n2)

	//	io.Copy(f, file)
	pic := FileName
	fmt.Println("file: ", header.Size)
	src, err := imaging.Open(pic)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	// Resize srcImage to size = 128x128px using the Lanczos filter.
	g := src.Bounds()

	// Get height and width

	if 4500000 < header.Size || 3061 < g.Dy() || 3061 < g.Dx() {
		height := float64(g.Dy())
		width := float64(g.Dx())

		perc := float64(4999999) / float64(header.Size)
		fmt.Printf(" %f %f \n perc= %f\n", width, height, perc)
		if 4500000 < header.Size {

			height = height * perc
			width = width * perc
		}

		fmt.Printf(" %f %f\n", width, height)

		for {
			if 3060 < height || 3060 < width {
				height = height * 0.9
				width = width * 0.9
			} else {
				break
			}
		}
		h := int(height)
		w := int(width)
		fmt.Printf(" %d %d\n", w, h)
		d := fmt.Sprintf("ffmpeg -y -i %s -vf scale=%d:%d %s", pic, w, h, pic)
		//d := fmt.Sprintf("ffmpeg -y -i %s -vf scale=512:512 %s", pic, pic)
		lsCmd := exec.Command("bash", "-c", d)
		_, err = lsCmd.Output()
		if err != nil {
			panic(err)
		}
	}
	defer file.Close()
	return pic
}
