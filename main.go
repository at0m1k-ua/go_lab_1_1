package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type MeasurementType struct {
	Name  string
	Label string
	Units string
	Value string
}

type TemplateModel struct {
	Measurements []MeasurementType
	CalcResult   []MeasurementType
	Error        error
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")

	m := TemplateModel{}
	m.Measurements = []MeasurementType{
		{Name: "hp", Label: "Hp", Units: "%"},
		{Name: "cp", Label: "Cp", Units: "%"},
		{Name: "sp", Label: "Sp", Units: "%"},
		{Name: "np", Label: "Np", Units: "%"},
		{Name: "op", Label: "Op", Units: "%"},
		{Name: "wp", Label: "Wp", Units: "%"},
		{Name: "ap", Label: "Ap", Units: "%"},
	}

	if r.Method != "POST" {
		_ = tmpl.Execute(w, m)
		return
	}

	err := r.ParseForm()
	if err != nil {
		return
	}

	for i := range m.Measurements {
		m.Measurements[i].Value = r.FormValue(m.Measurements[i].Name)
	}

	if res, err := Calculate(m.Measurements); err != nil {
		m.Error = err
	} else {
		m.CalcResult = res
	}

	_ = tmpl.Execute(w, m)
}

func Calculate(measurements []MeasurementType) ([]MeasurementType, error) {
	i := make(map[string]float64)
	var o []MeasurementType

	for _, m := range measurements {
		if m.Value == "" {
			return o, fmt.Errorf("поле \"%s\" не заповнене", m.Label)
		}

		if val, err := strconv.ParseFloat(m.Value, 64); err != nil {
			return o, fmt.Errorf("поле \"%s\" містить невірне значення", m.Label)
		} else {
			i[m.Name] = val
		}

	}

	hp := i["hp"]
	cp := i["cp"]
	sp := i["sp"]
	np := i["np"]
	op := i["op"]
	wp := i["wp"]
	ap := i["ap"]

	const delta = 0.01
	if math.Abs(hp+cp+sp+np+op+wp+ap-100) > delta {
		return o, fmt.Errorf("сума введених значень повинна дорівнювати 100")
	}

	kpc := 100 / (100 - wp)
	kpg := 100 / (100 - wp - ap)

	hc := hp * kpc
	cc := cp * kpc
	sc := sp * kpc
	nc := np * kpc
	oc := op * kpc
	ac := ap * kpc

	hg := hp * kpg
	cg := cp * kpg
	sg := sp * kpg
	ng := np * kpg
	og := op * kpg

	qrn := 339*cp + 1030*hp - 108.8*(op-sp) - 25*wp
	qsn := (qrn + 25*wp) * 100 / (100 - wp)
	qgn := (qrn + 25*wp) * 100 / (100 - wp - ap)

	o = []MeasurementType{
		{Name: "kpc", Label: "Qрс", Units: "", Value: fmt.Sprintf("%.2f", kpc)},
		{Name: "kpg", Label: "Qрг", Units: "", Value: fmt.Sprintf("%.2f", kpg)},
		{Name: "hc", Label: "Hc", Units: "%", Value: fmt.Sprintf("%.2f", hc)},
		{Name: "cc", Label: "Cc", Units: "%", Value: fmt.Sprintf("%.2f", cc)},
		{Name: "sc", Label: "Sc", Units: "%", Value: fmt.Sprintf("%.2f", sc)},
		{Name: "nc", Label: "Nc", Units: "%", Value: fmt.Sprintf("%.2f", nc)},
		{Name: "oc", Label: "Oc", Units: "%", Value: fmt.Sprintf("%.2f", oc)},
		{Name: "ac", Label: "Ac", Units: "%", Value: fmt.Sprintf("%.2f", ac)},
		{Name: "hg", Label: "Hг", Units: "%", Value: fmt.Sprintf("%.2f", hg)},
		{Name: "cg", Label: "Cг", Units: "%", Value: fmt.Sprintf("%.2f", cg)},
		{Name: "sg", Label: "Sг", Units: "%", Value: fmt.Sprintf("%.2f", sg)},
		{Name: "ng", Label: "Nг", Units: "%", Value: fmt.Sprintf("%.2f", ng)},
		{Name: "og", Label: "Oг", Units: "%", Value: fmt.Sprintf("%.2f", og)},
		{Name: "qrn", Label: "Qрн", Units: "КДж/кг", Value: fmt.Sprintf("%.2f", qrn)},
		{Name: "qsn", Label: "Qсн", Units: "КДж/кг", Value: fmt.Sprintf("%.2f", qsn)},
		{Name: "qgn", Label: "Qгн", Units: "КДж/кг", Value: fmt.Sprintf("%.2f", qgn)},
	}

	return o, nil
}

func main() {
	http.HandleFunc("/", IndexHandler)

	fmt.Println("Server is listening...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}
