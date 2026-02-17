package icon

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// GetBackgroundColor extracts dominant color from corners if autoBG is true and bgColor empty.
func GetBackgroundColor(bgColor string, autoBG bool, logoPath string) (color.RGBA, error) {
	var target color.RGBA
	if bgColor != "" {
		if len(bgColor) == 7 && bgColor[0] == '#' {
			r := parseHex(bgColor[1:3])
			g := parseHex(bgColor[3:5])
			b := parseHex(bgColor[5:7])
			target = color.RGBA{R: r, G: g, B: b, A: 255}
			return target, nil
		}
	}
	if autoBG {
		img, err := imaging.Open(logoPath)
		if err != nil {
			return target, fmt.Errorf("cannot open logo: %w", err)
		}
		bounds := img.Bounds()
		w, h := bounds.Dx(), bounds.Dy()
		if w == 0 || h == 0 {
			return target, fmt.Errorf("invalid image size")
		}
		corners := []image.Point{{0, 0}, {w - 1, 0}, {0, h - 1}, {w - 1, h - 1}}
		var rSum, gSum, bSum, count int
		for _, p := range corners {
			c := img.At(p.X, p.Y)
			r, g, b, _ := c.RGBA()
			rSum += int(r >> 8)
			gSum += int(g >> 8)
			bSum += int(b >> 8)
			count++
		}
		if count == 0 {
			return target, fmt.Errorf("no corner samples")
		}
		target = color.RGBA{R: uint8(rSum / count), G: uint8(gSum / count), B: uint8(bSum / count), A: 255}
	} else {
		target = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	}
	return target, nil
}

func parseHex(s string) uint8 {
	var v uint8
	fmt.Sscanf(s, "%02x", &v)
	return v
}

// IconPaths holds output file paths
type IconPaths struct {
	AdaptiveXML   string
	LegacyMDpi    string
	LegacyHdpi    string
	LegacyXhdpi   string
	LegacyXXhdpi  string
	LegacyXXXhdpi string
	NotifMDpi     string
	NotifHdpi     string
	NotifXhdpi    string
	NotifXXhdpi   string
	NotifXXXhdpi  string
}

// GenerateAppIcon creates the adaptive launcher icon and legacy PNGs
func GenerateAppIcon(logoPath, baseName, outputDir string, bg color.RGBA, dryRun bool) (IconPaths, error) {
	var paths IconPaths

	logo, err := imaging.Open(logoPath)
	if err != nil {
		return paths, fmt.Errorf("load logo: %w", err)
	}
	fg := imaging.Resize(logo, 512, 512, imaging.Lanczos)

	bgImg := image.NewRGBA(image.Rect(0, 0, 512, 512))
	draw.Draw(bgImg, bgImg.Bounds(), &image.Uniform{C: bg}, image.Point{}, draw.Src)

	combined := image.NewRGBA(bgImg.Bounds())
	draw.Draw(combined, bgImg.Bounds(), bgImg, image.Point{}, draw.Src)
	offset := (512 - fg.Bounds().Dx()) / 2
	draw.Draw(combined, fg.Bounds().Add(image.Pt(offset, offset)), fg, image.Point{}, draw.Over)

	if !dryRun {
		xmlDir := filepath.Join(outputDir, "app/src", baseName, "res/mipmap-anydpi-v26")
		if err := os.MkdirAll(xmlDir, 0755); err != nil {
			return paths, err
		}
		xmlPath := filepath.Join(xmlDir, "ic_launcher.xml")
		xml := generateAdaptiveXML("ic_launcher")
		if err := os.WriteFile(xmlPath, []byte(xml), 0644); err != nil {
			return paths, err
		}
		paths.AdaptiveXML = xmlPath
	}

	densities := []struct {
		name string
		size int
	}{
		{"mdpi", 48},
		{"hdpi", 72},
		{"xhdpi", 96},
		{"xxhdpi", 144},
		{"xxxhdpi", 192},
	}
	outputBase := filepath.Join(outputDir, "app/src", baseName, "res")
	for _, d := range densities {
		img := imaging.Resize(combined, d.size, d.size, imaging.Lanczos)
		dir := filepath.Join(outputBase, "mipmap-"+d.name)
		if !dryRun {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return paths, err
			}
			outPath := filepath.Join(dir, "ic_launcher.png")
			f, err := os.Create(outPath)
			if err != nil {
				return paths, err
			}
			defer f.Close()
			if err := png.Encode(f, img); err != nil {
				return paths, err
			}
			switch d.name {
			case "mdpi":
				paths.LegacyMDpi = outPath
			case "hdpi":
				paths.LegacyHdpi = outPath
			case "xhdpi":
				paths.LegacyXhdpi = outPath
			case "xxhdpi":
				paths.LegacyXXhdpi = outPath
			case "xxxhdpi":
				paths.LegacyXXXhdpi = outPath
			}
		}
	}

	return paths, nil
}

func generateAdaptiveXML(name string) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<adaptive-icon xmlns:android="http://schemas.android.com/apk/res/android">
    <background android:drawable="@color/ic_launcher_background"/>
    <foreground android:drawable="@mipmap/%s_foreground"/>
</adaptive-icon>`, name)
}

// GenerateNotificationIcon creates a transparent notification icon
func GenerateNotificationIcon(logoPath, baseName, outputDir string, dryRun bool) (string, error) {
	logo, err := imaging.Open(logoPath)
	if err != nil {
		return "", fmt.Errorf("load logo: %w", err)
	}
	densities := []struct{ name string; size int }{
		{"mdpi", 24},
		{"hdpi", 36},
		{"xhdpi", 48},
		{"xxhdpi", 72},
		{"xxxhdpi", 96},
	}
	outputBase := filepath.Join(outputDir, "app/src", baseName, "res")
	var lastPath string
	for _, d := range densities {
		img := imaging.Resize(logo, d.size, d.size, imaging.Lanczos)
		dir := filepath.Join(outputBase, "drawable-"+d.name)
		if !dryRun {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return "", err
			}
			outPath := filepath.Join(dir, "ic_notification_icon.webp")
			f, err := os.Create(outPath)
			if err != nil {
				return "", err
			}
			defer f.Close()
			// Encode as PNG for now (MVP)
			if err := png.Encode(f, img); err != nil {
				return "", err
			}
			lastPath = outPath
		}
	}
	return lastPath, nil
}