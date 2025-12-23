package server

import (
	"net/http"
)

func RegisterSignupHandler(mux *http.ServeMux) {
	// mux := &http.ServeMux{}
	mux.HandleFunc("POST /signup", PostSignupHandler)
}

func PostSignupHandler(w http.ResponseWriter, r *http.Request) {
	// var stringBuilder *strings.Builder = &strings.Builder{}
	// query := r.URL.Query()

	// stringBuilder.WriteString("SELECT * FROM products where 1=1 ")
	// // filter
	// filterID := query.Get("id")
	// filterName := query.Get("name")
	// applyFilter(stringBuilder, filterID, filterName)

	// // apply sort
	// sortBy := query.Get("sortBy")
	// sortOrder := query.Get("sortOrder")
	// applySort(stringBuilder, sortBy, sortOrder)

	// stringBuilder.WriteRune(';')
	// // fetch data from db
	// db := db.GetDB()
	// stmt, err := db.Prepare(stringBuilder.String())
	// if err != nil {
	// 	log.Println(err)
	// 	http.Error(w, "ERROR fetching data 101", http.StatusInternalServerError)
	// 	return
	// }

	// rows, err := stmt.Query()

	// if err != nil {
	// 	http.Error(w, "ERROR fetching data 102", http.StatusInternalServerError)
	// }

	// p := Product{}
	// resp := []Product{}
	// for rows.Next() {
	// 	rows.Scan(&p.Id, &p.Name, &p.Count)
	// 	resp = append(resp, p)
	// }
	// data := ProductResponse{
	// 	Status: http.StatusOK,
	// 	Count:  len(resp),
	// 	Data:   resp,
	// }

	// jsonData, err := json.Marshal(data)
	// if err != nil {
	// 	http.Error(w, "Error reading products", http.StatusInternalServerError)
	// } else {
	// 	w.Header().Set(constants.Header_ContentType, constants.ContentType_ApplicationJSON)
	// 	w.Write(jsonData)
	// }
}
