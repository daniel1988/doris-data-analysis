package doris

import (
	"fmt"
)

type Order struct {
	Field string `json:"field"`
	Desc  bool   `json:"desc"`
}

func BuildOrder(order Order) string {
	if order.Desc {
		return fmt.Sprintf("`%s` DESC", order.Field)
	}
	return fmt.Sprintf("`%s` ASC", order.Field)
}

func BuildOrders(orders []Order) []string {
	if len(orders) == 0 {
		return []string{}
	}
	orderSqls := make([]string, 0)
	for _, order := range orders {
		orderSqls = append(orderSqls, BuildOrder(order))
	}
	return orderSqls
}
