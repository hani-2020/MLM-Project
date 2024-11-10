from django import forms

#unit_field = ((1, 'currency'), ('2', 'percentage'))

class BaseDetailCollectionForm(forms.Form):
    joining_package_fee = forms.FloatField()
    additional_product_price = forms.FloatField(required=False, initial=0)
    number_of_users = forms.IntegerField(required=True)
    sponsor_bonus = forms.FloatField(required=False, initial=0)
    # sponsor_bonus_unit = forms.ChoiceField(choices=unit_field)
    # business_expenses_per_member = forms.FloatField()
    # business_expenses_unit = forms.ChoiceField(choices=unit_field)
    # tax_per_person = forms.FloatField()
    # tax_unit = forms.ChoiceField(choices=unit_field)

class BinaryForm(BaseDetailCollectionForm):
    binary_bonus = forms.FloatField(required=False, initial=0)
    # binary_bonus_unit = forms.ChoiceField(choices=unit_field)
    # binary_bonus_pairing_ratio = forms.ChoiceField(choices=((1, '1:1'),(2, '2:2')))
    matching_bonus_per_level = forms.CharField(max_length=255, required=False)
    # matching_bonus_unit = forms.ChoiceField(choices=unit_field)
    # matching_bonus_levels = forms.IntegerField()
    capping_amount = forms.FloatField(required=False, initial=0)
    capping_scope = forms.MultipleChoiceField(choices=((1, 'Binary bonus'), (2, 'Matching bonus'), (3, 'Sponsor bonus')),widget=forms.CheckboxSelectMultiple,required=False)