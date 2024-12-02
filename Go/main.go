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
	Children          []*Member
	Position          string
	Parent            *Member
	Left              int
	Right             int
	Level             int
	Sale              float64
	SponsorBonus      float64
	BinaryBonus       float64
	MatchingBonus     float64
	PoolBonus         float64
	LeftSales         float64
	RightSales        float64
	LeftCarryForward  float64
	RightCarryForward float64
}

// type MemberExport struct {
// 	ID                int     `json:"id"`
// 	LeftMember        int     `json:"left_member"`
// 	RightMember       int     `json:"right_member"`
// 	Children          []int   `json:"children"`
// 	Position          string  `json:"position"`
// 	Parent            int     `json:"parent"`
// 	Left              int     `json:"left"`
// 	Right             int     `json:"right"`
// 	Level             int     `json:"level"`
// 	Sale              float64 `json:"sale"`
// 	SponsorBonus      float64 `json:"sponsor_bonus"`
// 	BinaryBonus       float64 `json:"binary_bonus"`
// 	LeftSales         float64 `json:"left_sales"`
// 	RightSales        float64 `json:"right_sales"`
// 	LeftCarryForward  float64 `json:"left_carry_forward"`
// 	RightCarryForward float64 `json:"right_carry_forward"`
// 	MatchingBonus     float64 `json:"matching_bonus"`
// }

type ExportCycleData struct {
	NumberUsers   int     `json:"number_users"`
	Expense       float64 `json:"expense"`
	Revenue       float64 `json:"revenue"`
	Profit        float64 `json:"profit"`
	Cycle         int     `json:"cycle"`
	BinaryBonus   float64 `json:"binary_bonus"`
	MatchingBonus float64 `json:"matching_bonus"`
	SponsorBonus  float64 `json:"sponsor_bonus"`
	PoolBonus     float64 `json:"pool_bonus"`
	// MemberData    []MemberExport `json:"member_data"`
}

type ExportData struct {
	PlanType           string            `json:"plan_type"`
	TotalExpense       float64           `json:"total_expense"`
	TotalRevenue       float64           `json:"total_revenue"`
	//TotalProfit        float64           `json:"total_profit"`
	TotalCycles        int               `json:"total_cycles"`
	TotalBinaryBonus   float64           `json:"total_binary_bonus"`
	TotalSponsorBonus  float64           `json:"total_sponsor_bonus"`
	TotalMatchingBonus float64           `json:"total_matching_bonus"`
	TotalPoolBonus     float64           `json:"total_pool_bonus"`
	CycleData          []ExportCycleData `json:"cycle_data"`
}

var members []*Member
var current_id int = 1
var queue []*Member

func build_unilevel_tree(num_users int, joining_package_fee float64, additional_product_price float64, downlines_per_user int) {
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
		if current_id <= num_users && len(current_member.Children) < downlines_per_user-1 {
			child_member := &Member{
				ID:     current_id,
				Parent: current_member,
				Level:  current_member.Level + 1,
				Sale:   joining_package_fee + additional_product_price,
			}
			current_member.Children = append(current_member.Children, child_member)
			members = append(members, child_member)
			queue = append(queue, child_member)
			current_id = current_id + 1
		}
		if current_id <= num_users && len(current_member.Children) >= downlines_per_user-1 {
			child_member := &Member{
				ID:     current_id,
				Parent: current_member,
				Level:  current_member.Level + 1,
				Sale:   joining_package_fee + additional_product_price,
			}
			current_member.Children = append(current_member.Children, child_member)
			members = append(members, child_member)
			queue = append(queue, child_member)
			current_id = current_id + 1
			queue = queue[1:]
		}
	}
}

func build_binary_tree(num_users int, joining_package_fee float64, additional_product_price float64) {
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

func set_get_sponsor_bonus(sponsor_perc float64, capping_amount float64, capping_scope map[string]bool) {
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
				binaryBonus = math.Min(float64(member.LeftSales), float64(member.RightSales)) * binary_percentage / 100
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
	if capping_scope["2"] && matching_bonus > capping_amount {
		parent.MatchingBonus = capping_amount
	} else {
		parent.MatchingBonus = matching_bonus
	}
	iterant = iterant + 1
	parent = parent.Parent
	apply_matching_bonus(member, parent, matching_perc_list, iterant, capping_amount, capping_scope)
}

func set_get_pool_bonus(pool_perc float64, dist_no int, expense float64, revenue float64, sponsorBonus float64, binary_bonus float64, matching_bonus float64) float64{
	var pool_amount float64
	no_people := math.Min(float64(dist_no), float64(len(members)))
	profit := revenue - (expense + sponsorBonus + binary_bonus + matching_bonus)
	if profit < 0{
		pool_amount = 0
	}else{
		pool_amount = profit * pool_perc / 100
		for _, member := range members{
			if member.ID < int(no_people){
				member.PoolBonus = pool_amount/no_people
			}
		}
	}
	return pool_amount
}

func set_get_level_bonus(level_percs_list []float64, capping_amount float64, capping_scope map[string]bool) float64 {
	var total_bonus float64
	for _, member := range members {
		iterant := 0
		if member.Parent == nil {
			continue
		}
		parent := member.Parent
		apply_level_bonus(member, parent, level_percs_list, iterant, capping_amount, capping_scope)
	}
	for _, member := range members {
		total_bonus = total_bonus + member.SponsorBonus
	}
	return total_bonus
}

func apply_level_bonus(member *Member, parent *Member, level_percs_list []float64, iterant int, capping_amount float64, capping_scope map[string]bool) {
	if iterant >= len(level_percs_list) || parent == nil {
		return
	}
	sponsor_bonus := parent.SponsorBonus
	sponsor_bonus = sponsor_bonus + (member.Sale * level_percs_list[iterant] / 100)
	if capping_scope["4"] && sponsor_bonus > capping_amount {
		parent.SponsorBonus = capping_amount
	} else {
		parent.SponsorBonus = sponsor_bonus
	}
	iterant = iterant + 1
	parent = parent.Parent
	apply_level_bonus(member, parent, level_percs_list, iterant, capping_amount, capping_scope)
}

func set_get_uni_matching_bonus(matching_percs_list []float64, capping_amount float64, capping_scope map[string]bool) float64 {
	var total_bonus float64
	for _, member := range members {
		iterant := 0
		if member.Parent == nil {
			continue
		}
		parent := member.Parent
		apply_uni_matching_bonus(member, parent, matching_percs_list, iterant, capping_amount, capping_scope)
	}
	for _, member := range members {
		total_bonus = total_bonus + member.MatchingBonus
	}
	return total_bonus
}

func apply_uni_matching_bonus(member *Member, parent *Member, matching_perc_list []float64, iterant int, capping_amount float64, capping_scope map[string]bool) {
	if iterant >= len(matching_perc_list) || parent == nil {
		return
	}
	matching_bonus := parent.MatchingBonus
	matching_bonus = matching_bonus + (member.SponsorBonus * matching_perc_list[iterant] / 100)
	if capping_scope["2"] && matching_bonus > capping_amount {
		parent.MatchingBonus = capping_amount
	} else {
		parent.MatchingBonus = matching_bonus
	}
	iterant = iterant + 1
	parent = parent.Parent
	apply_uni_matching_bonus(member, parent, matching_perc_list, iterant, capping_amount, capping_scope)
}

func main() {
	//for binary plan
	http.HandleFunc("/binary-calc/", func(w http.ResponseWriter, r *http.Request) {
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
		expense_per_member := data["expenses_per_member"].(float64)
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
		members = []*Member{}
		current_id = 1
		queue = []*Member{}
		var totalExpense, totalRevenue, totalBinaryBonus, totalMatchingBonus, totalSponsorBonus, totalPoolBonus float64
		var cycleList []ExportCycleData
		cycle_num := 0
		for total_num_of_users > 0 {
			var binary_bonus, matching_bonus, pool_bonus float64
			cycle_num = cycle_num + 1
			for product := range product_order_list {
				number_of_users := int(productCatalogueMap[product_order_list[product]]["quantity"])
				joining_package_fee := productCatalogueMap[product_order_list[product]]["price"]
				if total_num_of_users >= number_of_users {
					build_binary_tree(number_of_users, joining_package_fee, additional_product_price)
					total_num_of_users = total_num_of_users - number_of_users
				} else {
					number_of_users = total_num_of_users
					build_binary_tree(number_of_users, joining_package_fee, additional_product_price)
					total_num_of_users = total_num_of_users - number_of_users
					break
				}
			}
			set_get_sponsor_bonus(sponsor_perc, capping_amount, cappingScopeMap)
			binary_bonus = set_get_binary_bonus(binaryBonusPairingRatios, binaryBonusRange, capping_amount, cappingScopeMap)
			matching_bonus = set_get_matching_bonus(matching_perc_list, capping_amount, cappingScopeMap)
			// var copiedMembers []MemberExport
			var revenue /*binaryBonus, matchingBonus*/, sponsorBonus float64
			for _, member := range members {
				revenue = revenue + member.Sale
				sponsorBonus = sponsorBonus + member.SponsorBonus
				member.MatchingBonus = 0
				member.SponsorBonus = 0
				// sponsorBonus = sponsorBonus + member.SponsorBonus
				// 	leftmember := -1
				// 	rightmember := -1
				// 	parent := -1
				// 	if member.LeftMember != nil {
				// 		leftmember = member.LeftMember.ID
				// 	}
				// 	if member.RightMember != nil {
				// 		rightmember = member.RightMember.ID
				// 	}
				// 	if member.Parent != nil {
				// 		parent = member.Parent.ID
				// 	}
				// 	copiedMember := MemberExport{
				// 		ID:                member.ID,
				// 		LeftMember:        leftmember,
				// 		RightMember:       rightmember,
				// 		Position:          member.Position,
				// 		Parent:            parent,
				// 		Left:              member.Left,
				// 		Right:             member.Right,
				// 		Level:             member.Level,
				// 		Sale:              member.Sale,
				// 		SponsorBonus:      member.SponsorBonus,
				// 		BinaryBonus:       member.BinaryBonus,
				// 		LeftSales:         member.LeftSales,
				// 		RightSales:        member.RightSales,
				// 		LeftCarryForward:  member.LeftCarryForward,
				// 		RightCarryForward: member.RightCarryForward,
				// 		MatchingBonus:     member.MatchingBonus,
				// 	}
				// 	copiedMembers = append(copiedMembers, copiedMember)
			}
			// expense := expense_per_member*float64(len(copiedMembers))
			expense := expense_per_member * float64(len(members))
			pool_perc := 3.0
			dist_no := 100
			pool_bonus = set_get_pool_bonus(pool_perc, dist_no, expense, revenue, sponsorBonus, binary_bonus, matching_bonus)
			totalExpense = totalExpense + expense
			totalRevenue = totalRevenue + revenue
			totalBinaryBonus = totalBinaryBonus + binary_bonus
			totalSponsorBonus = totalSponsorBonus + sponsorBonus
			totalMatchingBonus = totalMatchingBonus + matching_bonus
			totalPoolBonus = totalPoolBonus + pool_bonus
			exportCycleData := ExportCycleData{
				// NumberUsers:   len(copiedMembers),
				NumberUsers:   len(members),
				Expense:       expense,
				Revenue:       revenue,
				// Profit:        revenue - expense,
				Cycle:         cycle_num,
				BinaryBonus:   binary_bonus,
				SponsorBonus:  sponsorBonus,
				MatchingBonus: matching_bonus,
				PoolBonus: pool_bonus,
				// MemberData:    copiedMembers,
			}
			cycleList = append(cycleList, exportCycleData)
			// for _, member := range members {
			// 	member.MatchingBonus = 0
			// 	member.SponsorBonus = 0
			// }
		}
		exportData := ExportData{
			PlanType: "binary",
			TotalExpense:       totalExpense,
			TotalRevenue:       totalRevenue,
			// TotalProfit:        totalRevenue - totalExpense,
			TotalCycles:        cycle_num,
			TotalBinaryBonus:   totalBinaryBonus,
			TotalSponsorBonus:  totalSponsorBonus,
			TotalMatchingBonus: totalMatchingBonus,
			TotalPoolBonus: totalPoolBonus,
			CycleData:          cycleList,
		}
		response, err := json.Marshal(exportData)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})

	//for unilevel plan
	http.HandleFunc("/unilevel-calc/", func(w http.ResponseWriter, r *http.Request) {
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
		expense_per_member := data["expenses_per_member"].(float64)
		downlines_per_user := int(data["downlines_per_user"].(float64))
		additional_product_price := data["additional_product_price"].(float64)

		product_orders := data["product_order_list"].([]interface{})
		products_catalogue := data["products_catalogue"].(map[string]interface{})

		level_percs := data["level_bonus"].([]interface{})

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

		var level_percs_list []float64
		for _, v := range level_percs {
			level_percs_list = append(level_percs_list, v.(float64))
		}
		members = []*Member{}
		current_id = 1
		queue = []*Member{}
		var totalExpense, totalRevenue, totalMatchingBonus, totalSponsorBonus float64
		var cycleList []ExportCycleData
		cycle_num := 0
		for total_num_of_users > 0 {
			var level_bonus, matching_bonus float64
			cycle_num = cycle_num + 1
			for product := range product_order_list {
				number_of_users := int(productCatalogueMap[product_order_list[product]]["quantity"])
				joining_package_fee := productCatalogueMap[product_order_list[product]]["price"]
				if total_num_of_users >= number_of_users {
					build_unilevel_tree(number_of_users, joining_package_fee, additional_product_price, downlines_per_user)
					total_num_of_users = total_num_of_users - number_of_users
				} else {
					number_of_users = total_num_of_users
					build_unilevel_tree(number_of_users, joining_package_fee, additional_product_price, downlines_per_user)
					total_num_of_users = total_num_of_users - number_of_users
					break
				}
			}
			level_bonus = set_get_level_bonus(level_percs_list, capping_amount, cappingScopeMap)
			matching_bonus = set_get_uni_matching_bonus(matching_perc_list, capping_amount, cappingScopeMap)
			var revenue float64
			for _, member := range members {
				revenue = revenue + member.Sale
				member.MatchingBonus = 0
				member.SponsorBonus = 0
			}
			expense := expense_per_member * float64(len(members))
			totalExpense = totalExpense + expense
			totalRevenue = totalRevenue + revenue
			totalSponsorBonus = totalSponsorBonus + level_bonus
			totalMatchingBonus = totalMatchingBonus + matching_bonus
			exportCycleData := ExportCycleData{
				NumberUsers:   len(members),
				Expense:       expense,
				Revenue:       revenue,
				Cycle:         cycle_num,
				SponsorBonus:  level_bonus,
				MatchingBonus: matching_bonus,
			}
			cycleList = append(cycleList, exportCycleData)
		}
		exportData := ExportData{
			PlanType: "unilevel",
			TotalExpense:       totalExpense,
			TotalRevenue:       totalRevenue,
			TotalCycles:        cycle_num,
			TotalSponsorBonus:  totalSponsorBonus,
			TotalMatchingBonus: totalMatchingBonus,
			CycleData:          cycleList,
		}
		fmt.Println(exportData)
		response, err := json.Marshal(exportData)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	})
	http.ListenAndServe(":8080", nil)
}
