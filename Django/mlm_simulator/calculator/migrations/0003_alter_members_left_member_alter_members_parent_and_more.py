# Generated by Django 4.2.16 on 2024-11-10 19:50

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('calculator', '0002_rename_member_members'),
    ]

    operations = [
        migrations.AlterField(
            model_name='members',
            name='left_member',
            field=models.IntegerField(blank=True, null=True),
        ),
        migrations.AlterField(
            model_name='members',
            name='parent',
            field=models.IntegerField(blank=True, null=True),
        ),
        migrations.AlterField(
            model_name='members',
            name='right_member',
            field=models.IntegerField(blank=True, null=True),
        ),
    ]