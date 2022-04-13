package subs

import (
	"database/sql"

	"log"

	_ "github.com/lib/pq"
)

func Parse_json(aMap map[string]interface{}) string {
	// Предполагаю, что сообщения прияходят с малой периодичностью, поэтому для каждого запроса
	// отдельно создаем соединение с бд, в ином случае необходимо открыть соединение с подключением подписчика
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres dbname=testdb password=311Negative533 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	order_uid := aMap["order_uid"]
	main_table := map[string]interface{}{}
	for key, val := range aMap {

		switch val.(type) {
		case map[string]interface{}:
			if key == "delivery" {

				_, err := db.Exec("insert into delivery  values ($1, $2, $3, $4, $5, $6, $7, $8)",
					order_uid,
					val.(map[string]interface{})["name"],
					val.(map[string]interface{})["phone"],
					val.(map[string]interface{})["zip"],
					val.(map[string]interface{})["city"],
					val.(map[string]interface{})["address"],
					val.(map[string]interface{})["region"],
					val.(map[string]interface{})["email"])
				if err != nil {
					panic(err)
				}
			}

			if key == "payment" {
				_, err := db.Exec("insert into payment  values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
					order_uid,
					val.(map[string]interface{})["transaction"],
					val.(map[string]interface{})["request_id"],
					val.(map[string]interface{})["currency"],
					val.(map[string]interface{})["provider"],
					val.(map[string]interface{})["amount"],
					val.(map[string]interface{})["payment_dt"],
					val.(map[string]interface{})["bank"],
					val.(map[string]interface{})["delivery_cost"],
					val.(map[string]interface{})["goods_total"],
					val.(map[string]interface{})["custom_fee"])
				if err != nil {
					panic(err)
				}
			}

		case []interface{}:
			for _, val := range val.([]interface{}) {
				_, err := db.Exec("insert into items  values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
					order_uid,
					val.(map[string]interface{})["chrt_id"],
					val.(map[string]interface{})["track_number"],
					val.(map[string]interface{})["price"],
					val.(map[string]interface{})["rid"],
					val.(map[string]interface{})["name"],
					val.(map[string]interface{})["sale"],
					val.(map[string]interface{})["size"],
					val.(map[string]interface{})["total_price"],
					val.(map[string]interface{})["nm_id"],
					val.(map[string]interface{})["brand"],
					val.(map[string]interface{})["status"])
				if err != nil {
					panic(err)
				}
			}

		default:
			main_table[key] = val
		}

	}
	db.Exec("insert into main values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		main_table["order_uid"],
		main_table["track_number"],
		main_table["entry"],
		main_table["locale"],
		main_table["internal_signature"],
		main_table["customer_id"],
		main_table["delivery_service"],
		main_table["shardkey"],
		main_table["sm_id"],
		main_table["date_created"],
		main_table["oof_shard"],
	)

	return main_table["order_uid"].(string)
}
