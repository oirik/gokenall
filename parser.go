package gokenall

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	columnCount = 15
)

// JapanZipCode is a parsed line from ken_all.csv.
type JapanZipCode struct {
	JISCode                   string `json:"jis_code"`      // 全国地方公共団体コード（JIS X0401、X0402）………　半角数字
	OldZipCode                string `json:"old_zip_code"`  // （旧）郵便番号（5桁）………………………………………　半角数字
	ZipCode                   string `json:"zip_code"`      // 郵便番号（7桁）………………………………………　半角数字
	PrefKana                  string `json:"pref_kana"`     // 都道府県名　…………　半角カタカナ（コード順に掲載）　（注1）
	CityKana                  string `json:"city_kana"`     // 市区町村名　…………　半角カタカナ（コード順に掲載）　（注1）
	StreetKana                string `json:"street_kana"`   // 町域名　………………　半角カタカナ（五十音順に掲載）　（注1）
	Pref                      string `json:"pref"`          // 都道府県名　…………　漢字（コード順に掲載）　（注1,2）
	City                      string `json:"city"`          // 市区町村名　…………　漢字（コード順に掲載）　（注1,2）
	Street                    string `json:"street"`        // 町域名　………………　漢字（五十音順に掲載）　（注1,2）
	StreetDuplicateZipCodeFlg string `json:"-"`             // 一町域が二以上の郵便番号で表される場合の表示　（注3）　（「1」は該当、「0」は該当せず）
	NumberedSmallStreetFlg    string `json:"-"`             // 小字毎に番地が起番されている町域の表示　（注4）　（「1」は該当、「0」は該当せず）
	NumberedStreetFlg         string `json:"-"`             // 丁目を有する町域の場合の表示　（「1」は該当、「0」は該当せず）
	ZipCodeDuplicateStreetFlg string `json:"-"`             // 一つの郵便番号で二以上の町域を表す場合の表示　（注5）　（「1」は該当、「0」は該当せず）
	UpdateFlg                 string `json:"update_flg"`    // 更新の表示（注6）（「0」は変更なし、「1」は変更あり、「2」廃止（廃止データのみ使用））
	UpdateReason              string `json:"update_reason"` // 変更理由　（「0」は変更なし、「1」市政・区政・町政・分区・政令指定都市施行、「2」住居表示の実施、「3」区画整理、「4」郵便区調整等、「5」訂正、「6」廃止（廃止データのみ使用））
	PrefCode                  string `json:"pref_code"`     // <ken_allにはない追加項目> 都道府県コード(JIS X0401)
}

func parseCSV(line string, trim bool) (*JapanZipCode, error) {
	csvReader := csv.NewReader(strings.NewReader(line))
	cols, err := csvReader.Read()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read csv format: %s", line)
	}
	return parseArray(cols, trim)
}

func parseArray(cols []string, trim bool) (*JapanZipCode, error) {
	if len(cols) != columnCount {
		return nil, errors.New("Column count is wrong")
	}
	if trim {
		for i := range cols {
			cols[i] = strings.TrimSpace(cols[i])
		}
	}
	p := JapanZipCode{
		JISCode:                   cols[0],
		OldZipCode:                cols[1],
		ZipCode:                   cols[2],
		PrefKana:                  cols[3],
		CityKana:                  cols[4],
		StreetKana:                cols[5],
		Pref:                      cols[6],
		City:                      cols[7],
		Street:                    cols[8],
		StreetDuplicateZipCodeFlg: cols[9],
		NumberedSmallStreetFlg:    cols[10],
		NumberedStreetFlg:         cols[11],
		ZipCodeDuplicateStreetFlg: cols[12],
		UpdateFlg:                 cols[13],
		UpdateReason:              cols[14],
		PrefCode:                  cols[0][:2],
	}

	return &p, nil
}

func (p *JapanZipCode) revertCSV() string {
	cols := p.revertArray()
	for i := range cols {
		if i >= 1 && i <= 8 {
			cols[i] = fmt.Sprintf("\"%s\"", cols[i])
		}
	}
	return strings.Join(cols, ",")
}

func (p *JapanZipCode) revertArray() []string {
	return []string{
		p.JISCode,
		p.OldZipCode,
		p.ZipCode,
		p.PrefKana,
		p.CityKana,
		p.StreetKana,
		p.Pref,
		p.City,
		p.Street,
		p.StreetDuplicateZipCodeFlg,
		p.NumberedSmallStreetFlg,
		p.NumberedStreetFlg,
		p.ZipCodeDuplicateStreetFlg,
		p.UpdateFlg,
		p.UpdateReason,
	}
}

func (p *JapanZipCode) isMultiLineStart() bool {
	oi := strings.LastIndexAny(p.Street, "(（")
	if oi < 0 {
		return false
	}
	ci := strings.LastIndexAny(p.Street, ")）")
	if ci < 0 {
		return true
	}
	return oi > ci
}

func (p *JapanZipCode) isMultiLineEnd() bool {
	ci := strings.IndexAny(p.Street, ")）")
	if ci < 0 {
		return false
	}
	oi := strings.IndexAny(p.Street, "(（")
	if oi < 0 {
		return true
	}
	return ci < oi
}
