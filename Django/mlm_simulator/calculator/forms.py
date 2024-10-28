from django import forms

unit_field = ((1, 'currency'), ('2', 'percentage'))

class BaseDetailCollectionForm(forms.Form):
    joining_package_fee = forms.DecimalField()
    additional_product_price = forms.DecimalField()
    number_of_users = forms.IntegerField(required=True)
    sponsor_bonus = forms.DecimalField()
    sponsor_bonus_unit = forms.ChoiceField(choices=unit_field)
    business_expenses_per_member = forms.DecimalField()
    business_expenses_unit = forms.ChoiceField(choices=unit_field)
    tax_per_person = forms.DecimalField()
    tax_unit = forms.ChoiceField(choices=unit_field)

class BinaryForm(BaseDetailCollectionForm):
    binary_bonus = forms.DecimalField()
    binary_bonus_unit = forms.ChoiceField(choices=unit_field)
    binary_bonus_pairing_ratio = forms.ChoiceField(choices=((1, '1:1'),(2, '2:2')))
    matching_bonus = forms.DecimalField()
    matching_bonus_unit = forms.ChoiceField(choices=unit_field)
    matching_bonus_levels = forms.IntegerField()
    capping_amount = forms.DecimalField()
    capping_scope = forms.ChoiceField(choices=((1, 'Binary bonus'),(2, 'Matching bonus'),(3,'Sponsor bonus')))