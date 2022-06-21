package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"restapi-go/entity"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type ItemHandlerInterface interface {
	ItemHandler(w http.ResponseWriter, r *http.Request)
}

type ItemHandler struct {
	db *sql.DB
}

func NewItemHandler(db *sql.DB) *ItemHandler {
	return &ItemHandler{db: db}
}

var (
	db  *sql.DB
	err error
)

func (h *ItemHandler) ItemHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	item_id := params["item_id"]
	switch r.Method {
	case http.MethodGet:
		h.getItemHandler(w, r, item_id)
	case http.MethodPost:
		h.createItemsHandler(w, r)
	case http.MethodPut:
		h.UpdateOrderId(w, r, item_id)
	case http.MethodDelete:
		h.DeleteOrderHandler(w, r, item_id)
	}
}

func (h *ItemHandler) getItemHandler(w http.ResponseWriter, r *http.Request, item_id string) {
	ctx := context.Background()
	queryString := `select
		o.order_id as order_id
		,o.customer_name
		,o.ordered_at
		,json_agg(json_build_object(
			'lineItemId',i.item_id
			,'itemCode',i.item_code
			,'description',i.description
			,'quantity',i.quantity
			,'order_id',i.order_id
		)) as items
	from orders o join items i
	on o.order_id = i.order_id
	group by o.order_id`
	rows, err := h.db.QueryContext(ctx, queryString)
	if err != nil {
		fmt.Println("query row error", err)
	}
	defer rows.Close()

	var orders []*entity.Order
	for rows.Next() {
		var o entity.Order
		var itemsStr string
		if serr := rows.Scan(&o.Order_id, &o.Customer_name, &o.Ordered_at, &itemsStr); serr != nil {
			fmt.Println("Scan error", serr)
		}
		var items []entity.Item
		if err := json.Unmarshal([]byte(itemsStr), &items); err != nil {
			fmt.Errorf("Error when parsing items")
		} else {
			o.Item = append(o.Item, items...)
		}
		orders = append(orders, &o)
	}

	jsonData, _ := json.Marshal(&orders)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonData)

}

func (h *ItemHandler) createItemsHandler(w http.ResponseWriter, r *http.Request) {
	var newOrder entity.Order
	json.NewDecoder(r.Body).Decode(&newOrder)

	sqlStatement := `INSERT INTO orders (customer_name, ordered_at) VALUES ($1, $2) RETURNING order_id`
	ctx := context.Background()
	var id int
	err := h.db.QueryRowContext(ctx, sqlStatement, newOrder.Customer_name, newOrder.Ordered_at).Scan(&id)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(newOrder.Item); i++ {
		var items entity.Item
		items.Item_code = newOrder.Item[i].Item_code
		items.Description = newOrder.Item[i].Description
		items.Quantity = newOrder.Item[i].Quantity
		sqlStatement1 := `INSERT INTO items (item_code, description, quantity, order_id) VALUES ($1, $2, $3, $4) returning item_id`
		_, err := h.db.Exec(sqlStatement1, items.Item_code, items.Description, items.Quantity, id)
		if err != nil {
			panic(err)
		}
	}

	w.Write([]byte("Successfully created"))
}

func (h *ItemHandler) DeleteOrderHandler(w http.ResponseWriter, r *http.Request, item_id string) {
	sqlstament := `DELETE from items where order_id = $1`
	if idInt, err := strconv.Atoi(item_id); err == nil {
		fmt.Println("tesst", item_id)
		sqlstament1 := `DELETE from orders where order_id = $1`
		_, err := h.db.Exec(sqlstament, idInt)
		if err != nil {
			panic(err)
		}
		_, err = h.db.Exec(sqlstament1, idInt)
		if err != nil {
			panic(err)
		}
		w.Write([]byte("Successfully deleted"))
		return
	}

}

func (h *ItemHandler) UpdateOrderId(w http.ResponseWriter, r *http.Request, item_id string) {
	if item_id != "" {
		var newOrder entity.Order
		json.NewDecoder(r.Body).Decode(&newOrder)
		fmt.Println("tesst", newOrder)
		sqlstatment := `update orders set customer_name = $1 , ordered_at = $2 where order_id = $3`

		res, err := h.db.Exec(sqlstatment,
			newOrder.Customer_name,
			time.Now(),
			item_id,
		)
		if err != nil {
			panic(err)
		}
		for i := 0; i < len(newOrder.Item); i++ {
			var items entity.Item
			items.Item_id = newOrder.Item[i].Item_id
			items.Item_code = newOrder.Item[i].Item_code
			items.Description = newOrder.Item[i].Description
			items.Quantity = newOrder.Item[i].Quantity
			query := `update items set item_code = $1, description = $2, quantity = $3 where order_id = $4 and item_id = $5`

			_, err := h.db.Exec(query, items.Item_code, items.Description, items.Quantity, item_id, items.Item_id)
			if err != nil {
				panic(nil)
			}
		}
		count, err := res.RowsAffected()
		if err != nil {
			panic(err)
		}

		w.Write([]byte(fmt.Sprint("User update ", count)))
		return
	}
}
