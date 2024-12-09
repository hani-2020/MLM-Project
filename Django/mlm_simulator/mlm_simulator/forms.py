from django import forms

class BaseDetailCollectionForm(forms.Form):
    additional_product_price = forms.FloatField(required=False, initial=0)
    number_of_users = forms.IntegerField(required=True)
    sponsor_bonus = forms.FloatField(required=False, initial=0)
    product_names = forms.CharField(max_length=255, required=False)
    product_quantities = forms.CharField(max_length=255, required=False)
    product_prices = forms.CharField(max_length=255, required=False)
    capping_amount = forms.FloatField(required=False, initial=0)
    capping_scope = forms.MultipleChoiceField(choices=((1, 'Binary bonus'), (2, 'Matching bonus'), (3, 'Sponsor bonus'), (4, 'Level bonus')),widget=forms.CheckboxSelectMultiple,required=False)
    business_expenses_per_member = forms.FloatField(required=False, initial=1)
    matching_bonus_per_level = forms.CharField(max_length=255, required=False)
    pool_bonus = forms.FloatField(required=False, initial=0)
    pool_distribution = forms.IntegerField(required=False, initial=1)

class BinaryForm(BaseDetailCollectionForm):
    binary_bonus_range = forms.CharField(max_length=255, required=False)
    binary_bonus = forms.CharField(max_length=255, required=False)
    binary_bonus_pairing_ratio = forms.CharField(max_length=255, required=False)

class UnilevelForm(BaseDetailCollectionForm):
    downlines_per_user = forms.IntegerField(required=False)
    level_bonus_per_level = forms.CharField(max_length=255, required=False)