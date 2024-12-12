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
}

type ExportData struct {
	PlanType           string            `json:"plan_type"`
	TotalExpense       float64           `json:"total_expense"`
	TotalRevenue       float64           `json:"total_revenue"`
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

func set_get_sponsor_bonus(sponsor_perc float64, capping_amount float64, capping_scope map[string]bool, sponsor_bonus_as_amount bool) float64 {
	var sponsorBonus float64
	for _, member := range members {
		var sponsor_bonus float64
		if member.Parent != nil {
			if sponsor_bonus_as_amount != true {
				sponsor_bonus = member.Parent.SponsorBonus + (member.Sale * sponsor_perc / 100)
			} else {
				sponsor_bonus = member.Parent.SponsorBonus + sponsor_perc
			}
			if capping_scope["3"] && sponsor_bonus > capping_amount {
				member.Parent.SponsorBonus = capping_amount
			} else {
				member.Parent.SponsorBonus = sponsor_bonus
			}
		}
	}
	for _, member := range members {
		sponsorBonus = sponsorBonus + member.SponsorBonus
	}
	return sponsorBonus
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

func set_get_binary_bonus(binaryBonusPairingRatios map[string]int, binaryBonusRange []map[string]float64, capping_amount float64, capping_scope map[string]bool, binary_bonus_as_amount bool) float64 {
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
		left_bits := int(member.LeftSales) / binaryBonusPairingRatios["left"]
		right_bits := int(member.RightSales) / binaryBonusPairingRatios["right"]
		bits := math.Min(float64(left_bits), float64(right_bits))
		left_amount := bits * float64(binaryBonusPairingRatios["left"])
		right_amount := bits * float64(binaryBonusPairingRatios["right"])
		member.LeftCarryForward = member.LeftSales - left_amount
		member.RightCarryForward = member.RightSales - right_amount
		var binaryBonus float64
		for i := range binaryBonusRange {
			if bits >= binaryBonusRange[i]["min"] && bits <= binaryBonusRange[i]["max"] {
				binary_percentage := binaryBonusRange[i]["bonus"]
				if binary_bonus_as_amount != true {
					binaryBonus = math.Min(float64(left_amount), float64(right_amount)) * binary_percentage / 100
				} else {
					if math.Min(float64(left_amount), float64(right_amount)) == 0 {
						binaryBonus = 0
					} else {
						binaryBonus = binary_percentage
					}
				}
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

func set_get_matching_bonus(matching_percs_list []float64, capping_amount float64, capping_scope map[string]bool, matching_bonus_as_amount bool) float64 {
	var total_bonus float64
	for _, member := range members {
		iterant := 0
		if member.Parent == nil {
			continue
		}
		parent := member.Parent
		apply_matching_bonus(member, parent, matching_percs_list, iterant, capping_amount, capping_scope, matching_bonus_as_amount)
	}
	for _, member := range members {
		total_bonus = total_bonus + member.MatchingBonus
	}
	return total_bonus
}

func apply_matching_bonus(member *Member, parent *Member, matching_perc_list []float64, iterant int, capping_amount float64, capping_scope map[string]bool, matching_bonus_as_amount bool) {
	if iterant >= len(matching_perc_list) || parent == nil {
		return
	}
	matching_bonus := parent.MatchingBonus
	if matching_bonus_as_amount != true {
		matching_bonus = matching_bonus + (member.BinaryBonus * matching_perc_list[iterant] / 100)
	} else {
		if member.BinaryBonus == 0 {
			matching_bonus = matching_bonus + 0
		} else {
			matching_bonus = matching_bonus + matching_perc_list[iterant]
		}
	}
	if capping_scope["2"] && matching_bonus > capping_amount {
		parent.MatchingBonus = capping_amount
	} else {
		parent.MatchingBonus = matching_bonus
	}
	iterant = iterant + 1
	parent = parent.Parent
	apply_matching_bonus(member, parent, matching_perc_list, iterant, capping_amount, capping_scope, matching_bonus_as_amount)
}

func set_get_pool_bonus(pool_perc float64, dist_no int, expense float64, revenue float64, sponsorBonus float64, binary_bonus float64, matching_bonus float64) float64 {
	var pool_amount float64
	no_people := math.Min(float64(dist_no), float64(len(members)))
	profit := revenue - (expense + sponsorBonus + binary_bonus + matching_bonus)
	if profit < 0 {
		pool_amount = 0
	} else {
		pool_amount = profit * pool_perc / 100
		for _, member := range members {
			if member.ID < int(no_people+1) && member.ID > 1 {
				member.PoolBonus = pool_amount / no_people
			}
		}
	}
	return pool_amount
}

func set_get_level_bonus(level_percs_list []float64, capping_amount float64, capping_scope map[string]bool, level_bonus_as_amount bool) float64 {
	var total_bonus float64
	for _, member := range members {
		iterant := 0
		if member.Parent == nil {
			continue
		}
		parent := member.Parent
		apply_level_bonus(member, parent, level_percs_list, iterant, capping_amount, capping_scope, level_bonus_as_amount)
	}
	for _, member := range members {
		total_bonus = total_bonus + member.SponsorBonus
	}
	return total_bonus
}

func apply_level_bonus(member *Member, parent *Member, level_percs_list []float64, iterant int, capping_amount float64, capping_scope map[string]bool, level_bonus_as_amount bool) {
	if iterant >= len(level_percs_list) || parent == nil {
		return
	}
	sponsor_bonus := parent.SponsorBonus
	if level_bonus_as_amount != true {
		sponsor_bonus = sponsor_bonus + (member.Sale * level_percs_list[iterant] / 100)
	} else {
		sponsor_bonus = sponsor_bonus + level_percs_list[iterant]
	}
	if capping_scope["4"] && sponsor_bonus > capping_amount {
		parent.SponsorBonus = capping_amount
	} else {
		parent.SponsorBonus = sponsor_bonus
	}
	iterant = iterant + 1
	parent = parent.Parent
	apply_level_bonus(member, parent, level_percs_list, iterant, capping_amount, capping_scope, level_bonus_as_amount)
}

func set_get_uni_matching_bonus(matching_percs_list []float64, capping_amount float64, capping_scope map[string]bool, matching_bonus_as_amount bool) float64 {
	var total_bonus float64
	for _, member := range members {
		iterant := 0
		if member.Parent == nil {
			continue
		}
		parent := member.Parent
		apply_uni_matching_bonus(member, parent, matching_percs_list, iterant, capping_amount, capping_scope, matching_bonus_as_amount)
	}
	for _, member := range members {
		total_bonus = total_bonus + member.MatchingBonus
	}
	return total_bonus
}

func apply_uni_matching_bonus(member *Member, parent *Member, matching_perc_list []float64, iterant int, capping_amount float64, capping_scope map[string]bool, matching_bonus_as_amount bool) {
	if iterant >= len(matching_perc_list) || parent == nil {
		return
	}
	matching_bonus := parent.MatchingBonus
	if matching_bonus_as_amount != true {
		matching_bonus = matching_bonus + (member.SponsorBonus * matching_perc_list[iterant] / 100)
	} else {
		if member.SponsorBonus == 0 {
			matching_bonus = matching_bonus + 0
		} else {
			matching_bonus = matching_bonus + matching_perc_list[iterant]
		}
	}
	if capping_scope["2"] && matching_bonus > capping_amount {
		parent.MatchingBonus = capping_amount
	} else {
		parent.MatchingBonus = matching_bonus
	}
	iterant = iterant + 1
	parent = parent.Parent
	apply_uni_matching_bonus(member, parent, matching_perc_list, iterant, capping_amount, capping_scope, matching_bonus_as_amount)
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
		sponsor_bonus_as_amount := data["sponsor_bonus_as_amount"].(bool)

		binary_bonus_pairing_ratios := data["binary_bonus_pairing_ratios"].(map[string]interface{})
		binary_bonus_range := data["binary_bonus_range"].([]interface{})
		binary_bonus_as_amount := data["binary_bonus_as_amount"].(bool)

		matching_percs := data["matching_bonus_list"].([]interface{})
		matching_bonus_as_amount := data["matching_bonus_as_amount"].(bool)

		capping_amount := data["capping_amount"].(float64)
		rawCappingScope := data["capping_scope"].([]interface{})
		cappingScopeMap := make(map[string]bool)

		pool_perc := data["pool_bonus"].(float64)
		pool_distribution := data["pool_distribution"].(float64)

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
			var binary_bonus, matching_bonus, pool_bonus, sponsorBonus, revenue float64
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
			sponsorBonus = set_get_sponsor_bonus(sponsor_perc, capping_amount, cappingScopeMap, sponsor_bonus_as_amount)
			sponsorBonus = sponsorBonus - members[0].SponsorBonus
			binary_bonus = set_get_binary_bonus(binaryBonusPairingRatios, binaryBonusRange, capping_amount, cappingScopeMap, binary_bonus_as_amount)
			binary_bonus = binary_bonus - members[0].BinaryBonus
			matching_bonus = set_get_matching_bonus(matching_perc_list, capping_amount, cappingScopeMap, matching_bonus_as_amount)
			matching_bonus = matching_bonus - members[0].MatchingBonus
			for _, member := range members {
				revenue = revenue + member.Sale
				member.MatchingBonus = 0
				member.SponsorBonus = 0
			}
			expense := expense_per_member * float64(len(members))
			pool_bonus = set_get_pool_bonus(pool_perc, int(pool_distribution), expense, revenue, sponsorBonus, binary_bonus, matching_bonus)
			totalExpense = totalExpense + expense
			totalRevenue = totalRevenue + revenue
			totalBinaryBonus = totalBinaryBonus + binary_bonus
			totalSponsorBonus = totalSponsorBonus + sponsorBonus
			totalMatchingBonus = totalMatchingBonus + matching_bonus
			totalPoolBonus = totalPoolBonus + pool_bonus
			exportCycleData := ExportCycleData{
				NumberUsers:   len(members),
				Expense:       expense,
				Revenue:       revenue,
				Cycle:         cycle_num,
				BinaryBonus:   binary_bonus,
				SponsorBonus:  sponsorBonus,
				MatchingBonus: matching_bonus,
				PoolBonus:     pool_bonus,
			}
			cycleList = append(cycleList, exportCycleData)
		}
		exportData := ExportData{
			PlanType:           "binary",
			TotalExpense:       totalExpense,
			TotalRevenue:       totalRevenue,
			TotalCycles:        cycle_num,
			TotalBinaryBonus:   totalBinaryBonus,
			TotalSponsorBonus:  totalSponsorBonus,
			TotalMatchingBonus: totalMatchingBonus,
			TotalPoolBonus:     totalPoolBonus,
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
		level_bonus_as_amount := data["level_bonus_as_amount"].(bool)

		matching_percs := data["matching_bonus_list"].([]interface{})
		matching_bonus_as_amount := data["matching_bonus_as_amount"].(bool)

		capping_amount := data["capping_amount"].(float64)
		rawCappingScope := data["capping_scope"].([]interface{})
		cappingScopeMap := make(map[string]bool)

		pool_perc := data["pool_bonus"].(float64)
		pool_distribution := data["pool_distribution"].(float64)

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
		var totalExpense, totalRevenue, totalMatchingBonus, totalSponsorBonus, totalPoolBonus float64
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
			level_bonus = set_get_level_bonus(level_percs_list, capping_amount, cappingScopeMap, level_bonus_as_amount)
			level_bonus = level_bonus - members[0].SponsorBonus
			matching_bonus = set_get_uni_matching_bonus(matching_perc_list, capping_amount, cappingScopeMap, matching_bonus_as_amount)
			matching_bonus = matching_bonus - members[0].MatchingBonus
			var revenue float64
			for _, member := range members {
				revenue = revenue + member.Sale
				member.MatchingBonus = 0
				member.SponsorBonus = 0
			}
			expense := expense_per_member * float64(len(members))
			pool_bonus := set_get_pool_bonus(pool_perc, int(pool_distribution), expense, revenue, level_bonus, 0, matching_bonus)
			totalExpense = totalExpense + expense
			totalRevenue = totalRevenue + revenue
			totalSponsorBonus = totalSponsorBonus + level_bonus
			totalMatchingBonus = totalMatchingBonus + matching_bonus
			totalPoolBonus = totalPoolBonus + pool_bonus
			exportCycleData := ExportCycleData{
				NumberUsers:   len(members),
				Expense:       expense,
				Revenue:       revenue,
				Cycle:         cycle_num,
				SponsorBonus:  level_bonus,
				MatchingBonus: matching_bonus,
				PoolBonus:     pool_bonus,
			}
			cycleList = append(cycleList, exportCycleData)
		}
		exportData := ExportData{
			PlanType:           "unilevel",
			TotalExpense:       totalExpense,
			TotalRevenue:       totalRevenue,
			TotalCycles:        cycle_num,
			TotalSponsorBonus:  totalSponsorBonus,
			TotalMatchingBonus: totalMatchingBonus,
			TotalPoolBonus:     totalPoolBonus,
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
