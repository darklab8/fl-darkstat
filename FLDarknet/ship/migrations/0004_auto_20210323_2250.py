# Generated by Django 3.1.7 on 2021-03-23 15:50

from django.db import migrations


class Migration(migrations.Migration):

    dependencies = [
        ('ship', '0003_remove_ship_infocard'),
    ]

    operations = [
        migrations.RemoveField(
            model_name='ship',
            name='distance_render',
        ),
        migrations.RemoveField(
            model_name='ship',
            name='docking_camera',
        ),
        migrations.RemoveField(
            model_name='ship',
            name='solar_radius',
        ),
    ]