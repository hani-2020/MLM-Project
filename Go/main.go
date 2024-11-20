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
var current_id int = 1
var queue []*Member

func build_tree(num_users int, joining_package_fee float64, additional_product_price float64) {
	num_users = current_id + num_users - 1
	if num_users <= 0 {
		return
	}
	if current_id == 1 {
		root := &Member{ID: current_id, Parent: nil, Level: 1}
		members = append(members, root)
		queue = []*Member{root}
		current_id = current_id + 1
		num_users = num_users + 1
	}
	for current_id <= num_users {
		current_member := queue[0]
		if current_id <= num_users && current_member.LeftMember == nil {
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
		if current_id <= num_users && current_member.RightMember == nil {
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
			queue = queue[1:]
		}
	}
}

func set_get_sponsor_bonus(sponsor_perc float64, capping_amount float64, capping_scope map[string]bool, starting_id int) float64 {
	var total_bonus float64 = 0
	for _, member := range members {
		if member.ID >= starting_id && member.Parent != nil {
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

func set_get_binary_bonus(binaryBonusPairingRatios map[string]int, binaryBonusRange []map[string]float64, capping_amount float64, capping_scope map[string]bool) float64 {
	var total_bonus float64 = 0
	for _, member := range members {
		left_sales := 0.0
		right_sales := 0.0
		if member.LeftMember != nil {
			left_sales = traverse(member.LeftMember)
			member.LeftSales = left_sales + member.LeftCarryForward
		}
		if member.RightMember != nil {
			right_sales = traverse(member.RightMember)
			member.RightSales = right_sales + member.RightCarryForward
		}
		left_bits := int(left_sales) / binaryBonusPairingRatios["left"]
		right_bits := int(right_sales) / binaryBonusPairingRatios["right"]
		bits := math.Min(float64(left_bits), float64(right_bits))
		left_amount := bits * float64(binaryBonusPairingRatios["left"])
		right_amount := bits * float64(binaryBonusPairingRatios["right"])
		member.LeftCarryForward = member.LeftSales - left_amount
		member.RightCarryForward = member.RightSales - right_amount
		var binaryBonus float64
		for i := range binaryBonusRange {
			if bits >= binaryBonusRange[i]["min"] && bits <= binaryBonusRange[i]["max"] {
				binary_percentage := binaryBonusRange[i]["bonus"]
				binaryBonus = math.Min(float64(member.LeftSales), float64(member.RightSales))*binary_percentage/100
				break
			}
		}
		if capping_scope["1"] && binaryBonus > capping_amount {
			member.BinaryBonus = capping_amount
		} else {
			member.BinaryBonus = binaryBonus
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

// func num_cycles(number_of_users float64, products_catalogue map[string]map[string]float64) float64 {
// 	cycles := 0.0
// 	for number_of_users > 0 {
// 		for _, v1 := range products_catalogue {
// 			for key := range v1 {
// 				if key == "quantity" {
// 					number_of_users = number_of_users - v1["quantity"]
// 					cycles = cycles + 1
// 				}
// 			}
// 		}
// 	}
// 	return cycles
// }

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

		total_num_of_users := int(data["number_of_users"].(float64))
		additional_product_price := data["additional_product_price"].(float64)

		product_orders := data["product_order_list"].([]interface{})
		products_catalogue := data["products_catalogue"].(map[string]interface{})

		sponsor_perc := data["sponsor_bonus"].(float64)

		binary_bonus_pairing_ratios := data["binary_bonus_pairing_ratios"].(map[string]interface{})
		binary_bonus_range := data["binary_bonus_range"].([]interface{})

		matching_percs := data["matching_bonus_list"].([]interface{})

		capping_amount := data["capping_amount"].(float64)
		rawCappingScope := data["capping_scope"].([]interface{})
		cappingScopeMap := make(map[string]bool)

		var product_order_list []string

		for _, v := range rawCappingScope {
			cappingScopeMap[v.(string)] = true
		}

		for product_order := range product_orders {
			product_order_list = append(product_order_list, product_orders[product_order].(string))
		}

		productCatalogueMap := make(map[string]map[string]float64)
		for k, v := range products_catalogue {
			product_map := v.(map[string]interface{})
			detail_map := make(map[string]float64)
			for k2, v2 := range product_map {
				detail_map[k2] = v2.(float64)
			}
			productCatalogueMap[k] = detail_map
		}

		var matching_perc_list []float64
		for _, v := range matching_percs {
			matching_perc_list = append(matching_perc_list, v.(float64))
		}

		binaryBonusPairingRatios := make(map[string]int)
		for key, value := range binary_bonus_pairing_ratios {
			binaryBonusPairingRatios[key] = int(value.(float64))
		}

		binaryBonusRange := []map[string]float64{}
		for _, item := range binary_bonus_range {
			rangeMap := item.(map[string]interface{})
			convertedMap := map[string]float64{}
			for key, value := range rangeMap {
				convertedMap[key] = float64(value.(float64))
			}
			binaryBonusRange = append(binaryBonusRange, convertedMap)
		}
		var sponsor_bonus, binary_bonus, matching_bonus float64
		members = []*Member{}
		current_id = 1
		queue = []*Member{}
		var cycles_data [][]*Member
		var cycle_start_ids []int = []int{0}
		cycle_num := 0
		fmt.Println("###################start#######################")
		for total_num_of_users > 0 {
			cycle_num = cycle_num + 1
			for product := range product_order_list {
				number_of_users := int(productCatalogueMap[product_order_list[product]]["quantity"])
				joining_package_fee := productCatalogueMap[product_order_list[product]]["price"]
				if total_num_of_users >= number_of_users {
					build_tree(number_of_users, joining_package_fee, additional_product_price)
					total_num_of_users = total_num_of_users - number_of_users
				} else {
					number_of_users = total_num_of_users
					build_tree(number_of_users, joining_package_fee, additional_product_price)
					total_num_of_users = total_num_of_users - number_of_users
					break
				}
			}
			cycle_start_ids = append(cycle_start_ids, len(members))
			starting_id := cycle_start_ids[cycle_num-1]
			sponsor_bonus = set_get_sponsor_bonus(sponsor_perc, capping_amount, cappingScopeMap, starting_id)
			binary_bonus = set_get_binary_bonus(binaryBonusPairingRatios, binaryBonusRange, capping_amount, cappingScopeMap)
			matching_bonus = set_get_matching_bonus(matching_perc_list, capping_amount, cappingScopeMap)
			var copiedMembers []*Member
			for _, member := range members {
				copiedMember := &Member{
					ID:                member.ID,
					LeftMember:        member.LeftMember,
					RightMember:       member.RightMember,
					Position:          member.Position,
					Parent:            member.Parent,
					Left:              member.Left,
					Right:             member.Right,
					Level:             member.Level,
					Sale:              member.Sale,
					SponsorBonus:      member.SponsorBonus,
					BinaryBonus:       member.BinaryBonus,
					LeftSales:         member.LeftSales,
					RightSales:        member.RightSales,
					LeftCarryForward:  member.LeftCarryForward,
					RightCarryForward: member.RightCarryForward,
					MatchingBonus:     member.MatchingBonus,
				}
				copiedMembers = append(copiedMembers, copiedMember)
			}
			cycles_data = append(cycles_data, copiedMembers)
			for _, member := range members {
				member.MatchingBonus = 0
			}
		}
		fmt.Println("###################end#######################")
		cycle := 0
		for _, members := range cycles_data {
			cycle = cycle + 1
			fmt.Println("Cycle:", cycle)
			for _, member := range members {
				fmt.Println("id:", member.ID)
				if member.Parent != nil {
					fmt.Println("Parent Member:", member.Parent.ID)
				} else {
					fmt.Println("Parent Member: nil")
				}
				if member.LeftMember != nil {
					fmt.Println("Left Member:", member.LeftMember.ID)
				} else {
					fmt.Println("Left Member: nil")
				}
				if member.RightMember != nil {
					fmt.Println("Right Member:", member.RightMember.ID)
				} else {
					fmt.Println("Right Member: nil")
				}
				fmt.Println("JPF:", member.Sale)
				fmt.Println("Left Carry:", member.LeftCarryForward)
				fmt.Println("Right Carry:", member.RightCarryForward)
				fmt.Println("Left sales:", member.LeftSales)
				fmt.Println("Right sales:", member.RightSales)
				fmt.Println("Sponsor Bonus:", member.SponsorBonus)
				fmt.Println("Binary Bonus:", member.BinaryBonus)
				fmt.Println("Matching Bonus:", member.MatchingBonus)
			}
		}
		fmt.Println("##################")
		fmt.Println(sponsor_bonus)
		fmt.Println(binary_bonus)
		fmt.Println(matching_bonus)
		fmt.Println("##################")
	})
	http.ListenAndServe(":8080", nil)
}
