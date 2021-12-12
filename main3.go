package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

var normalJson = `
	{
		"a1":"1",
		"a2":{"b1":"123"}
	}
`

var normalMap = make(map[string]interface{}, 2)

type param struct {
	Required bool
	Min      int
	Max      int
}

var testObject = map[string]param{
	"a1":    param{Required: true, Min: 0, Max: 10},
	"a2.b1": param{Required: true, Min: 1, Max: 7},
}

func main() {
	json.Unmarshal([]byte(normalJson), &normalMap)
	TestSimpleItem()
}

// 正常系のjsonをmapにする
// mapからテスト対象の項目を取り込んで、
// 取り込んだparamでvalition
// 行う

// テスト対象項目を探し、
// 自動生成した値でセットして、テスト実施

// 必須の場合
// 	reqired true
// 任意の場合
// 	reqired false

// min　0の場合
// 	データなし　同値　正常
// 	1桁のデータ作成　超過　正常
// min 1の場合
// 	未満　なしのデータ　　異常
// 	同値　1桁のデータ作成　　正常
// 	超過　2桁のデータ作成	正常
// mxn 0の場合
// 	データなし　同値　正常
// 	1桁のデータ作成　超過　正常
// mxn 10の場合
// 	未満　9桁のデータ　　正常
// 	同値　10桁のデータ作成　　正常
// 	超過　11桁のデータ作成	異常

// min:0 ,max:0
// min:1, max:0
// min:0, max:10
// min:1, max:10
// min:1, max:1
func executeValidation() {
	// 正常系確認
	TestNormal()
	// 確認
	TestSimpleItem()
	return
}

func TestNormal() {
	// 正常系のjsonをmapにする
	return
}

func TestSimpleItem() {
	// mapからテスト対象の項目を取り込んで、
	for k, p := range testObject {
		// 項目分解
		objName := sliceObject(k)
		// paramでテストデータ作成
		// validationData := generateJsonTestData(objName, p)
		fmt.Println(p)
		m := searchMap(objName, &normalMap, 0)

		fmt.Println(m)
		// validation(validationData.testJsonData)
		// assert.Equal(t, validationData.hasErr, true)
	}

	return
}

type fieldName struct {
	splitedName []string
	name        string
}

func sliceObject(obj string) fieldName {
	splitedName := strings.Split(obj, ".")
	return fieldName{
		splitedName: splitedName,
		name:        obj,
	}
}

type checkData struct {
	testJsonData string
	hasErr       bool
}
type expectData struct {
	input  interface{}
	hasErr bool
}

// func generatorTestJsonData(objName fieldName, p param) (result []checkData) {
// 	// paramでテストデータ作成
// 	testData := generateDataString(p)
// 	for _, v := range testData {
// 		// 生成したテストデータをobjnameにセット
// 		result = append(result, (objName, v))
// 	}
// 	// 生成したjsonを返却
// 	return result
// }

// 値のタイプ
// paramでテストデータを生成
// json 生成

func generateJsonTestData(objName fieldName, p param) (result []checkData) {
	currentObjName := objName.splitedName[0]

	// oneLayer
	for key, _ := range normalMap {
		if key == currentObjName {
			testData := generateData(normalMap[key], p)
			for _, v := range testData {
				normalMap[key] = v.input
				//json生成
				testJsonData, err := json.Marshal(normalMap)
				if err != nil {
					fmt.Println(err)
				}
				result = append(result, checkData{testJsonData: string(testJsonData), hasErr: v.hasErr})
			}
			break
		}
	}

	return result
}

// func setValue()
func searchMap(objName fieldName, m map[string]interface{}, currentIndex int) map[string]interface{} {
	currentName := objName.splitedName[currentIndex]

	if hasField(currentName, m) && isLastField(currentIndex, len(objName.splitedName)-1) {
		return m
	}

	// NextName := objName.splitedName[currentIndex+1]
	return searchMap(objName, m[currentName], currentIndex+1)
}

// func getNextIndexMap(m, NextName) *map[string]interface{} {

// }
func isLastField(currentIndex, endIndex int) bool {
	return currentIndex == endIndex
}

func hasField(name string, m map[string]interface{}) bool {
	for k, _ := range m {
		if name == k {
			return true
		}
	}
	return false
}

func generateData(value interface{}, p param) (result []expectData) {
	switch value.(type) {
	case string:
		return generateDataString(p)
	case int64:
		return generateDataInt64(p)
	// case map[string]interface{}:
	// 	return generateDataObject(p)
	default:
		return result
	}
}

func generateDataString(p param) (result []expectData) {
	// min:0 ,max:0
	if p.Min == 0 && p.Min == p.Max {
		// 	データなし　同値　正常
		result = append(result, expectData{input: "", hasErr: false})
		// 	任意データ作成
		result = append(result, expectData{input: stringDataGenerator(1), hasErr: false})

	} else if p.Required && p.Min != 0 && p.Max == 0 {
		// min:1, max:0
		// min:2, max:0
		result = append(result, expectData{input: stringDataGenerator(p.Min - 1), hasErr: true})
		result = append(result, expectData{input: stringDataGenerator(p.Min), hasErr: false})
		result = append(result, expectData{input: stringDataGenerator(p.Min + 1), hasErr: false})
	} else if p.Min == 0 && p.Max != 0 {
		// min:0, max:10
		result = append(result, expectData{input: "", hasErr: false})
		result = append(result, expectData{input: stringDataGenerator(p.Max - 1), hasErr: false})
		result = append(result, expectData{input: stringDataGenerator(p.Max), hasErr: false})
		result = append(result, expectData{input: stringDataGenerator(p.Max + 1), hasErr: true})
	} else if p.Required && p.Min != 0 && p.Max == p.Min {
		// min:1, max:1
		// min:5, max:5
		result = append(result, expectData{input: stringDataGenerator(p.Min - 1), hasErr: true})
		result = append(result, expectData{input: stringDataGenerator(p.Min), hasErr: false})
		result = append(result, expectData{input: stringDataGenerator(p.Min + 1), hasErr: true})
	} else if p.Required && p.Min != 0 && p.Max != p.Min && p.Min < p.Max {
		// min:1, max:10
		// min:3, max:4
		result = append(result, expectData{input: stringDataGenerator(p.Min - 1), hasErr: true})
		result = append(result, expectData{input: stringDataGenerator(p.Min), hasErr: false})
		result = append(result, expectData{input: stringDataGenerator(p.Max - 1), hasErr: false})
		result = append(result, expectData{input: stringDataGenerator(p.Max), hasErr: false})
		result = append(result, expectData{input: stringDataGenerator(p.Max + 1), hasErr: true})
	} else {
		fmt.Println("登録間違い")
	}
	return result
}

func generateDataInt64(p param) (result []expectData) {
	// min:0 ,max:0
	if p.Min == 0 && p.Min == p.Max {
		// 	任意データ作成
		result = append(result, expectData{input: 0, hasErr: false})
		result = append(result, expectData{input: 1, hasErr: false})
		// } else if p.Min != 0 && p.Max == 0 {
		// 	// min:1, max:0
		// 	// min:2, max:0
		// 	result = append(result, expectData{input: int64DataGenerator(p.Min) + 1, hasErr: true})
		// 	result = append(result, expectData{input: stringDataGenerator(p.Min), hasErr: false})
		// 	result = append(result, expectData{input: stringDataGenerator(p.Min + 1), hasErr: false})
	} else if p.Min == 0 && p.Max != 0 {
		// min:0, max:3　⇒ 998,999,1000
		result = append(result, expectData{input: 0, hasErr: false})
		result = append(result, expectData{input: int64DataGenerator(p.Max) - 1, hasErr: false})
		result = append(result, expectData{input: int64DataGenerator(p.Max), hasErr: false})
		result = append(result, expectData{input: int64DataGenerator(p.Max) + 1, hasErr: true})
		// } else if p.Min != 0 && p.Max == p.Min {
		// 	// min:1, max:1
		// 	// min:5, max:5
		// 	result = append(result, expectData{input: stringDataGenerator(p.Min - 1), hasErr: true})
		// 	result = append(result, expectData{input: stringDataGenerator(p.Min), hasErr: false})
		// 	result = append(result, expectData{input: stringDataGenerator(p.Min + 1), hasErr: true})
		// } else if p.Min != 0 && p.Max != p.Min && p.Min < p.Max {
		// 	// min:1, max:10
		// 	// min:3, max:4
		// 	result = append(result, expectData{input: stringDataGenerator(p.Min - 1), hasErr: true})
		// 	result = append(result, expectData{input: stringDataGenerator(p.Min), hasErr: false})
		// 	result = append(result, expectData{input: stringDataGenerator(p.Max - 1), hasErr: false})
		// 	result = append(result, expectData{input: stringDataGenerator(p.Max), hasErr: false})
		// 	result = append(result, expectData{input: stringDataGenerator(p.Max + 1), hasErr: true})
	} else {
		fmt.Println("数字確認には仕様対象外のパターンがあります")
	}
	return result
}

func int64DataGenerator(digit int) int64 {
	return int64(math.Pow10(digit)) - 1
}

// stringDataGenerator
func stringDataGenerator(digit int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		fmt.Println("can not generator string Data ")
		return ""
	}

	var result string
	for _, v := range b {
		result += string(letters[int(v)%len(letters)])
	}

	return result
}
