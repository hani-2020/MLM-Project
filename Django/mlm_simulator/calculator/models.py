from django.db import models

class Members(models.Model):
    user_id = models.IntegerField(primary_key=True)
    left_member = models.IntegerField(null=True, blank=True)
    right_member = models.IntegerField(null=True, blank=True)
    position = models.CharField(max_length=10, null=True, blank=True)
    parent = models.IntegerField(null=True, blank=True)
    level = models.IntegerField(null=True, blank=True)
    ##################hits pythons recursions limit from 500 members:BEGIN##############################
    # left = models.IntegerField(null=True, blank=True)
    # right = models.IntegerField(null=True, blank=True)
    ##################hits pythons recursions limit from 500 members:END##############################
    sale = models.FloatField(null=True, blank=True)
    sponsor_bonus = models.FloatField(null=True, blank=True)
    binary_bonus = models.FloatField(null=True, blank=True)
    left_sales = models.FloatField(null=True, blank=True)
    right_sales = models.FloatField(null=True, blank=True)
    carry_forward = models.FloatField(null=True, blank=True)
    matching_bonus = models.FloatField(null=True, blank=True)