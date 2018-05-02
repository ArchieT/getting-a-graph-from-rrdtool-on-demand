package main

import (
	"os"
	"fmt"
	//"io/ioutil"
	"net/http"
	"github.com/ziutek/rrd"
	"time"
	"log"
	"strings"
)

type Parameters struct {
	Start, End    time.Time
	Width, Height uint
	Step          uint
	Title, VLabel string
}

type LineDef struct {
	Width uint8
	Red   uint8
	Green uint8
	Blue  uint8
}

func (l *LineDef) Color() string {
	return fmt.Sprintf("#%02x%02x%02x", l.Red, l.Green, l.Blue)
}

func (l *LineDef) No() bool {
	return l.Width == 0
}

func (l *LineDef) Yes() bool {
	return !l.No()
}

type DefParameters struct {
	Name       string
	Average    LineDef
	AverageMax LineDef
	AverageMin LineDef
	Min        LineDef
	Max        LineDef
}

type Def struct {
	RRDFile string
	Params  DefParameters
}

type DefProto struct {
	Name string
	Type string
	LineDef
}

func Graph(params Parameters, defs []Def) (rrd.GraphInfo, []byte, error) {
	a := rrd.NewGrapher()
	a.SetSize(params.Width, params.Height)
	a.SetTitle(params.Title)
	a.SetVLabel(params.VLabel)
	for _, d := range defs {
		if d.Params.Average.Yes() {
			a.Def(d.Params.Name+"a", d.RRDFile, d.Params.Name, "AVERAGE")
			a.Line(float32(d.Params.Average.Width), d.Params.Name+"a", d.Params.Average.Color())
		}
		if d.Params.Min.Yes() {
			a.Def(d.Params.Name+"i", d.RRDFile, d.Params.Name, "MIN")
			a.Line(float32(d.Params.Min.Width), d.Params.Name+"i", d.Params.Min.Color())
		}
		if d.Params.AverageMin.Yes() {
			a.Def(d.Params.Name+"ai", d.RRDFile, d.Params.Name, "AVERAGE", "reduce=MIN")
			a.Line(float32(d.Params.AverageMin.Width), d.Params.Name+"ai", d.Params.AverageMin.Color())
		}
		if d.Params.AverageMax.Yes() {
			a.Def(d.Params.Name+"ax", d.RRDFile, d.Params.Name, "AVERAGE", "reduce=MAX")
			a.Line(float32(d.Params.AverageMax.Width), d.Params.Name+"ax", d.Params.AverageMax.Color())
		}
		if d.Params.Max.Yes() {
			a.Def(d.Params.Name+"x", d.RRDFile, d.Params.Name, "MAX")
			a.Line(float32(d.Params.Max.Width), d.Params.Name+"x", d.Params.Max.Color())
		}
	}
	if params.Step != 0 {
		a.AddOptions("-S " + fmt.Sprint(params.Step))
	}
	return a.Graph(params.Start,params.End)
}

func StringOfLowerLetters(s string) bool {
	for c := range s {
		if c < 'a' {
			return false
		}
		if c > 'z' {
			return false
		}
	}
	return true
}

func ParsingLineArg(s string) (string, string, LineDef, bool) {
	var d LineDef
	var w string
	fmt.Sscanf(s, "LINE%2d:%s#%02x%02x%02x", &(d.Width), w, &(d.Red), &(d.Green), &(d.Blue))
	var i byte = 0
	var a byte = 0
	if w[len(w)-2] == 'a' {
		switch w[len(w)-1] {
		case 'x':
			a = 'a'
			i = 'x'
		case 'i':
			a = 'a'
			i = 'i'
		default:
			return "", "", d, false
		}
	} else {
		switch w[len(w)-1] {
		case 'a':
			i = 'a'
		case 'i':
			i = 'a'
		case 'x':
			i = 'x'
		default:
			return "", "", d, false
		}
	}
	var n string
	var t string
	if a == 0 {
		n = w[:len(w)-1]
		t = string(i)
	} else if a == 'a' {
		n = w[:len(w)-2]
		t = string(a) + string(i)
	} else {
		n = ""
		t = ""
	}
	return n, t, d, StringOfLowerLetters(n)
}

func mergeProto(ps []DefProto, files map[string]string) (ou []Def) {
	ou = make([]Def, 0, 3)
	for _, p := range ps {
		var found *DefParameters = nil
		for _, o := range ou {
			if o.RRDFile == files[p.Name] {
				found = &(o.Params)
				break
			}
		}
		if found == nil {
			emptyLineDef := LineDef{0, 0, 0, 0}
			ou = append(ou,
				Def{
					files[p.Name],
					DefParameters{
						p.Name,
						emptyLineDef,
						emptyLineDef,
						emptyLineDef,
						emptyLineDef,
						emptyLineDef,
					},
				})
			found = &(ou[len(ou)-1].Params)
		}
		switch p.Type {
		case "a":
			found.Average = p.LineDef
		case "ax":
			found.AverageMax = p.LineDef
		case "ai":
			found.AverageMin = p.LineDef
		case "i":
			found.Min = p.LineDef
		case "x":
			found.Max = p.LineDef
		default:
			panic("panicbotencasetuw185lkjfdsgdsgyyy" + p.Type + "ppp")
		}
	}
	return
}

func main() {
	argsWithoutProg := os.Args[1:]
	files := make(map[string]string)
	key := ""
	for _, v := range argsWithoutProg {
		if key == "" {
			key = v
		} else {
			files[key] = v
			key = ""
		}
	}
	http.HandleFunc("/",
		func(w http.ResponseWriter, r *http.Request) {
			theUnjustPath := strings.Split(r.URL.Path, "/")
			thePath := make([]string, 0, 10)
			for _, pael := range theUnjustPath {
				if pael != "" {
					thePath = append(thePath, pael)
				}
			}
			if thePath[0] == "deliver" {
				timely := strings.Split(thePath[1], "_")
				secondss := timely[0]
				steppp := "nostep"
				if len(timely) > 1 {
					steppp = timely[1]
				}
				widthss := thePath[2]
				heigss := thePath[3]
				var secs, wid, hei, ste uint
				fmt.Sscanf(secondss, "%d", &secs)
				fmt.Sscanf(widthss, "%d", &wid)
				fmt.Sscanf(heigss, "%d", &hei)
				if steppp == "nostep" {
					ste = 0
				} else {
					fmt.Sscanf(steppp, "%d", &ste)
				}
				if secs == 0 || secs > (100*365*24*3600) || wid > 2050 || hei > 2050 || wid == 0 || hei == 0 {
					panic("panic220lkjfdsgkljfd" + secondss + ",ii," + widthss + ",oo," + heigss + "ppppppfdsgsdhgsv")
				}
				what := thePath[len(thePath)-1]
				whats := strings.Split(what, "_")
				ou := make([]DefProto, 0, 10)
				for _, wha := range whats {
					n, t, d, g := ParsingLineArg(wha)
					if files[n] == "" || !g {
						panic("thishandlefunclksagfdsaglkh")
					}
					ou = append(ou, DefProto{n, t, d})
				}
				_, b, e := Graph(
					Parameters{
						time.Now().Add(time.Duration(-secs) * time.Second),
						time.Now(),
						wid, hei,
						ste,
						"Temperatura z ostatnich " + secondss + "s step:" + steppp,
						"Temperatura w °C",
					}, mergeProto(ou, files))
				w.Write(b)
				if e != nil {
					panic("tutajpanicaaaaaa256aaaa")
				}
			} else {
				fmt.Fprintln(w,
					"<a href=\"deliver/360000_900/1000/420/LINE2:tempa#000000\">aaa</a>")
			}
		})
	log.Fatal(http.ListenAndServe(":8085", nil))
}