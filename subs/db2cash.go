package subs

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var Cash_data map[string]interface{}

// Загрузка данных из бд в глобальную переменную Cash_data для каждого уникального
// order_id, на случай если имеются одинаковые строки в бд
func Load_cash() map[string]interface{} {
	Cash_data = make(map[string]interface{})
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=testdb password=311Negative533 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("select distinct order_uid from main where order_uid != ''")
	if err != nil {
		log.Println("Wrong query")
	}
	var order_uid []string
	for rows.Next() {
		var data string
		rows.Scan(&data)
		order_uid = append(order_uid, data)

	}

	for _, val := range order_uid {
		make_json(val, db)
	}
	return Cash_data
}

func make_json(order_uid string, db *sql.DB) {

	type items_t struct {
		order_uid    string
		chrt_id      int
		track_number string
		price        int
		rid          string
		name         string
		sale         int
		size         string
		total_price  int
		nm_id        int
		brand        string
		status       int
	}

	type delivery_t struct {
		order_uid string
		name      string
		phone     string
		zip       string
		city      string
		address   string
		region    string
		email     string
	}

	type payment_t struct {
		order_uid     string
		transaction   string
		request_id    string
		currency      string
		provider      string
		amount        int
		payment_dt    int
		bank          string
		delivery_cost int
		goods_total   int
		custom_fee    int
	}

	type json struct {
		order_uid          string
		track_number       string
		entry              string
		delivery           delivery_t
		payment            payment_t
		items              []items_t
		locale             string
		internal_signature string
		customer_id        string
		delivery_service   string
		shardkey           string
		sm_id              int
		date_created       string
		oof_shard          string
	}
	rows_items, err := db.Query("select * from items where order_uid = $1", order_uid)
	if err != nil {
		log.Println("Wrong query")
	}
	var albom_items = []items_t{}

	for rows_items.Next() {
		var alb_items items_t
		if err := rows_items.Scan(
			&alb_items.order_uid,
			&alb_items.chrt_id,
			&alb_items.track_number,
			&alb_items.price,
			&alb_items.rid,
			&alb_items.name,
			&alb_items.sale,
			&alb_items.size,
			&alb_items.total_price,
			&alb_items.nm_id,
			&alb_items.brand,
			&alb_items.status); err != nil {
			log.Print("Failed to convert data into string")
		}
		albom_items = append(albom_items, alb_items)
	}
	rows_del, err := db.Query("select * from delivery where order_uid = $1", order_uid)
	if err != nil {
		log.Println("Wrong query")
	}

	var alb_del delivery_t
	for rows_del.Next() {

		if err := rows_del.Scan(
			&alb_del.order_uid,
			&alb_del.name,
			&alb_del.phone,
			&alb_del.zip,
			&alb_del.city,
			&alb_del.address,
			&alb_del.region,
			&alb_del.email,
		); err != nil {
			log.Print("Failed to convert data into string")
		}
	}

	rows_pay, err := db.Query("select * from payment where order_uid = $1", order_uid)
	if err != nil {
		log.Println("Wrong query")
	}
	var alb_pay payment_t
	for rows_pay.Next() {

		if err := rows_pay.Scan(
			&alb_pay.order_uid,
			&alb_pay.transaction,
			&alb_pay.request_id,
			&alb_pay.currency,
			&alb_pay.provider,
			&alb_pay.amount,
			&alb_pay.payment_dt,
			&alb_pay.bank,
			&alb_pay.delivery_cost,
			&alb_pay.goods_total,
			&alb_pay.custom_fee,
		); err != nil {
			log.Print("Failed to convert data into string")
		}
	}

	rows_main, err := db.Query("select * from main where order_uid = $1", order_uid)
	if err != nil {
		log.Println("Wrong query")
	}

	var alb_main json

	for rows_main.Next() {

		if err := rows_main.Scan(
			&alb_main.order_uid,
			&alb_main.track_number,
			&alb_main.entry,
			&alb_main.locale,
			&alb_main.internal_signature,
			&alb_main.customer_id,
			&alb_main.delivery_service,
			&alb_main.shardkey,
			&alb_main.sm_id,
			&alb_main.date_created,
			&alb_main.oof_shard,
		); err != nil {
			log.Print("Failed to convert data into string")
		}
		alb_main.delivery = alb_del
		alb_main.items = albom_items
		alb_main.payment = alb_pay

	}
	Cash_data[order_uid] = alb_main
}
