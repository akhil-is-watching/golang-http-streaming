package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/video", func(c *fiber.Ctx) error {
		params := c.Queries()
		videoFile, err := os.Open(fmt.Sprintf("videos/%s.mp4", params["ques"]))
		if err != nil {
			panic(err)
		}

		fileInfo, err := videoFile.Stat()
		if err != nil {
			return err
		}

		c.Set("Content-Type", "video/mp4")
		c.Set("Content-Length", strconv.Itoa(int(fileInfo.Size())))
		c.Set("Accept-Ranges", "bytes")

		if c.Get("Range") != "" {
			rangeHeader := c.Get("Range")
			rangeValue := rangeHeader[len("bytes="):]
			start, end := parseRange(rangeValue, int(fileInfo.Size()))

			c.Set("Content-Range", "bytes "+strconv.Itoa(start)+"-"+strconv.Itoa(end)+"/"+strconv.Itoa(int(fileInfo.Size())))
			c.Status(fiber.StatusPartialContent)

			videoFile.Seek(int64(start), 0)

			buffer := make([]byte, end-start+1)
			videoFile.Read(buffer)
			c.Send(buffer)
		} else {
			c.SendStream(videoFile)
		}

		return nil
	})

	err := app.Listen(":3000")
	if err != nil {
		panic(err)
	}
}

func parseRange(rangeHeader string, fileSize int) (int, int) {
	var start, end int

	// Assuming the range header is in the format "bytes=start-end"
	parts := strings.Split(rangeHeader, "-")
	if len(parts) == 2 {
		fmt.Sscanf(parts[0], "bytes=%d", &start)
		fmt.Sscanf(parts[1], "%d", &end)
	}

	if end == 0 || end > fileSize-1 {
		end = fileSize - 1
	}

	return start, end
}
