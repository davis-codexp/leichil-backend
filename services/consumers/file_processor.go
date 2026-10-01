package consumers

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"strings"

	"github.com/gofiber/fiber/v2/log"
)

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
