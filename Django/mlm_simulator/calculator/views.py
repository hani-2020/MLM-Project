from django.shortcuts import render
from django.views import View
from .forms import BinaryForm

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

    def set_and_get_sponsor_bonus(self, sponsor_percentage):
        total_bonus = 0
        for member in self.members:
            right_bonus = 0
            left_bonus = 0
            if member.right_member and member.right_member.sale:
                right_bonus = member.right_member.sale * sponsor_percentage/100
            if member.left_member and member.left_member.sale:
                left_bonus = member.left_member.sale * sponsor_percentage/100
            member.sponsor_bonus = right_bonus + left_bonus
            total_bonus = total_bonus + member.sponsor_bonus
        return total_bonus

    def set_and_get_binary_bonus(self, binary_percentage):
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
            member.binary_bonus = min(left_sales, right_sales) * binary_percentage/100
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
    
    def set_and_get_matching_bonus(self, matching_percentages):
        for member in self.members:
            iterant = 0
            if not member.parent:
                continue
            parent = member.parent
            if parent.matching_bonus is None:
                parent.matching_bonus = 0
            self.apply_matching_bonus(member, parent, matching_percentages, iterant)
        sum = 0
        for member in self.members:
            if member.matching_bonus is None:
                continue
            sum = sum + member.matching_bonus
        return sum

    def apply_matching_bonus(self, member, parent, matching_percentages, iterant):
        if iterant >= len(matching_percentages) or parent is None:
            return
        parent.matching_bonus = parent.matching_bonus + (float(member.binary_bonus)*float(matching_percentages[iterant])/100)
        iterant = iterant + 1
        parent = parent.parent
        self.apply_matching_bonus(member, parent, matching_percentages, iterant)

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
                  f'carry forward: {current_member.carry_forward if current_member.carry_forward else None}, ',
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
            joining_package_fee = form.cleaned_data['joining_package_fee']
            additional_product_price = form.cleaned_data['additional_product_price']
            sponsor_bonus = form.cleaned_data['sponsor_bonus']
            binary_bonus = form.cleaned_data['binary_bonus']
            matching_bonus_string = form.cleaned_data['matching_bonus_per_level']
            matching_bonus_list = [float(value) for value in matching_bonus_string.split(",")]
        tree = Tree(number_of_users, joining_package_fee, additional_product_price)
        sponsor_bonus = tree.set_and_get_sponsor_bonus(sponsor_bonus)
        binary_bonus = tree.set_and_get_binary_bonus(binary_bonus)
        matching_bonus = tree.set_and_get_matching_bonus(matching_bonus_list)
        tree.display_tree()
        pass
        
            


