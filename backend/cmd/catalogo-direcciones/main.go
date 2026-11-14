// Comando de desarrollo: genera internal/database/catalogos/direcciones.csv.gz
// con los municipios (INEGI) y los asentamientos con código postal (SEPOMEX)
// de los estados elegidos. La API lo carga al arrancar (database.CargarDirecciones).
//
//	go run ./cmd/catalogo-direcciones -inegi AGEEML_utf8.csv -sepomex hoja1.csv,hoja2.csv
//
// Entradas:
//   - INEGI: catálogo único de municipios (CSV UTF-8 con CVE_ENT, CVE_MUN, NOM_MUN).
//   - SEPOMEX: una hoja por estado exportada a CSV UTF-8 (d_codigo, d_asenta, …).
//
// Salida (CSV, una línea por registro, comprimido con gzip):
//
//	M,<cve_ent>,<cve_mun>,<nombre>                                  municipio
//	A,<cve_ent>,<cve_mun>,<id_asenta>,<cp>,<nombre>,<tipo>,<zona>,<ciudad>   asentamiento
package main

import (
	"compress/gzip"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

func main() {
	inegi := flag.String("inegi", "", "CSV UTF-8 del catálogo de municipios de INEGI")
	sepomex := flag.String("sepomex", "", "CSV UTF-8 de SEPOMEX, separados por coma (una hoja por estado)")
	salida := flag.String("salida", "internal/database/catalogos/direcciones.csv.gz", "archivo de salida")
	flag.Parse()
	if *inegi == "" || *sepomex == "" {
		flag.Usage()
		os.Exit(2)
	}
	if err := generar(*inegi, strings.Split(*sepomex, ","), *salida); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type clave struct{ ent, mun string }

func generar(rutaInegi string, rutasSepomex []string, salida string) error {
	// 1. Asentamientos de SEPOMEX: definen qué estados se incluyen.
	var asentamientos [][]string
	estados := map[string]bool{}
	municipiosUsados := map[clave]bool{}
	vistos := map[string]bool{}
	for _, ruta := range rutasSepomex {
		filas, err := leer(ruta)
		if err != nil {
			return err
		}
		col := columnas(filas[0])
		for _, f := range filas[1:] {
			ent, mun := f[col["c_estado"]], f[col["c_mnpio"]]
			id := f[col["id_asenta_cpcons"]]
			k := ent + mun + id
			if vistos[k] {
				return fmt.Errorf("%s: asentamiento repetido %s", ruta, k)
			}
			vistos[k] = true
			estados[ent] = true
			municipiosUsados[clave{ent, mun}] = true
			asentamientos = append(asentamientos, []string{"A", ent, mun, id,
				f[col["d_codigo"]], limpio(f[col["d_asenta"]]), limpio(f[col["d_tipo_asenta"]]),
				limpio(f[col["d_zona"]]), limpio(f[col["d_ciudad"]])})
		}
	}

	// 2. Municipios de INEGI de esos estados (nombres oficiales).
	filas, err := leer(rutaInegi)
	if err != nil {
		return err
	}
	col := columnas(filas[0])
	var municipios [][]string
	conocidos := map[clave]bool{}
	for _, f := range filas[1:] {
		ent, mun := f[col["CVE_ENT"]], f[col["CVE_MUN"]]
		if !estados[ent] {
			continue
		}
		conocidos[clave{ent, mun}] = true
		municipios = append(municipios, []string{"M", ent, mun, limpio(f[col["NOM_MUN"]])})
	}
	for k := range municipiosUsados {
		if !conocidos[k] {
			return fmt.Errorf("el municipio %s-%s de SEPOMEX no está en INEGI", k.ent, k.mun)
		}
	}

	// Orden estable: el archivo generado no cambia si los datos no cambian.
	sort.Slice(municipios, func(i, j int) bool { return strings.Join(municipios[i], ",") < strings.Join(municipios[j], ",") })
	sort.Slice(asentamientos, func(i, j int) bool {
		a, b := asentamientos[i], asentamientos[j]
		return a[1]+a[2]+a[3] < b[1]+b[2]+b[3]
	})

	out, err := os.Create(salida)
	if err != nil {
		return err
	}
	defer out.Close()
	gz, _ := gzip.NewWriterLevel(out, gzip.BestCompression)
	w := csv.NewWriter(gz)
	if err := w.WriteAll(append(municipios, asentamientos...)); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	fmt.Printf("estados %d · municipios %d · asentamientos %d → %s\n", len(estados), len(municipios), len(asentamientos), salida)
	return nil
}

func leer(ruta string) ([][]string, error) {
	f, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.FieldsPerRecord = -1
	var filas [][]string
	for {
		fila, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%s: %w", ruta, err)
		}
		if len(fila) > 1 && strings.TrimSpace(strings.Join(fila, "")) != "" { // ignora líneas vacías
			filas = append(filas, fila)
		}
	}
	if len(filas) < 2 {
		return nil, fmt.Errorf("%s: sin datos", ruta)
	}
	return filas, nil
}

// columnas: nombre → índice (quita el BOM que deja Excel en la primera).
func columnas(enc []string) map[string]int {
	m := map[string]int{}
	for i, c := range enc {
		m[strings.TrimPrefix(strings.TrimSpace(c), "\ufeff")] = i
	}
	return m
}

func limpio(s string) string { return strings.Join(strings.Fields(s), " ") }
