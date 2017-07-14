package utils

import (
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"strings"
	"time"
)

func Selector(q ...string) (r bson.M) {
	if len(q) < 1 {
		r = nil
		return
	}
	r = make(bson.M, len(q))
	for _, s := range q {
		r[s] = 1
	}
	return
}

func Deselector(q ...string) (r bson.M) {
	if len(q) < 1 {
		r = nil
		return
	}
	r = make(bson.M, len(q))
	for _, s := range q {
		r[s] = 0
	}
	return
}

func SelDeSel(sel []string, desel []string) (r bson.M) {
	if len(sel)+len(desel) < 1 {
		r = nil
		return
	}
	r = make(bson.M, len(sel)+len(desel))
	for _, s := range sel {
		r[s] = 1
	}
	for _, s := range desel {
		r[s] = 0
	}
	return
}

//
//func UpdateBsonFromMap(mapModel map[string]interface{}) (data bson.M){
//	data = bson.M{}
//	for key, value := range mapModel {
//		//var er error
//		//var rInt int64
//		//rInt, er = strconv.ParseInt(value, 10, 64)
//		//if er == nil {
//		//	data[key] = rInt
//		//	continue
//		//}
//		//var rBool bool
//		//rBool, er = strconv.ParseBool(value)
//		//if er == nil {
//		//	data[key]= rBool
//		//	continue
//		//}
//		//var rFloat float64
//		//rFloat, er = strconv.ParseFloat(value, 64)
//		//if er == nil {
//		//	data[key]= rFloat
//		//	continue
//		//}
//		data[key] = value
//	}
//	data = bson.M{"$set":data}
//	return
//}

func GetBsonFindArray(and []map[string]string, or []map[string]string) (query bson.M) {
	query = bson.M{}
	andArray := []bson.M{}
	for _, obj := range and {
		for key, value := range obj {
			var er error
			var rInt int64
			var opr string = ""
			if strings.HasPrefix(value, ">=") {
				values := strings.Split(value, ">=")
				opr = "$gte"
				value = values[1]
			} else if strings.HasPrefix(value, ">") {
				values := strings.Split(value, ">")
				opr = "$gt"
				value = values[1]
			} else if strings.HasPrefix(value, "<=") {
				values := strings.Split(value, "<=")
				opr = "$lte"
				value = values[1]
			} else if strings.HasPrefix(value, "<") {
				values := strings.Split(value, "<")
				opr = "$lt"
				value = values[1]
			} else if strings.HasPrefix(value, "!=") {
				values := strings.Split(value, "!=")
				opr = "$ne"
				value = values[1]
			} else if strings.HasPrefix(value, "==") {
				values := strings.Split(value, "==")
				opr = "$eq"
				value = values[1]
			}

			layout := "2006-01-02T15:04:05.999Z"
			t, err := time.Parse(layout, value)
			if err == nil {
				if opr == "" {
					andArray = append(andArray, bson.M{key: t})
				} else {
					andArray = append(andArray, bson.M{key: bson.M{opr: t}})
				}
				continue
			}

			if strings.HasPrefix(value, "in:") {
				values := strings.Split(value, ":")
				values = strings.Split(values[1], ",")
				andArray = append(andArray, bson.M{key: bson.M{"$in": values}})
				continue
			}

			rInt, er = strconv.ParseInt(value, 10, 64)
			if er == nil {
				if opr == "" {
					andArray = append(andArray, bson.M{key: rInt})
				} else {
					andArray = append(andArray, bson.M{key: bson.M{opr: rInt}})
				}
				continue
			}
			var rBool bool
			rBool, er = strconv.ParseBool(value)
			if er == nil {
				if opr == "" {
					andArray = append(andArray, bson.M{key: rBool})
				} else {
					andArray = append(andArray, bson.M{key: bson.M{opr: rBool}})
				}
				continue
			}
			var rFloat float64
			rFloat, er = strconv.ParseFloat(value, 64)
			if er == nil {
				if opr == "" {
					andArray = append(andArray, bson.M{key: rFloat})
				} else {
					andArray = append(andArray, bson.M{key: bson.M{opr: rFloat}})
				}
				continue
			}
			if strings.HasPrefix(value, "ObjectId(") && strings.HasSuffix(value, ")") {
				value = strings.Split(strings.SplitAfter(value, "(")[1], ")")[0]
				if bson.IsObjectIdHex(value) {
					andArray = append(andArray, bson.M{key: bson.ObjectIdHex(value)})
					continue
				}
			}
			if strings.HasPrefix(value, "!") {
				value = strings.Split(value, "!")[1]
				andArray = append(andArray, bson.M{key: bson.M{"$regex": value, "$options": "i"}})
			} else {
				if opr == "" {
					andArray = append(andArray, bson.M{key: bson.M{"$regex": value}})
				} else {
					andArray = append(andArray, bson.M{key: bson.M{opr: value}})
				}
			}
		}
	}

	orArray := []bson.M{}
	for _, obj := range or {
		for key, value := range obj {
			var er error
			var rInt int64
			var opr string = ""
			if strings.HasPrefix(value, ">") {
				values := strings.Split(value, ">")
				opr = "$gt"
				value = values[1]
			} else if strings.HasPrefix(value, ">=") {
				values := strings.Split(value, ">=")
				opr = "$gte"
				value = values[1]
			} else if strings.HasPrefix(value, "<") {
				values := strings.Split(value, "<")
				opr = "$lt"
				value = values[1]
			} else if strings.HasPrefix(value, "<=") {
				values := strings.Split(value, "<=")
				opr = "$lte"
				value = values[1]
			} else if strings.HasPrefix(value, "!=") {
				values := strings.Split(value, "!=")
				opr = "$ne"
				value = values[1]
			} else if strings.HasPrefix(value, "==") {
				values := strings.Split(value, "==")
				opr = "$eq"
				value = values[1]
			}

			layout := "2006-01-02T15:04:05.999Z"
			t, err := time.Parse(layout, value)
			if err == nil {
				if opr == "" {
					orArray = append(orArray, bson.M{key: t})
				} else {
					orArray = append(orArray, bson.M{key: bson.M{opr: t}})
				}
				continue
			}

			if strings.HasPrefix(value, "in:") {
				values := strings.Split(value, ":")
				values = strings.Split(values[1], ",")
				orArray = append(orArray, bson.M{key: bson.M{"$in": values}})
				continue
			}

			rInt, er = strconv.ParseInt(value, 10, 64)
			if er == nil {
				if opr == "" {
					orArray = append(orArray, bson.M{key: rInt})
				} else {
					orArray = append(orArray, bson.M{key: bson.M{opr: rInt}})
				}
				continue
			}
			var rBool bool
			rBool, er = strconv.ParseBool(value)
			if er == nil {
				if opr == "" {
					orArray = append(orArray, bson.M{key: rBool})
				} else {
					orArray = append(orArray, bson.M{key: bson.M{opr: rBool}})
				}
				continue
			}
			var rFloat float64
			rFloat, er = strconv.ParseFloat(value, 64)
			if er == nil {
				if opr == "" {
					orArray = append(orArray, bson.M{key: rFloat})
				} else {
					orArray = append(orArray, bson.M{key: bson.M{opr: rFloat}})
				}
				continue
			}
			if strings.HasPrefix(value, "ObjectId(") && strings.HasSuffix(value, ")") {
				value = strings.Split(strings.SplitAfter(value, "(")[1], ")")[0]
				if bson.IsObjectIdHex(value) {
					orArray = append(orArray, bson.M{key: bson.ObjectIdHex(value)})
					continue
				}
			}
			if strings.HasPrefix(value, "!") {
				value = strings.Split(value, "!")[1]
				orArray = append(orArray, bson.M{key: bson.M{"$regex": value, "$options": "i"}})
			} else {
				if opr == "" {
					orArray = append(orArray, bson.M{key: bson.M{"$regex": value}})
				} else {
					orArray = append(orArray, bson.M{key: bson.M{opr: value}})
				}
			}
		}
	}

	if len(andArray) > 0 && len(orArray) > 0 {
		query = bson.M{"$and": []bson.M{{"$and": andArray}, {"$or": orArray}}}
	} else if len(andArray) > 0 && len(orArray) == 0 {
		query = bson.M{"$and": andArray}
	} else if len(andArray) == 0 && len(orArray) > 0 {
		query = bson.M{"$or": orArray}
	}
	return
}
