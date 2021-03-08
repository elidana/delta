package main

import (
	"io"
	"text/template"
	"time"
)

var infoTemplate string = `*     GAMIT/GLOBK station.info.geonet
* 
*     Generated by GeoNet on {{update}}
*     Send questions, comments or concerns to info@geonet.org.nz
*
*SITE  Station Name      Session Start      Session Stop       Ant Ht   HtCod  Ant N    Ant E    Receiver Type         Vers                  SwVer  Receiver SN           Antenna Type     Dome   Antenna SN
{{range $i := . -}}
{{printf " %-5s" $i.Site -}}
{{trunc $i.Name 16 | printf " %-16s" -}}
{{at $i.Start | printf "  %s" -}}
{{at $i.End   | printf "  %s" -}}
{{printf " %8.4f" $i.Height -}}
{{printf "  %-5s" $i.HeightCod -}}
{{printf " %8.4f" $i.North -}}
{{printf " %8.4f" $i.East -}}
{{printf "  %-20s" $i.ReceiverType -}}
{{printf "  %-21s" $i.Version -}}
{{printf " %5s" $i.Software -}}
{{printf "  %-20s" $i.ReceiverSerial -}}
{{printf "  %-15s" $i.AntennaType -}}
{{printf "  %-6s" $i.Radome -}}
{{printf " %s" $i.AntennaSerial}}
{{end -}}
`

type Info struct {
	Site string
	Name string

	Start time.Time
	End   time.Time

	Height float64
	North  float64
	East   float64

	Radome    string
	HeightCod string

	Version  string
	Software string

	ReceiverType   string
	ReceiverSerial string
	AntennaType    string
	AntennaSerial  string
}

func (i Info) Like(info Info) bool {
	switch {
	case i.Site != info.Site:
		return false
	case i.Name != info.Name:
		return false
	case i.Height != info.Height:
		return false
	case i.North != info.North:
		return false
	case i.East != info.East:
		return false
	case i.Radome != info.Radome:
		return false
	case i.HeightCod != info.HeightCod:
		return false
	case i.Version != info.Version:
		return false
	case i.Software != info.Software:
		return false
	case i.ReceiverType != info.ReceiverType:
		return false
	case i.ReceiverSerial != info.ReceiverSerial:
		return false
	case i.AntennaType != info.AntennaType:
		return false
	case i.AntennaSerial != info.AntennaSerial:
		return false
	default:
		return true
	}
}

// remove overlaps of stream information not related to firmware
func squashInfo(info ...Info) []Info {

	if !(len(info) > 0) {
		return nil
	}

	var list []Info
	running := info[0]
	for i := 1; i < len(info); i++ {
		switch {
		case running.Like(info[i]):
			running.End = info[i].End
		default:
			list = append(list, running)
			running = info[i]
		}
	}
	list = append(list, running)

	return list
}

// encode stream information in the "info" format.
func encodeInfo(wr io.Writer, update time.Time, info ...Info) error {

	var funcMap template.FuncMap = template.FuncMap{
		"update": func() string {
			return update.UTC().Format("2006-01-02 15:04 MST")
		},
		"at": func(at time.Time) string {
			if at.After(time.Now()) {
				return "9999 999 23 59 30"
			}
			return at.UTC().Format("2006 002 15 04 05")
		},
		"trunc": func(s string, n int) string {
			if len(s) > n {
				return s[0:n]
			}
			return s
		},
	}

	tmpl, err := template.New("info").Funcs(funcMap).Parse(infoTemplate)
	if err != nil {
		return err
	}

	if err := tmpl.Execute(wr, info); err != nil {
		return err
	}

	return nil
}
