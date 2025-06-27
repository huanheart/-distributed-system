package ffmpeg

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// 统计本地音视频的时长
// CountDuration 返回本地音视频文件的时长
func CountDuration(filePath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("ffprobe 执行失败: %v, 详情: %s", err, stderr.String())
	}

	// 去掉换行符并转换为 float64
	durationStr := strings.TrimSpace(out.String())
	durationFloat, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, fmt.Errorf("转换为 float64 失败: %v", err)
	}

	return durationFloat, nil
}
