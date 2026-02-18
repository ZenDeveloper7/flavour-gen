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

// GetBackgroundColor extracts dominant color from logo corners for adaptive icon background
func GetBackgroundColor(bgColor string, autoBG bool, logoPath string) (color.RGBA, error) {
	var target color.RGBA

	// Use provided bgColor if given
	if bgColor != "" && len(bgColor) == 7 && bgColor[0] == '#' {
		r := parseHex(bgColor[1:3])
		g := parseHex(bgColor[3:5])
		b := parseHex(bgColor[5:7])
		return color.RGBA{R: r, G: g, B: b, A: 255}, nil
	}

	// Auto-detect from logo corners
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

		// Sample 4 corners
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
		target = color.RGBA{
			R: uint8(rSum / count),
			G: uint8(gSum / count),
			B: uint8(bSum / count),
			A: 255,
		}
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

// IconPaths holds output file paths for generated icons
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
	AppLogo       string
}

// GenerateAppIcon creates launcher icons: adaptive XML + legacy PNGs + 512x512 app_logo
func GenerateAppIcon(logoPath, baseName, outputDir string, bg color.RGBA, dryRun bool) (IconPaths, error) {
	var paths IconPaths

	logo, err := imaging.Open(logoPath)
	if err != nil {
		return paths, fmt.Errorf("load logo: %w", err)
	}

	// Resize logo for foreground
	fg := imaging.Resize(logo, 512, 512, imaging.Lanczos)

	// Create background
	bgImg := image.NewRGBA(image.Rect(0, 0, 512, 512))
	draw.Draw(bgImg, bgImg.Bounds(), &image.Uniform{C: bg}, image.Point{}, draw.Src)

	// Combine foreground + background for launcher icon
	combined := image.NewRGBA(bgImg.Bounds())
	draw.Draw(combined, bgImg.Bounds(), bgImg, image.Point{}, draw.Src)
	offset := (512 - fg.Bounds().Dx()) / 2
	draw.Draw(combined, fg.Bounds().Add(image.Pt(offset, offset)), fg, image.Point{}, draw.Over)

	if !dryRun {
		// Save 512x512 app_logo.png in drawable
		drawableDir := filepath.Join(outputDir, "app/src", baseName, "res/drawable")
		if err := os.MkdirAll(drawableDir, 0755); err != nil {
			return paths, err
		}
		appLogoPath := filepath.Join(drawableDir, "app_logo.png")
		if f, err := os.Create(appLogoPath); err == nil {
			png.Encode(f, fg)
			f.Close()
			paths.AppLogo = appLogoPath
		}

		// Create adaptive icon XML
		xmlDir := filepath.Join(outputDir, "app/src", baseName, "res/mipmap-anydpi-v26")
		os.MkdirAll(xmlDir, 0755)
		xmlPath := filepath.Join(xmlDir, "ic_launcher.xml")
		os.WriteFile(xmlPath, []byte(generateAdaptiveXML("ic_launcher")), 0644)
		paths.AdaptiveXML = xmlPath
	}

	// Generate legacy launcher icons at different densities
	densities := []struct {
		name string
		size int
	}{
		{"mdpi", 48}, {"hdpi", 72}, {"xhdpi", 96}, {"xxhdpi", 144}, {"xxxhdpi", 192},
	}
	outputBase := filepath.Join(outputDir, "app/src", baseName, "res")
	for _, d := range densities {
		img := imaging.Resize(combined, d.size, d.size, imaging.Lanczos)
		dir := filepath.Join(outputBase, "mipmap-"+d.name)
		if !dryRun {
			os.MkdirAll(dir, 0755)
			outPath := filepath.Join(dir, "ic_launcher.png")
			if f, err := os.Create(outPath); err == nil {
				png.Encode(f, img)
				f.Close()
				setLegacyPath(&paths, d.name, outPath)
			}
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
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<adaptive-icon xmlns:android="http://schemas.android.com/apk/res/android">
    <background android:drawable="@color/ic_launcher_background"/>
    <foreground android:drawable="@mipmap/%s_foreground"/>
</adaptive-icon>`, name)
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

	densities := []struct{ name string; size int }{
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
