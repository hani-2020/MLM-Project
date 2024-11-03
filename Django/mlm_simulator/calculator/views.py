from django.shortcuts import render
from django.views import View
from .forms import BinaryForm
import math

class Member:

    def __init__(self, id, left_value, right_value, parent):
        self.id = id
        self.left_member = None
        self.right_member = None
        self.parent = parent
        self.left_value = left_value
        self.right_value = right_value
    
    def __str__(self):
        return f"Member(id={self.id}, range=({self.left_value}, {self.right_value}))"

    def __repr__(self):
        return self.__str__()

class Tree:

    def __init__(self):
        self.root = None

    def createTree(self, number_of_members):
        if number_of_members <= 0:
            return None
        if self.root is None:
            right_value = 2 * number_of_members
            self.root = Member(1, 1, right_value, None) 
            self._add_members(self.root, number_of_members - 1)
        return self.root

    def _add_members(self, node, remaining_members):
        if remaining_members <= 1:
            return
        left_id = 2 * node.id
        right_id = 2 * node.id + 1

        mid_value = (node.left_value + node.right_value) // 2
        node.left_member = Member(left_id, node.left_value + 1, mid_value, node)
        node.right_member = Member(right_id, mid_value + 1, node.right_value - 1, node)
        self._add_members(node.left_member, remaining_members // 2)
        self._add_members(node.right_member, remaining_members // 2)

    def __str__(self):
        return self._display_tree(self.root)

    def _display_tree(self, node, level=0):
        if node is None:
            return ""
        result = " " * (level * 4) + str(node) + "\n"
        result += self._display_tree(node.left_member, level + 1)
        result += self._display_tree(node.right_member, level + 1)
        return result

    

class Calculator(View):

    template_name = 'calculator.html'
    
    def get(self, request, *args, **kwargs):
        form = BinaryForm()
        context = {'form':form}
        return render(request, self.template_name, context)
    
    def post(self, request, *args, **kwargs):
        number_of_users = int(request.POST.get('number_of_users'))
        tree = Tree()
        tree.createTree(number_of_users)
        print(tree)
        pass
        
            


