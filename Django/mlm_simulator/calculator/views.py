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
        self.left = None
        self.right = None

class Tree:

    def __init__(self, num_members):
        self.root = None
        self.num_members = num_members
        self.build_tree()
        self.sum = 1
        self.assign_left_right(self.root)

    def build_tree(self):
        if self.num_members <= 0:
            return
        self.root = Member(id=1, parent=None)
        queue = [self.root]
        current_id = 2 
        while current_id <= self.num_members:
            current_member = queue.pop(0)
            if current_id <= self.num_members:
                left_child = Member(id=current_id, parent=current_member)
                left_child.position = 'Left'
                current_member.left_member = left_child
                queue.append(left_child)
                current_id += 1
            if current_id <= self.num_members:
                right_child = Member(id=current_id, parent=current_member)
                right_child.position = 'Right'
                current_member.right_member = right_child
                queue.append(right_child)
                current_id += 1

    def assign_left_right(self, node):
        if node is None:
            return
        node.left = self.sum
        self.sum += 1
        if node.left_member:
            self.assign_left_right(node.left_member)
        node.right = self.sum
        self.sum += 1 
        if node.right_member:
            self.assign_left_right(node.right_member)
        

    def display_tree(self):
        queue = [self.root]
        while queue:
            current_member = queue.pop(0)
            print(f'Member ID: {current_member.id}, ',
                  f'Left member: {current_member.left_member.id if current_member.left_member else None}, ',
                  f'Right member: {current_member.right_member.id if current_member.right_member else None}, ',
                  f'left: {current_member.left if current_member.left else None}, ',
                  f'right: {current_member.right if current_member.right else None}, ',
                  f'Parent: {current_member.parent.id if current_member.parent else None}')
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
        number_of_users = int(request.POST.get('number_of_users'))
        tree = Tree(number_of_users)
        tree.display_tree()
        pass
        
            


