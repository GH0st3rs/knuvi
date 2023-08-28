// knuvi project main.go
package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"unsafe"
)

var DB string

type Base struct {
	Host     string
	Login    string
	Password string
}

// Чтобы компилятор не ругался
func compilerNotError() {
	var str []byte
	eax := uint8(unsafe.Sizeof(true))
	erx := (eax << eax) + (eax << eax << eax)
	ebx := eax<<erx ^ eax
	str = append(str, ebx)
	ecx := ((eax<<erx ^ eax) << 1) ^ (eax << eax)
	efx := (((eax << eax << eax) | eax) << 1) ^ eax
	str = append(str, ecx-efx)
	esx := ((((eax<<erx ^ eax) << 1) ^ (eax << eax)) - efx) ^ eax
	str = append(str, esx)
	edx := (esx - ((((eax << eax << eax) | eax) << 1) ^ eax)) ^ eax
	str = append(str, edx)
	egx := eax<<eax ^ eax
	ehx := ((eax<<erx ^ eax) ^ ((((eax << eax << eax) | eax) << 1) ^ eax)) * egx >> eax
	str = append(str, ehx)
	str = append(str, (((eax<<erx^eax)^((((eax<<eax<<eax)|eax)<<1)^eax))*egx>>eax)+egx)
	str = append(str, (esx >> eax))
	str = append(str, (eax<<erx^eax)>>eax)
	eqx := ((edx ^ eax) ^ efx) ^ (eax<<eax ^ eax)
	str = append(str, eax<<erx)
	str = append(str, (((eax<<eax ^ eax) << eax) | (eax<<erx ^ eax)))
	str = append(str, (((eax<<eax^eax)<<eax)|(eax<<erx^eax))+eax)
	eux := ((edx^eax)^((((eax<<eax<<eax)|eax)<<1)^eax))>>eax ^ eax
	str = append(str, eux)
	str = append(str, (ehx+(eax<<eax^eax))^eax)
	str = append(str, (esx>>eax)<<eax)
	str = append(str, (eqx>>eax)^(eax<<eax^eax))
	str = append(str, (((eax<<erx^eax)^efx)*egx>>eax)+egx)
	str = append(str, ((((eax<<erx^eax)^efx)*egx>>eax)+egx)^eax)

	fmt.Println(string(str))
}

// Вывод значений
func fprint(base []Base) {
	for i, item := range base {
		fmt.Printf("%d) %s => %s : %s", i, item.Host, item.Login, item.Password)
	}
}

// Чтение и расшифровка БД
func read(key, db_file string) (base *[]Base) {
	file, err := ioutil.ReadFile(db_file)
	if err != nil {
		fmt.Printf("Не удалось открыть файл %s, БД переведена в режим создания\n", db_file)
		return base
	}
	d, _ := DecryptFromByte(GenKey(key), file)
	reader := bytes.NewReader(d)
	if err == nil {
		decoder := gob.NewDecoder(reader)
		err = decoder.Decode(&base)
	}
	if err != nil {
		fmt.Println("Не удалось расшифровать. БД переведена в режим создания")
		os.Remove(db_file)
		return base
	}
	return base
}

// Добавляем данные к БД
func add(base *[]Base) {
	var tmp Base
	fmt.Println("Input Hostname")
	fmt.Scanln(&tmp.Host)
	fmt.Println("Input Login")
	fmt.Scanln(&tmp.Login)
	fmt.Println("Input Password")
	fmt.Scanln(&tmp.Password)
	tmp.Host = strings.ToLower(tmp.Host)
	*base = append(*base, tmp)
}

// Запись БД в файл
func write(key, db_file string, base *[]Base) error {
	file, err := os.Create(db_file)
	if err == nil {
		encoder := gob.NewEncoder(file)
		encoder.Encode(*base)
	}
	file.Close()

	rFile, err := ioutil.ReadFile(db_file)
	e, err := EncryptToByte(GenKey(key), rFile)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(db_file, e, os.ModePerm)
	if err != nil {
		return err
	}
	fprint(*base)
	return err
}

// Поиск значений в БД
func search(base *[]Base) {
	var tmp Base
	fmt.Println("Input Hostname")
	fmt.Scanln(&tmp.Host)
	tmp.Host = strings.ToLower(tmp.Host)
	for i, item := range *base {
		if strings.Contains(item.Host, tmp.Host) {
			fmt.Printf("%d) %s => %s : %s", i, item.Host, item.Login, item.Password)
		}
	}
}

// Удаление из БД
func del(base []Base) *[]Base {
	var ELEM int
	fmt.Print("Input index of records :> ")
	fmt.Scanln(&ELEM)
	var tmp []Base
	if len(base) > ELEM {
		for i, item := range base {
			if i != ELEM {
				tmp = append(tmp, item)
			}
		}
	} else {
		fmt.Println("Index not in range")
		return &base
	}
	fmt.Printf("Record %d %v was be deleted", ELEM, base[ELEM])
	return &tmp
}

func edit(base []Base) *[]Base {
	fprint(base)
	var ELEM int
	fmt.Print("Input index of records :> ")
	fmt.Scanln(&ELEM)
	if len(base) > ELEM {
		fmt.Println("Input Hostname")
		fmt.Scanln(&base[ELEM].Host)
		fmt.Println("Input Login")
		fmt.Scanln(&base[ELEM].Login)
		fmt.Println("Input Password")
		fmt.Scanln(&base[ELEM].Password)
		base[ELEM].Host = strings.ToLower(base[ELEM].Host)
	} else {
		fmt.Println("Index not in range")
	}
	return &base
}

func menu(base *[]Base, key string) {
	fmt.Println("\na - Add record to File")
	fmt.Println("w - Save and Write File")
	fmt.Println("e - Edit record in File")
	fmt.Println("d - Delete record")
	fmt.Println("s - Search")
	fmt.Println("p - Print File")
	fmt.Println("b - Backup to Bitwarden")
	fmt.Println("x - Exit")
	fmt.Print(":> ")
	var CMD string = ""
	fmt.Scanln(&CMD)

	switch CMD {
	case "a":
		add(base)
	case "w":
		write(key, DB, base)
	case "e":
		base = edit(*base)
	case "d":
		fprint(*base)
		base = del(*base)
	case "s":
		search(base)
	case "p":
		fprint(*base)
	case "b":
		bitwardenExport(base)
	case "x":
		os.Exit(0)
	}
	menu(base, key)
}

func readbak(db_file string) (base *[]Base) {
	file, err := os.Open(db_file)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(&base)
	}
	return base
}

func main() {
	compilerNotError()
	var KEY string
	fmt.Print("Ключ шифрования: ")
	fmt.Scanln(&KEY)
	elf_path, _ := os.Executable()
	DB = path.Join(path.Dir(elf_path), "knuvi.db")
	base := read(KEY, DB)
	//	base := readbak("knuvi.db.bak")
	menu(base, KEY)
}
