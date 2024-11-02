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

    def createTree(self, number_of_members, member=None, id=1, left_value=1, right_value=None):
        if number_of_members <= 0:
            return None
        if self.root is None:
            right_value = 2 * number_of_members
            self.root = Member(id, left_value, right_value, None)
            member = self.root
        else:
            member = Member(id, left_value, right_value, member.parent) 
        mid_value = (left_value + right_value) // 2
        if number_of_members > 1:
            member.left_member = self.createTree(
                number_of_members=number_of_members - 1,
                member=member,
                id=id * 2,
                left_value=left_value,
                right_value=mid_value
            )
        if number_of_members > 1:
            member.right_member = self.createTree(
                number_of_members=number_of_members - 1,
                member=member,
                id=id * 2 + 1,
                left_value=mid_value+1,
                right_value=right_value
            )
        return member

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
        pass
        
            


