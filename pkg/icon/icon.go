package icon

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
)

// createCircularIcon creates a circular icon with white background and transparent corners
func createCircularIcon(img image.Image, size int) image.Image {
	// Resize to target size (the logo with white background)
	resized := imaging.Resize(img, size, size, imaging.Lanczos)

	// Create output image with transparent background (RGBA)
	output := image.NewRGBA(image.Rect(0, 0, size, size))

	// Calculate center and radius (slightly smaller to ensure circle fits)
	centerX := float64(size) / 2
	centerY := float64(size) / 2
	radius := float64(size) / 2

	// Copy pixels - white background inside circle, transparent outside
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) - centerX
			dy := float64(y) - centerY
			dist := math.Sqrt(dx*dx + dy*dy)

			if dist <= radius {
				// Inside circle - copy pixel from resized image (which has white background)
				pixel := resized.At(x, y)
				r, g, b, a := pixel.RGBA()
				output.SetRGBA(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})
			} else {
				// Outside circle - transparent
				output.SetRGBA(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}

	return output
}

// GetBackgroundColor returns white background by default (or custom color if provided)
func GetBackgroundColor(bgColor string, autoBG bool, logoPath string) (color.RGBA, error) {
	var target color.RGBA

	// Use provided bgColor if given
	if bgColor != "" && len(bgColor) == 7 && bgColor[0] == '#' {
		r := parseHex(bgColor[1:3])
		g := parseHex(bgColor[3:5])
		b := parseHex(bgColor[5:7])
		return color.RGBA{R: r, G: g, B: b, A: 255}, nil
	}

	// Default to white background
	target = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	return target, nil
}

func parseHex(s string) uint8 {
	var v uint8
	fmt.Sscanf(s, "%02x", &v)
	return v
}

// IconPaths holds output file paths for generated icons
type IconPaths struct {
	AdaptiveXML   string
	LegacyMDpi    string
	LegacyHdpi    string
	LegacyXhdpi   string
	LegacyXXhdpi  string
	LegacyXXXhdpi string
	Foreground512 string // 512x512 foreground
	Playstore     string // ic_launcher-playstore.png
	NotifMDpi     string
	NotifHdpi     string
	NotifXhdpi    string
	NotifXXhdpi   string
	NotifXXXhdpi  string
	AppLogo       string
}

// GenerateAppIcon creates launcher icons: adaptive XML + legacy PNGs + playstore icon
func GenerateAppIcon(logoPath, baseName, outputDir string, bg color.RGBA, dryRun bool) (IconPaths, error) {
	var paths IconPaths

	logo, err := imaging.Open(logoPath)
	if err != nil {
		return paths, fmt.Errorf("load logo: %w", err)
	}

	// Resize logo for foreground (transparent background)
	fg := imaging.Resize(logo, 512, 512, imaging.Lanczos)

	// Create white background
	bgImg := image.NewRGBA(image.Rect(0, 0, 512, 512))
	draw.Draw(bgImg, bgImg.Bounds(), &image.Uniform{C: bg}, image.Point{}, draw.Src)

	// Combine foreground + background for launcher icon
	combined := image.NewRGBA(bgImg.Bounds())
	draw.Draw(combined, bgImg.Bounds(), bgImg, image.Point{}, draw.Src)
	offset := (512 - fg.Bounds().Dx()) / 2
	draw.Draw(combined, fg.Bounds().Add(image.Pt(offset, offset)), fg, image.Point{}, draw.Over)

	resDir := filepath.Join(outputDir, "app/src", baseName, "res")

	if !dryRun {
		// Save 512x512 app_logo.png in drawable (white background)
		drawableDir := filepath.Join(resDir, "drawable")
		os.MkdirAll(drawableDir, 0755)
		appLogoPath := filepath.Join(drawableDir, "app_logo.png")
		if f, err := os.Create(appLogoPath); err == nil {
			png.Encode(f, combined)
			f.Close()
			paths.AppLogo = appLogoPath
		}

		// Create adaptive icon XML for ic_launcher
		xmlDir := filepath.Join(resDir, "mipmap-anydpi-v26")
		os.MkdirAll(xmlDir, 0755)

		// ic_launcher.xml
		launcherXML := generateAdaptiveXML("ic_launcher")
		os.WriteFile(filepath.Join(xmlDir, "ic_launcher.xml"), []byte(launcherXML), 0644)
		paths.AdaptiveXML = filepath.Join(xmlDir, "ic_launcher.xml")

		// ic_launcher_round.xml
		roundXML := generateAdaptiveXML("ic_launcher_round")
		os.WriteFile(filepath.Join(xmlDir, "ic_launcher_round.xml"), []byte(roundXML), 0644)
	}

	// Generate legacy launcher icons at different densities
	densities := []struct {
		name string
		size int
	}{
		{"mdpi", 48}, {"hdpi", 72}, {"xhdpi", 96}, {"xxhdpi", 144}, {"xxxhdpi", 192},
	}

	// Resize foreground for each density
	foregroundImages := make(map[string]image.Image)
	for _, d := range densities {
		img := imaging.Resize(fg, d.size, d.size, imaging.Lanczos)
		foregroundImages[d.name] = img
	}

	for _, d := range densities {
		// Combined icon (with background)
		img := imaging.Resize(combined, d.size, d.size, imaging.Lanczos)
		dir := filepath.Join(resDir, "mipmap-"+d.name)
		if !dryRun {
			os.MkdirAll(dir, 0755)
			// ic_launcher.png
			outPath := filepath.Join(dir, "ic_launcher.png")
			if f, err := os.Create(outPath); err == nil {
				png.Encode(f, img)
				f.Close()
				setLegacyPath(&paths, d.name, outPath)
			}
			// ic_launcher_foreground.png (transparent foreground only)
			fgPath := filepath.Join(dir, "ic_launcher_foreground.png")
			if f, err := os.Create(fgPath); err == nil {
				png.Encode(f, foregroundImages[d.name])
				f.Close()
			}
		}
	}

	// Generate ic_launcher_round - for API 25+ (adaptive icon with round shape)
	for _, d := range densities {
		// Create circular icon
		circularImg := createCircularIcon(combined, d.size)
		dir := filepath.Join(resDir, "mipmap-"+d.name)
		if !dryRun {
			roundPath := filepath.Join(dir, "ic_launcher_round.png")
			if f, err := os.Create(roundPath); err == nil {
				png.Encode(f, circularImg)
				f.Close()
			}
		}
	}

	// Create ic_launcher-playstore.png (512x512 with white bg, same as app_logo)
	// This goes in the root of src/<flavor>/ (same level as google-services.json)
	if !dryRun {
		srcDir := filepath.Join(outputDir, "app/src", baseName)
		playstorePath := filepath.Join(srcDir, "ic_launcher-playstore.png")
		if f, err := os.Create(playstorePath); err == nil {
			png.Encode(f, combined)
			f.Close()
			paths.Playstore = playstorePath
		}
	}

	return paths, nil
}

func setLegacyPath(paths *IconPaths, density, path string) {
	switch density {
	case "mdpi":
		paths.LegacyMDpi = path
	case "hdpi":
		paths.LegacyHdpi = path
	case "xhdpi":
		paths.LegacyXhdpi = path
	case "xxhdpi":
		paths.LegacyXXhdpi = path
	case "xxxhdpi":
		paths.LegacyXXXhdpi = path
	}
}

func generateAdaptiveXML(name string) string {
	// Use the name parameter for the foreground drawable
	foreground := "@drawable/ic_launcher_foreground"
	if name != "ic_launcher" {
		foreground = "@drawable/ic_launcher_foreground"
	}
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<adaptive-icon xmlns:android="http://schemas.android.com/apk/res/android">
    <background android:drawable="@color/ic_launcher_background"/>
    <foreground android:drawable="%s"/>
</adaptive-icon>`, foreground)
}

// GenerateNotificationIcon creates white-on-transparent notification icons
func GenerateNotificationIcon(logoPath, baseName, outputDir string, dryRun bool) (string, error) {
	logo, err := imaging.Open(logoPath)
	if err != nil {
		return "", fmt.Errorf("load logo: %w", err)
	}

	// Convert to grayscale then to white on transparent
	logo = imaging.Grayscale(logo)
	logo = imaging.Resize(logo, 512, 512, imaging.Lanczos)

	// Create white logo on transparent background
	whiteLogo := image.NewRGBA(logo.Bounds())
	for y := 0; y < logo.Bounds().Dy(); y++ {
		for x := 0; x < logo.Bounds().Dx(); x++ {
			px := logo.At(x, y)
			r, g, b, _ := px.RGBA()
			// Non-white pixels become white, white becomes transparent
			if r < 65535 || g < 65535 || b < 65535 {
				whiteLogo.SetRGBA(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				whiteLogo.SetRGBA(x, y, color.RGBA{255, 255, 255, 0})
			}
		}
	}

	densities := []struct {
		name string
		size int
	}{
		{"mdpi", 24}, {"hdpi", 36}, {"xhdpi", 48}, {"xxhdpi", 72}, {"xxxhdpi", 96},
	}
	outputBase := filepath.Join(outputDir, "app/src", baseName, "res")
	var lastPath string

	for _, d := range densities {
		img := imaging.Resize(whiteLogo, d.size, d.size, imaging.Lanczos)
		dir := filepath.Join(outputBase, "drawable-"+d.name)
		if !dryRun {
			os.MkdirAll(dir, 0755)
			outPath := filepath.Join(dir, "ic_notification_icon.png")
			if f, err := os.Create(outPath); err == nil {
				png.Encode(f, img)
				f.Close()
				lastPath = outPath
			}
		}
	}
	return lastPath, nil
}
