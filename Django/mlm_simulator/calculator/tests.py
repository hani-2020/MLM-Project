from django.test import TestCase
from django.urls import reverse
from .forms import BinaryForm
from .views import Tree

class CalculatorViewTests(TestCase):

    def test_get_calculator_view(self):
        response = self.client.get(reverse('binary-calculator'))
        self.assertEqual(response.status_code, 200)
        self.assertTemplateUsed(response, 'calculator.html')
        self.assertIsInstance(response.context['form'], BinaryForm)

    def test_post_calculator_view(self):
        data = {
            'number_of_users': 15,
            'joining_package_fee': 900,
            'additional_product_price': 100,
            'sponsor_bonus': 10,
            'binary_bonus': 10,
            'matching_bonus_per_level': "10,5",
            'capping_amount': 10,
            'capping_scope': ["1","2","3"],
        }
        response = self.client.post(reverse('binary-calculator'), data)
        self.assertEqual(response.status_code, 200)
        self.assertTemplateUsed(response, 'result.html')
        self.assertIn('sponsor_bonus', response.context)
        self.assertIn('binary_bonus', response.context)
        self.assertIn('matching_bonus', response.context)
        self.assertIsInstance(response.context['sponsor_bonus'], float)
        self.assertIsInstance(response.context['binary_bonus'], float)
        self.assertIsInstance(response.context['matching_bonus'], float)


class TreeTests(TestCase):

    def setUp(self):
        self.tree = Tree(num_members=15, package_price=900, additional_product_price=100)

    def test_build_tree(self):
        self.assertEqual(len(self.tree.members), 15)
        self.assertEqual(self.tree.root.id, 1)

    def test_set_and_get_sponsor_bonus(self):
        sponsor_bonus = self.tree.set_and_get_sponsor_bonus(sponsor_percentage=10, capping_amount=10000, capping_scope="3")
        self.assertIsInstance(sponsor_bonus, float)
        self.assertEqual(sponsor_bonus, 1400)

    def test_set_and_get_binary_bonus(self):
        binary_bonus = self.tree.set_and_get_binary_bonus(binary_percentage=10, capping_amount=10000, capping_scope="1")
        self.assertIsInstance(binary_bonus, float)
        self.assertEqual(binary_bonus, 1700)

    def test_set_and_get_matching_bonus(self):
        self.tree.set_and_get_binary_bonus(binary_percentage=10, capping_amount=10000, capping_scope="1")
        matching_bonus = self.tree.set_and_get_matching_bonus(
            matching_percentages=[10, 5],
            capping_amount=10000,
            capping_scope="2"
        )
        self.assertIsInstance(matching_bonus, float)
        self.assertEqual(matching_bonus, 120)
    
