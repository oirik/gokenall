package gokenall

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/csv"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"golang.org/x/text/width"
)

const (
	kenAllSiteURL = "https://www.post.japanpost.jp/zipcode/dl/kogaki-zip.html"
	kenAllFileURL = "https://www.post.japanpost.jp/zipcode/dl/kogaki/zip/ken_all.zip"
)

// Download downloads ken_all file from japanpost website.
// The file on website is zip archived.
// If extract flag sets true, the file is decompressed to csv file.
func Download(w io.Writer, extract bool) error {
	resp, err := http.Get(kenAllFileURL)
	if err != nil {
		return errors.Wrapf(err, "failed to download url: %s", kenAllFileURL)
	}
	defer resp.Body.Close()

	if !extract {
		if _, err = io.Copy(w, resp.Body); err != nil {
			return errors.Wrap(err, "failed to copy from reader to writer")
		}
		return nil
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read http body")
	}

	zipReader, err := zip.NewReader(bytes.NewReader(bodyBytes), resp.ContentLength)
	if err != nil {
		return errors.Wrap(err, "failed to allocate reader")
	}
	if len(zipReader.File) != 1 {
		return errors.Errorf("downloaded zip file does not contain 1 file but %d files", len(zipReader.File))
	}

	srcFile, err := zipReader.File[0].Open()
	if err != nil {
		return errors.Wrapf(err, "failed to open the decompress file in zip file: %s", zipReader.File[0].Name)
	}
	defer srcFile.Close()

	if _, err = io.Copy(w, srcFile); err != nil {
		return errors.Wrap(err, "failed to copy from decompress file in zip file to writer")
	}

	return nil
}

// Updated checks whether the file on japanpost website is updated or not
// by comparing with the date text in website.
func Updated(compareDate time.Time) (result bool, updatedDate time.Time, retErr error) {
	resp, err := http.Get(kenAllSiteURL)
	if err != nil {
		retErr = errors.Wrapf(err, "failed to download url: %s", kenAllSiteURL)
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		retErr = errors.Wrapf(err, "failed to read contents of url: %s", kenAllSiteURL)
		return
	}

	reg := regexp.MustCompile(`<small>(\d+年\d+月\d+日)更新</small>`)
	matches := reg.FindStringSubmatch(string(bytes))
	if len(matches) == 0 {
		retErr = errors.New("not found updated date string in website")
		return
	}
	updatedDate, err = time.Parse("2006年1月2日", matches[1])
	if err != nil {
		retErr = errors.Wrapf(err, "failed to parse updated date string: %s", matches[1])
		return
	}

	result = compareDate.Before(updatedDate)
	return
}

// NormalizeOption is the condition flags at normalize.
type NormalizeOption uint

const (
	// NormalizeWidth is set if you want to align the letter format.
	NormalizeWidth NormalizeOption = 1 << iota
	// NormalizeUTF8 is set if you want to convert from sjis to UTF8.
	NormalizeUTF8
	// NormalizeTrim is set if you want to trim the text.
	NormalizeTrim
	// bitsNormalizeOption is a number of normalize options.
	bitsNormalizeOption = iota
	// NoNormalizeOption represents no flag is set for normalize.
	NoNormalizeOption = NormalizeOption(0)
	// AllNormalizeOption represents all flags are set for normalize.
	AllNormalizeOption = NormalizeOption(1<<bitsNormalizeOption - 1)
	// DefaultNormalizeOption represents default flags are set for normalize.
	DefaultNormalizeOption = AllNormalizeOption
)

// Normalize make original ken_all texts easy to use.
// See detail information in https://github.com/oirik/gokenall.
// Optionaly change width / encoding / trim. (default true for all)
func Normalize(r io.Reader, w io.Writer, option NormalizeOption) error {
	var inputLines, outputLines int

	csvReader := csv.NewReader(transform.NewReader(r, japanese.ShiftJIS.NewDecoder()))
	csvReader.ReuseRecord = true

	var writer *bufio.Writer
	if option&NormalizeWidth == 0 {
		if option&NormalizeUTF8 == 0 {
			writer = bufio.NewWriter(transform.NewWriter(w, japanese.ShiftJIS.NewEncoder()))
		} else {
			writer = bufio.NewWriter(w)
		}
	} else {
		if option&NormalizeUTF8 == 0 {
			writer = bufio.NewWriter(transform.NewWriter(w, transform.Chain(norm.NFD, width.Fold, norm.NFC, japanese.ShiftJIS.NewEncoder())))
		} else {
			writer = bufio.NewWriter(transform.NewWriter(w, transform.Chain(norm.NFD, width.Fold, norm.NFC)))
		}
	}

	normer := newNormalizer()

	for {

		inputLines++

		var input *JapanZipCode
		{
			cols, err := csvReader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				return errors.Wrap(err, "failed to read csv")
			}
			input, err = parseArray(cols, (option&NormalizeTrim != 0))
			if err != nil {
				return errors.Wrapf(err, "failed to parse data: input-line=%d", inputLines)
			}
		}

		normer.push(input)
		for normer.canPop() {

			output := normer.pop()
			outputCSV := output.revertCSV()

			if outputLines > 0 {
				outputCSV = "\n" + outputCSV
			}
			outputLines++

			_, err := writer.WriteString(outputCSV)
			if err != nil {
				return errors.Wrapf(err, "failed to write string to output: input-line=%d output-line=%d", inputLines, outputLines)
			}
		}
	}
	writer.Flush()
	return nil
}

// Parse parses input csv texts to JapanZipCode data structure.
func Parse(r io.Reader) ([]*JapanZipCode, error) {
	var inputLines int
	list := []*JapanZipCode{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		inputLines++
		p, err := parseCSV(scanner.Text(), false)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse line: input-line=%d", inputLines)
		}
		list = append(list, p)
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to read scanner: input-line=%d", inputLines)
	}
	return list, nil
}
