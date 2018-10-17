package gokenall

import (
	"reflect"
	"testing"
)

func TestJapanZipCode_isMultiLineStart(t *testing.T) {
	tests := []struct {
		name string
		p    *JapanZipCode
		want bool
	}{
		{"", &JapanZipCode{Street: ""}, false},
		{"", &JapanZipCode{Street: "a"}, false},

		{"hankaku", &JapanZipCode{Street: "()"}, false},
		{"hankaku", &JapanZipCode{Street: "a(a)a"}, false},
		{"hankaku", &JapanZipCode{Street: ")"}, false},
		{"hankaku", &JapanZipCode{Street: "a)"}, false},
		{"hankaku", &JapanZipCode{Street: "("}, true},
		{"hankaku", &JapanZipCode{Street: "a("}, true},
		{"hankaku", &JapanZipCode{Street: "(a"}, true},
		{"hankaku", &JapanZipCode{Street: "a(a"}, true},
		{"hankaku", &JapanZipCode{Street: "()(a"}, true},

		{"zenkaku", &JapanZipCode{Street: "（）"}, false},
		{"zenkaku", &JapanZipCode{Street: "a（a）a"}, false},
		{"zenkaku", &JapanZipCode{Street: "）"}, false},
		{"zenkaku", &JapanZipCode{Street: "a）"}, false},
		{"zenkaku", &JapanZipCode{Street: "（"}, true},
		{"zenkaku", &JapanZipCode{Street: "a（"}, true},
		{"zenkaku", &JapanZipCode{Street: "（a"}, true},
		{"zenkaku", &JapanZipCode{Street: "a（a"}, true},
		{"zenkaku", &JapanZipCode{Street: "（）（a"}, true},

		{"hankaku+zenkaku", &JapanZipCode{Street: "(）"}, false},
		{"hankaku+zenkaku", &JapanZipCode{Street: "a(a）a"}, false},
		{"hankaku+zenkaku", &JapanZipCode{Street: "）"}, false},
		{"hankaku+zenkaku", &JapanZipCode{Street: "a）"}, false},
		{"hankaku+zenkaku", &JapanZipCode{Street: "("}, true},
		{"hankaku+zenkaku", &JapanZipCode{Street: "a("}, true},
		{"hankaku+zenkaku", &JapanZipCode{Street: "(a"}, true},
		{"hankaku+zenkaku", &JapanZipCode{Street: "a(a"}, true},
		{"hankaku+zenkaku", &JapanZipCode{Street: "(）(a"}, true},

		{"zenkaku+hankaku", &JapanZipCode{Street: "（)"}, false},
		{"zenkaku+hankaku", &JapanZipCode{Street: "a（a)a"}, false},
		{"zenkaku+hankaku", &JapanZipCode{Street: ")"}, false},
		{"zenkaku+hankaku", &JapanZipCode{Street: "a)"}, false},
		{"zenkaku+hankaku", &JapanZipCode{Street: "（"}, true},
		{"zenkaku+hankaku", &JapanZipCode{Street: "a（"}, true},
		{"zenkaku+hankaku", &JapanZipCode{Street: "（a"}, true},
		{"zenkaku+hankaku", &JapanZipCode{Street: "a（a"}, true},
		{"zenkaku+hankaku", &JapanZipCode{Street: "（)（a"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.isMultiLineStart(); got != tt.want {
				t.Errorf("JapanZipCode.isMultiLineStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJapanZipCode_isMultiLineEnd(t *testing.T) {
	tests := []struct {
		name string
		p    *JapanZipCode
		want bool
	}{
		{"", &JapanZipCode{Street: ""}, false},
		{"", &JapanZipCode{Street: "a"}, false},

		{"zenkaku", &JapanZipCode{Street: "（）"}, false},
		{"zenkaku", &JapanZipCode{Street: "a（a）a"}, false},
		{"zenkaku", &JapanZipCode{Street: "（"}, false},
		{"zenkaku", &JapanZipCode{Street: "（a"}, false},
		{"zenkaku", &JapanZipCode{Street: "）"}, true},
		{"zenkaku", &JapanZipCode{Street: "a）"}, true},
		{"zenkaku", &JapanZipCode{Street: "）a"}, true},
		{"zenkaku", &JapanZipCode{Street: "a）a"}, true},
		{"zenkaku", &JapanZipCode{Street: "）（）a"}, true},

		{"hankaku+zenkaku", &JapanZipCode{Street: "(）"}, false},
		{"hankaku+zenkaku", &JapanZipCode{Street: "a(a）a"}, false},
		{"hankaku+zenkaku", &JapanZipCode{Street: "("}, false},
		{"hankaku+zenkaku", &JapanZipCode{Street: "(a"}, false},
		{"hankaku+zenkaku", &JapanZipCode{Street: "）"}, true},
		{"hankaku+zenkaku", &JapanZipCode{Street: "a）"}, true},
		{"hankaku+zenkaku", &JapanZipCode{Street: "）a"}, true},
		{"hankaku+zenkaku", &JapanZipCode{Street: "a）a"}, true},
		{"hankaku+zenkaku", &JapanZipCode{Street: "）(）a"}, true},

		{"zenkaku+hankaku", &JapanZipCode{Street: "（)"}, false},
		{"zenkaku+hankaku", &JapanZipCode{Street: "a（a)a"}, false},
		{"zenkaku+hankaku", &JapanZipCode{Street: "（"}, false},
		{"zenkaku+hankaku", &JapanZipCode{Street: "（a"}, false},
		{"zenkaku+hankaku", &JapanZipCode{Street: ")"}, true},
		{"zenkaku+hankaku", &JapanZipCode{Street: "a)"}, true},
		{"zenkaku+hankaku", &JapanZipCode{Street: ")a"}, true},
		{"zenkaku+hankaku", &JapanZipCode{Street: "a)a"}, true},
		{"zenkaku+hankaku", &JapanZipCode{Street: ")（)a"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.isMultiLineEnd(); got != tt.want {
				t.Errorf("JapanZipCode.isMultiLineEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseCSV(t *testing.T) {
	type args struct {
		line string
		trim bool
	}
	tests := []struct {
		name    string
		args    args
		want    *JapanZipCode
		wantErr bool
	}{
		{"", args{`01101,"060  ","0600007","ﾎｯｶｲﾄﾞｳ","ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ","ｷﾀ7ｼﾞｮｳﾆｼ","北海道","札幌市中央区","北七条西",0,0,1,0,0,0`, false}, &JapanZipCode{
			JISCode:                   "01101",
			OldZipCode:                "060  ",
			ZipCode:                   "0600007",
			PrefKana:                  "ﾎｯｶｲﾄﾞｳ",
			CityKana:                  "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ",
			StreetKana:                "ｷﾀ7ｼﾞｮｳﾆｼ",
			Pref:                      "北海道",
			City:                      "札幌市中央区",
			Street:                    "北七条西",
			StreetDuplicateZipCodeFlg: "0",
			NumberedSmallStreetFlg:    "0",
			NumberedStreetFlg:         "1",
			ZipCodeDuplicateStreetFlg: "0",
			UpdateFlg:                 "0",
			UpdateReason:              "0",
			PrefCode:                  "01",
		}, false},
		{"", args{`01101,"060  ","0600007","ﾎｯｶｲﾄﾞｳ","ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ","ｷﾀ7ｼﾞｮｳﾆｼ","北海道","札幌市中央区","北七条西",1,2,3,4,5,6`, true}, &JapanZipCode{
			JISCode:                   "01101",
			OldZipCode:                "060",
			ZipCode:                   "0600007",
			PrefKana:                  "ﾎｯｶｲﾄﾞｳ",
			CityKana:                  "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ",
			StreetKana:                "ｷﾀ7ｼﾞｮｳﾆｼ",
			Pref:                      "北海道",
			City:                      "札幌市中央区",
			Street:                    "北七条西",
			StreetDuplicateZipCodeFlg: "1",
			NumberedSmallStreetFlg:    "2",
			NumberedStreetFlg:         "3",
			ZipCodeDuplicateStreetFlg: "4",
			UpdateFlg:                 "5",
			UpdateReason:              "6",
			PrefCode:                  "01",
		}, false},
		{"", args{`01101,"060  ","0600007","ﾎｯｶｲﾄﾞｳ","ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ","ｷﾀ7ｼﾞｮｳﾆｼ","北海道","札幌市中央区","北七条西",0,0,1,0,0`, false}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCSV(tt.args.line, tt.args.trim)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseArray(t *testing.T) {
	type args struct {
		cols []string
		trim bool
	}
	tests := []struct {
		name    string
		args    args
		want    *JapanZipCode
		wantErr bool
	}{
		{"", args{[]string{"01101", "060  ", "0600007", "ﾎｯｶｲﾄﾞｳ", "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ", "ｷﾀ7ｼﾞｮｳﾆｼ", "北海道", "札幌市中央区", "北七条西", "0", "0", "1", "0", "0", "0"}, false}, &JapanZipCode{
			JISCode:                   "01101",
			OldZipCode:                "060  ",
			ZipCode:                   "0600007",
			PrefKana:                  "ﾎｯｶｲﾄﾞｳ",
			CityKana:                  "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ",
			StreetKana:                "ｷﾀ7ｼﾞｮｳﾆｼ",
			Pref:                      "北海道",
			City:                      "札幌市中央区",
			Street:                    "北七条西",
			StreetDuplicateZipCodeFlg: "0",
			NumberedSmallStreetFlg:    "0",
			NumberedStreetFlg:         "1",
			ZipCodeDuplicateStreetFlg: "0",
			UpdateFlg:                 "0",
			UpdateReason:              "0",
			PrefCode:                  "01",
		}, false},
		{"", args{[]string{"01101", "060  ", "0600007", "ﾎｯｶｲﾄﾞｳ", "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ", "ｷﾀ7ｼﾞｮｳﾆｼ", "北海道", "札幌市中央区", "北七条西", "1", "2", "3", "4", "5", "6"}, true}, &JapanZipCode{
			JISCode:                   "01101",
			OldZipCode:                "060",
			ZipCode:                   "0600007",
			PrefKana:                  "ﾎｯｶｲﾄﾞｳ",
			CityKana:                  "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ",
			StreetKana:                "ｷﾀ7ｼﾞｮｳﾆｼ",
			Pref:                      "北海道",
			City:                      "札幌市中央区",
			Street:                    "北七条西",
			StreetDuplicateZipCodeFlg: "1",
			NumberedSmallStreetFlg:    "2",
			NumberedStreetFlg:         "3",
			ZipCodeDuplicateStreetFlg: "4",
			UpdateFlg:                 "5",
			UpdateReason:              "6",
			PrefCode:                  "01",
		}, false},
		{"", args{[]string{"01101", "060  ", "0600007", "ﾎｯｶｲﾄﾞｳ", "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ", "ｷﾀ7ｼﾞｮｳﾆｼ", "北海道", "札幌市中央区", "北七条西", "0", "0", "1", "0", "0"}, false}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseArray(tt.args.cols, tt.args.trim)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJapanZipCode_revertCSV(t *testing.T) {
	tests := []struct {
		name string
		p    *JapanZipCode
		want string
	}{
		{"", &JapanZipCode{
			JISCode:                   "01101",
			OldZipCode:                "060  ",
			ZipCode:                   "0600007",
			PrefKana:                  "ﾎｯｶｲﾄﾞｳ",
			CityKana:                  "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ",
			StreetKana:                "ｷﾀ7ｼﾞｮｳﾆｼ",
			Pref:                      "北海道",
			City:                      "札幌市中央区",
			Street:                    "北七条西",
			StreetDuplicateZipCodeFlg: "1",
			NumberedSmallStreetFlg:    "2",
			NumberedStreetFlg:         "3",
			ZipCodeDuplicateStreetFlg: "4",
			UpdateFlg:                 "5",
			UpdateReason:              "6",
			PrefCode:                  "01",
		}, `01101,"060  ","0600007","ﾎｯｶｲﾄﾞｳ","ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ","ｷﾀ7ｼﾞｮｳﾆｼ","北海道","札幌市中央区","北七条西",1,2,3,4,5,6`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.revertCSV(); got != tt.want {
				t.Errorf("JapanZipCode.revertCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJapanZipCode_revertArray(t *testing.T) {
	tests := []struct {
		name string
		p    *JapanZipCode
		want []string
	}{
		{"", &JapanZipCode{
			JISCode:                   "01101",
			OldZipCode:                "060  ",
			ZipCode:                   "0600007",
			PrefKana:                  "ﾎｯｶｲﾄﾞｳ",
			CityKana:                  "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ",
			StreetKana:                "ｷﾀ7ｼﾞｮｳﾆｼ",
			Pref:                      "北海道",
			City:                      "札幌市中央区",
			Street:                    "北七条西",
			StreetDuplicateZipCodeFlg: "1",
			NumberedSmallStreetFlg:    "2",
			NumberedStreetFlg:         "3",
			ZipCodeDuplicateStreetFlg: "4",
			UpdateFlg:                 "5",
			UpdateReason:              "6",
			PrefCode:                  "01",
		}, []string{"01101", "060  ", "0600007", "ﾎｯｶｲﾄﾞｳ", "ｻｯﾎﾟﾛｼﾁｭｳｵｳｸ", "ｷﾀ7ｼﾞｮｳﾆｼ", "北海道", "札幌市中央区", "北七条西", "1", "2", "3", "4", "5", "6"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.revertArray(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JapanZipCode.revertArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
