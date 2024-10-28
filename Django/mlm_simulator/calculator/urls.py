from django.urls import path
from . import views

urlpatterns = [
    path('binary/', views.Calculator.as_view(), name='binary-calculator')
]