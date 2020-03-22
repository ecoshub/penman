package penman

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type ReadLine struct {
	file   *os.File
	active bool
	offset int
	scan   *bufio.Scanner
}

/*
usage;
rl := Reader(dir)
temp := rl.Next()
for temp != nil {
	fmt.Printf(temp)
	temp = rl.Next()
}
*/

func Reader(dir string) (*ReadLine, error) {
	_, filedir := SplitDir(dir)
	file, err := os.Open(filedir)
	if err != nil {
		return nil, err
	}
	rl := ReadLine{file: file, scan: bufio.NewScanner(file), active: true}
	return &rl, nil
}

func (r *ReadLine) Next() []byte {
	if !r.active {
		return nil
	}
	if r.scan.Scan() {
		buffer := r.scan.Bytes()
		r.offset += len(buffer) + len(NewLine())
		return buffer
	}
	r.Close()
	return nil
}

func (r *ReadLine) Close() {
	r.active = false
	r.file.Close()
}

func Read(dir string) []byte {
	dir = PreProcess(dir)
	_, filedir := SplitDir(dir)
	buff, err := ioutil.ReadFile(filedir)
	if err != nil {
		fmt.Printf("Read File Error:%v\n", err)
	} else {
		return buff
	}
	return []byte{}
}

func SRead(dir string) string {
	return string(Read(dir))
}

func ReadAt(dir string, offset int64, length int64) []byte {
	dir = PreProcess(dir)
	f, err := os.Open(dir)
	if err != nil {
		fmt.Println("File Open Error:", err)
	} else {
		defer f.Close()
	}
	data := make([]byte, length)
	_, err = f.Seek(offset, 0)
	if err != nil {
		fmt.Println("Seeker Error:", err)
	}
	_, err = f.Read(data)
	if err != nil {
		fmt.Println("Read Error:", err)
	}
	return data
}

// Write
// if file exist append end
// curr prefix not lower-upper key senstive
// dir: curr\new_folder\new_text.txt is current directory
// desk prefix not lower-upper key senstive
// dir: desk\new_folder\new_text.txt is desktop directory
func Write(dir string, buff []byte) {
	dir = PreProcess(dir)
	newdir, newfile := SplitDir(dir)
	err := os.MkdirAll(newdir, os.ModePerm)
	if err != nil {
		fmt.Println("Make Directory Error:", err)
	} else {
		if IsFileExist(newfile) {
			// apppend
			appendFile(newfile, buff)
		} else {
			// create
			writeFile(newfile, buff)
		}
	}
}

// main write function
func writeFile(filedir string, buffer []byte) {
	err := ioutil.WriteFile(filedir, buffer, os.ModePerm)
	if err != nil {
		fmt.Printf("File Write Error:%v\n", err)
	}
}

// main append function
func appendFile(filedir string, buff []byte) {
	f, err := os.OpenFile(filedir, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Printf("File Open Error:%v\n", err)
	}
	defer f.Close()
	if _, err = f.Write(buff); err != nil {
		fmt.Printf("File Write Error:%v\n", err)
	}
}

// string write
func SWrite(dir string, data string) {
	Write(dir, []byte(data))
}

//string over write
func SOWrite(dir string, data string) {
	OWrite(dir, []byte(data))
}

// string writeln
func SWriteln(dir string, data string) {
	Write(dir, []byte(data+NewLine()))
}

// ReWrite
// if file exist over write else create and write
// curr prefix not lower-upper key senstive
// dir: curr\new_folder\new_text.txt is current directory
// desk prefix not lower-upper key senstive
// dir: desk\new_folder\new_text.txt is desktop directory
func OWrite(dir string, buff []byte) {
	dir = PreProcess(dir)
	newdir, newfile := SplitDir(dir)
	err := os.MkdirAll(newdir, os.ModePerm)
	if err != nil {
		fmt.Println("Make Directory Error:", err)
	} else {
		writeFile(newfile, buff)
	}
}

func GetLineHas(dir, key string) (int64, int) {
	dir = PreProcess(dir)
	file := SRead(dir)
	tokens := strings.Split(file, NewLine())
	count := int64(0)
	lennl := len(NewLine())
	for _, v := range tokens {
		if strings.Contains(v, key) {
			return int64(count), len(v)
		}
		count += int64(len(v) + lennl)
	}
	return int64(-1), 0
}

func GetLineHasAll(dir, key string) ([]int64, []int) {
	dir = PreProcess(dir)
	file := SRead(dir)
	tokens := strings.Split(file, NewLine())
	count := int64(0)
	lennl := len(NewLine())
	offsets := make([]int64, 0, 1024)
	lens := make([]int, 0, 1024)
	for _, v := range tokens {
		if strings.Contains(v, key) {
			offsets = append(offsets, int64(count))
			lens = append(lens, len(v))
		}
		count += int64(len(v) + lennl)
	}
	return offsets, lens
}

func UpdateLine(dir, key, newval string) {
	dir = PreProcess(dir)
	file := SRead(dir)
	tokens := strings.Split(file, NewLine())
	for i, v := range tokens {
		if strings.Contains(v, key) {
			tokens[i] = newval
		}
	}
	SOWrite(dir, strings.Join(tokens, NewLine()))
}

func UpdateLineWithOffset(dir string, offset int64, length int, newval string) {
	dir = PreProcess(dir)
	file := SRead(dir)
	tokens := strings.Split(file, NewLine())
	count := int64(0)
	lennl := len(NewLine())
	for i, v := range tokens {
		if count == offset {
			tokens[i] = newval
		}
		count += int64(len(v) + lennl)
	}
	SOWrite(dir, strings.Join(tokens, NewLine()))
}

func DeleteLineWithOffset(dir string, offset int64, length int) {
	dir = PreProcess(dir)
	file := SRead(dir)
	tokens := strings.Split(file, NewLine())
	count := int64(0)
	lennl := len(NewLine())
	for i, v := range tokens {
		if count == offset {
			tokens[len(tokens)-1], tokens[i] = tokens[i], tokens[len(tokens)-1]
		}
		count += int64(len(v) + lennl)
	}
	SOWrite(dir, strings.Join(tokens[:len(tokens)-1], NewLine()))
}
