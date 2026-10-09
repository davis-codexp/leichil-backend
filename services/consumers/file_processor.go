package consumers

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	_ "github.com/gen2brain/heic"
	"github.com/gofiber/fiber/v2/log"
	"github.com/nfnt/resize"
	"github.com/skrashevich/go-webp"
)

func RunImageProcessorV2(file_name string) {
	var media_path = os.Getenv("MEDIA_PATH")
	if media_path == "" {
		log.Error("Invalid Media Path")
		return
	}
	var image_path string = filepath.Join(media_path, file_name)
	const limit = 720
	info, err := os.Stat(image_path)
	if err != nil {
		log.Error("Unable to get image properties: ", err.Error())
		return
	}
	var file_size_kb int64 = info.Size() / 1000
	file, err := os.Open(image_path)
	if err != nil {
		log.Error("Unable to open image: ", err.Error())
		return
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		log.Error("Failed to decode image: ", err.Error())
		return
	}
	var bounds image.Rectangle = img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	var longest = max(width, height)
	var processed image.Image = img
	if file_size_kb > 500 && longest > limit {
		var scale int = (longest + limit - 1) / limit
		var target_width int = width / scale
		var target_height int = height / scale
		target_height = max(target_height, 1)
		target_width = max(target_width, 1)
		processed = resize.Resize(uint(target_width), uint(target_height), img, resize.Lanczos3)
	}
	var ext string = filepath.Ext(file_name)
	var base_name = strings.TrimSuffix(file_name, ext)
	var output_path string = filepath.Join(media_path, base_name+".webp")
	dest, err := os.Create(output_path)
	if err != nil {
		log.Error("Failed to generate result file: ", err.Error())
		return
	}
	defer dest.Close()
	if err = webp.Encode(dest, processed, &webp.Options{
		Lossy:   true,
		Quality: 80,
	}); err != nil {
		log.Error("Native webp compression failed: ", err.Error())
		return
	}

	_ = file.Close()
	if err := os.Remove(image_path); err != nil {
		log.Error("Failed cleanup of original file: ", err.Error())
		return
	}
	log.Info("Image compression and conversion completed")
}

func GetMimeType(buffer []byte, filename string) string {
	var mime_type = http.DetectContentType(buffer)
	if mime_type != "application/octet-stream" {
		return mime_type
	}
	if is_heif(buffer) {
		return "image/heic"
	}
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".heic", ".heif", ".hif":
		return "image/heic"
	}
	return mime_type
}

func is_heif(buffer []byte) bool {
	if len(buffer) > 12 || string(buffer[4:8]) != "ftyp" {
		return false
	}
	switch string(buffer[8:12]) {
	case "heic", "heix", "heim", "heis", "mif1", "hevc", "hevx", "hevm":
		return true
	default:
		return false
	}
}

func RunImageProcessor(file_name string) {
	media_path := os.Getenv("MEDIA_PATH")
	image_path := media_path + "/" + file_name
	const limit = 720
	factor := 0
	scale_factor := 1

	my_file, err := os.Stat(image_path)
	if err != nil {
		log.Error("Unable to get image path: ", err.Error())
		return
	}

	file_size := my_file.Size() / 1000
	file, err := os.Open(image_path)
	if err != nil {
		log.Error("Unable to open image: ", err.Error())
		return
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", image_path, err)
	}

	if image.Width >= image.Height {
		factor = image.Width
	} else {
		factor = image.Height
	}

	if file_size > 100 && factor > 720 {
		scale_factor = (factor / limit) + 1
	}

	width := image.Width / scale_factor
	height := image.Height / scale_factor
	scale := fmt.Sprintf("scale=%v:%v", width, height)
	output_file := strings.Split(file_name, ".")[0] + ".webp"
	output_path := media_path + "/" + output_file

	cmd := exec.Command("ffmpeg", "-i", image_path, "-vf", scale, output_path)
	_, err = cmd.Output()
	if err != nil {
		log.Error("Could not run ffmpeg: ", err.Error())
		return
	}

	cmd = exec.Command("rm", image_path)
	_, err = cmd.Output()
	if err != nil {
		log.Error("Could not delete file", err.Error())
		return
	}
}
