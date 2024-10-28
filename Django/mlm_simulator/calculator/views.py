from django.shortcuts import render
from django.views import View
from .forms import BinaryForm

class Calculator(View):

    template_name = 'calculator.html'

    def get(self, request, *args, **kwargs):
        form = BinaryForm()
        context = {'form':form}
        return render(request, self.template_name, context)
