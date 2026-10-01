package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type BitrateProfile struct {
	Resolution  string
	VideoBit    string
	MaxVideoBit string
	BufSize     string
	AudioBit    string
	Bandwidth   int
}

/*
var profile_map = map[string]BitrateProfile{
	"480p":  {Resolution: "854x480", VideoBit: "1500k", AudioBit: "128k"},
	"720p":  {Resolution: "1280x720", VideoBit: "2500k", AudioBit: "128k"},
	"1080p": {Resolution: "1920x1080", VideoBit: "5000k", AudioBit: "192k"},
}
*/

var profile_map = map[string]BitrateProfile{
	"480p": {
		Resolution: "854x480", VideoBit: "800k", MaxVideoBit: "1200k",
		BufSize: "1600k", AudioBit: "128k", Bandwidth: 800000,
	},
	"720p": {
		Resolution: "1280x720", VideoBit: "1800k", MaxVideoBit: "2500k",
		BufSize: "3600k", AudioBit: "128k", Bandwidth: 1800000,
	},
	"1080p": {
		Resolution: "1920x1080", VideoBit: "3200k", MaxVideoBit: "4500k",
		BufSize: "6400k", AudioBit: "192k", Bandwidth: 3200000,
	},
}

func sanitize_data[T string](data *T) T {
	var value T
	if data == nil {
		return value
	}
	return *data
}

func main() {
	/* Sample cmd -> ./mdx-transcoder -input=<path_to_raw_file> -qualities=<comma_separated_qualities>  */
	var input_path *string = flag.String("input", "", "Absolute path to the source raw video file")
	var qualities_input *string = flag.String("qualities", "480p,720p,1080p", "Comma separated string of target video profiles")
	flag.Parse()
	fmt.Println(*input_path)
	if sanitize_data(input_path) == "" {
		log.Fatalf("\nError: Missing input file path")
	}
	if _, err := os.Stat(*input_path); os.IsNotExist(err) {
		log.Fatalf("\nError: The source file could not be found in the given path")
	}

	var asset_id string = uuid.New().String()
	var parent_dir string = filepath.Join(".", "videos", asset_id)
	fmt.Printf("\nInitializing Transcode operation for Asset ID: %s", asset_id)
	fmt.Printf("\nCreating Parent Directory at: %s\n", parent_dir)
	if err := os.MkdirAll(parent_dir, 0755); err != nil {
		log.Fatalf("\nError: Failed to create parent directory: %v", err)
	}
	var target_qualities []string = strings.Split(sanitize_data(qualities_input), ",")
	var variant_playlist []string
	var wg sync.WaitGroup

	for _, quality := range target_qualities {
		var quality_key string = strings.TrimSpace(strings.ToLower(quality))
		profile, verified := profile_map[quality_key]
		if !verified {
			fmt.Printf("Warning: Quality Profile %s is unrecognized. Skipping.\n", quality)
			continue
		}

		wg.Add(1)
		go func(key string, bp BitrateProfile) {
			defer wg.Done()
			var sub_dir string = filepath.Join(parent_dir, key)
			if err := os.MkdirAll(sub_dir, 0755); err != nil {
				log.Fatalf("Failed to create sub directory: %v", err)
			}
			var playlist_name string = fmt.Sprintf("%s/endpoint.m3u8", key)
			variant_playlist = append(variant_playlist, playlist_name)
			var playlist_path string = filepath.Join(sub_dir, "endpoint.m3u8")
			var segment_file_pattern = filepath.Join(sub_dir, "segment_%03d.ts")

			/*
				var ffmeg_args = []string{
					"-y", "-i", *input_path,
					"-vf", fmt.Sprintf("scale=%s", bp.Resolution),
					"-c:v", "libx264", "-profile:v", "high", "-preset", "medium",
					"-b:v", bp.VideoBit, "-maxrate", bp.MaxVideoBit, "-bufsize", bp.BufSize,
					"-c:a", "aac", "-b:a", bp.AudioBit, "-hls_time", "10",
					"-hls_playlist_type", "vod",
					"-hls_segment_type", "fmp4",
					"-hls_segment_filename", segment_file_pattern,
					playlist_path,
				}
			*/

			var ffmeg_args = []string{
				"-y",
				"-hwaccel", "videotoolbox", // Enables GPU-accelerated hardware decoding
				"-i", *input_path,
				"-vf", fmt.Sprintf("scale=%s", bp.Resolution),
				"-c:v", "h264_videotoolbox", // Switches from CPU (libx264) to M1 GPU encoder
				"-g", "240",
				"-realtime", "0",
				"-profile:v", "high",
				"-b:v", bp.VideoBit,
				"-maxrate", bp.MaxVideoBit,
				"-bufsize", bp.BufSize,
				"-c:a", "aac",
				"-b:a", bp.AudioBit,
				"-hls_time", "10",
				"-hls_playlist_type", "vod",
				"-hls_segment_type", "mpegts",
				"-hls_segment_filename", segment_file_pattern,
				playlist_path,
			}

			var cmd = exec.Command("ffmpeg", ffmeg_args...)
			if err := cmd.Run(); err != nil {
				log.Fatalf("\nError: FFMPEG encoding exception on profile %s: %v", key, err)
			}
			fmt.Printf("Completed variant asset block: %s\n", key)
		}(quality_key, profile)
	}
	wg.Wait()

	/*----------------GENERATE THE MASTER PLAYLIST .m3u8 FILE---------------------*/
	if len(variant_playlist) > 0 {
		var master_playlist_path string = filepath.Join(parent_dir, "master.m3u8")
		fmt.Printf("Generating Master Playlist")
		master_file, err := os.Create(master_playlist_path)
		if err != nil {
			log.Fatalf("Error: Failed to create Master Playlist: %v", err)
		}
		defer master_file.Close()
		_, _ = master_file.WriteString("#EXTM3U\n")
		for _, val := range target_qualities {
			var key string = strings.TrimSpace(strings.ToLower(val))
			profile, verified := profile_map[key]
			if !verified {
				continue
			}

			var line string = fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s\n%s/endpoint.m3u8\n",
				profile.Bandwidth,
				profile.Resolution,
				key,
			)

			_, _ = master_file.WriteString(line)
		}
	}
	fmt.Printf("\n Transcoding completed")
}
