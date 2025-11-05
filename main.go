package main

import (
	"errors"
	"flag"
	"image"
	"image/draw"
	"image/png"
	"io/fs"
	"os"
)

func main() {
	// Парсим флаги
	flag.Parse()

	// Проверяем кол-во позиционных аргументов
	if len(flag.Args()) < 2 {
		println("Usage: veilsteg [flags] <what_to_hide> <where_to_hide.png>")
		os.Exit(2)
	}

	// Проверяем файл на чтение
	for _, filePath := range flag.Args() {
		// Пытаемся получить доступ к файлу
		err := checkPath(filePath)

		// Производим проверку на ошибки
		switch {
		case errors.Is(err, fs.ErrNotExist):
			println("File", filePath, "do not exists")
			os.Exit(1)
		case err == nil:
			continue
		default:
			println("Reading error:\n", err.Error())
			os.Exit(1)
		}
	}

	// Записываем пути до нагрузки и контейнера для сокрытия
	payloadPath, coverPNGPath := flag.Arg(0), flag.Arg(1)

	// Проверяем являются ли они файлом
	if !isFile(payloadPath) {
		println(payloadPath, "is not a file")
	}
	if !isFile(coverPNGPath) {
		println(coverPNGPath, "is not a file")
	}

	// Открываем изображение-контейнер
	in, err := os.Open(coverPNGPath)
	if err != nil {
		println("Error reading file to hide:\n", err.Error())
		os.Exit(1)
	}
	defer in.Close()

	// Декодируем изображение
	src, err := png.Decode(in)
	if err != nil {
		println("Image decoding error:\n", err.Error())
		os.Exit(1)
	}

	//
	b := src.Bounds()
	rgba := image.NewRGBA(b)
	draw.Draw(rgba, b, src, b.Min, draw.Src)

	//
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			// Читаем исходный цвет
			r, g, bl, a := rgba.At(x, y).RGBA()

			// RGBA() возвращает 16-битные компоненты [0..65535], приводим к 8-битным
			R := uint8(r >> 8)
			G := uint8(g >> 8)
			B := uint8(bl >> 8)
			A := uint8(a >> 8)

			println(R, G, B, A)
		}
	}

}
