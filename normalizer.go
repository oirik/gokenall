package gokenall

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/width"
)

type normalizer struct {
	inputs               []*JapanZipCode
	outputs              []*JapanZipCode
	clearStreetReg       *regexp.Regexp
	streetBlessReg       *regexp.Regexp
	streetKanaBlessReg   *regexp.Regexp
	streetInnerBlessReg1 *regexp.Regexp
	streetInnerBlessReg2 *regexp.Regexp
	streetInnerBlessReg3 *regexp.Regexp
	streetInnerBlessReg4 *regexp.Regexp
	streetInnerBlessReg5 *regexp.Regexp
}

const (
	clearStreetRegExp       = `(^以下に掲載がない場合$|の次に番地がくる場合$|.+一円$)`
	streetBlessRegExp       = `（([^）]+)）$`
	streetKanaBlessRegExp   = `\(([^\)]+)\)$`
	streetInnerBlessRegExp1 = `^(その他|地階・階層不明|.*を除く)$`
	streetInnerBlessRegExp2 = `^([０１２３４５６７８９]+)階$`
	streetInnerBlessRegExp3 = `^([０１２３４５６７８９]+)～([０１２３４５６７８９]+)(丁目|番地|番)$`
	streetInnerBlessRegExp4 = `^([０１２３４５６７８９、]+)(丁目|番地|番)$`
	streetInnerBlessRegExp5 = `^[^「」～－０１２３４５６７８９]+$`
)

const (
	maxDevideNum = 99
)

var zen2hanMap = map[string]string{
	"丁目": "ﾁｮｳﾒ",
	"番地": "ﾊﾞﾝﾁ",
	"番":  "ﾊﾞﾝ",
}

func newNormalizer() *normalizer {
	return &normalizer{
		inputs:               make([]*JapanZipCode, 0),
		outputs:              make([]*JapanZipCode, 0),
		clearStreetReg:       regexp.MustCompile(clearStreetRegExp),
		streetBlessReg:       regexp.MustCompile(streetBlessRegExp),
		streetKanaBlessReg:   regexp.MustCompile(streetKanaBlessRegExp),
		streetInnerBlessReg1: regexp.MustCompile(streetInnerBlessRegExp1),
		streetInnerBlessReg2: regexp.MustCompile(streetInnerBlessRegExp2),
		streetInnerBlessReg3: regexp.MustCompile(streetInnerBlessRegExp3),
		streetInnerBlessReg4: regexp.MustCompile(streetInnerBlessRegExp4),
		streetInnerBlessReg5: regexp.MustCompile(streetInnerBlessRegExp5),
	}
}

func (normer *normalizer) push(p *JapanZipCode) {
	normer.inputs = append(normer.inputs, p)
	normer.normalize()
}

func (normer *normalizer) canPop() bool {
	return len(normer.outputs) > 0
}

func (normer *normalizer) pop() *JapanZipCode {
	if len(normer.outputs) == 0 {
		return nil
	}
	ret := normer.outputs[0]
	normer.outputs = normer.outputs[1:]
	return ret
}

func (normer *normalizer) normalize() {
	if len(normer.inputs) == 0 {
		return
	}

	input := normer.inputs[0]
	if input.isMultiLineStart() {
		if ok := normer.normalizeMulti(); !ok {
			return
		}
		input = normer.inputs[0]
	}
	normer.inputs = normer.inputs[1:]

	outputs := normer.normalizeStreet(input)

	normer.outputs = append(normer.outputs, outputs...)
	return
}

func (normer *normalizer) normalizeStreet(input *JapanZipCode) []*JapanZipCode {
	outputs := []*JapanZipCode{input}

	if normer.clearStreetReg.FindString(outputs[0].Street) != "" {
		outputs[0].Street = ""
		outputs[0].StreetKana = ""
	} else if matches := normer.streetBlessReg.FindStringSubmatch(outputs[0].Street); matches != nil {

		var innerBless, innerBlessKana string
		innerBless = matches[1]

		matchesKana := normer.streetKanaBlessReg.FindStringSubmatch(outputs[0].StreetKana)
		if len(matchesKana) > 1 {
			innerBlessKana = matchesKana[1]
		}

		outputs[0].Street = normer.streetBlessReg.ReplaceAllString(outputs[0].Street, "")
		outputs[0].StreetKana = normer.streetKanaBlessReg.ReplaceAllString(outputs[0].StreetKana, "")

		if normer.streetInnerBlessReg1.FindString(innerBless) != "" { // `^(その他|地階・階層不明|.*を除く)$`
			// Do nothing
		} else if matches := normer.streetInnerBlessReg2.FindStringSubmatch(innerBless); matches != nil { // `^([０１２３４５６７８９]+)階$`
			outputs[0].Street += matches[1] + "階"
			outputs[0].StreetKana += width.Narrow.String(matches[1]) + "ｶｲ"
		} else if matches := normer.streetInnerBlessReg3.FindStringSubmatch(innerBless); matches != nil { // `^([０１２３４５６７８９]+)～([０１２３４５６７８９]+)(丁目|番地|番)$`
			start := zenkaku2Int(matches[1])
			end := zenkaku2Int(matches[2])
			if (end - start + 1) <= maxDevideNum {
				for i := start; i <= end; i++ {
					ad := *outputs[0]
					ad.Street += int2Zenkaku(i) + matches[3]
					ad.StreetKana += fmt.Sprintf("%d%s", i, zen2hanMap[matches[3]])
					outputs = append(outputs, &ad)
				}
				outputs = outputs[1:]
			}
		} else if matches := normer.streetInnerBlessReg4.FindStringSubmatch(innerBless); matches != nil { // `^([０１２３４５６７８９、]+)(丁目|番地|番)$`
			splits := strings.Split(matches[1], "、")
			for _, s := range splits {
				ad := *outputs[0]
				ad.Street += s + matches[2]
				ad.StreetKana += fmt.Sprintf("%d%s", zenkaku2Int(s), zen2hanMap[matches[2]])
				outputs = append(outputs, &ad)
			}
			outputs = outputs[1:]
		} else if normer.streetInnerBlessReg5.FindString(innerBless) != "" { // `^[^「」～－０１２３４５６７８９]+$`
			splits := strings.Split(innerBless, "、")
			splitsKana := strings.Split(innerBlessKana, "､")
			if len(splits) == len(splitsKana) {
				for i := range splits {
					ad := *outputs[0]
					ad.Street += splits[i]
					ad.StreetKana += splitsKana[i]
					outputs = append(outputs, &ad)
				}
				outputs = outputs[1:]
			}
		} else {
			//fmt.Println(innerBless)
		}
	}

	return outputs
}

func (normer *normalizer) normalizeMulti() bool {
	endIndex := 0
	for i, p := range normer.inputs {
		if p.isMultiLineEnd() {
			endIndex = i
			break
		}
	}
	if endIndex == 0 {
		return false
	}

	var street, streetKana strings.Builder
	for i := 0; i <= endIndex; i++ {
		street.WriteString(normer.inputs[i].Street)
		streetKana.WriteString(normer.inputs[i].StreetKana)
	}

	normer.inputs[endIndex].Street = street.String()
	normer.inputs[endIndex].StreetKana = streetKana.String()

	normer.inputs = normer.inputs[endIndex:]

	return true
}

func zenkaku2Int(t string) int {
	i, err := strconv.ParseInt(width.Narrow.String(t), 10, 32)
	if err != nil {
		panic(err)
	}
	return int(i)
}

func int2Zenkaku(i int) string {
	return width.Widen.String(strconv.FormatInt(int64(i), 10))
}
