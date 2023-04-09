package util

import (
	"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"image"
	_ "image/jpeg"
	"mime/multipart"
)

// RotateImage 苹果手机拍照的图片，会有方向属性Orientation，经过Decode和Encode，编码处理后，方向属性会丢失，导致图片被旋转
func RotateImage(img image.Image, ori int) image.Image {
	switch ori {
	case 6: //90度图片旋转
		img = rotate90(img)
	case 3:
		img = rotate180(img)
	case 8:
		img = rotate270(img)
	}
	return img
}

// 旋转90度
func rotate90(m image.Image) image.Image {
	rotate := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dy(), m.Bounds().Dx()))
	// 矩阵旋转
	for x := m.Bounds().Min.Y; x < m.Bounds().Max.Y; x++ {
		for y := m.Bounds().Max.X - 1; y >= m.Bounds().Min.X; y-- {
			//  设置像素点
			rotate.Set(m.Bounds().Max.Y-x, y, m.At(y, x))
		}
	}
	return rotate
}

// 旋转180度
func rotate180(m image.Image) image.Image {
	rotate := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dx(), m.Bounds().Dy()))
	// 矩阵旋转
	for x := m.Bounds().Min.X; x < m.Bounds().Max.X; x++ {
		for y := m.Bounds().Min.Y; y < m.Bounds().Max.Y; y++ {
			//  设置像素点
			rotate.Set(m.Bounds().Max.X-x, m.Bounds().Max.Y-y, m.At(x, y))
		}
	}
	return rotate
}

// 旋转270度
func rotate270(m image.Image) image.Image {
	rotate := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dy(), m.Bounds().Dx()))
	// 矩阵旋转
	for x := m.Bounds().Min.Y; x < m.Bounds().Max.Y; x++ {
		for y := m.Bounds().Max.X - 1; y >= m.Bounds().Min.X; y-- {
			// 设置像素点
			rotate.Set(x, m.Bounds().Max.X-y, m.At(y, x))
		}
	}
	return rotate
}

// 方向判断
func ReadOrientation(file multipart.File) int {
	//file, err := os.Open(filename)
	//if err != nil {
	//	fmt.Println("failed to open file, err: ", err)
	//	return 0
	//}
	//defer file.Close()
	x, err := exif.Decode(file)
	if err != nil {
		fmt.Println("failed to decode file, err: ", err)
		return 0
	}

	orientation, err := x.Get(exif.Orientation)
	if err != nil {
		fmt.Println("failed to get orientation, err: ", err)
		return 0
	}
	orientVal, err := orientation.Int(0)
	if err != nil {
		fmt.Println("failed to convert type of orientation, err: ", err)
		return 0
	}

	fmt.Println("the value of photo orientation is :", orientVal)
	return orientVal
}
