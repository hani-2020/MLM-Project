package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
)

type Member struct {
	ID                int
	LeftMember        *Member
	RightMember       *Member
	Position          string
	Parent            *Member
	Left              int
	Right             int
	Level             int
	Sale              float64
	SponsorBonus      float64
	BinaryBonus       float64
	LeftSales         float64
	RightSales        float64
	LeftCarryForward  float64
	RightCarryForward float64
	MatchingBonus     float64
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

func traverse(member *Member) float64 {
	if member == nil {
		return 0
	}
	currentSales := member.Sale
	leftSales := traverse(member.LeftMember)
	rightSales := traverse(member.RightMember)
	return currentSales + leftSales + rightSales
}

func set_get_binary_bonus(binary_percentage float64, capping_amount float64, capping_scope map[string]bool) float64 {
	var total_bonus float64 = 0
	for _, member := range members {
		left_sales := 0.0
		right_sales := 0.0
		if member.LeftMember != nil {
			left_sales = traverse(member.LeftMember)
			member.LeftSales = left_sales
		}
		if member.RightMember != nil {
			right_sales = traverse(member.RightMember)
			member.RightSales = right_sales
		}
		binaryBonus := math.Min(float64(left_sales), float64(right_sales)) * binary_percentage / 100
		if capping_scope["1"] && binaryBonus > capping_amount {
			member.BinaryBonus = capping_amount
		} else {
			member.BinaryBonus = binaryBonus
		}
		carry_forward := left_sales - right_sales
		if carry_forward > 0 {
			member.LeftCarryForward = carry_forward
		} else if carry_forward <= 0 {
			member.RightCarryForward = carry_forward
		}
		total_bonus = total_bonus + member.BinaryBonus
	}
	return total_bonus
}

func set_get_matching_bonus(matching_percs_list []float64, capping_amount float64, capping_scope map[string]bool) float64 {
	var total_bonus float64
	for _, member := range members {
		iterant := 0
		if member.Parent == nil {
			continue
		}
		parent := member.Parent
		apply_matching_bonus(member, parent, matching_percs_list, iterant, capping_amount, capping_scope)
	}
	for _, member := range members {
		total_bonus = total_bonus + member.MatchingBonus
	}
	return total_bonus
}

func apply_matching_bonus(member *Member, parent *Member, matching_perc_list []float64, iterant int, capping_amount float64, capping_scope map[string]bool) {
	if iterant >= len(matching_perc_list) || parent == nil {
		return
	}
	matching_bonus := parent.MatchingBonus
	matching_bonus = matching_bonus + (member.BinaryBonus * matching_perc_list[iterant] / 100)
	if capping_scope["2"] && parent.MatchingBonus > capping_amount {
		parent.MatchingBonus = capping_amount
	} else {
		parent.MatchingBonus = matching_bonus
	}
	iterant = iterant + 1
	parent = parent.Parent
	apply_matching_bonus(member, parent, matching_perc_list, iterant, capping_amount, capping_scope)
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
		binary_perc := data["binary_bonus"].(float64)
		matching_percs := data["matching_bonus_list"].([]interface{})
		var matching_perc_list []float64
		for _, v := range matching_percs {
			matching_perc_list = append(matching_perc_list, v.(float64))
		}
		capping_amount := data["capping_amount"].(float64)
		rawCappingScope := data["capping_scope"].([]interface{})
		cappingScopeMap := make(map[string]bool)
		for _, v := range rawCappingScope {
			cappingScopeMap[v.(string)] = true
		}
		build_tree(number_of_users, joining_package_fee, additional_product_price)
		sponsor_bonus := set_get_sponsor_bonus(sponsor_perc, capping_amount, cappingScopeMap)
		binary_bonus := set_get_binary_bonus(binary_perc, capping_amount, cappingScopeMap)
		matching_bonus := set_get_matching_bonus(matching_perc_list, capping_amount, cappingScopeMap)
		for _, member := range members {
			fmt.Println("#################")
			fmt.Println("id:", member.ID)
			fmt.Println("sponsor bonus:", member.SponsorBonus)
			fmt.Println("binary bonus:", member.BinaryBonus)
			fmt.Println("matching bonus:", member.MatchingBonus)
			fmt.Println("#################")
		}
		fmt.Println(sponsor_bonus)
		fmt.Println(binary_bonus)
		fmt.Println(matching_bonus)
	})
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
