"model for table about space ships"
from django.db import models

# Create your models here.


class Ship(models.Model):
    "Any related stats about space ship information"
    class Meta:
        verbose_name_plural = "ships"

    nickname = models.CharField(
        max_length=50, db_index=True, blank=True, null=True)

    ids_name = models.IntegerField(db_index=True, blank=True, null=True)
    name = models.CharField(
        max_length=50, db_index=True, blank=True, null=True)

    ids_info = models.IntegerField(db_index=True, blank=True, null=True)

    info_name = models.CharField(
        max_length=100, db_index=True, blank=True, null=True)

    mass = models.FloatField(blank=True, null=True)
    hold_size = models.IntegerField(blank=True, null=True)
    linear_drag = models.FloatField(blank=True, null=True)
    max_bank_angle = models.IntegerField(db_index=True, blank=True, null=True)
    camera_angular_acceleration = models.FloatField(
        blank=True, null=True, verbose_name='ang speed')

    camera_horizontal_turn_angle = models.IntegerField(
        blank=True, null=True, verbose_name='hor ang')

    camera_vertical_turn_up_angle = models.IntegerField(
        blank=True, null=True, verbose_name='turn up')

    camera_vertical_turn_down_angle = models.IntegerField(
        blank=True, null=True, verbose_name='turn down')

    camera_turn_look_ahead_slerp_amount = models.FloatField(
        blank=True, null=True, verbose_name='look ahead')

    hit_pts = models.IntegerField(blank=True, null=True)
    nudge_force = models.FloatField(blank=True, null=True)
    strafe_force = models.IntegerField(blank=True, null=True)
    strafe_power_usage = models.IntegerField(
        blank=True, null=True, verbose_name='strafe usage')
    explosion_resistance = models.FloatField(
        blank=True, null=True, verbose_name='exp res')

    ids_info1 = models.IntegerField(blank=True, null=True)
    ids_info2 = models.IntegerField(blank=True, null=True)
    ids_info3 = models.IntegerField(blank=True, null=True)
    ship_class = models.IntegerField(blank=True, null=True)
    nanobot_limit = models.IntegerField(
        blank=True, null=True, verbose_name='nanobots')
    shield_battery_limit = models.IntegerField(
        blank=True, null=True, verbose_name='batteries')

    # different original name: type
    typeof = models.CharField(
        max_length=50, db_index=True, blank=True, null=True)
