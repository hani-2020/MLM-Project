package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Member struct {
	ID            int
	LeftMember    *Member
	RightMember   *Member
	Position      string
	Parent        *Member
	Left 		  int
	Right		  int
	Level         int
	Sale          float64
	SponsorBonus  float64
	BinaryBonus   float64
	LeftSales     float64
	RightSales    float64
	CarryForward  float64
	MatchingBonus float64
}

var members []*Member

func build_tree(num_users int, joining_package_fee float64, additional_product_price float64) {
	members = []*Member{}
	if num_users <= 0 {
		return
	}
	root := &Member{ID: 1, Parent: nil, Level: 1}
	members = append(members, root)
	queue := []*Member{root}
	current_id := 2
	for current_id <= num_users {
		current_member := queue[0]
		queue = queue[1:]
		if current_id <= num_users {
			current_member.LeftMember = &Member{
				ID:       current_id,
				Parent:   current_member,
				Level:    current_member.Level + 1,
				Sale:     joining_package_fee + additional_product_price,
				Position: "Left",
			}
			members = append(members, current_member.LeftMember)
			queue = append(queue, current_member.LeftMember)
			current_id = current_id + 1
		}
		if current_id <= num_users {
			current_member.RightMember = &Member{
				ID:       current_id,
				Parent:   current_member,
				Level:    current_member.Level + 1,
				Sale:     joining_package_fee + additional_product_price,
				Position: "Right",
			}
			members = append(members, current_member.RightMember)
			queue = append(queue, current_member.RightMember)
			current_id = current_id + 1
		}
	}
}

func set_get_sponsor_bonus(sponsor_perc float64, capping_amount float64, capping_scope map[string]bool) float64 {
	var total_bonus float64 = 0
	for _, member := range members {
		if member.Parent != nil {
			sponsor_bonus := member.Parent.SponsorBonus + (member.Sale * sponsor_perc / 100)
			if capping_scope["3"] && sponsor_bonus > capping_amount {
				member.Parent.SponsorBonus = capping_amount
			} else {
				member.Parent.SponsorBonus = sponsor_bonus
			}
		}
	}
	for _, member := range members {
		total_bonus = total_bonus + member.SponsorBonus
	}
	return total_bonus
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var data map[string]interface{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(body, &data); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		number_of_users := int(data["number_of_users"].(float64))
		joining_package_fee := data["joining_package_fee"].(float64)
		additional_product_price := data["additional_product_price"].(float64)
		sponsor_perc := data["sponsor_bonus"].(float64)
		capping_amount := data["capping_amount"].(float64)
		rawCappingScope, ok := data["capping_scope"].([]interface{})
		if !ok {
			fmt.Println("Error: capping_scope is not a list")
			return
		}
		var capping_scope []string
		for _, v := range rawCappingScope {
			str, ok := v.(string)
			if ok {
				capping_scope = append(capping_scope, str)
			} else {
				fmt.Println("Error: non-string value found in capping_scope")
				return
			}
		}
		cappingScopeMap := make(map[string]bool)
		for _, v := range capping_scope {
			cappingScopeMap[v] = true
		}
		build_tree(number_of_users, joining_package_fee, additional_product_price)
		sponsor_bonus := set_get_sponsor_bonus(sponsor_perc, capping_amount, cappingScopeMap)
		for _, member := range members {
			fmt.Println("#################")
			fmt.Println("id:", member.ID)
			fmt.Println("sponsor bonus:", member.SponsorBonus)
			fmt.Println("#################")
		}
		fmt.Println(sponsor_bonus)
	})
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
