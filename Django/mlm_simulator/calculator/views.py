from django.shortcuts import render
from django.views import View
from .forms import BinaryForm

class Calculator(View):

    template_name = 'calculator.html'

    def get(self, request, *args, **kwargs):
        form = BinaryForm()
        context = {'form':form}
        return render(request, self.template_name, context)
    
    def post(self, request, *args, **kwargs):
        number_of_users = int(request.POST.get('number_of_users'))
        users_dict = {}
        for user_id in range(number_of_users):
            users_dict[user_id] = {'parent_id':None,
                                  'position':None,
                                  'bonus':None
                                }
        
            


