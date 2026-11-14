package database

import (
	"reflect"
	"testing"
)

func TestDividir(t *testing.T) {
	script := `-- comentario
CREATE TABLE A (ID INTEGER);

INSERT INTO A VALUES (1);
SET TERM ^ ;
CREATE TRIGGER T FOR A BEFORE INSERT AS
BEGIN
  NEW.ID = 1;
END^
SET TERM ; ^
DROP TABLE A;
`
	got, err := dividir(script)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"CREATE TABLE A (ID INTEGER)",
		"INSERT INTO A VALUES (1)",
		"CREATE TRIGGER T FOR A BEFORE INSERT AS\nBEGIN\n  NEW.ID = 1;\nEND",
		"DROP TABLE A",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("sentencias:\n got  %q\n want %q", got, want)
	}
}

func TestDividirSinTerminar(t *testing.T) {
	if _, err := dividir("CREATE TABLE A (ID INTEGER)"); err == nil {
		t.Fatal("esperaba error por sentencia sin terminar")
	}
}

func TestMigracionesEmbebidas(t *testing.T) {
	ms, err := cargar()
	if err != nil {
		t.Fatal(err)
	}
	if len(ms) == 0 || ms[0].version != "0001" {
		t.Fatalf("esperaba 0001 como primera migración, obtuve %+v", ms)
	}
}
