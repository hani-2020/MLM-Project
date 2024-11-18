from django.shortcuts import render
from django.views import View
from .forms import BinaryForm
from .models import Members
from django.http import JsonResponse
import requests

class Member:

    def __init__(self, id, parent):
        self.id = id
        self.left_member = None
        self.right_member = None
        self.position = None
        self.parent = parent
        self.level = None
        ##################hits pythons recursions limit from 500 members:BEGIN##############################
        # self.left = None
        # self.right = None
        ##################hits pythons recursions limit from 500 members:END##############################
        self.sale = None
        self.sponsor_bonus = None
        self.binary_bonus = None
        self.left_sales = None
        self.right_sales = None
        self.carry_forward = None
        self.matching_bonus = None

class Tree:

    def __init__(self, num_members, package_price, additional_product_price):
        self.root = None
        self.num_members = num_members
        self.members = []
        self.build_tree()
        ##################hits pythons recursions limit from 500 members:BEGIN##############################
        # self.sum = 1
        # self.root.left = self.sum
        #self.assign_left_right(self.root)
        ##################hits pythons recursions limit from 500 members:END##############################
        self.set_member_sales(package_price, additional_product_price)

    def build_tree(self):
        if self.num_members <= 0:
            return
        self.root = Member(id=1, parent=None)
        self.root.level = 1
        self.members.append(self.root)
        queue = [self.root]
        current_id = 2 
        while current_id <= self.num_members:
            current_member = queue.pop(0)
            if current_id <= self.num_members:
                left_child = Member(id=current_id, parent=current_member)
                if left_child.level != left_child.parent.level + 1:
                    left_child.level = left_child.parent.level + 1
                left_child.position = 'Left'
                current_member.left_member = left_child
                queue.append(left_child)
                self.members.append(left_child)
                current_id += 1
            if current_id <= self.num_members:
                right_child = Member(id=current_id, parent=current_member)
                if right_child.level != right_child.parent.level + 1:
                    right_child.level = right_child.parent.level + 1
                right_child.position = 'Right'
                current_member.right_member = right_child
                queue.append(right_child)
                self.members.append(right_child)
                current_id += 1

    def set_member_sales(self, package_price, additional_product_price):
        for member in self.members:
            if member.id != 1:
                member.sale = package_price + (additional_product_price or 0)

    def set_and_get_sponsor_bonus(self, sponsor_percentage, capping_amount, capping_scope):
        total_bonus = 0
        for member in self.members:
            right_bonus = 0
            left_bonus = 0
            if member.right_member and member.right_member.sale:
                right_bonus = member.right_member.sale * sponsor_percentage/100
            if member.left_member and member.left_member.sale:
                left_bonus = member.left_member.sale * sponsor_percentage/100
            sponsor_bonus = right_bonus + left_bonus
            if '3' in capping_scope and sponsor_bonus>capping_amount:
                member.sponsor_bonus = capping_amount
            else:
                member.sponsor_bonus = sponsor_bonus
            total_bonus = total_bonus + member.sponsor_bonus
        return total_bonus

    def set_and_get_binary_bonus(self, binary_percentage, capping_amount, capping_scope):
        total_bonus = 0
        for member in self.members:
            left_sales = 0
            right_sales = 0
            if member.left_member:
                left_sales = self.traverse(member.left_member)
                member.left_sales = left_sales
            if member.right_member:
                right_sales = self.traverse(member.right_member)
                member.right_sales = right_sales
            binary_bonus = min(left_sales, right_sales) * binary_percentage/100
            if "1" in capping_scope and binary_bonus>capping_amount:
                member.binary_bonus = capping_amount
            else:
                member.binary_bonus = binary_bonus
            carrry_forward = left_sales - right_sales
            if member.left_member and carrry_forward>0:
                member.left_member.carry_forward = carrry_forward
            elif member.right_member and carrry_forward<0:
                member.right_member.carryforward = -1 * carrry_forward
            total_bonus = total_bonus + member.binary_bonus
        return total_bonus

    def traverse(self, node):
        if not node:
            return 0
        current_sales = node.sale if node.sale is not None else 0
        left_sales = self.traverse(node.left_member)
        right_sales = self.traverse(node.right_member)
        return current_sales + left_sales + right_sales
    
    def set_and_get_matching_bonus(self, matching_percentages, capping_amount, capping_scope):
        for member in self.members:
            iterant = 0
            if not member.parent:
                continue
            parent = member.parent
            if parent.matching_bonus is None:
                parent.matching_bonus = 0
            self.apply_matching_bonus(member, parent, matching_percentages, iterant, capping_amount, capping_scope)
        sum = 0
        for member in self.members:
            if member.matching_bonus is None:
                continue
            sum = sum + float(member.matching_bonus)
        return sum

    def apply_matching_bonus(self, member, parent, matching_percentages, iterant, capping_amount, capping_scope):
        if iterant >= len(matching_percentages) or parent is None:
            return
        matching_bonus = parent.matching_bonus
        if member.binary_bonus is None:
            member.binary_bonus = 0
        matching_bonus = matching_bonus + (member.binary_bonus*matching_percentages[iterant]/100)
        if "2" in capping_scope and parent.matching_bonus>capping_amount:
            parent.matching_bonus = capping_amount
        else:
            parent.matching_bonus = matching_bonus
        iterant = iterant + 1
        parent = parent.parent
        self.apply_matching_bonus(member, parent, matching_percentages, iterant, capping_amount, capping_scope)

    ##################hits pythons recursions limit from 500 members:BEGIN##############################
    # def assign_left_right(self, node):
    #     if not node.parent and node.left and node.right:
    #         return
    #     if not node.left:
    #         node.left = self.sum
    #     self.sum = self.sum + 1
    #     if node.left_member and not node.left_member.left:
    #         self.assign_left_right(node.left_member)
    #     if node.right_member and not node.right_member.left:
    #         self.assign_left_right(node.right_member)
    #     if not (node.left_member or node.right_member):
    #         node.right = self.sum
    #         self.assign_left_right(node.parent)
    #     if not node.right:
    #         node.right = self.sum
    #         if node.parent:
    #             self.assign_left_right(node.parent)
    ##################hits pythons recursions limit from 500 members:END##############################

    def display_tree(self):
        queue = [self.root]
        while queue:
            current_member = queue.pop(0)
            print(f'Member ID: {current_member.id}, ',
                  f'Sponsor Bonus: {current_member.sponsor_bonus if current_member.sponsor_bonus else None}, ',
                  f'Binary Bonus: {current_member.binary_bonus if current_member.binary_bonus else None}, '
                  f'Matching Bonus: {current_member.matching_bonus if current_member.matching_bonus else None}, '
                  )
            if current_member.left_member:
                queue.append(current_member.left_member)
            if current_member.right_member:
                queue.append(current_member.right_member)

class Calculator(View):

    template_name = 'calculator.html'
    
    def get(self, request, *args, **kwargs):
        form = BinaryForm()
        context = {'form':form}
        return render(request, self.template_name, context)
    
    def post(self, request, *args, **kwargs):
        form = BinaryForm(request.POST)
        if form.is_valid():
            number_of_users = form.cleaned_data['number_of_users']
            # joining_package_fee = form.cleaned_data['joining_package_fee']
            additional_product_price = form.cleaned_data['additional_product_price']

            product_names = form.cleaned_data['product_names']
            product_prices = form.cleaned_data['product_prices']
            product_quantities = form.cleaned_data['product_quantities']

            sponsor_bonus = form.cleaned_data['sponsor_bonus']

            binary_bonus_range = form.cleaned_data['binary_bonus_range']
            binary_bonus = form.cleaned_data['binary_bonus']
            binary_bonus_pairing_ratio = form.cleaned_data['binary_bonus_pairing_ratio']

            matching_bonus_string = form.cleaned_data['matching_bonus_per_level']

            capping_amount = form.cleaned_data['capping_amount']
            capping_scope = form.cleaned_data['capping_scope']
    
        product_names_list = product_names.split(", ") if product_names else []
        product_prices_list = [float(price) for price in product_prices.split(",")] if product_prices else []
        product_quantities_list = [int(quantity) for quantity in product_quantities.split(",")] if product_quantities else []
        products_catalogue = {}
        for i in range(len(product_names_list)):
            products_catalogue[product_names_list[i]] = {
                'price': product_prices_list[i],
                'quantity': product_quantities_list[i]
            }
        
        left, right = map(int, binary_bonus_pairing_ratio.split(":"))
        binary_bonus_pairing_ratio_dict = {"left": left, "right": right}
        binary_bonus_list = [float(perc) for perc in binary_bonus.split(",")]
        parts = binary_bonus_range.split(",")
        binary_bonus_dict = []
        for i in range(len(parts)):
            if "-" in parts[i]:
                min_val, max_val = map(int, parts[i].split("-"))
                binary_bonus_dict.append({"min": min_val, "max": max_val, "bonus": binary_bonus_list[i]})
            else:
                binary_bonus_dict.append({"min": int(parts[i]), "max": 10**100,"bonus": binary_bonus_list[i]})
    
        if matching_bonus_string:
            matching_bonus_list = [float(value) for value in matching_bonus_string.split(",")]
        else:
            matching_bonus_list = [0]

        if not capping_scope or not capping_amount:
            capping_amount = 10**100
        ##################go stuff:BEGIN#########################
        input = {
            'number_of_users': number_of_users,
            # 'joining_package_fee': joining_package_fee,
            'additional_product_price': additional_product_price,
            'product_order_list': product_names_list,
            'products_catalogue': products_catalogue,
            'sponsor_bonus': sponsor_bonus,
            'binary_bonus_pairing_ratios': binary_bonus_pairing_ratio_dict,
            'binary_bonus_range': binary_bonus_dict,
            'binary_bonus': binary_bonus,
            'matching_bonus_list': matching_bonus_list,
            'capping_amount': capping_amount,
            'capping_scope': capping_scope
        }
        self.send_to_go(input)
        ###################go stuff:END#####################
        # tree = Tree(number_of_users, joining_package_fee, additional_product_price)
        # sponsor_bonus = tree.set_and_get_sponsor_bonus(sponsor_bonus, capping_amount, capping_scope)
        # binary_bonus = tree.set_and_get_binary_bonus(binary_bonus, capping_amount, capping_scope)
        # matching_bonus = tree.set_and_get_matching_bonus(matching_bonus_list, capping_amount, capping_scope)
        # tree.display_tree()
        # self.store_in_db(tree.members)
        context = {}
        # context = {
        #     'sponsor_bonus':sponsor_bonus,
        #     'binary_bonus':binary_bonus,
        #     "matching_bonus":matching_bonus
        # }
        return render(request, 'result.html', context)
    
    def store_in_db(self, tree_members):
        Members.objects.all().delete()
        members_to_create = []
        for member in tree_members:
            members_to_create.append(
                Members(
                    user_id=member.id,
                    left_member=member.left_member.id if member.left_member else None,
                    right_member=member.right_member.id if member.right_member else None,
                    position=member.position,
                    parent=member.parent.id if member.parent else None,
                    level=member.level,
                    sale=member.sale,
                    sponsor_bonus=member.sponsor_bonus,
                    binary_bonus=member.binary_bonus,
                    left_sales=member.left_sales,
                    right_sales=member.right_sales,
                    carry_forward=member.carry_forward,
                    matching_bonus=member.matching_bonus
                )
            )
        Members.objects.bulk_create(members_to_create)

    def send_to_go(self, input):
        try:
            response = requests.post("http://localhost:8080", json=input)
            response_data = response.json()
            return JsonResponse(response_data)
        except requests.exceptions.RequestException as e:
            return JsonResponse({"error": str(e)}, status=500)