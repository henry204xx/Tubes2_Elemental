package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	baseURL := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	resp, err := http.Get(baseURL)
	if err != nil {
		fmt.Println("Error fetching page:", err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Error loading HTML:", err)
		return
	}

	os.MkdirAll("downloaded_images", os.ModePerm)

	iterator := 0
	skipped := false

	doc.Find("table.list-table.col-list.icon-hover").Each(func(i int, table *goquery.Selection) {
		if iterator == 1 && !skipped {
			skipped = true
			return
		}

		tierPath := fmt.Sprintf("downloaded_images/tier_%d", iterator)
		os.MkdirAll(tierPath, os.ModePerm)
		iterator++

		table.Find("tbody tr").Each(func(j int, row *goquery.Selection) {
			if j == 0 {
				return
			}

			img := row.Find("td img").First()
			if img.Length() == 0 {
				return
			}

			imgURL, exists := img.Attr("data-src")
			if !exists {
				imgURL, exists = img.Attr("src")
			}
			if !exists || imgURL == "" {
				return
			}

			if strings.HasPrefix(imgURL, "//") {
				imgURL = "https:" + imgURL
			}

			u, err := url.Parse(imgURL)
			if err != nil {
				fmt.Println("Invalid URL:", imgURL)
				return
			}

			parts := strings.Split(u.Path, "/")
			filename := "unknown"
			for _, part := range parts {
				if strings.HasSuffix(part, ".svg") || strings.HasSuffix(part, ".png") {
					filename = part
					break
				}
			}
			if filename == "unknown" && len(parts) >= 5 {
				filename = parts[len(parts)-5]
			}

			fmt.Println("Downloading:", filename)

			originalURL := imgURL
			if strings.Contains(imgURL, "scale-to-width-down") || strings.Contains(imgURL, "revision") {
				if revIndex := indexOf(parts, "revision"); revIndex != -1 {
					svgPath := strings.Join(parts[:revIndex], "/")
					originalURL = fmt.Sprintf("%s://%s%s", u.Scheme, u.Host, svgPath)
				}
			}

			savePath := path.Join(tierPath, filename)

			if err := downloadImage(originalURL, savePath); err != nil {
				fmt.Println("Failed to download original URL:", err)
				if err := downloadImage(imgURL, savePath); err != nil {
					fmt.Println("Failed to download fallback URL:", err)
				} else {
					fmt.Println("Successfully saved with fallback URL:", savePath)
				}
			} else {
				fmt.Println("Saved to:", savePath)
			}
		})
	})

	fmt.Println("Download process finished!")
}

func indexOf(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return -1
}

func downloadImage(imageURL string, filePath string) error {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", imageURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
