from django.urls import path
from . import views

urlpatterns = [
    path('binary/', views.BinaryCalculator.as_view(), name='binary-calculator'),
    path('unilevel/', views.UnilevelCalculator.as_view(), name='unilevel-calculator')
]