from django.shortcuts import render, redirect
from django.utils.decorators import method_decorator
from django.views.decorators.csrf import csrf_exempt
from django.views import View
from .forms import BinaryForm, UnilevelForm
from django.http import JsonResponse
import os
import json
from cryptography.fernet import Fernet
import requests
from urllib.parse import quote

empty_json = {}

def send_to_go(input, link, url):
        global empty_json
        domain = url.split(':')
        try:
            response = requests.post(domain[0] + ":" + domain[1] + ":8080/" + link + "/", json=input)
            response_data = response.json()
            empty_json = response_data
            return response_data
        except requests.exceptions.RequestException as e:
            return JsonResponse({"error": str(e)}, status=500)

def Home(request):
    return render(request, 'interface.html')

class BinaryCalculator(View):

    @method_decorator(csrf_exempt)
    def dispatch(self, *args, **kwargs):
        return super().dispatch(*args, **kwargs)

    def post(self, request, *args, **kwargs):
        inputData = {
            'number_of_users': int(request.POST.get('number_of_users')),
            'business_expenses_per_member': float(request.POST.get('business_expenses_per_member')) if request.POST.get('business_expenses_per_member')!="" else 0,
            'additional_product_price': float(request.POST.get('additional_product_price')),
            'product_prices': ','.join(request.POST.getlist('product_prices')).rstrip(','),
            'product_names': ','.join(request.POST.getlist('product_names')).rstrip(','),
            'product_quantities': ','.join(request.POST.getlist('product_quantities')).rstrip(','),
            'sponsor_bonus': float(request.POST.get('sponsor_bonus')),
            'binary_bonus_range': ','.join(request.POST.getlist('binary_bonus_range')),
            'binary_bonus': ','.join(request.POST.getlist('binary_bonus')),
            'binary_bonus_pairing_ratio': ','.join(request.POST.getlist('binary_bonus_pairing_ratio')),
            'matching_bonus_per_level': ','.join(request.POST.getlist('matching_bonus_per_level')),
            'capping_amount': request.POST.get('capping_amount'), 
            'pool_bonus': float(request.POST.get('pool_bonus')),
            'pool_distribution': int(request.POST.get('pool_distribution'))
        }
        capping_scope = request.POST.getlist('capping_scope')
        if capping_scope != ['']:
            inputData['capping_scope'] = list(map(int, capping_scope[0].split(',')))
        form = BinaryForm(inputData)
        if form.is_valid():
            number_of_users = form.cleaned_data['number_of_users']
            expenses_per_member = form.cleaned_data['business_expenses_per_member']
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

            pool_bonus = form.cleaned_data['pool_bonus']
            pool_distribution = form.cleaned_data['pool_distribution']
        product_names_list = product_names.split(",") if product_names else []
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
        input = {
            'number_of_users': number_of_users,
            'expenses_per_member': expenses_per_member,
            'additional_product_price': additional_product_price,
            'product_order_list': product_names_list,
            'products_catalogue': products_catalogue,
            'sponsor_bonus': sponsor_bonus,
            'binary_bonus_pairing_ratios': binary_bonus_pairing_ratio_dict,
            'binary_bonus_range': binary_bonus_dict,
            'matching_bonus_list': matching_bonus_list,
            'capping_amount': capping_amount,
            'capping_scope': capping_scope,
            'pool_bonus': pool_bonus,
            'pool_distribution': pool_distribution
        }
        url = request.build_absolute_uri()
        send_to_go(input, "binary-calc", url)
        return redirect("result")
        
class UnilevelCalculator(View):

    @method_decorator(csrf_exempt)
    def dispatch(self, *args, **kwargs):
        return super().dispatch(*args, **kwargs)
    
    def post(self, request, *args, **kwargs):
        inputData = {
            'number_of_users': int(request.POST.get('number_of_users')),
            'business_expenses_per_member': float(request.POST.get('business_expenses_per_member')) if request.POST.get('business_expenses_per_member')!="" else 0,
            'additional_product_price': float(request.POST.get('additional_product_price')),
            'product_prices': ','.join(request.POST.getlist('product_prices')).rstrip(','),
            'product_names': ','.join(request.POST.getlist('product_names')).rstrip(','),
            'product_quantities': ','.join(request.POST.getlist('product_quantities')).rstrip(','),
            'level_bonus_per_level': ','.join(request.POST.getlist('matching_bonus_per_level')),
            'matching_bonus_per_level': ','.join(request.POST.getlist('matching_bonus_per_level')),
            'capping_amount': request.POST.get('capping_amount'), 
            "downlines_per_user": int(request.POST.get('downlines_per_user')),
            'pool_bonus': float(request.POST.get('pool_bonus')),
            'pool_distribution': int(request.POST.get('pool_distribution'))
        }
        capping_scope = request.POST.getlist('capping_scope')
        if capping_scope != ['']:
            inputData['capping_scope'] = list(map(int, capping_scope[0].split(',')))
        form = UnilevelForm(inputData)
        if form.is_valid():
            number_of_users = form.cleaned_data['number_of_users']
            downlines_per_user = form.cleaned_data['downlines_per_user']
            expenses_per_member = form.cleaned_data['business_expenses_per_member']
            additional_product_price = form.cleaned_data['additional_product_price']

            product_names = form.cleaned_data['product_names']
            product_prices = form.cleaned_data['product_prices']
            product_quantities = form.cleaned_data['product_quantities']

            level_bonus_string = form.cleaned_data['level_bonus_per_level']

            matching_bonus_string = form.cleaned_data['matching_bonus_per_level']

            capping_amount = form.cleaned_data['capping_amount']
            capping_scope = form.cleaned_data['capping_scope']

            pool_bonus = form.cleaned_data['pool_bonus']
            pool_distribution = form.cleaned_data['pool_distribution']
        product_names_list = product_names.split(",") if product_names else []
        product_prices_list = [float(price) for price in product_prices.split(",")] if product_prices else []
        product_quantities_list = [int(quantity) for quantity in product_quantities.split(",")] if product_quantities else []
        products_catalogue = {}
        for i in range(len(product_names_list)):
            products_catalogue[product_names_list[i]] = {
                'price': product_prices_list[i],
                'quantity': product_quantities_list[i]
            }

        level_bonus_list = [float(value) for value in level_bonus_string.split(",")]
        if matching_bonus_string:
            matching_bonus_list = [float(value) for value in matching_bonus_string.split(",")]
        else:
            matching_bonus_list = [0]
        if not capping_scope or not capping_amount:
            capping_amount = 10**100

        input = {
            'number_of_users': number_of_users,
            'expenses_per_member': expenses_per_member,
            'downlines_per_user': downlines_per_user,
            'additional_product_price': additional_product_price,
            'product_order_list': product_names_list,
            'products_catalogue': products_catalogue,
            'level_bonus': level_bonus_list,
            'matching_bonus_list': matching_bonus_list,
            'capping_amount': capping_amount,
            'capping_scope': capping_scope,
            'pool_bonus': pool_bonus,
            'pool_distribution': pool_distribution
        }
        url = request.build_absolute_uri()
        send_to_go(input, "unilevel-calc", url)
        return redirect("result")
    
class Result(View):
    def get(self, request, *args, **kwargs):
        context = empty_json
        return render(request, 'result.html', context)