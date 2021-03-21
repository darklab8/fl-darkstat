from django.apps import AppConfig

class MainConfig(AppConfig):
    name = 'main'
    def ready(self):
        # username = 'guest'
        # try:
        #     user = User.objects.get(username=username)
        # except User.DoesNotExist:
        #     user = User(username=username)
        #     u.set_password('guest')
        #     user.is_staff = True
        #     user.save()

        #user = authenticate(username='guest', password='guest')
        #return user
        pass
        #main_loop()