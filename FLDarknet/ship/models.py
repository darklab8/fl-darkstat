from django.db import models

# Create your models here.
class Ship(models.Model):
    class Meta:
        verbose_name_plural = "ships"
    
    nickname = models.CharField(max_length=50, db_index=True, blank=True, null=True)

    ids_name = models.IntegerField(db_index=True, blank=True, null=True)
    ids_info = models.IntegerField(db_index=True, blank=True, null=True)

    name = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    infocard = models.CharField(max_length=500, db_index=True, blank=True, null=True)

    typeof = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    mass = models.FloatField(blank=True, null=True)
    hold_size = models.IntegerField(blank=True, null=True)
    linear_drag = models.FloatField(blank=True, null=True)
    max_bank_angle = models.IntegerField(db_index=True, blank=True, null=True)
    camera_angular_acceleration = models.FloatField(blank=True, null=True)
    camera_horizontal_turn_angle = models.IntegerField(blank=True, null=True)
    camera_vertical_turn_up_angle = models.IntegerField(blank=True, null=True)
    camera_vertical_turn_down_angle = models.IntegerField(blank=True, null=True)
    camera_turn_look_ahead_slerp_amount = models.FloatField(blank=True, null=True)
    hit_pts = models.IntegerField(blank=True, null=True)
    nudge_force = models.FloatField(blank=True, null=True)
    strafe_force = models.IntegerField(blank=True, null=True)
    strafe_power_usage = models.IntegerField(blank=True, null=True)
    bay_door_anim = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    bay_doors_open_snd = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    bay_doors_close_snd = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    HP_bay_surface = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    HP_bay_external = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    HP_tractor_source = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    explosion_resistance = models.FloatField(blank=True, null=True)
    ids_info1 = models.IntegerField(blank=True, null=True)
    ids_info2 = models.IntegerField(blank=True, null=True)
    ids_info3 = models.IntegerField(blank=True, null=True)
    ship_class = models.IntegerField(blank=True, null=True)
    nanobot_limit = models.IntegerField(blank=True, null=True)
    shield_battery_limit = models.IntegerField(blank=True, null=True)
    docking_camera = models.IntegerField(blank=True, null=True)
    distance_render = models.IntegerField(blank=True, null=True)
    nomad = models.BooleanField(blank=True, null=True)
    loadout = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    solar_radius = models.IntegerField(blank=True, null=True)
    shape_name = models.CharField(max_length=50, db_index=True, blank=True, null=True)
    destructible = models.BooleanField(blank=True, null=True)
    phantom_physics = models.BooleanField(blank=True, null=True)
        
